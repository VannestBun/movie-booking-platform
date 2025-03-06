package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/database"
)

type Booking struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	MovieID     string    `json:"movie_id"`
	SeatNumber  string    `json:"seat_number"`
	BookingTime time.Time `json:"booking_time"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (cfg *apiConfig) handlerBookingsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		UserID      uuid.UUID `json:"user_id"`
		MovieID     uuid.UUID `json:"movie_id"`
		SeatNumber  string    `json:"seat_number"`
		BookingTime time.Time `json:"booking_time"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	// Validate request data
	if params.UserID == uuid.Nil || params.MovieID == uuid.Nil || params.SeatNumber == "" {
		respondWithError(w, http.StatusBadRequest, "Missing required fields", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	movie, err := cfg.db.GetMovie(r.Context(), params.MovieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Movie not found", err)
		return
	}

	// Check if seat is available

	// ---????----

	booking, err := cfg.db.CreateBooking(r.Context(), database.CreateBookingParams{
		UserID:      user.ID,
		MovieID:     movie.ID,
		SeatNumber:  params.SeatNumber,
		BookingTime: params.BookingTime,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create booking", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}
