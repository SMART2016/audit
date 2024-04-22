package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
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
		// Wrap the standard ResponseWriter with our custom writer
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 OK if not set explicitly
		}
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

		logMsg = fmt.Sprintf("RequestId: %s, CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: %d.\n",
			requestID,
			creds.Username,
			getUserRole(creds.Username),
			serviceId,
			r.Method+":"+r.RequestURI,
			r.RemoteAddr,
			r.UserAgent(),
			time.Now().Format(time.RFC3339),
			wrappedWriter.statusCode,
		)
		publishEventLogs(serviceId, logMsg)
	})
}

// AuthMiddleware checks the JWT token and authorizes users
func AuthorizationMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func getClaimsAndTokenFromAuthzHeader(r *http.Request) (*Claims, *jwt.Token) {
	tokenString := r.Header.Get("Authorization")
	claims := &Claims{}

	token, _ := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return claims, token
}

func hasAPIAccess(role string, r *http.Request) bool {
	fmt.Println("role =", role, "   URI =", r.RequestURI, "  Perms=", rolePermissions[role], "  API Perms= ", apiPermissions[role])
	var permitted bool
	//First check if the user has permission for the API
	if permittedApis, ok := apiPermissions[role]; ok {
		if _, ok := permittedApis[r.RequestURI]; ok {
			permitted = true
		} else {
			permitted = false
		}
	}
	return permitted
}
