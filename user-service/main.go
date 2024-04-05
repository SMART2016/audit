package main

import (
	"log"
	"net/http"
)

var users = []User{}

func main() {
	log.Fatal(http.ListenAndServe(":8080", Router{}.getRoutes()))
}
