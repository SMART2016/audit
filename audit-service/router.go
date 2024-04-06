package main

import (
	"github.com/gorilla/mux"
)

type Router struct{}

func (r Router) getRoutes() *mux.Router {
	router := mux.NewRouter()
	searchEventsRouter := router.PathPrefix("/audit-service/v1").Subrouter()
	searchEventsRouter.HandleFunc("/logevents", AuthMiddleware(QueryValidatorMiddleware(submitQuery))).Methods("POST")
	searchEventsRouter.HandleFunc("/unsafe/logevents", QueryValidatorMiddleware(submitQuery)).Methods("POST")
	searchEventsRouter.HandleFunc("/health", AuthMiddleware(health)).Methods("GET")

	authRouter := router.PathPrefix("/auth-service/v1").Subrouter()
	authRouter.HandleFunc("/register", AuthMiddleware(RegisterNewUserHandler)).Methods("POST")
	authRouter.HandleFunc("/login", AuthMiddleware(LoginHandler)).Methods("POST")
	authRouter.HandleFunc("/health", AuthMiddleware(health)).Methods("GET")

	return router
}
