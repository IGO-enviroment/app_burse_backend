package profile_repo

import (
	"app_burse_backend/internal/domain"
	"app_burse_backend/internal/service"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetById(id int) (*domain.Profile, error) {
	var profile domain.Profile

	err := service.GetById(r.db, id, domain.ProfileTable, &profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}
