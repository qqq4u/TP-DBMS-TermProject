package models

// easyjson -all ./internal/models/user.go

type User struct {
	ID       int    `json:"-"`
	Nickname string `json:"nickname,omitempty"`
	Fullname string `json:"fullname"`
	About    string `json:"about,omitempty"`
	Email    string `json:"email"`
}
