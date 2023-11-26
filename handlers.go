package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	Password string `json:"password"`
	Hash     string `json:"hash"`
}

func hashPassword(w http.ResponseWriter, r *http.Request) {
	var p password

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
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
	// w.Write(hash)

}

// expects both the password and the hash
// returns status code 200 if the password is correct
func verifyPassword(w http.ResponseWriter, r *http.Request) {
	var p password

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if p.Hash == "" || p.Password == "" {
		http.Error(w, "Password and hash are required", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(p.Hash), []byte(p.Password))
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
