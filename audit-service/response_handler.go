package main

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

// responseWriter is a wrapper around http.ResponseWriter that allows us to capture
// the response status code and body, for logging purposes.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       bytes.Buffer
}

// Write captures the body data and writes it to the underlying ResponseWriter.
func (rw *responseWriter) Write(data []byte) (int, error) {
	rw.body.Write(data)
	return rw.ResponseWriter.Write(data)
}

// WriteHeader captures the status code and delegates to the underlying ResponseWriter.
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// logResponseMiddleware logs the response status code and body, then sends the response back to the client.
func logResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the standard ResponseWriter with our custom writer
		wrappedWriter := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default to 200 OK if not set explicitly
		}

		// Call the next handler with our wrapped writer
		next.ServeHTTP(wrappedWriter, r)

		userName, role, err := getUserClaims(*wrappedWriter, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// After the handler has written to the response, log the details
		requestID := getRequestId()
		serviceId := getServiceId(r.RequestURI)

		logMsg := fmt.Sprintf("RequestId: %s, CurrentUser: %s,Role: %s, System: %s, Action: %s, IP: %s, Agent: %s, Time: %s, Status: %d.\n",
			requestID,
			userName,
			role,
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

func getUserClaims(w responseWriter, r *http.Request) (string, string, error) {
	userName := ""
	var err error
	role := ""
	if strings.HasSuffix(r.RequestURI, "/login") {
		userName, err = getUserDetailsForlogin(w)
		if err != nil {
			http.Error(w.ResponseWriter, err.Error(), http.StatusInternalServerError)
			return "", "", err
		}
		role = getUserRole(userName)
	} else {
		claims, _ := getClaimsAndTokenFromAuthzHeader(r)
		userName = claims.Username
		role = claims.Role
	}
	return userName, role, nil
}

func getUserDetailsForlogin(w responseWriter) (string, error) {
	claims := &Claims{}
	jwt.ParseWithClaims(string(w.body.Bytes()), claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	return claims.Username, nil
}
