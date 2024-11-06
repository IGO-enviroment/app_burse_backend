package users_repository

import (
	"app_burse_backend/internal/domain"
	"app_burse_backend/internal/service"
	"app_burse_backend/pkg/postgres"
	"fmt"
)

type Repository struct {
	db postgres.Database
}

func NewRepository(db postgres.Database) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetById(id int) (*domain.User, error) {
	var user domain.User

	err := service.GetByField(r.db, "id", id, domain.UserTable, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := service.GetByField(r.db, domain.UserTable, "email", email, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *Repository) Create(fields []service.FieldDB) (int, error) {
	var id int

	mapper := service.NewMapperFields(fields...)
	err := r.db.QueryRowx(
		fmt.Sprintf("INSERT INTO users %s VALUES %s RETURNING id", mapper.GetColumnNames(), mapper.GetPlaceholders()),
		mapper.GetValues()...,
	).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
