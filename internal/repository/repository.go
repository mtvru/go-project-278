package repository

import "github.com/mtvru/go-project-278/internal/db"

type Repository struct {
	*db.Queries
}

func New(dbtx db.DBTX) *Repository {
	return &Repository{Queries: db.New(dbtx)}
}
