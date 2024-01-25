package models

import (
	"errors"
	"strings"

	"github.com/dauid64/super_chat_backend/src/security"
)

type User struct {
	ID uint64 `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (user *User) Prepare(stage string) error {
	if err := user.validate(stage); err != nil {
		return err
	}

	if err := user.format(stage); err != nil {
		return err
	}

	return nil
}

func (user *User) validate(stage string) error {
	if user.Email == "" {
		return errors.New("O nome é obrigatório e não pode estar em branco")
	}
	if stage == "cadastro" && user.Password == "" {
		return errors.New("a senha é obrigatório e não pode estar em branco")
	}

	return nil
}

func (user *User) format(stage string) error {
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	if stage == "cadastro" {
		passwordWithHash, err := security.Hash(user.Password)
		if err != nil {
			return err
		}

		user.Password = string(passwordWithHash)
	}

	return nil
}