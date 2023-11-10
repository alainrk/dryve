package app

import (
	"dryve/internal/config"
	"dryve/internal/service"
)

type App struct {
	Config      config.Config
	FileService service.FileService
}

func NewApp(config config.Config) *App {
	return &App{
		Config: config,
	}
}

func (a *App) WithFileService(s service.FileService) *App {
	a.FileService = s
	return a
}
