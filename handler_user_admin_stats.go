package main

import (
	"net/http"
	"time"

	"github.com/google/uuid"
)

type UserAdminStats struct {
	TotalMovies   int64  `json:"total_movies"`
	TotalUsers    int64  `json:"total_users"`
	TotalBookings int64  `json:"total_bookings"`
	TotalRevenue  string `json:"total_revenue"`
}

type TopFiveMovies struct {
	ID            uuid.UUID `json:"ID"`
	Title         string    `json:"title"`
	TotalBookings int64     `json:"total_bookings"`
	TotalRevenue  float64    `json:"total_revenue"`
}

type ShowtimeOccupancyResponse struct {
	Title          string    `json:"title"`
	StartTime      time.Time `json:"start_time"`
	OccupiedSeats  int64     `json:"occupied_seats"`
	AvailableSeats int64     `json:"available_seats"`
}

func (cfg *apiConfig) handlerUserAdminStats(w http.ResponseWriter, r *http.Request) {
	dbUserAdmin, err := cfg.db.GetUserAdminStats(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve user admin cinema stats", err)
		return
	}

	respondWithJSON(w, http.StatusOK, UserAdminStats{
		TotalMovies:   dbUserAdmin.TotalMovies,
		TotalUsers:    dbUserAdmin.TotalUsers,
		TotalBookings: dbUserAdmin.TotalBookings,
		TotalRevenue:  dbUserAdmin.TotalRevenue,
	})

}

func (cfg *apiConfig) handlerUserAdminStatsTopFiveMovies(w http.ResponseWriter, r *http.Request) {
	topFiveMovies, err := cfg.db.GetTopFiveMovies(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve top five movies", err)
		return
	}

	movies := []TopFiveMovies{}
	for _, topFiveMovie := range topFiveMovies {
		movies = append(movies, TopFiveMovies{
			ID:            topFiveMovie.ID,
			Title:         topFiveMovie.Title,
			TotalBookings: topFiveMovie.TotalBookings,
			TotalRevenue:  float64(topFiveMovie.TotalBookings) * 10.99,
		})
	}
	respondWithJSON(w, http.StatusOK, movies)
}

func (cfg *apiConfig) handlerUserAdminStatsShowtimeOccupancy(w http.ResponseWriter, r *http.Request) {
	movieTitleString := r.PathValue("movieTitle")

	dbShowtimeOccupancy, err := cfg.db.GetShowtimeOccupancy(r.Context(), movieTitleString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve showtime occupancy data", err)
		return
	}

	totalSeats := 96

	showtimeOccupancy := []ShowtimeOccupancyResponse{}
	for _, showtimeResponse := range dbShowtimeOccupancy {
		availableSeats := int64(totalSeats) - showtimeResponse.OccupiedSeats
		showtimeOccupancy = append(showtimeOccupancy, ShowtimeOccupancyResponse{
			Title:          showtimeResponse.Title.String,
			StartTime:      showtimeResponse.StartTime,
			OccupiedSeats:  showtimeResponse.OccupiedSeats,
			AvailableSeats: availableSeats,
		})
	}
	respondWithJSON(w, http.StatusOK, showtimeOccupancy)
}
