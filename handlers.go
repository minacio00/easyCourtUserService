package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/minacio00/easyCourtUserService/database"
	"github.com/minacio00/easyCourtUserService/models"
	"github.com/minacio00/easyCourtUserService/services"
	"golang.org/x/crypto/bcrypt"
)

func setTokenHeader(w http.ResponseWriter, token string) {
	w.Header().Set("Authorization", "Bearer "+token)
}

// On success returns the hashed password as a json, a jwt token as a cookie
func hashPassword(w http.ResponseWriter, r *http.Request) {
	p := &models.Credentials{}
	requestData := struct {
		Password string `json:"password"`
	}{}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.Password = requestData.Password
	defer r.Body.Close()
	if p.Password == "" {
		http.Error(w, "empty password", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := generateJWT(p.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	setTokenHeader(w, token)

	// returns the hashed password as json
	encoder := json.NewEncoder(w)
	p.Password = string(hash)
	encoder.Encode(p)

}

// returns status code 202 if the password is correct and a jwt token as a cookie
func signing(w http.ResponseWriter, r *http.Request) {
	var p models.Credentials
	var user models.Tenant

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	if p.Password == "" || p.Email == "" {
		http.Error(w, "Password and Email are requeried", http.StatusBadRequest)
		return
	}
	err = database.Db.First(&user, "email = ?", p.Email).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(p.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	token, err := generateJWT(user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	setTokenHeader(w, token)
	w.WriteHeader(http.StatusAccepted)

}

func Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header not found", http.StatusUnauthorized)
		return
	}

	// Check if the header has the expected "Bearer " prefix
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	// Extract the token string (remove the "Bearer " prefix)
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	// Now tokenString contains the JWT token, and you can use it as needed
	fmt.Println("JWT Token:", tokenString)

	// Other response logic...
	w.WriteHeader(http.StatusOK)

	authService := services.NewAuthenticatorService(database.MongoClient)

	// Check if the token is blacklisted
	blacklisted, err := authService.IsTokenBlacklisted(tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !blacklisted {
		// Blacklist the token if it's not already expired
		valid, err := authService.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusInternalServerError)
			return
		}

		//if the token is already expired just remove the cookie
		if !valid {
			setTokenHeader(w, "")
			w.WriteHeader(http.StatusNoContent)
			return
		}

		//else blacklist the cookie
		err = authService.BlacklistToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		setTokenHeader(w, "")
		w.WriteHeader(http.StatusNoContent)

	}

	// Respond with success message or appropriate status code
	setTokenHeader(w, "")
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w, "Logged out successfully")
}

func GetAllBlacklistedTokens(w http.ResponseWriter, r *http.Request) {
	authService := services.NewAuthenticatorService(database.MongoClient)

	tokens, err := authService.GetAllBlacklistedTokens()
	if err != nil {
		log.Fatal(err)
	}
	tokensJson, _ := json.Marshal(tokens)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(tokensJson)

}

//todo: implement refresh
//todo: write code to see all blacklisted
