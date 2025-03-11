package main

import (
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
		Title:           updateData.Title,
		Description:     updateData.Description,
		DurationMinutes: updateData.DurationMinutes,
		PosterImageUrl:  updateData.PosterImageUrl,
		TrailerVideoUrl: updateData.TrailerVideoUrl,
	})

	w.WriteHeader(http.StatusNoContent)

}
