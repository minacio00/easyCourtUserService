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
	addr := os.Getenv("ADDR")
	if addr == "" {
		log.Fatal("could not read the addres from the dockerfile")
	}
	database.Connectdb()
	viper.SetConfigFile(".env")
	viper.ReadInConfig()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/hashPassword", hashPassword)
	r.Post("/api/signing", signing)
	fmt.Println("server is running on :8081")
	log.Fatal(http.ListenAndServe(addr, r))

}
