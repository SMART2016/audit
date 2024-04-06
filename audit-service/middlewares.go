package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"net/http"
)

// AuthMiddleware checks the JWT token and authorizes users
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		tokenString := r.Header.Get("Authorization")
		claims := &Claims{}

		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		//TODO: populate with user role system and api
		logRequestDetails(r, "AUDIT_HEALTH", "audit-service")
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		var query map[string]interface{}
		json.NewDecoder(r.Body).Decode(&query)
		queryType := query["type"].(string)

		if !hasAccessToSystem(claims.Role, queryType, r) {
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
