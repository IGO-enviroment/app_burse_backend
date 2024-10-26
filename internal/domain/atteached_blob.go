package domain

type AttactedBlob struct {
	ID int `json:"id"`

	BlobID int `json:"blob_id"`

	ModelType string `json:"model_type"`
	ModelID   int    `json:"model_id"`
}
