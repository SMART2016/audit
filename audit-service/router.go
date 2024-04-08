package main

import (
	"github.com/gorilla/mux"
)

type Router struct{}

func (r Router) getRoutes() *mux.Router {
	router := mux.NewRouter()
	router.Use(logResponseMiddleware)
	searchEventsRouter := router.PathPrefix("/audit-service/v1").Subrouter()
	searchEventsRouter.HandleFunc("/logevents", AuthMiddleware(LoggingMiddleware(QueryValidatorMiddleware(SubmitQuery)))).Methods("POST")
	//searchEventsRouter.HandleFunc("/unsafe/logevents", QueryValidatorMiddleware(SubmitQuery)).Methods("POST")
	searchEventsRouter.HandleFunc("/Health", AuthMiddleware(LoggingMiddleware(Health))).Methods("GET")

	authRouter := router.PathPrefix("/auth-service/v1").Subrouter()
	authRouter.HandleFunc("/register", AuthMiddleware(LoggingMiddleware(RegisterNewUserHandler))).Methods("POST")
	//authRouter.HandleFunc("/unsafe/register", RegisterNewUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", LoginLoggerMiddleware(LoginHandler)).Methods("POST")
	authRouter.HandleFunc("/Health", AuthMiddleware(LoggingMiddleware(Health))).Methods("GET")
	authRouter.HandleFunc("/users", AuthMiddleware(LoggingMiddleware(GetUsersHandler))).Methods("GET")

	return router
}
