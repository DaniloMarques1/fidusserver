package main

import (
	"net/http"

	"github.com/danilomarques1/fidusserver/handlers"
)

type FidusServer struct {
	mux  *http.ServeMux
	port string
}

func NewFidusServer(port string) *FidusServer {
	return &FidusServer{mux: http.NewServeMux(), port: port}
}

func (f *FidusServer) Start() error {
	f.mux.HandleFunc("POST /fidus/master/register", handlers.CreateMaster)

	f.mux.HandleFunc("PUT /auth", func(w http.ResponseWriter, r *http.Request) {
		// TODO: validation and return a jwt token
	})

	return http.ListenAndServe(f.port, f.mux)
}
