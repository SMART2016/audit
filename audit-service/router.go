package main

import (
	"github.com/gorilla/mux"
)

type Router struct{}

func (r Router) getRoutes() *mux.Router {
	router := mux.NewRouter()
	//router.Use(logResponseMiddleware)
	searchEventsRouter := router.PathPrefix("/audit-service/v1").Subrouter()
	searchEventsRouter.HandleFunc("/logevents", LoggingMiddleware(AuthorizationMiddleware(QueryValidatorMiddleware(SubmitQuery)))).Methods("POST")
	//searchEventsRouter.HandleFunc("/unsafe/logevents", QueryValidatorMiddleware(SubmitQuery)).Methods("POST")
	searchEventsRouter.HandleFunc("/health", LoggingMiddleware(AuthorizationMiddleware(Health))).Methods("GET")

	authRouter := router.PathPrefix("/auth-service/v1").Subrouter()
	authRouter.HandleFunc("/register", LoggingMiddleware(AuthorizationMiddleware(RegisterNewUserHandler))).Methods("POST")
	//authRouter.HandleFunc("/unsafe/register", RegisterNewUserHandler).Methods("POST")
	authRouter.HandleFunc("/login", BasicAuthMiddleware(LoginHandler)).Methods("POST")
	authRouter.HandleFunc("/health", LoggingMiddleware(AuthorizationMiddleware(Health))).Methods("GET")
	authRouter.HandleFunc("/users", LoggingMiddleware(AuthorizationMiddleware(GetUsersHandler))).Methods("GET")

	return router
}
