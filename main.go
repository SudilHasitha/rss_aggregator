package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Hello, World!")
	godotenv.Load()
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,   // Maximum value not ignored by any of major browsers
		AllowCredentials: false, // Allows cookies to be sent
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/readiness", readinessHandler)
	v1Router.Get("/error", handlerErr)

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}
	fmt.Printf("Starting server on port %s\n", portString)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
