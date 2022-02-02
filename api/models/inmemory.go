package models

// InMemoryRequest holds the data coming from user
type InMemoryRequest struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}
