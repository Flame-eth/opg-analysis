package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("Port is not set")
	}

	router := chi.NewRouter()

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Server running on ppprt %s", portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal("Error starting server")
	}
	
}