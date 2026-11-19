package users

import (
	"auth/internal/httpx"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type usersStore interface {
	Get(ctx context.Context, id int) (*User, error)
	Insert(ctx context.Context, newUser *User) error
}

type Handler struct {
	UsersStore usersStore
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "couldn't decode json", http.StatusBadRequest)
		return
	}

	newUser := &User{}
	newUser.FullName = input.FullName
	newUser.Email = input.Email
	newUser.Username = input.Username

	if err := newUser.SetPassword(input.Password); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.UsersStore.Insert(r.Context(), newUser); err != nil {
		switch {
		case errors.Is(err, ErrDuplicateUsername):
			http.Error(w, "duplicate username", http.StatusConflict)
		case errors.Is(err, ErrDuplicateEmail):
			http.Error(w, "duplicate email", http.StatusConflict)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	header := http.Header{}
	header.Set("Location", fmt.Sprintf("/v1/users/%v", newUser.ID))

	err := httpx.WriteJSON(w, httpx.Envelope{"user": newUser}, http.StatusCreated, header)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
