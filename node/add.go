package node

import (
	"context"
	"fmt"
	mainlog "log"
	"nkn-server/block"
	"nkn-server/script"
	"nkn-server/xtime"
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
	altNode := false

	block.NodesMutex.Lock()
	generationId := GetGenerationId()
	if generationId > GetGenerationsCount() {
		altNode = true
	} else {
		NewNode(ip, generationId)
	}
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

	var strScript string
	if altNode {
		strScript = script.GetString("/1/alt.sh", generationId, makeLog)
	} else {
		strScript = script.GetString("/1/install.sh", generationId, makeLog)
	}

	if config.NodeNum > 1 && !altNode {
		strScript += script.GetString("/1/add.sh", generationId, makeLog)
	}

	output, err := session.CombinedOutput(strScript)
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
	now := xtime.ToStr(time.Now())
	newNode := DataNode{
		ID:              primitive.NewObjectID(),
		Ip:              ip,
		Generation:      generationId,
		Height:          0,
		Version:         "-",
		WorkTime:        "-",
		MinedEver:       0,
		MinedToday:      0,
		NodeStatus:      "OFFLINE",
		LastBlockNumber: 0,
		Created:         now,
		LastUpdate:      now,
		LastOfflineTime: now,
	}

	result, err := db.NodeCollection.InsertOne(context.TODO(), newNode)
	if err != nil {
		log.MyLog.Println(err)
	}

	log.MyLog.Println(result)
}
