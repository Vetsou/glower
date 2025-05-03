package mocks

import (
	"glower/database/model"
	"net/http"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoMock struct{ mock.Mock }

func (r *AuthRepoMock) InsertUser(user model.User) error {
	args := r.Called(user)
	return args.Error(0)
}

func (r *AuthRepoMock) GetUser(email string) (model.User, error) {
	args := r.Called(email)
	return args.Get(0).(model.User), args.Error(1)
}

const (
	WrongCredentialsMsg = "Invalid email or password."
	WrongFormatMsg      = "Invalid form data. Please fill all required fields."

	CorrectEmail    = "test@example.com"
	CorrectPassword = "Password123"
)

func GetTestUser() model.User {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(CorrectPassword), bcrypt.DefaultCost)

	user := model.User{
		FirstName:    "FirstNameTest",
		LastName:     "LastNameTest",
		Email:        CorrectEmail,
		PasswordHash: hashedPass,
	}
	return user
}

type LoginTestData struct {
	Email, Password string
	ExpectedCode    int
	ExpectedMsg     string
}

func GetLoginDataTestCases() []LoginTestData {
	return []LoginTestData{
		{CorrectEmail, CorrectPassword, http.StatusOK, ""},
		{CorrectEmail, "wrongPass", http.StatusUnauthorized, WrongCredentialsMsg},
		{"", "", http.StatusBadRequest, WrongFormatMsg},
		{CorrectEmail, "", http.StatusBadRequest, WrongFormatMsg},
		{"", CorrectPassword, http.StatusBadRequest, WrongFormatMsg},
		{"test_example.com", CorrectPassword, http.StatusBadRequest, WrongFormatMsg},
		{"testexample.com", CorrectPassword, http.StatusBadRequest, WrongFormatMsg},
		{"test@example", CorrectPassword, http.StatusBadRequest, WrongFormatMsg},
	}
}

type RegisterTestData struct {
	FirstName, LastName       string
	Email                     string
	Password, ConfirmPassword string
	ExpectedCode              int
	ExpectedMsg               string
}

func GetRegisterDataTestCases() []RegisterTestData {
	return []RegisterTestData{
		{
			"fn", "ls", "test@example.com", "password123", "",
			http.StatusBadRequest, WrongFormatMsg,
		},
		{
			"fn", "ls", "test@example.com", "", "password123",
			http.StatusBadRequest, WrongFormatMsg,
		},
		{
			"fn", "", "test@example.com", "password123", "password123",
			http.StatusBadRequest, WrongFormatMsg,
		},
		{
			"", "ls", "test@example.com", "password123", "password123",
			http.StatusBadRequest, WrongFormatMsg,
		},
		{
			"veryveryveryveryveryveryveryveryveryverylongfirstname",
			"ls", "test@example.com", "password", "password",
			http.StatusBadRequest,
			"Invalid form data: first name cannot be longer than 50 characters",
		},
		{
			"firstname", "veryveryveryveryveryveryveryveryveryverylonglastname",
			"test@example.com", "password", "password",
			http.StatusBadRequest,
			"Invalid form data: last name cannot be longer than 50 characters",
		},
		{
			"fs", "ls",
			"veryveryveryveryveryveryveryveryveryveryveryveryverylongtest@example.com",
			"password", "password",
			http.StatusBadRequest,
			"Invalid form data: email cannot be longer than 70 characters",
		},
		{
			"fs", "ls", "test@example.com",
			"veryveryveryveryveryveryveryveryveryveryveryveryverylongpassword",
			"password",
			http.StatusBadRequest,
			"Invalid form data: password cannot be longer than 60 characters",
		},
		{
			"fs", "ls", "test@example.com",
			"password",
			"veryveryveryveryveryveryveryveryveryveryveryveryverylongpassword",
			http.StatusBadRequest,
			"Invalid form data: confirm password cannot be longer than 60 characters",
		},
		{
			"fs", "ls", "test@example.com",
			"password",
			"password_diff",
			http.StatusBadRequest,
			"Invalid form data: passwords must match",
		},
		{
			"fs", "ls", "test@example.com",
			"password",
			"password",
			http.StatusBadRequest,
			"Password is too weak. Please use at least 10 characters",
		},
		{
			"fs", "ls", "test@example.com",
			"Qazxc12345!@#$%",
			"Qazxc12345!@#$%",
			http.StatusOK,
			"",
		},
	}
}
