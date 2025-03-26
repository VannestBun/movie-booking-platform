package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/vannestbun/movie-booking/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	jwtSecret      string
	s3Client       *s3.Client
	s3Bucket       string
	s3Region       string
}

func main() {
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}
	dbQueries := database.New(dbConn)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	s3Bucket := os.Getenv("S3_BUCKET")
	if s3Bucket == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}

	s3Region := os.Getenv("S3_REGION")
	if s3Region == "" {
		log.Fatal("S3_REGION environment variable is not set")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	client := s3.NewFromConfig(awsCfg)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		jwtSecret:      jwtSecret,
		s3Client:       client,
		s3Bucket:       s3Bucket,
		s3Region:       s3Region,
	}

	mux := http.NewServeMux()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	mux.HandleFunc("GET /api/users/{userID}", apiCfg.handlerUserBookings)
	mux.HandleFunc("GET /api/users/admin/stats", apiCfg.handlerUserAdminStats)
	mux.HandleFunc("GET /api/users/admin/stats/top-five-movies", apiCfg.handlerUserAdminStatsTopFiveMovies)
	mux.HandleFunc("GET /api/users/admin/stats/{movieTitle}", apiCfg.handlerUserAdminStatsShowtimeOccupancy)

	mux.HandleFunc("POST /api/movies", apiCfg.handlerMoviesCreate)
	mux.HandleFunc("GET /api/movies", apiCfg.handlerMoviesRetrieve)
	mux.HandleFunc("GET /api/movies/{movieID}", apiCfg.handlerMoviesGet)
	mux.HandleFunc("PUT /api/movies/{movieID}", apiCfg.handlerMoviesUpdate)
	mux.HandleFunc("DELETE /api/movies/{movieID}", apiCfg.handlerMoviesDelete)

	mux.HandleFunc("GET /api/showtimes/{showtimeID}", apiCfg.handlerShowtimeCreate)

	mux.HandleFunc("POST /api/bookings", apiCfg.handlerBookingsCreate)
	mux.HandleFunc("DELETE /api/bookings/{bookingID}", apiCfg.handlerBookingsDelete)

	mux.HandleFunc("GET /api/booking-seats/{movieID}/{startTime}", apiCfg.handlerBookedSeatsGet)
	mux.HandleFunc("POST /api/booking-seats", apiCfg.handlerBookedSeatsCreate)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	// ------- FRONT END ROUTE -------

	fs := apiCfg.middlewareMetricsInc(http.FileServer(http.Dir("./frontend")))
	mux.Handle("/", fs)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: c.Handler(mux),
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(srv.ListenAndServe())
}
