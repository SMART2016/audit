package main

import (
	"log"
	"net/http"
	"time"
)

func logRequestDetails(r *http.Request, action string, serviceId string) {
	log.Printf("CurrentUser: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s\n",
		"Dipanjan",
		serviceId,
		action,
		r.RemoteAddr,
		r.UserAgent(),
		time.Now().Format(time.RFC3339))
}
