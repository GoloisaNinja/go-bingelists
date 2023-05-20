package config

var Repo *Repository

type Repository struct {
	Config *AppConfig
}

func NewRepo(config *AppConfig) *Repository {
	return &Repository{
		Config: config,
	}
}

func NewAppConfiguration(r *Repository) {
	Repo = r
}
