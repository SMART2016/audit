package main

import "github.com/gorilla/mux"

type Router struct{}

func (r Router) getRoutes() *mux.Router {
	router := mux.NewRouter()
	userRouter := router.PathPrefix("/user-service/v1").Subrouter()
	userRouter.HandleFunc("/users", createUser).Methods("POST")
	userRouter.HandleFunc("/users/{id}", updateUser).Methods("PUT")
	userRouter.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	userRouter.HandleFunc("/users/{id}", getUser).Methods("GET")
	userRouter.HandleFunc("/health", health).Methods("GET")

	return router
}
