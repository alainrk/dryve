package app

import (
	"dryve/internal/config"
	"dryve/internal/service"
)

type App struct {
	config      config.Config
	fileService service.FileService
}

func NewApp(config config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) WithFileService(s service.FileService) *App {
	a.fileService = s
	return a
}
