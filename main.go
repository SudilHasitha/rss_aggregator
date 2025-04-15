package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/SudilHasitha/rss_aggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq" // PostgreSQL driver
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	fmt.Println("Server is starting...")

	// feed, err := urlToFeed("https://lexfridman.com/feed/podcast/")
	// if err != nil {
	// 	log.Fatal("Failed to parse feed:", err)
	// }
	// fmt.Println("Feed", feed.Channel.Items)

	godotenv.Load()
	portString := os.Getenv("PORT")

	if portString == "" {
		log.Fatal("PORT environment variable not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	queries := database.New(conn)
	if err != nil {
		log.Fatal("Failed to create queries:", err)
	}

	apiConfig := apiConfig{
		DB: queries,
	}

	// Start scraping in a separate goroutine
	go startScraping(queries, 15, 120*time.Second)

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
	v1Router.Post("/users", apiConfig.createUserrHandler)
	v1Router.Get("/users", apiConfig.authMiddleware(apiConfig.getUserHandler))
	v1Router.Post("/feeds", apiConfig.authMiddleware(apiConfig.createFeedHandler))
	v1Router.Get("/feeds", apiConfig.getFeedHandler)
	v1Router.Post("/feed_follows", apiConfig.authMiddleware(apiConfig.createFeedFollowHandler))
	v1Router.Get("/feed_follows", apiConfig.authMiddleware(apiConfig.getFeedFollowsHandler))
	v1Router.Delete("/feed_follows/{feedFollowID}", apiConfig.authMiddleware(apiConfig.deleteFeedFollowHandler))
	v1Router.Get("/posts", apiConfig.authMiddleware(apiConfig.getPostsForUserHandler))

	router.Mount("/v1", v1Router)

	server := &http.Server{
		Addr:    ":" + portString,
		Handler: router,
	}
	fmt.Printf("Starting server on port %s\n", portString)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
