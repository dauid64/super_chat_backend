package models

type User struct {
	ID uint64 `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}
