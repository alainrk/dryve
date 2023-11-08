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

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var configFile = "./config.json"

func main() {
	rand.Seed(time.Now().UnixNano())

	if f := os.Getenv("CONFIG_FILE"); f != "" {
		configFile = f
	}
	config := config.NewConfig(configFile)

	db, err := repository.NewDB(config.Database)
	if err != nil {
		fmt.Printf("database initialization failed with err %v\n", err)
	}

	dao := repository.NewDAO(db)

	app := app.NewApp(config).
		WithFileService(service.NewFileService(dao))

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/files", func(r chi.Router) {
		r.Get("/{id}", app.GetFile)
	})

	http.ListenAndServe(fmt.Sprintf(":%d", config.HTTP.Port), r)
}
