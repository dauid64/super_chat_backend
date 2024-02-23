package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/dauid64/super_chat_backend/src/authetication"
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

func TestCreateUserWithEmailNull(t *testing.T) {
	userForCreate := models.User{
		Email:    "",
		Password: "123456",
	}
	userForCreateJson, err := json.Marshal(userForCreate)
	if err != nil {
		t.Errorf("Erro ao codificar usuário para JSON")
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/usuarios",
		bytes.NewBuffer(userForCreateJson),
	)

	CreateUser(rr, req)

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

	if responseBody.Erro != "O e-mail é obrigatório e não pode estar em branco" {
		t.Errorf("Resposta inesperada (%s)", responseBody.Erro)
	}
}

func TestCreateUserWithPasswordNull(t *testing.T) {
	userForCreate := models.User{
		Email:    "test@gmail.com",
		Password: "",
	}
	userForCreateJson, err := json.Marshal(userForCreate)
	if err != nil {
		t.Errorf("Erro ao codificar usuário para JSON")
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodPost,
		"/usuarios",
		bytes.NewBuffer(userForCreateJson),
	)

	CreateUser(rr, req)

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

	if responseBody.Erro != "a senha é obrigatório e não pode estar em branco" {
		t.Errorf("Resposta inesperada (%s)", responseBody.Erro)
	}
}

func TestSearchUserCorrect(t *testing.T) {
	mock, gormdb := database.DbMock(t)

	userPassword, err := security.Hash("12345678")
	if err != nil {
		t.Errorf("Erro ao gerar senha (%s)", err)
	}

	userFind := models.User{
		Email:    "test@gmail.com",
		Password: string(userPassword),
	}

	rowsInsert := sqlmock.NewRows([]string{"id"}).AddRow(1)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "users" (.+) VALUES (.+)`).WillReturnRows(rowsInsert)
	mock.ExpectCommit()

	result := gormdb.Create(&userFind)
	if result.Error != nil {
		t.Errorf("ocorreu um erro ao criar o usuário (%s)", result.Error)
	}

	token, err := authetication.CreateToken(1)
	if err != nil {
		t.Errorf("Erro ao gerar token (%s)", err)
	}

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/usuarios",
		nil,
	)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	rowsSelect := sqlmock.NewRows([]string{"id", "email", "password"}).
		AddRow(1, "test@gmail.com", string(userPassword))

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id != $1 AND "users"."deleted_at" IS NULL`)).WithArgs(1).WillReturnRows(rowsSelect)

	SearchUsers(rr, req)

	if http.StatusOK != rr.Code {
		t.Errorf("Status code inesperado %d", rr.Code)
	}

	response := rr.Result()
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("Erro ao ler resposta do servidor (%s)", err)
	}

	var userResponse []models.User
	err = json.Unmarshal(data, &userResponse)
	if err != nil {
		t.Errorf("Erro ao converter corpo da resposta (%s)", err)
	}

	if userFind.ID != userResponse[0].ID {
		t.Errorf("Resposta inesperada (%s)", string(data))
	}
}

func TestSearchUserUnauthorized(t *testing.T) {
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(
		http.MethodGet,
		"/usuarios",
		nil,
	)
	req.Header.Add("Authorization", "Bearer token_invalido")

	SearchUsers(rr, req)

	if http.StatusUnauthorized != rr.Code {
		t.Errorf("Status code inesperado %d", rr.Code)
	}
}
