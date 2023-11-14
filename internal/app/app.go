package app

import (
	"dryve/internal/config"
	"dryve/internal/service"
)

type App struct {
	Config       config.Config
	FileService  service.FileService
	UserService  service.UserService
	EmailService service.EmailService
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

func (a *App) WithUserService(s service.UserService) *App {
	a.UserService = s
	return a
}

func (a *App) WithEmailService(s service.EmailService) *App {
	a.EmailService = s
	return a
}
