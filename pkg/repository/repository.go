package repository

import (
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	PullRequest *PrRepository
	Team        *TeamRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Team:        NewTeam(db),
		PullRequest: NewPR(db),
	}
}
