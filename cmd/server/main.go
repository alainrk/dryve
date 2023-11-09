package main

import (
	"dryve/internal/app"
	"dryve/internal/config"
	"dryve/internal/repository"
	"dryve/internal/service"
	"fmt"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var defaultConfigPath = "./config.json"

func main() {
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
		WithFileService(service.NewFileService(dao))

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

	r.Use(middleware.Logger)
	r.Use(middleware.RedirectSlashes)

	r.Route("/files", func(r chi.Router) {
		r.Post("/", app.UploadFile)
		// r.Get("/{id}", app.GetFile)
		// r.Delete("/{id}", app.DeleteFile)
		// r.Delete("/range/{from_year}-{from_month}-{from_day}/{to_year}-{to_month}-{to_day}", app.DeleteFiles)
	})

	return r
}
