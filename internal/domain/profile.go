package domain

const ProfileTable = "profiles"

type Profile struct {
	ID int `json:"id"`

	UserID int `json:"user_id"`
}
