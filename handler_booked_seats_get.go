package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/database"
)

type BookingSeat struct {
	ID        uuid.UUID
	BookingID uuid.UUID
	SeatCode  string
}

func (cfg *apiConfig) handlerBookedSeatsGet(w http.ResponseWriter, r *http.Request) {

	startTimeString := r.PathValue("startTime")
	startTime, err := time.Parse("3:04 PM", startTimeString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't parse time", err)
		return
	}

	movieIDString := r.PathValue("movieID")
	movieID, err := uuid.Parse(movieIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	// Get show time to verify it exist
	_, err = cfg.db.GetShowtimeByMovieID(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get showtime", err)
		return
	}

	seatCode, err := cfg.db.GetBookedSeats(r.Context(), database.GetBookedSeatsParams{
		StartTime: startTime,
		MovieID:   movieID,
	})
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Booked seats not found", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, seatCode)

}

func (cfg *apiConfig) handlerBookedSeatsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		BookingID uuid.UUID `json:"booking_id"`
		SeatCode  string    `json:"seat_code"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	bookedSeat, err := cfg.db.CreateBookingSeat(r.Context(), database.CreateBookingSeatParams{
		BookingID: params.BookingID,
		SeatCode:  params.SeatCode,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create booking", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, BookingSeat{
		ID:        bookedSeat.ID,
		BookingID: bookedSeat.BookingID,
		SeatCode:  bookedSeat.SeatCode,
	})

}
