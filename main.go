package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"nkn-server/config"
	"nkn-server/db"
	"nkn-server/handlers"
	"nkn-server/log"
	"nkn-server/node"
	"nkn-server/node2"
)

func main() {
	config.Start()
	log.Start()
	defer log.Stop()

	db.Start()
	defer db.Stop()

	go log.Timer()

	go node.UpdateBase()
	if config.NodeNum > 1 {
		go node2.UpdateBase()
	}

	initRouter()
}

func initRouter() {
	log.MyLog.Println("Server will start at http://localhost:9999/")

	router := mux.NewRouter()

	if config.NodeNum > 1 {
		node2Router := router.PathPrefix("/node2").Subrouter()
		//node2Router.HandleFunc("", handlers.Nodes2Get).Methods("GET")
		node2Router.HandleFunc("/add", handlers.Node2Add).Methods("POST")
		node2Router.HandleFunc("/make", handlers.Node2Make).Methods("POST")
		node2Router.HandleFunc("/delete", handlers.Node2Delete).Methods("POST")
	}

	nodeRouter := router.PathPrefix("/node").Subrouter()
	nodeRouter.HandleFunc("", handlers.NodesGet).Methods("GET")
	nodeRouter.HandleFunc("/add", handlers.NodeAdd).Methods("POST")
	nodeRouter.HandleFunc("/make", handlers.NodeMake).Methods("POST")
	nodeRouter.HandleFunc("/delete", handlers.NodeDelete).Methods("POST")

	genRouter := router.PathPrefix("/generation").Subrouter()
	genRouter.HandleFunc("/count", handlers.GenerationCount).Methods("GET")

	fs := http.FileServer(http.Dir("./public/"))
	router.PathPrefix("/generations/").Handler(fs)

	router.PathPrefix("/").HandlerFunc(handlers.CatchAll)

	srv := &http.Server{
		Addr:     config.ServerAddr,
		ErrorLog: log.MyLog,
		Handler:  router,
	}

	err := srv.ListenAndServe()
	log.MyLog.Fatal(err)
}
