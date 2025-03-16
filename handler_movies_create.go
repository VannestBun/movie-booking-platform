package main

import (
	"crypto/rand"
	"strings"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/vannestbun/movie-booking/internal/auth"
	"github.com/vannestbun/movie-booking/internal/database"
)

type Movie struct {
	ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	DurationMinutes int32     `json:"duration_minutes"`
	PosterImageUrl  string    `json:"poster_image_url"`
	TrailerVideoUrl string    `json:"trailer_video_url"`
	Rating          string    `json:"rating"`
	Genre           string    `json:"genre"`
	Director        string    `json:"director"`
	Casts           []string  `json:"casts"`
}

func (cfg *apiConfig) handlerMoviesCreate(w http.ResponseWriter, r *http.Request) {
	// Set upload limit
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB limit for images

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
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found", err)
		return
	}

	if user.UserRole != "admin" {
		respondWithError(w, http.StatusUnauthorized, "User not authorized", err)
		return
	}

	// Parse the multipart form data
	err = r.ParseMultipartForm(10 << 20) // 10MB limit in memory
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't parse form", err)
		return
	}

	// Extract form fields
	title := r.FormValue("title")
	description := r.FormValue("description")
	durationMinutesInt, err := strconv.Atoi(r.FormValue("duration_minutes"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid duration", err)
		return
	}
	rating := r.FormValue("rating")
	genre := r.FormValue("genre")
	director := r.FormValue("director")
	castsInput := r.FormValue("casts")
	casts := strings.Split(castsInput, ",")
	// Trim whitespace from each name
	for i := range casts {
		casts[i] = strings.TrimSpace(casts[i])
	}

	// Convert int to int32
	durationMinutes := int32(durationMinutesInt)

	// Create a placeholder for image URL - we'll fill this after upload
	posterImageUrl := ""
	trailerVideoUrl := r.FormValue("trailer_video_url") // Keeping this as a URL for now

	// Get the image file from form
	file, header, err := r.FormFile("poster_image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No image file provided", err)
		return
	}
	defer file.Close()

	// Validate that it's an image file
	mediaType, _, err := mime.ParseMediaType(header.Header.Get("Content-Type"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Content-Type", err)
		return
	}
	if mediaType != "image/jpeg" && mediaType != "image/png" {
		respondWithError(w, http.StatusBadRequest, "Invalid file type", nil)
		return
	}

	// Determine the appropriate file extension based on mediaType
	var fileExt string
	if mediaType == "image/jpeg" {
		fileExt = ".jpeg"
	} else if mediaType == "image/png" {
		fileExt = ".png"
	}

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "movie-poster-*"+fileExt)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create temporary file", err)
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy the uploaded file to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not write file to disk", err)
		return
	}

	// Reset file pointer to beginning
	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not reset file pointer", err)
		return
	}

	// Generate a unique key for S3
	fileExt = filepath.Ext(header.Filename)
	randomBytes := make([]byte, 16)
	_, err = rand.Read(randomBytes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate random filename", err)
		return
	}
	fileKey := fmt.Sprintf("%x%s", randomBytes, fileExt)
	_, err = cfg.s3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(cfg.s3Bucket),
		Key:         aws.String(fileKey),
		Body:        tempFile,
		ContentType: aws.String(mediaType),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error uploading file to S3", err)
		return
	}

	// Construct the S3 URL
	posterImageUrl = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		cfg.s3Bucket, cfg.s3Region, fileKey)

	movie, err := cfg.db.CreateMovie(r.Context(), database.CreateMovieParams{
		Title:           title,
		Description:     description,
		DurationMinutes: durationMinutes,
		PosterImageUrl:  posterImageUrl,
		TrailerVideoUrl: trailerVideoUrl,
		Rating:          rating,
		Genre:           genre,
		Director:        director,
		Casts:           casts,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create movie", err)
		return
	}

	showtimes := []string{"10:30", "14:15", "18:00", "21:30"}
	layout := "15:04" // 24-hour time format
	
	for _, timeStr := range showtimes {
		t, err := time.Parse(layout, timeStr)
        if err != nil {
			respondWithError(w, http.StatusInternalServerError, "cannot convert time", err)
            return
        }
		_, err = cfg.db.CreateShowtime(r.Context(), database.CreateShowtimeParams{
			MovieID: movie.ID,
			StartTime: t,
		})
	}

	respondWithJSON(w, http.StatusCreated, Movie{
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
}
