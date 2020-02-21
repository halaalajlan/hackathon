package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/halaalajlan/hackathon/logger"
	mid "github.com/halaalajlan/hackathon/middleware"
)

type Server struct {
	server *http.Server
}

func NewServer() *Server {
	defaultServer := &http.Server{
		ReadTimeout: 10 * time.Second,
		Addr:        "127.0.0.1:3333",
	}
	as := &Server{server: defaultServer}
	as.registerRoutes()
	return as
}
func (as *Server) Start() {

	log.Infof("Starting  server at http://%s", "127.0.0.1:3333")
	log.Fatal(as.server.ListenAndServe())
}

// Shutdown attempts to gracefully shutdown the server.
func (as *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return as.server.Shutdown(ctx)
}

func (as *Server) registerRoutes() {
	router := mux.NewRouter()

	router.Use(mid.RequireAPIKey)
	router.HandleFunc("/api/login", as.Login)
	router.HandleFunc("/api/patient", as.Patients)
	handler := handlers.CombinedLoggingHandler(log.Writer(), router)
	as.server.Handler = handler
}
