package main

import (
	"database/sql"
	"net/http"

	"github.com/google/uuid"
)

type userBookingInfo struct {
	Email           string
	StartTime       sql.NullTime
	Title           sql.NullString
	Description     sql.NullString
	DurationMinutes sql.NullInt32
	PosterImageUrl  sql.NullString
	TrailerVideoUrl sql.NullString
	Rating          sql.NullString
	Genre           sql.NullString
	Director        sql.NullString
	Casts           []string
	BookingID       uuid.NullUUID
	SeatCode        []sql.NullString
}

func (cfg *apiConfig) handlerUserBookings(w http.ResponseWriter, r *http.Request) {
	userIDString := r.PathValue("userID")
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
		return
	}

	dbUserBookingInfo, err := cfg.db.GetUserBookingInfo(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get user booking info", err)
		return
	}

	userBookingInfos := []userBookingInfo{}

	for _, booking := range dbUserBookingInfo {
		seatCodes, err := cfg.db.GetBookedSeatsByBookingID(r.Context(), booking.BookingID.UUID)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't get seat code", err)
			return
		}

		// Convert []string to []sql.NullString
		seatCodesNull := make([]sql.NullString, len(seatCodes))
		for i, seat := range seatCodes {
			seatCodesNull[i] = sql.NullString{String: seat, Valid: true}
		}

		userBookingInfos = append(userBookingInfos, userBookingInfo{
			Email:           booking.Email,
			StartTime:       booking.StartTime,
			Title:           booking.Title,
			Description:     booking.Description,
			DurationMinutes: booking.DurationMinutes,
			PosterImageUrl:  booking.PosterImageUrl,
			TrailerVideoUrl: booking.TrailerVideoUrl,
			Rating:          booking.Rating,
			Genre:           booking.Genre,
			Director:        booking.Director,
			Casts:           booking.Casts,
			BookingID:       booking.BookingID,
			SeatCode:        seatCodesNull,
		})
	}

	respondWithJSON(w, http.StatusOK, userBookingInfos)
}
