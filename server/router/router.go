package router

import (
	"fmt"
	"log"
	"net/http"
	_ "status-page-monitor/internal/database"
	_ "status-page-monitor/internal/database/models"
	"status-page-monitor/server/handlers"
	"status-page-monitor/server/middleware"

	"github.com/gorilla/mux"
)

func InitRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/", helloWorldHandler).Methods("GET")

	apiRouter := r.PathPrefix("/api").Subrouter()

	adminApiRouter := apiRouter.PathPrefix("/admin").Subrouter()
	publicApiRouter := apiRouter.PathPrefix("/").Subrouter()

	// GET api/admin/servers
	adminApiRouter.HandleFunc("/servers", middleware.Chain(handlers.GetAllServersHandler, middleware.Auth(), middleware.Logging())).Methods("GET")

	// GET api/admin/servers/:id
	adminApiRouter.HandleFunc("/servers/{id}", middleware.Chain(handlers.GetServerById, middleware.Auth(), middleware.Logging())).Methods("GET")

	// POST api/admin/servers
	adminApiRouter.HandleFunc("/servers", middleware.Chain(handlers.CreateServer, middleware.Auth(), middleware.Logging())).Methods("POST")

	// PUT api/admin/servers/:id
	adminApiRouter.HandleFunc("/servers/{id}", middleware.Chain(handlers.EditServer, middleware.Auth(), middleware.Logging())).Methods("PUT")

	//POST api/login
	publicApiRouter.HandleFunc("/login", middleware.Chain(handlers.LoginHandler, middleware.Logging()))

	log.Fatal(http.ListenAndServe(":8000", r))
}

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, "Hello world")
}

func PlugHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Plug")
}
