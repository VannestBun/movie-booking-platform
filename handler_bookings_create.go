package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/database"
)

type Booking struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  time.Time
	ShowtimeID uuid.UUID
}

func (cfg *apiConfig) handlerBookingsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		UserID   string   `json:"user_id"`
		Showtime string   `json:"showtime"`
		MovieId  string   `json:"movie_id"`
		SeatCode []string `json:"seat_code"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}
	userID, err := uuid.Parse(params.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id format", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	// Don't forget AUTH USING THE TOKENS

	movieID, err := uuid.Parse(params.MovieId)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie_id format", err)
		return
	}

	startTime, err := time.Parse("3:04 PM", params.Showtime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse time", err)
		return
	}

	showtime, err := cfg.db.GetShowtimeByMovieAndStartTime(r.Context(), database.GetShowtimeByMovieAndStartTimeParams{
		MovieID:   movieID,
		StartTime: startTime,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Showtime not found", err)
		return
	}

	booking, err := cfg.db.CreateBooking(r.Context(), database.CreateBookingParams{
		UserID:     user.ID,
		ShowtimeID: showtime.ID,
	})

	// this should also automatically booked the approrpiate seats
	for _, seat_code := range params.SeatCode {
		_, err := cfg.db.CreateBookingSeat(r.Context(), database.CreateBookingSeatParams{
			BookingID: booking.ID,
			SeatCode:  seat_code,
		})
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Unable to create booking seat", err)
			return
		}
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create booking", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, Booking{
		ID:         booking.ID,
		UserID:     booking.UserID,
		CreatedAt:  booking.CreatedAt,
		UpdatedAt:  booking.UpdatedAt,
		ShowtimeID: booking.ShowtimeID,
	})

}

// so first you check the auth
// POST req with params of showtimeID, make booking, then post sseat with this newly made booking id
