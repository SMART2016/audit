package main

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
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

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello User")
	logRequestDetails(r, "AUDIT_HEALTH", "audit-service")
	json.NewEncoder(w).Encode("I am Healthy")
}

func submitQuery(w http.ResponseWriter, r *http.Request) {
	esClient := getNewElasticsearchClient()

	var query map[string]interface{}
	json.NewDecoder(r.Body).Decode(&query)
	b, err := json.Marshal(query)
	if err != nil {
		panic(err)
	}

	queryStr := string(b)
	resp, err := esClient.submitQuery(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Print(string(resp))
	logRequestDetails(r, "AUDIT_SEARCH", "audit-service")
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func createSearchPayload(w http.ResponseWriter, queryValues url.Values) ([]byte, error, bool) {
	fmt.Println(queryValues)
	requestBody, err := json.Marshal("")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil, nil, true
	}
	return requestBody, err, false
}
