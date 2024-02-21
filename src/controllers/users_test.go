package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/dauid64/super_chat_backend/src/database"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type UserBody struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
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
	mock := database.DbMock(t)
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
}
