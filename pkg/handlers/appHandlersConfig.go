package handlers

import "go-bingelists/pkg/config"

var Repo *Repository

type Repository struct {
	Config *config.AppConfig
}

func New(config *config.AppConfig) *Repository {
	return &Repository{
		Config: config,
	}
}

func ConfigNewHandlers(r *Repository) {
	Repo = r
}
