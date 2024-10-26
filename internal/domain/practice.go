package domain

const PracticeTable = "practices"

type Practice struct {
	ID int `json:"id"`

	Name string `json:"name"`

	Description string `json:"description"`

	Duration int `json:"duration"`

	InstitutionId int `json:"institutionId"`
	EnterpriseId  int `json:"enterpriseId"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type PracticeFields struct {
	ID            int
	Name          string
	Description   string
	Duration      int
	InstitutionId int
	EnterpriseId  int
	CreatedAt     string
	UpdatedAt     string
}
