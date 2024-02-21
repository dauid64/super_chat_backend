package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/dauid64/super_chat_backend/src/database"
	"github.com/dauid64/super_chat_backend/src/models"
	"github.com/dauid64/super_chat_backend/src/security"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
	"gorm.io/gorm"
)

type UserBody struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type ResponseBody struct {
	Erro string `json:"erro,omitempty"`
}

func TestCreateUserWithNoneBody(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/usuarios", nil)

	CreateUser(rr, req)

	if http.StatusBadRequest != rr.Code {
		t.Errorf("Status code inesperado %d", rr.Code)
	}
}

func TestCreateUserCorrect(t *testing.T) {
	mock, _ := database.DbMock(t)
	expectedSQL := regexp.QuoteMeta("INSERT INTO \"users\" (\"created_at\",\"updated_at\",\"deleted_at\",\"email\",\"password\")")
	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(expectedSQL).WillReturnRows(rows)
	mock.ExpectCommit()

	userForCreate := []byte(`{
		"email": "test@gmail.com",
		"password": "123456"
	}`)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/usuarios",
		bytes.NewBuffer(userForCreate),
	)

	CreateUser(rr, req)

	if http.StatusCreated != rr.Code {
		t.Errorf("Status code inesperado %d", rr.Code)
	}

	response := rr.Result()
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Erro ao ler resposta do servidor (%s)", err)
	}

	var userResponse models.User
	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		t.Errorf("Erro ao converter corpo da resposta (%s)", err)
	}

	if userResponse.ID != 1 || userResponse.Email != "test@gmail.com" {
		t.Errorf("Resposta inesperada (%+v)", userResponse)
	}
}

func TestCreateUserWithExistEmail(t *testing.T) {
	mock, gormdb := database.DbMock(t)

	rows := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WillReturnRows(rows)
	mock.ExpectCommit()

	userPassword, err := security.Hash("12345678")
	if err != nil {
		t.Errorf("Erro ao gerar senha (%s)", err)
	}

	existingUser := models.User{
		Email:    "test@gmail.com",
		Password: string(userPassword),
	}

	result := gormdb.Create(&existingUser)
	if result.Error != nil {
		t.Errorf("ocorreu um erro ao criar o usuário (%s)", result.Error)
	}

	userForCreate := []byte(`{
		"email": "test@gmail.com",
		"password": "12345678"
	}`)

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WillReturnError(gorm.ErrDuplicatedKey)
	mock.ExpectRollback()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/usuarios",
		bytes.NewBuffer(userForCreate),
	)

	CreateUser(rr, req)

	if http.StatusInternalServerError != rr.Code {
		t.Errorf("Status code inesperado %d", rr.Code)
	}

	response := rr.Result()
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Erro ao ler resposta do servidor (%s)", err)
	}

	var responseBody ResponseBody
	err = json.Unmarshal(data, &responseBody)
	if err != nil {
		t.Errorf("Erro ao converter corpo da resposta (%s)", err)
	}

	if responseBody.Erro != "e-mail já cadastrado" {
		t.Errorf("Resposta inesperada (%s)", responseBody.Erro)
	}
}
