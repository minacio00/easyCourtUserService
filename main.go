package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/minacio00/easyCourtUserService/database"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = viper.GetString("ADDR")
		if addr == "" {
			log.Fatal("could not read the address from the dockerfile")
		}
	}
	database.Connectdb()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/hashPassword", hashPassword)
	r.Post("/api/signing", signing)
	r.Post("/api/logout", Logout)
	//test this endpoint
	r.Get("/api/", GetAllBlacklistedTokens)
	fmt.Printf("server is running on %v", addr)
	log.Fatal(http.ListenAndServe(":"+addr, r))

}
