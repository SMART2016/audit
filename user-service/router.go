package main

import "github.com/gorilla/mux"

type Router struct{}

func (r Router) getRoutes() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	router.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	router.HandleFunc("/user/{id}", getUser).Methods("GET")
	router.HandleFunc("/health", health).Methods("GET")

	return router
}
