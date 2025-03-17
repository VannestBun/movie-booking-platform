package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerBookingsDelete(w http.ResponseWriter, r *http.Request) {
	bookingIDString := r.PathValue("bookingID")
	bookingID, err := uuid.Parse(bookingIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid booking ID", err)
		return
	}

	err = cfg.db.DeleteSeats(r.Context(), bookingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete seats", err)
		return
	}

	err = cfg.db.DeleteBooking(r.Context(), bookingID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete booking", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}