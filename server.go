package main

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/danilomarques1/fidusserver/apierror"
	"github.com/danilomarques1/fidusserver/handlers"
	"github.com/danilomarques1/fidusserver/response"
	"github.com/danilomarques1/fidusserver/token"
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

	f.router.Group(func(passwordRouter chi.Router) {
		passwordRouter.Use(AuthMiddleware)
		passwordRouter.Post("/fidus/password/store", handlers.StorePassword)
	})

	log.Printf("Server running at %v\n", f.port)
	return http.ListenAndServe(f.port, f.router)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if len(authHeader) == 0 {
			response.Error(w, apierror.Forbidden())
			return
		}
		splitedAuthHeader := strings.Split(authHeader, " ")
		if len(splitedAuthHeader) < 2 {
			response.Error(w, apierror.Forbidden())
			return
		}
		tokenStr := splitedAuthHeader[1]
		masterId, err := token.ParseToken(tokenStr)
		if err != nil {
			response.Error(w, apierror.Forbidden())
			return
		}

		ctx := context.WithValue(r.Context(), "masterId", masterId) // TODO: get the master id
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
