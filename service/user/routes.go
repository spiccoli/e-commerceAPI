package user

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/spiccoli/e-commerceAPI/service/auth"
	"github.com/spiccoli/e-commerceAPI/types"
	"github.com/spiccoli/e-commerceAPI/utils"
	"net/http"
)

type Handler struct {
	store types.UserStore
}

func NewHandler(store types.UserStore) *Handler {
	return &Handler{store: store}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/login", h.handleLogin).Methods(http.MethodPost)
	router.HandleFunc("/register", h.handleRegister).Methods(http.MethodPost)
}

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) handleRegister(w http.ResponseWriter, r *http.Request) {
	//Json payload
	var payload types.RegisterUser
	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
	}

	// validate payload
	if err := utils.Validate.Struct(payload); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("validation errors: %v", validationErrors))
			return
		}
		// If it's not a validation error, handle it as a generic error
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("unexpected error: %w", err))
		return
	}

	// check if user already exists
	_, err := h.store.GetUserByEmail(payload.Email)
	if err == nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("user with mail %s already exists", payload.Email))
		return
	}
	hashedPassword, err := auth.HashedPassword(payload.Password)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// if user does not exist, create user
	err = h.store.CreateUser(types.User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
		Email:     payload.Email,
		Password:  hashedPassword,
	})
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}
	err = utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "user created"})
}
