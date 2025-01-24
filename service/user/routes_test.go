package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/spiccoli/e-commerceAPI/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserServiceHandler(t *testing.T) {
	userStore := &mockUserStore{}
	handler := NewHandler(userStore)
	t.Run("Should fail with 400 if email in payload is invalid", func(t *testing.T) {
		payload := types.RegisterUser{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "invalid", // Invalid email
			Password:  "password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)

		// Assert that a 400 Bad Request is returned for an invalid email
		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code 400, got %d", rr.Code)
		}
	})

	t.Run("Should correctly register the user", func(t *testing.T) {
		payload := types.RegisterUser{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "asd@gmail.com",
			Password:  "password",
		}
		marshalled, _ := json.Marshal(payload)

		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		router := mux.NewRouter()
		router.HandleFunc("/register", handler.handleRegister).Methods(http.MethodPost)
		router.ServeHTTP(rr, req)
		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code 201, got %d", rr.Code)
		}
	})
}

type mockUserStore struct{}

func (m *mockUserStore) GetUserByEmail(email string) (*types.User, error) {
	return nil, fmt.Errorf("user not found")
}
func (m *mockUserStore) GetUserByID(id int) (*types.User, error) {
	return nil, nil
}
func (m *mockUserStore) CreateUser(user types.User) error {
	return nil
}
