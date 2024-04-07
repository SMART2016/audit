package main

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"time"
)

const (
	ROLE_ADMIN = "admin"
	ROLE_USER  = "user"
)

// username:password storage
var users = map[string]string{}

// username:role storage
var userRoles = map[string]string{}

// Role to permitted attribute mapping , could be done better
var rolePermissions = map[string][]string{
	ROLE_ADMIN: {"auth-service", "audit-service"}, // Admins can access any system
	ROLE_USER:  {"auth-service", "audit-service"}, // Regular users can only access monitoring-service logs
}

// Role to Permitted API mapping, Method could also be added
var apiPermissions = map[string]map[string]string{
	ROLE_ADMIN: {"/audit-service/v1/health": "audit_svc_health", "/audit-service/v1/logevents": "audit_svc_log_event", "/auth-service/v1/register": "auth_svc_registration", "/auth-service/v1/login": "auth_svc_login", "/auth-service/v1/users": "auth_svc_users_list", "/auth-service/v1/health": "auth_svc_health"}, // Admins can access any system
	ROLE_USER:  {"/audit-service/v1/health": "audit_svc_health", "/auth-service/v1/login": "auth_svc_login", "/audit-service/v1/logevents": "audit_svc_log_event"},                                                                                                                                                      // Regular users can only access monitoring-service logs
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
/**
POST http://localhost:9191/auth-service/v1/register
payload: {
	"username":"admin",
    "password":"admin",
	"role":"admin"
}
*/
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

/*
*
GET http://localhost:9191/auth-service/v1/users
*/
func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(userRoles)
	w.WriteHeader(http.StatusOK)
}

/**
POST http://localhost:9191/auth-service/v1/login
payload: {
	"username":"admin",
    "password":"admin"
}
*/
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
	expirationTime := time.Now().Add(3000 * time.Minute)
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

func getUserRole(username string) string {
	return userRoles[username]
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

// Used to identify the role to attribute permission
func checkAttributeAccess(role string, system interface{}) (bool, map[string]interface{}) {
	AttributeFilterMap := map[string]interface{}{}
	permitted := false
	AttributeFilterMap["System.keyword"] = system
	if system != nil {
		if permittedSystems, ok := rolePermissions[role]; ok {
			for _, s := range permittedSystems {
				if s == system {
					permitted = true
					return permitted, AttributeFilterMap
				} else {
					permitted = false
				}
			}
		}
	}
	return permitted, AttributeFilterMap
}