package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// Logger for Login middleware
func LoginLoggerMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		serviceId := getServiceId(r.RequestURI)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}

		var query map[string]interface{}
		if err = json.Unmarshal(body, &query); err != nil {
			fmt.Println(err)
		}
		if r.RequestURI == "/auth-service/v1/register" {
			userRole := getUserRole(query["username"].(string))
			if strings.EqualFold(userRole, "") {
				if query["role"] == ROLE_ADMIN {
					userRole = ROLE_ADMIN
				} else {
					userRole = ROLE_USER
				}
			}
		}
		log.Printf("CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: Initiated\n",
			query["username"],
			getUserRole(query["username"].(string)),
			serviceId,
			r.Method+":"+r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Now().Format(time.RFC3339))
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	}
}

// loggingMiddleware checks the JWT token and authorizes users
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		claims, _ := getClaimsAndTokenFromAuthzHeader(r)
		serviceId := getServiceId(r.RequestURI)
		log.Printf("CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: Initiated\n",
			claims.Username,
			claims.Role,
			serviceId,
			r.Method+":"+r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Now().Format(time.RFC3339))

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r)
	}
}

// AuthMiddleware checks the JWT token and authorizes users
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		claims, token := getClaimsAndTokenFromAuthzHeader(r)
		if err != nil || token == nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		hasAccess := hasAPIAccess(claims.Role, r)
		if !hasAccess {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", claims.Username)
		ctx = context.WithValue(ctx, "role", claims.Role)

		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		next.ServeHTTP(w, r.WithContext(ctx))
		// Pass the username and role to the next handler

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
