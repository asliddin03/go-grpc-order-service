package app

import "log"

type App struct{}

func New() *App {
	return &App{}
}

func (a *App) Run() error {
	log.Println("order server started")

	return nil
}
