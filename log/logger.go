package log

import (
	mainlog "log"
	"nkn-server/config"
	"os"
)

var MyLog *mainlog.Logger
var UpdateLog *mainlog.Logger

var serverFile, updateFile *os.File

func Start() {
	logDir := config.DirRoot + "/logdir/"

	file := config.ServerLog
	if file == "" {
		file = "server.log"
	}

	serverFile, err := os.OpenFile(logDir+file, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		mainlog.Println(err)
	}

	MyLog = mainlog.New(serverFile, "SERVER\t", mainlog.Ldate|mainlog.Ltime|mainlog.Lmicroseconds|mainlog.Llongfile)
	MyLog.Println("Logger MyLog started")

	file = config.UpdateLog
	if file == "" {
		file = "update.log"
	}

	updateFile, erru := os.OpenFile(logDir+file, os.O_RDWR|os.O_CREATE, 0666)
	if erru != nil {
		mainlog.Println(erru)
	}

	UpdateLog = mainlog.New(updateFile, "UPDATE\t", mainlog.Ldate|mainlog.Ltime|mainlog.Lmicroseconds|mainlog.Llongfile)
	UpdateLog.Println("Logger UpdateLog started")
}

func Stop() {
	serverFile.Close()
	updateFile.Close()
}
