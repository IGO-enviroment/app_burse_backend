package users_repository

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

func (r *Repository) GetById(id int) (*domain.User, error) {
	var user domain.User

	err := service.GetById(r.db, id, domain.UserTable, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
