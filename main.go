package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"nkn-server/config"
	"nkn-server/db"
	"nkn-server/handlers"
	"nkn-server/log"
)

func main() {
	config.Start()
	log.Start()
	defer log.Stop()

	db.Start()
	defer db.Stop()

	initRouter()
}

/*
Routes:
router.GET("/", handlers.PermissionDenied)
router.POST("/", handlers.NodeIpPOST(db))

router.POST("/delete", handlers.Delete(db))

router.GET("/api", handlers.ApiGET(db))
router.GET("/usage", handlers.GetGenerationNumber)

router.GET("/generations/:fileName", handlers.GetGeneration)

router.GET("/my-nodes", handlers.MyNodesGET) html-page
*/
func initRouter() {
	log.MyLog.Println("Server will start at http://localhost:9999/")

	router := mux.NewRouter()
	nodeRouter := router.PathPrefix("/node").Subrouter()
	nodeRouter.HandleFunc("", handlers.NodesGet).Methods("GET")
	nodeRouter.HandleFunc("/add", handlers.NodeAdd).Methods("POST")
	nodeRouter.HandleFunc("/make", handlers.NodeMake).Methods("POST")
	nodeRouter.HandleFunc("/{id:[0-9]+}", handlers.NodeDelete).Methods("DELETE")

	genRouter := router.PathPrefix("/generation").Subrouter()
	genRouter.HandleFunc("/count", handlers.GenerationCount).Methods("GET")

	fs := http.FileServer(http.Dir("./public/"))
	router.PathPrefix("/generations/").Handler(fs)

	router.PathPrefix("/").HandlerFunc(handlers.CatchAll)

	/*http.Handle("/", ghandlers.CombinedLoggingHandler(os.Stdout, http.FileServer(http.Dir("."))))
	log.Fatal(http.ListenAndServe(":9999", loggingHandler(router)))*/

	srv := &http.Server{
		Addr:     config.ServerAddr,
		ErrorLog: log.MyLog,
		Handler:  router,
	}

	err := srv.ListenAndServe()
	log.MyLog.Fatal(err)
}

/*func loggingHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		h.ServeHTTP(w, r)
	})
}*/
