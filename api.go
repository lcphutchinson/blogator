package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"github.com/google/uuid"
	"github.com/lcphutchinson/database"
)

type apiConfig struct {
	DB *database.Queries
}

func (cfg *apiConfig) mdAuth(w http.ResponseWriter, r *http.Request) (database.User, bool) {
	auths := r.Header["Authorization"]
	for _, auth := range auths {
		authType, token, ok := strings.Cut(auth, " ")
		if !ok || authType != "ApiKey" {
			continue
		}
		user, err := cfg.DB.ReadUser(r.Context(), token)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Not Found")
			return database.User{}, false
		}
		return user, true
	}
	respondWithError(w, http.StatusUnauthorized, "Unauthorized")
	return database.User{}, false
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	body, _ := json.Marshal(payload)
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(body)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJSON(w, code, map[string]string{"error": msg})
}

func bootMux(config apiConfig, mux *http.ServeMux) {
	mux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("GET /v1/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
	})

	mux.HandleFunc("GET /v1/feeds", func(w http.ResponseWriter, r *http.Request) {
		feeds, err := config.DB.ListFeeds(r.Context())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		respondWithJSON(w, http.StatusOK, feeds)
	})

	mux.HandleFunc("POST /v1/feeds", func(w http.ResponseWriter, r *http.Request) {
		caller, ok := config.mdAuth(w, r)
		if !ok {
			return
		}
		params := database.CreateFeedParams {
			UserID:	caller.ID,
		}
		rData, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = json.Unmarshal(rData, &params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		newFeed, err := config.DB.CreateFeed(r.Context(), params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		respondWithJSON(w, http.StatusOK, newFeed)
	})

	mux.HandleFunc("DELETE /v1/feed_follows/*", func(w http.ResponseWriter, r *http.Request) {
		target, _ := strings.CutPrefix(r.URL.Path, "/v1/feed_follows/")
		targetID, err := uuid.Parse(target)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Bad Request")
			return
		}
		caller, ok := config.mdAuth(w, r)
		if !ok {
			return
		}
		params := database.RemoveFollowParams{
			FeedID: targetID,
			UserID: caller.ID,
		}
		_, err = config.DB.RemoveFollow(r.Context(), params)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Not Found")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	
	mux.HandleFunc("GET /v1/feed_follows", func(w http.ResponseWriter, r *http.Request) {
		caller, ok := config.mdAuth(w, r)
		if !ok {
			return
		}
		userFollows, err := config.DB.ListUserFollows(r.Context(), caller.ID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		respondWithJSON(w, http.StatusOK, userFollows)
	})

	mux.HandleFunc("POST /v1/feed_follows", func(w http.ResponseWriter, r *http.Request) {
		caller, ok := config.mdAuth(w, r)
		if !ok {
			return
		}
		params := database.CreateFollowParams{
			UserID: caller.ID,
		}
		rData, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = json.Unmarshal(rData, &params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		newFollow, err := config.DB.CreateFollow(r.Context(), params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		respondWithJSON(w, http.StatusOK, newFollow)
	})

	mux.HandleFunc("GET /v1/users", func(w http.ResponseWriter, r *http.Request) {
		user, ok := config.mdAuth(w, r)
		if ok {
			respondWithJSON(w, http.StatusOK, user)
		}
	})

	mux.HandleFunc("POST /v1/users", func(w http.ResponseWriter, r *http.Request) {
		params := struct{
			Name	string `json:"name"`
		}{}
		rBody, err := io.ReadAll(r.Body)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		err = json.Unmarshal(rBody, &params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		newUser, err := config.DB.CreateUser(context.TODO(), params.Name)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		respondWithJSON(w, http.StatusOK, newUser)
	})

	mux.HandleFunc("GET /v1/posts", func(w http.ResponseWriter, r *http.Request) {
		caller, ok := config.mdAuth(w, r)
		if !ok {
			return
		}
		var limit int32
		limstr := r.URL.Query().Get("limit")
		_, err := fmt.Sscanf(limstr, "%d", &limit)
		if err != nil {
			limit = 20
		}
		params := database.GetPostsByUserParams{
			UserID: caller.ID,
			Limit: limit,
		}
		posts, err := config.DB.GetPostsByUser(r.Context(), params)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		respondWithJSON(w, http.StatusOK, posts)
	})
}
