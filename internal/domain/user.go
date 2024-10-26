package domain

import "time"

const UserTable = "users"

type User struct {
	ID int `json:"id"`

	Email string `json:"email"`

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
