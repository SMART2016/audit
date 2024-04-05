package main

import (
	"log"
	"net/http"
	"time"
)

func logRequestDetails(r *http.Request, action string) {
	log.Printf("Action: %s, IP: %s, Agent: %s, Time: %s\n",
		action,
		r.RemoteAddr,
		r.UserAgent(),
		time.Now().Format(time.RFC3339))
}
