package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMoviesGet(w http.ResponseWriter, r *http.Request) {
	movieIDString := r.PathValue("movieID")
	movieID, err := uuid.Parse(movieIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	dbMovie, err := cfg.db.GetMovie(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get movie", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Movie{
		ID:              dbMovie.ID,
		CreatedAt:       dbMovie.CreatedAt,
		UpdatedAt:       dbMovie.UpdatedAt,
		Title:           dbMovie.Title,
		Description:     dbMovie.Description,
		DurationMinutes: dbMovie.DurationMinutes,
		PosterImageUrl:  dbMovie.PosterImageUrl,
		TrailerVideoUrl: dbMovie.TrailerVideoUrl,
	})
}

func (cfg *apiConfig) handlerMoviesRetrieve(w http.ResponseWriter, r *http.Request) {
	dbMovies, err := cfg.db.GetMovies(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve movies", err)
		return
	}

	movies := []Movie{}
	for _, dbMovie := range dbMovies {
		movies = append(movies, Movie{
			ID:              dbMovie.ID,
			CreatedAt:       dbMovie.CreatedAt,
			UpdatedAt:       dbMovie.UpdatedAt,
			Title:           dbMovie.Title,
			Description:     dbMovie.Description,
			DurationMinutes: dbMovie.DurationMinutes,
			PosterImageUrl:  dbMovie.PosterImageUrl,
			TrailerVideoUrl: dbMovie.TrailerVideoUrl,
		})
	}
	
	respondWithJSON(w, http.StatusOK, movies)
}
