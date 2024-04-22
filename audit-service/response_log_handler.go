package main

import (
	"bytes"
	"net/http"
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
