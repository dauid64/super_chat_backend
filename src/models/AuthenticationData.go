package models

type AuthenticationData struct {
	User User `json:"user"`
	Token string `json:"token"`
}