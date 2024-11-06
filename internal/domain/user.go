package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

const UserTable = "users"

type User struct {
	ID int `json:"id"`

	Email          string `json:"email"`
	DigestPassword string `json:"-"`

	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	MiddleName string `json:"middle_name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

func (u *User) SetPassword(password string) error {
	hash, err := u.HashPassword(password)
	if err != nil {
		return err
	}

	u.DigestPassword = hash
	return nil
}
