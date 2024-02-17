package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/1-ashraful-islam/boot.dev-projects/18.BlogAggregator/internal/database"
	"github.com/1-ashraful-islam/boot.dev-projects/18.BlogAggregator/internal/scrapper"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB     *database.Queries
	Logger *log.Logger
}

func (cfg *apiConfig) ScrapeFeeds(ctx context.Context, t time.Duration, n int32) {
	ticker := time.NewTicker(t)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Printf("Returning from ScrapeFeeds")
			return
		case <-ticker.C:
			feeds, err := cfg.DB.GetNextFeedsToFetch(ctx, n)
			if err != nil {
				cfg.Logger.Printf("Failed to get feeds to fetch: %+v", err)
				continue
			}
			if len(feeds) == 0 {
				log.Println("No feeds to fetch. Sleeping...")
				continue
			}

			var wg sync.WaitGroup
			for _, feed := range feeds {
				wg.Add(1)
				go func(feed database.Feed) {
					defer wg.Done()

					err := scrapper.ScrapeFeed(ctx, cfg.DB, feed)
					if err != nil {
						cfg.Logger.Printf("Failed to scrape feed: %+v", err)
					}
				}(feed)

			}
			wg.Wait()

		}
	}

}

type authedHandler func(w http.ResponseWriter, r *http.Request, u database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get api key from header
		authHeader := r.Header.Get("Authorization")

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			respondWithError(w, http.StatusUnauthorized, "Invalid authorization header. Expected format: 'Bearer <API key>'")
			return
		}
		apiKey := authHeaderParts[1]

		if apiKey == "" {
			respondWithError(w, http.StatusUnauthorized, "API key is missing")
			return
		}

		user, err := cfg.DB.GetUser(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid API key: "+apiKey)
			return
		}

		//call the handler with the authenticated user
		handler(w, r, user)
	}
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, u database.User) {
	respondWithJSON(w, http.StatusOK, u)
}

func (cfg *apiConfig) handlerUsersPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var u struct {
			Name string `json:"name"`
		}

		if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
			cfg.Logger.Printf("Failed to decode request body: %+v", err)
			respondWithError(w, http.StatusBadRequest, "Invalid request payload. Please provide a valid JSON object")
			return
		}
		defer r.Body.Close()

		if u.Name == "" {
			respondWithError(w, http.StatusBadRequest, "Name is required")
			return
		}

		//check if user already exists
		if _, err := cfg.DB.GetUserByName(r.Context(), u.Name); err == nil {
			respondWithError(w, http.StatusBadRequest, "User already exists")
			return
		}

		// Create a new user
		createdUser, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
			ID:        uuid.New(),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Name:      u.Name,
		})

		if err != nil {
			cfg.Logger.Printf("Failed to create user: %+v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		respondWithJSON(w, http.StatusCreated, createdUser)

	}
}

func (cfg *apiConfig) handlerFeedsPost(w http.ResponseWriter, r *http.Request, u database.User) {
	var f struct {
		Title       string `json:"title"`
		URL         string `json:"url"`
		Description string `json:"description"`
	}

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		cfg.Logger.Printf("Failed to decode request body: %+v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload. Please provide a valid JSON object")
		return
	}
	defer r.Body.Close()

	if f.Title == "" || f.URL == "" || f.Description == "" {
		respondWithError(w, http.StatusBadRequest, "title, url, and description fields are required")
		return
	}

	// check if feed already exists
	if _, err := cfg.DB.GetFeedByURL(r.Context(), f.URL); err == nil {
		respondWithError(w, http.StatusBadRequest, "Feed already exists")
		return
	}

	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      u.ID,
		Url:         f.URL,
		Title:       f.Title,
		Description: f.Description,
	})

	if err != nil {

		cfg.Logger.Printf("Failed to create feed: %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed")
		return
	}

	// create feed_follow for the user
	feed_follow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    u.ID,
	})

	if err != nil {
		cfg.Logger.Printf("Failed to create feed follow: %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed follow")
		return
	}
	var result = struct {
		Feed       database.Feed       `json:"feed"`
		FeedFollow database.FeedFollow `json:"feed_follow"`
	}{
		Feed:       feed,
		FeedFollow: feed_follow,
	}

	respondWithJSON(w, http.StatusCreated, result)
}

