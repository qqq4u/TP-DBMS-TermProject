package models

// easyjson -all ./internal/models/error.go

type Error struct {
	Message string `json:"message"`
}
