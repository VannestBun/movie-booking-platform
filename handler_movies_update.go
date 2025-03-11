package main

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/database"
)

func (cfg *apiConfig) handlerMoviesUpdate(w http.ResponseWriter, r *http.Request) {

	movieIDString := r.PathValue("movieID")
	movieID, err := uuid.Parse(movieIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	currentMovie, err := cfg.db.GetMovie(r.Context(), movieID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to fetch current movie", http.StatusInternalServerError)
		return
	}

	var updateData struct {
		Title           string `json:"title" validate:"required"`
		Description     string `json:"description,omitempty"`
		DurationMinutes int32  `json:"duration_minutes" validate:"gte=0"`
		PosterImageUrl  string `json:"poster_image_url,omitempty"`
		TrailerVideoUrl string `json:"trailer_video_url,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}

	err = cfg.db.UpdateMovie(r.Context(), database.UpdateMovieParams{
		ID:              movieID,
		Title:           ifEmptyUse(updateData.Title, currentMovie.Title),
		Description:     ifEmptyUse(updateData.Description, currentMovie.Description),
		DurationMinutes: ifZeroUse(updateData.DurationMinutes, currentMovie.DurationMinutes),
		PosterImageUrl:  ifEmptyUse(updateData.PosterImageUrl, currentMovie.PosterImageUrl),
		TrailerVideoUrl: ifEmptyUse(updateData.TrailerVideoUrl, currentMovie.TrailerVideoUrl),
	})

	w.WriteHeader(http.StatusNoContent)

}

func ifEmptyUse(newValue, currentValue string) string {
    if newValue == "" {
        return currentValue
    }
    return newValue
}

func ifZeroUse(newValue, currentValue int32) int32 {
    if newValue == 0 {
        return currentValue
    }
    return newValue
}