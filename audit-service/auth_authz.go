package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
)

var users = map[string]string{}     // username:password storage
var userRoles = map[string]string{} // username:role storage
var rolePermissions = map[string][]string{
	ROLE_ADMIN: {"user-service", "monitoring-service"}, // Admins can access any system
	ROLE_USER:  {"monitoring-service"},                 // Regular users can only access monitoring-service logs
}

var apiPermissions = map[string]map[string]string{
	ROLE_ADMIN: {"/audit-service/v1/logevents": "audit_svc_log_event", "/auth-service/v1/register": "user_svc_registration", "/auth-service/v1/login": "user_svc_login"}, // Admins can access any system
	ROLE_USER:  {"/auth-service/v1/login": "user_svc_login", "/audit-service/v1/logevents": "audit_svc_log_event"},                                                       // Regular users can only access monitoring-service logs
}

// JWT secret (store securely in a real application)
var jwtKey = []byte("my_secret_key")

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role,omitempty"` // Optional, defaults to "user"
}

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// RegisterHandler allows new users to register
func RegisterNewUserHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// In a real app, you should securely hash the password
	users[creds.Username] = creds.Password
	if creds.Role == ROLE_ADMIN {
		userRoles[creds.Username] = ROLE_ADMIN
	} else {
		userRoles[creds.Username] = ROLE_USER
	}

	w.Write([]byte("User registered successfully"))
}

// LoginHandler authenticates users and returns a JWT
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassword, ok := users[creds.Username]

	// Validate user credentials
	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Set expiration time for the token
	expirationTime := time.Now().Add(30 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		Role:     userRoles[creds.Username],
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(tokenString))
}

// hasAccessToSystem checks if the role has access to the specified system
func hasAccessToSystem(role, system string, r *http.Request) bool {
	//First check if the user has permission for the API
	if permittedApis, ok := apiPermissions[role]; ok {
		if _, ok := permittedApis[r.RequestURI]; ok {
			return true
		}
	}

	//First check if the user has permission to get data for a specific type
	if permittedSystems, ok := rolePermissions[role]; ok {
		for _, s := range permittedSystems {
			if s == system { // Admins have access to all systems
				return true
			}
		}
	}

	return false
}
