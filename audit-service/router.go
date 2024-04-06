package main

import "github.com/gorilla/mux"

type Router struct{}

func (r Router) getRoutes() *mux.Router {
	router := mux.NewRouter()
	searchEventsRouter := router.PathPrefix("/audit-service/v1").Subrouter()
	searchEventsRouter.HandleFunc("/logevents", submitQuery).Methods("POST")
	searchEventsRouter.HandleFunc("/health", health).Methods("GET")
	return router
}
