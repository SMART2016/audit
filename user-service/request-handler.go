package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	service_id = "user-service"
)

var users = []User{}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)
	users = append(users, user)
	logRequestDetails(r, "CREATE", service_id)
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updatedUser User
	json.NewDecoder(r.Body).Decode(&updatedUser)

	for i, user := range users {
		if user.ID == id {
			users[i] = updatedUser
			logRequestDetails(r, "UPDATE", service_id)
			json.NewEncoder(w).Encode(updatedUser)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	for i, user := range users {
		if user.ID == id {
			users = append(users[:i], users[i+1:]...)
			logRequestDetails(r, "DELETE", service_id)
			w.WriteHeader(http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	for _, user := range users {
		if user.ID == id {
			logRequestDetails(r, "READ", service_id)
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func health(w http.ResponseWriter, r *http.Request) {
	logRequestDetails(r, "HEALTH", service_id)
	json.NewEncoder(w).Encode("I am Healthy")
}
