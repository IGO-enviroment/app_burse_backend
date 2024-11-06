package domain

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const UserTable = "users"

type User struct {
	ID int `json:"id" db:"id"`

	Email          string `json:"email" db:"email"`
	DigestPassword string `json:"-" db:"digest_password"`

	FirstName  sql.NullString `json:"first_name" db:"first_name"`
	LastName   sql.NullString `json:"last_name" db:"last_name"`
	MiddleName sql.NullString `json:"middle_name" db:"middle_name"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
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
