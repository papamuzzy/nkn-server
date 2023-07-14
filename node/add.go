package node

import (
	"context"
	"fmt"
	mainlog "log"
	"nkn-server/block"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/ssh"
	"nkn-server/config"
	"nkn-server/db"
	"nkn-server/log"
)

func Add(ip string, generationId int) {
	block.NodesMutex.Lock()
	defer block.NodesMutex.Unlock()

	filter := bson.D{{"ip", ip}}
	count, err := db.NodeCollection.CountDocuments(context.TODO(), filter)
	if err != nil {
		log.MyLog.Println(err)
	}

	if count > 0 {
		log.MyLog.Println("there's such ip")
	} else {
		log.MyLog.Println("there no node with such ip")
		NewNode(ip, generationId)
	}
}

func Make(ip string) {
	file := config.DirRoot + "/logdir/Make_" + time.Now().Format("2006-01-02_15-04-05") + ".log"

	mFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		mainlog.Println(err)
	}
	defer mFile.Close()

	makeLog := mainlog.New(mFile, "SERVER\t", mainlog.Ldate|mainlog.Ltime|mainlog.Lmicroseconds|mainlog.Llongfile)
	makeLog.Println("Logger MyLog started, IP ", ip)

	start := time.Now()

	block.NodesMutex.Lock()
	generationId := GetGenerationId()
	NewNode(ip, generationId)
	block.NodesMutex.Unlock()

	conf := &ssh.ClientConfig{
		User: config.SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SshPassw),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:22", ip), conf)
	if err != nil {
		makeLog.Println(err)
	}
	makeLog.Println("SSH client -- YES!")

	session, err := client.NewSession()
	if err != nil {
		makeLog.Println(err)
	}
	makeLog.Println("SSH session -- YES!")

	script := fmt.Sprintf(`
		#!/bin/bash
		apt update -y
		apt purge needrestart -y
		apt-mark hold linux-image-generic linux-headers-generic openssh-server snapd
		apt upgrade -y
		apt -y install unzip vnstat htop screen mc
		
		username="nkn"
		benaddress="NKNKKevYkkzvrBBsNnmeTVf2oaTW3nK6Hu4K"
		config="https://nknrus.ru/config.tar"
		keys="http://5.180.183.19:9999/generations/%d.tar"
		
		useradd -m -p "pass" -s /bin/bash "$username" > /dev/null 2>&1
		usermod -a -G sudo "$username" > /dev/null 2>&1
		
		printf "Downloading........................................... "
		cd /home/$username > /dev/null 2>&1
		wget --quiet --continue --show-progress https://commercial.nkn.org/downloads/nkn-commercial/linux-amd64.zip > /dev/null 2>&1
		printf "DONE!\n"
		
		printf "Installing............................................ "
		unzip linux-amd64.zip > /dev/null 2>&1
		mv linux-amd64 nkn-commercial > /dev/null 2>&1
		chown -c $username:$username nkn-commercial/ > /dev/null 2>&1
		/home/$username/nkn-commercial/nkn-commercial -b $benaddress -d /home/$username/nkn-commercial/ -u $username install > /dev/null 2>&1
		printf "DONE!\n"
		printf "sleep 180"
		
		sleep 180
		
		DIR="/home/$username/nkn-commercial/services/nkn-node/"
		
		systemctl stop nkn-commercial.service > /dev/null 2>&1
		sleep 20
		cd $DIR > /dev/null 2>&1
		rm wallet.json > /dev/null 2>&1
		rm wallet.pswd > /dev/null 2>&1
		rm config.json > /dev/null 2>&1
		rm -Rf ChainDB > /dev/null 2>&1
		wget -O - "$keys" -q --show-progress | tar -xf -
		wget -O - "$config" -q --show-progress | tar -xf -
		chown -R $username:$username wallet.* > /dev/null 2>&1
		chown -R $username:$username config.* > /dev/null 2>&1
		printf "Downloading.......................................... DONE!\n"
		systemctl start nkn-commercial.service > /dev/null 2>&1
	`, generationId)

	output, err := session.CombinedOutput(script)
	if err != nil {
		makeLog.Println(err)
	}
	makeLog.Println(string(output))

	client.Close()
	session.Close()

	makeLog.Println("ran script successfully")

	total := time.Now().Sub(start).Seconds()

	makeLog.Println("Create node IP ", ip, " total time ", total, " sec")
}

func NewNode(ip string, generationId int) {
	newNode := DataNode{
		ID:              primitive.NewObjectID(),
		Ip:              ip,
		Generation:      generationId,
		Height:          "-",
		Version:         "-",
		WorkTime:        "-",
		MinedEver:       0,
		MinedToday:      0,
		NodeStatus:      "OFFLINE",
		LastBlockNumber: 0,
		LastUpdate:      "-",
		LastOfflineTime: time.Now().Format("2006-01-02 15:04:05"),
	}

	result, err := db.NodeCollection.InsertOne(context.TODO(), newNode)
	if err != nil {
		log.MyLog.Println(err)
	}

	log.MyLog.Println(result)
}
