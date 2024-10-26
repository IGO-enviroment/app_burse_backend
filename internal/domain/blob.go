package domain

type Blob struct {
	ID int `json:"id"`

	Url string `json:"url"`

	FileName string `json:"file_name"`

	Size int64 `json:"size"`

	AttactedID int `json:"attached_id"`
}
