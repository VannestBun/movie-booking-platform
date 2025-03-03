package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/database"
)

type Movie struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationMinutes int32     `json:"duration_minutes"`
	PosterImageUrl  string    `json:"poster_image_url"`
	TrailerVideoUrl string    `json:"trailer_video_url"`
}

func (cfg *apiConfig) handlerMoviesCreate(w http.ResponseWriter, r *http.Request) {
	
	type parameters struct {
		Title           string `json:"title"`
		Description     string `json:"description"`
		DurationMinutes int32  `json:"duration_minutes"`
		PosterImageUrl  string `json:"poster_image_url"`
		TrailerVideoUrl string `json:"trailer_video_url"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	movie, err := cfg.db.CreateMovie(r.Context(), database.CreateMovieParams{
		Title:           params.Title,
		Description:     params.Description,
		DurationMinutes: params.DurationMinutes,
		PosterImageUrl:  params.PosterImageUrl,
		TrailerVideoUrl: params.TrailerVideoUrl,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Movie{
		ID:              movie.ID,
		Title:           movie.Title,
		Description:     movie.Description,
		DurationMinutes: movie.DurationMinutes,
		PosterImageUrl:  movie.PosterImageUrl,
		TrailerVideoUrl: movie.TrailerVideoUrl,
	})
}
