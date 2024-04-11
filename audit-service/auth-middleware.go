package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// contextKey is a value for use with context.WithValue.
// It's a good practice to define keys as a custom type to avoid collisions.
type contextKey string

const credentialsContextKey = contextKey("credentials")

// extractCredentials decodes and extracts username and password from the Authorization header.
func extractCredentials(authHeader string) (*Credentials, error) {
	// The part after "Basic "
	encodedCreds := strings.TrimPrefix(authHeader, "Basic ")

	// Decode from Base64
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedCreds)
	if err != nil {
		return nil, err
	}

	// Convert to string and split to get username and password
	parts := strings.SplitN(string(decodedBytes), ":", 2)
	if len(parts) != 2 {
		return nil, err // You might want to use a more descriptive error here
	}

	return &Credentials{Username: parts[0], Password: parts[1]}, nil
}

// BasicAuthMiddleware decodes credentials and attaches them to the request context.
func BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Basic ") {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		creds, err := extractCredentials(authHeader)
		if err != nil {
			http.Error(w, "Failed to decode credentials", http.StatusUnauthorized)
			return
		}
		requestID := getRequestId()
		serviceId := getServiceId(r.RequestURI)
		logMsg := fmt.Sprintf("RequestId: %s, CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: Initiated\n",
			requestID,
			creds.Username,
			getUserRole(creds.Username),
			serviceId,
			r.Method+":"+r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Now().Format(time.RFC3339))
		//fmt.Println(logMsg)
		publishEventLogs(serviceId, logMsg)
		// Attach the credentials to the request context
		ctx := context.WithValue(r.Context(), credentialsContextKey, creds)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