func (cfg *apiConfig) handlerFeedsGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		feeds, err := cfg.DB.GetFeeds(r.Context())
		if err != nil {
			cfg.Logger.Printf("Failed to get feeds: %+v", err)
			respondWithError(w, http.StatusInternalServerError, "Failed to get feeds")
			return
		}

		respondWithJSON(w, http.StatusOK, feeds)
	}
}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, u database.User) {
	feed_follows, err := cfg.DB.GetFeedFollowsByUser(r.Context(), u.ID)
	if err != nil {
		cfg.Logger.Printf("Failed to get feed_follows: %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to get feed_follows")
		return
	}

	respondWithJSON(w, http.StatusOK, feed_follows)
}

func (cfg *apiConfig) handlerFeedFollowsPost(w http.ResponseWriter, r *http.Request, u database.User) {
	var ff struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&ff); err != nil {
		cfg.Logger.Printf("Failed to decode request body: %+v", err)
		respondWithError(w, http.StatusBadRequest, "Invalid request payload. Please provide a valid JSON object")
		return
	}
	defer r.Body.Close()

	if ff.FeedID == uuid.Nil {
		respondWithError(w, http.StatusBadRequest, "feed_id is required")
		return
	}

	// check if feed exists
	if _, err := cfg.DB.GetFeedByID(r.Context(), ff.FeedID); err != nil {
		respondWithError(w, http.StatusBadRequest, "Feed does not exist")
		return
	}

	// check if feed_follow already exists
	if _, err := cfg.DB.GetFeedFollows(r.Context(), database.GetFeedFollowsParams{FeedID: ff.FeedID, UserID: u.ID}); err == nil {
		respondWithError(w, http.StatusBadRequest, "Feed follow already exists")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    ff.FeedID,
		UserID:    u.ID,
	})

	if err != nil {
		cfg.Logger.Printf("Failed to create feed follow: %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, feedFollow)
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, u database.User) {
	feedFollowID := chi.URLParam(r, "feed_follows_id")
	if feedFollowID == "" {
		respondWithError(w, http.StatusBadRequest, "feed_follows_id is required")
		return
	}

	ffID, err := uuid.Parse(feedFollowID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid feed_follows_id")
		return
	}

	// check if feed_follow exists
	if _, err := cfg.DB.GetFeedFollows(r.Context(), database.GetFeedFollowsParams{FeedID: ffID, UserID: u.ID}); err != nil {
		cfg.Logger.Printf("Failed to get feed follow for feed_id %v: %+v", ffID, err)
		respondWithError(w, http.StatusBadRequest, "Feed follow does not exist")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		FeedID: ffID,
		UserID: u.ID,
	})
	if err != nil {
		cfg.Logger.Printf("Failed to delete feed follow: %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})

}

func main() {

	// Open the log file
	logFile, err := os.OpenFile("logs/blog-aggregator.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file:", err)
	}

	// Create a multi writer
	multiWriter := io.MultiWriter(logFile, os.Stdout)

	// Create the logger
	logger := log.New(multiWriter, "", log.Lshortfile|log.LstdFlags)

	// Load .env file
	if err := godotenv.Load(".env"); err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	// Connect to the database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatalf(errors.Wrap(err, "could not connect to the database").Error())
	}
	defer db.Close()

	// Check the connection
	err = db.Ping()
	if err != nil {
		logger.Fatalf(errors.Wrap(err, "could not ping the database").Error())
	}

	dbQueries := database.New(db)

	// Create a new instance of the API config
	apiConfig := &apiConfig{
		DB:     dbQueries,
		Logger: logger,
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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// test the scrapper
	go apiConfig.ScrapeFeeds(ctx, 1*time.Minute, 10)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen and Serve returned err: %v", err)
		}
	}()

	<-ctx.Done()
	log.Printf("Got interrupt signal: %s, shutting down...", ctx.Err())

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("HTTP server shutdown failed: %v", err)
	}

	log.Println("HTTP server shutdown complete")

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

	log.Println("CORS permissions are very permissive. Tighten before deploying to production")

	return cors.Handler(corsOptions)
}

func v1Router(apiConfig *apiConfig) http.Handler {
	r := chi.NewRouter()

	r.Get("/readiness", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Get("/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	})

	r.Post("/users", apiConfig.handlerUsersPost())
	r.Get("/users", apiConfig.middlewareAuth(apiConfig.handlerUsersGet))

	r.Post("/feeds", apiConfig.middlewareAuth(apiConfig.handlerFeedsPost))
	r.Get("/feeds", apiConfig.handlerFeedsGet())

	r.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedFollowsPost))
	r.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.handlerFeedFollowsGet))
	r.Delete("/feed_follows/{feed_follows_id}", apiConfig.middlewareAuth(apiConfig.handlerFeedFollowsDelete))

	return r
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON response: %+v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to marshal JSON response")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
