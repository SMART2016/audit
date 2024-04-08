package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
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
	json.NewEncoder(w).Encode("I am Healthy")
}

/*
*
POST http://localhost:9191/audit-service/v1/logevents

	Payload: {
	  "type":"user-service",
	  "es_query":{"query": {
	    "range": {
	      "time": {
	        "gte": "now-24h",
	        "lte": "now"
	      }
	    }
	  }
	  }
	}
*/
func submitQuery(w http.ResponseWriter, r *http.Request) {
	esClient := getNewElasticsearchClient()
	var query map[string]interface{}

	json.NewDecoder(r.Body).Decode(&query)
	fmt.Println("query", query)
	b, err := json.Marshal(query["es_query"].(map[string]interface{}))
	if err != nil {
		panic(err)
	}
	queryStr := string(b)

	claims, _ := getClaimsAndTokenFromAuthzHeader(r)
	queryType := query["type"]

	//First check if the user has permission to get data for a specific type
	permitted, AttributeFilterMap := checkAttributeAccess(claims.Role, queryType)
	if !permitted {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//Except admin all should be able to see events of there own no one elses.
	if !strings.EqualFold(claims.Role, ROLE_ADMIN) {
		AttributeFilterMap["CurrentUser"] = claims.Username
	}

	queryStr, err = generateAndFilterQuery(queryStr, AttributeFilterMap)
	resp, err := esClient.SubmitQuery(queryStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logRequestDetails(r, "AUDIT_SEARCH", "audit-service")
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}
