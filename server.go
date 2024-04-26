package main

import (
	"log"
	"net/http"

	"github.com/danilomarques1/fidusserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type FidusServer struct {
	router chi.Router
	port   string
}

func NewFidusServer(port string) *FidusServer {
	return &FidusServer{router: chi.NewRouter(), port: port}
}

func (f *FidusServer) Start() error {
	f.router.Use(middleware.Logger)

	f.router.Post("/fidus/master/register", handlers.CreateMaster)
	f.router.Post("/fidus/master/authenticate", handlers.AuthenticateMaster)

	log.Printf("Server running at %v\n", f.port)
	return http.ListenAndServe(f.port, f.router)
}
