package main

import (
	"encoding/json"
	"net/http"

	"github.com/minacio00/easyCourtUserService/database"
	"golang.org/x/crypto/bcrypt"
)

func generateTokenCookie(token string) *http.Cookie {
	cookie := &http.Cookie{
		Name:     "jwt_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
	return cookie
}
func setTokenCookie(w http.ResponseWriter, token string) {
	cookie := generateTokenCookie(token)
	http.SetCookie(w, cookie)
}
func hashPassword(w http.ResponseWriter, r *http.Request) {
	var p Credentials

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	encoder := json.NewEncoder(w)
	p.Password = string(hash)
	encoder.Encode(p)
	token, err := generateJWT(p.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// w.Write(hash)
	json.NewEncoder(w).Encode(map[string]string{"token": token})

}

// returns status code 200 if the password is correct
func signing(w http.ResponseWriter, r *http.Request) {
	var p Credentials
	var user Tenant

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
	setTokenCookie(w, token)
	// json.NewEncoder(w).Encode(map[string]string{"token": token})

	w.WriteHeader(http.StatusAccepted)

}

//todo: implement refresh
//todo: set the jwt as a cookie
//todo: blacklist tokens
