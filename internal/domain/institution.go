package domain

type Institution struct {
	ID int `json:"id"`

	Name string `json:"name"`

	Address string `json:"address"`

	PhoneNumber string `json:"phone_number"`

	Email string `json:"email"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
