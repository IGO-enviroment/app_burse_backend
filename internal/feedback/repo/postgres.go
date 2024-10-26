package feedback_repository

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create() error {
	_, err := r.db.Exec("INSERT INTO feedback (email) VALUES ($1)")
	return err
}

func (r *Repository) CreateForProfile() error {
	_, err := r.db.Exec("INSERT INTO feedback (profile_id) VALUES ($1)")
	return err
}
