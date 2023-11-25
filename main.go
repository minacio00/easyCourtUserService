package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/api/hashPassword", hashPassword)
	r.Post("/api/verifyPassword", verifyPassword)

	fmt.Println("server is running on :8081")
	log.Fatal(http.ListenAndServe(":8081", r))

}
