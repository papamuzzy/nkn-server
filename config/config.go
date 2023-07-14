package config

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var IsDebug bool
var IsTest bool
var DirRoot string
var MongoUri string
var MongoBase string
var MongoCollection string
var Ip string
var SshUser string
var SshPassw string
var ServerLog string
var UpdateLog string
var ServerAddr string

func Start() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Some error occured. Err: %s", err)
	}

	IsDebug = os.Getenv("DEBUG") == "1"
	IsTest = os.Getenv("TEST") == "1"

	ServerLog = os.Getenv("SERVER_LOG")
	UpdateLog = os.Getenv("UPDATE_LOG")

	MongoUri = os.Getenv("MONGO_URI")
	MongoBase = os.Getenv("MONGO_BASE")
	MongoCollection = os.Getenv("MONGO_COLLECTION")

	SshUser = os.Getenv("SSH_USER")
	SshPassw = os.Getenv("SSH_PASSWORD")

	Ip = os.Getenv("IP")

	/*minProfit, _ = strconv.ParseFloat(os.Getenv("MIN_PROFIT"), 64)
	maxLevel, _ = strconv.Atoi(os.Getenv("MAX_LEVEL"))
	fiat = strings.Split(os.Getenv("FIAT"), ",")*/

	getRoot()
	getFlags()
}

func getRoot() {
	if IsDebug {
		DirRoot, _ = os.Getwd()
		fmt.Println(DirRoot)
	} else {
		_, callerFile, _, _ := runtime.Caller(0)
		DirRoot = filepath.Dir(callerFile)
	}
}

func getFlags() {
	addr := flag.String("addr", ":9999", "Сетевой адрес веб-сервера")

	flag.Parse()

	ServerAddr = *addr
	//currencies = flag.Args()
}
