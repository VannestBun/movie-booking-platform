package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/auth"
	"github.com/vannestbun/movie-booking/internal/database"
)

func (cfg *apiConfig) handlerMoviesUpdate(w http.ResponseWriter, r *http.Request) {
    // Limit request size
    r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB limit

    token, err := auth.GetBearerToken(r.Header)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't find JWT", err)
        return
    }

    userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
    if err != nil {
        respondWithError(w, http.StatusUnauthorized, "Couldn't validate JWT", err)
        return
    }

    user, err := cfg.db.GetUser(r.Context(), userID)
    if err != nil || user.UserRole != "admin" {
        respondWithError(w, http.StatusUnauthorized, "User not authorized", err)
        return
    }

    // Parse form data
    err = r.ParseMultipartForm(10 << 20)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Couldn't parse form", err)
        return
    }

    movieIDStr := r.FormValue("id")
    movieID, err := uuid.Parse(movieIDStr)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid movie ID", err)
        return
    }

    movie, err := cfg.db.GetMovie(r.Context(), movieID)
    if err != nil {
        respondWithError(w, http.StatusNotFound, "Movie not found", err)
        return
    }

    // Update fields
    if title := r.FormValue("title"); title != "" {
        movie.Title = title
    }
    if description := r.FormValue("description"); description != "" {
        movie.Description = description
    }
    if durationStr := r.FormValue("duration_minutes"); durationStr != "" {
        if duration, err := strconv.Atoi(durationStr); err == nil {
            movie.DurationMinutes = int32(duration)
        }
    }
    if rating := r.FormValue("rating"); rating != "" {
        movie.Rating = rating
    }
    if genre := r.FormValue("genre"); genre != "" {
        movie.Genre = genre
    }
    if director := r.FormValue("director"); director != "" {
        movie.Director = director
    }
    if castsInput := r.FormValue("casts"); castsInput != "" {
        movie.Casts = strings.Split(castsInput, ",")
        for i := range movie.Casts {
            movie.Casts[i] = strings.TrimSpace(movie.Casts[i])
        }
    }
    if trailerURL := r.FormValue("trailer_video_url"); trailerURL != "" {
        movie.TrailerVideoUrl = trailerURL
    }

    // Handle optional image upload
    file, header, err := r.FormFile("poster_image")
    if err == nil {
        defer file.Close()
        mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
        if err == nil && (mediaType == "image/jpeg" || mediaType == "image/png") {
            fileExt := ".jpeg"
            if mediaType == "image/png" {
                fileExt = ".png"
            }
            tempFile, err := os.CreateTemp("", "movie-poster-*"+fileExt)
            if err == nil {
                defer os.Remove(tempFile.Name())
                defer tempFile.Close()
                io.Copy(tempFile, file)
                tempFile.Seek(0, io.SeekStart)

                randomBytes := make([]byte, 16)
                rand.Read(randomBytes)
                fileKey := fmt.Sprintf("%x%s", randomBytes, fileExt)

                _, err = cfg.s3Client.PutObject(r.Context(), &s3.PutObjectInput{
                    Bucket:      aws.String(cfg.s3Bucket),
                    Key:         aws.String(fileKey),
                    Body:        tempFile,
                    ContentType: aws.String(mediaType),
                })
                if err == nil {
                    movie.PosterImageUrl = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
                        cfg.s3Bucket, cfg.s3Region, fileKey)
                }
            }
        }
    }

    err = cfg.db.UpdateMovie(r.Context(), database.UpdateMovieParams{
        ID:              movie.ID,
        Title:           movie.Title,
        Description:     movie.Description,
        DurationMinutes: movie.DurationMinutes,
        PosterImageUrl:  movie.PosterImageUrl,
        TrailerVideoUrl: movie.TrailerVideoUrl,
        Rating:          movie.Rating,
        Genre:           movie.Genre,
        Director:        movie.Director,
        Casts:           movie.Casts,
    })
    if err != nil {
        respondWithError(w, http.StatusInternalServerError, "Couldn't update movie", err)
        return
    }

    respondWithJSON(w, http.StatusOK, movie)
}
