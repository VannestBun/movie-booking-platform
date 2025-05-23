package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerMoviesDelete(w http.ResponseWriter, r *http.Request) {
	movieIDString := r.PathValue("movieID")
	movieID, err := uuid.Parse(movieIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
		return
	}

	err = cfg.db.DeleteShowtime(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot delete showtime", err)
		return
	}

	err = cfg.db.DeleteMovie(r.Context(), movieID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't delete movie", err)
		return
	}

	// cannot delete when there is booking
	
	w.WriteHeader(http.StatusNoContent)
}

