package domain

import "time"

type Enterprise struct {
	ID int `json:"id"`

	Name string `json:"name"`

	PhoneNumber string `json:"phone_number"`

	Email string `json:"email"`

	Address string `json:"address"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
