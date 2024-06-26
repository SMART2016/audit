package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// loggingMiddleware checks the JWT token and authorizes users
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Wrap the standard ResponseWriter with our custom writer
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 OK if not set explicitly
		}

		claims, _ := getClaimsAndTokenFromAuthzHeader(r)
		serviceId := getServiceId(r.RequestURI)
		requestID := getRequestId()
		logMsg := fmt.Sprintf("RequestId: %s, CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: Initiated\n",
			requestID,
			claims.Username,
			claims.Role,
			serviceId,
			r.Method+":"+r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Now().Format(time.RFC3339))
		//fmt.Println(logMsg)
		publishEventLogs(serviceId, logMsg)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(r.Context(), "requestId", requestID)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(wrappedWriter, r.WithContext(ctx))

		logMsg = fmt.Sprintf("RequestId: %s, CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: %d.\n",
			requestID,
			claims.Username,
			claims.Role,
			serviceId,
			r.Method+":"+r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Now().Format(time.RFC3339),
			wrappedWriter.statusCode,
		)
		publishEventLogs(serviceId, logMsg)
	}
}

// QueryValidatorMiddleware checks the JWT token and authorizes users
func QueryValidatorMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		var query map[string]interface{}
		if err = json.Unmarshal(body, &query); err != nil {
			fmt.Println(err)
		}
		if _, ok := query["type"]; !ok {
			fmt.Println("type missing in payload")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if _, ok := query["es_query"]; !ok {
			fmt.Println("es_query missing in payload")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)

	}
}

// Extracts service name from request URL
func getServiceId(path string) string {
	// Split the string based on the "/" separator
	parts := strings.Split(path, "/")

	// Iterate over the parts to find the first non-empty substring
	for _, part := range parts {
		if part != "" {
			return part // Return the first non-empty part
		}
	}

	// Return an empty string if no non-empty part was found
	return ""
}

// generates a unique uuid along with current time
func getRequestId() string {
	// Generate a timestamp
	timestamp := time.Now().UTC().Format("20060102-150405.000")

	// Generate a UUID
	randomUUID := uuid.New().String()

	// Concatenate the timestamp with the UUID
	requestID := fmt.Sprintf("%s-%s", timestamp, randomUUID)

	return requestID
}
