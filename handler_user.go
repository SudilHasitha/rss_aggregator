package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SudilHasitha/rss_aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) createUserrHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		Name:      params.Name,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to create user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusCreated, databaseUserToAPIUser(user))
}

func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, http.StatusOK, databaseUserToAPIUser(user))
}

func (apiCfg *apiConfig) getPostsForUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPostsForUser(r.Context(), user.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't get posts for user: %v", err))
		return
	}
	respondWithJSON(w, http.StatusOK, databasePostsToAPIPosts(posts))
}
