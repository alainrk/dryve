package main

import (
	"dryve/internal/app"
	"dryve/internal/config"
	"dryve/internal/repository"
	"dryve/internal/service"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/httprate"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var defaultConfigPath = "./config.json"

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Load configuration
	if f := os.Getenv("CONFIG_FILE"); f != "" {
		defaultConfigPath = f
	}
	config := config.NewConfig(defaultConfigPath)

	// Initialize database
	db, err := repository.NewDB(config.Database)
	if err != nil {
		fmt.Printf("database initialization failed with err %v\n", err)
	}

	// Register data access objects
	dao := repository.NewDAO(db)

	// Create application and register services
	app := app.NewApp(config).
		WithFileService(service.NewFileService(dao, config.Storage.Path)).
		WithUserService(service.NewUserService(dao)).
		WithEmailService(service.NewMockEmailService(config.Email))

	// Create and setup middlewares and routes
	r := setupRouter(app)

	// Start server
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.HTTP.Port), r)
	if err != nil {
		fmt.Printf("server failed with err %v\n", err)
	}
}

// setupRouter creates and setups middlewares and routes
func setupRouter(app *app.App) *chi.Mux {
	r := chi.NewRouter()

	// Match request paths with a trailing slash and redirect to the same path without.
	r.Use(middleware.RedirectSlashes)
	// Set a few useful out-of-the-box middlewares.
	r.Use(middleware.Logger)
	// Provides the original IP in the request.
	r.Use(middleware.RealIP)
	// Recovers from panic and return 500 instead.
	r.Use(middleware.Recoverer)
	// Provides a unique ID formed by the a process ID and request ID.
	r.Use(middleware.RequestID)
	// Strip ending slashes.
	r.Use(middleware.StripSlashes)

	// Not enabled middlewares for now.
	// // Setup basic CORS.
	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
	// 	AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	// }))
	// // Gzip compression for who accepts it.
	// r.Use(middleware.Compress())
	// // Throttling if needed.
	// r.Use(middleware.ThrottleWithOpts(middleware.ThrottleOpts{}))

	// Public routes for authentication.
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", app.Login)
		r.Post("/register", app.Register)
	})

	// Protected routes (only requiring JWT)
	r.Group(func(r chi.Router) {
		// Search JWT, verify/deny, populate context with user information.
		r.Use(app.JWTMiddleware)

		r.Get("/user/verify/1", app.EmailVerifyStep1)
	})

	// Public routes
	// Public route for email verification
	r.Get("/user/verify/2/email/{id}/{code}", app.EmailVerifyStep2)
	// Public route for keepalive
	r.Get("/healthcheck", app.Healthcheck)

	// Protected routes (after JWT)
	r.Group(func(r chi.Router) {
		// Search JWT, verify/deny, populate context with user information.
		r.Use(app.JWTMiddleware)
		// Use JWT Claims to populate context (user, role, etc.)
		r.Use(app.AuthMiddleware)

		// TODO: Here the protected routes.
		//
		//
	})

	r.Route("/files", func(r chi.Router) {
		r.Get("/{id}", app.GetFile)
		r.Get("/range/{from}/{to}", app.SearchFilesByDateRange)

		// Protect the risky endpoints with a basic rate limiter
		// It needs to be pulled out when horizontal scaling is needed
		r.Group(func(r chi.Router) {
			r.Use(httprate.LimitByIP(app.Config.Limits.FileEndpointsRateLimit, 1*time.Minute))
			r.Post("/", app.UploadFile)
			r.Get("/{id}/download", app.DownloadFile)
			r.Delete("/{id}", app.DeleteFile)
			r.Delete("/range/{from}/{to}", app.DeleteFiles)
		})
	})

	return r
}
