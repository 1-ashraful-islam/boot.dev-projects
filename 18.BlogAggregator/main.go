package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/1-ashraful-islam/boot.dev-projects/18.BlogAggregator/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {
	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// Connect to database
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic("Could not open database: " + err.Error())
	}
	defer db.Close()

	dbQueries := database.New(db)

	// Create a new instance of the API config
	apiConfig := &apiConfig{
		DB: dbQueries,
	}

	port := os.Getenv("PORT")
	fmt.Println("Server is running on port", port)

	// start the router with go-chi
	r := chi.NewRouter()
	// cors
	r.Use(middlewareCors())

	r.Mount("/v1", v1Router(apiConfig))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world from root"))
	})

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("Server error", err)
	}

}

func middlewareCors() func(next http.Handler) http.Handler {
	corsOptions := cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	return cors.Handler(corsOptions)
}

func v1Router(apiConfig *apiConfig) http.Handler {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello world from v1 router"))
	})

	r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	})

	return r
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
