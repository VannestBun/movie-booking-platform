package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/database"
)

type Showtime struct {
	ID        uuid.UUID
	MovieID   uuid.UUID
	StartTime time.Time
}

func (cfg *apiConfig) handlerShowtimeCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		MovieID   uuid.UUID `json:"user_id"`
		StartTime time.Time `json:"showtime_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	movie, err := cfg.db.GetMovie(r.Context(), params.MovieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Movie not found", err)
		return
	}

	showtime, err := cfg.db.CreateShowtime(r.Context(), database.CreateShowtimeParams{
		MovieID:   movie.ID,
		StartTime: params.StartTime,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create showtime", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Showtime{
		ID:        showtime.ID,
		MovieID:   showtime.MovieID,
		StartTime: showtime.StartTime,
	})
}
