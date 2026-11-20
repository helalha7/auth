package tokens

import (
	"auth/internal/httpx"
	"auth/internal/users"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type usersStore interface {
	GetByUsername(ctx context.Context, username string) (*users.User, error)
}

type tokensStore interface {
	Insert(ctx context.Context, token *Token) error
}

type Handler struct {
	UsersStore  usersStore
	TokensStore tokensStore
}

func (h *Handler) CreateAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return
	}

	user, err := h.UsersStore.GetByUsername(r.Context(), input.Username)
	if err != nil {
		switch {
		case errors.Is(err, users.ErrUserNotFound):
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	if err := user.PasswordMatches(input.Password); err != nil {
		return
	}

	token := GenerateToken(user.ID, 10*time.Minute, Authenticaton)
	h.TokensStore.Insert(r.Context(), token)

	err = httpx.WriteJSON(w, httpx.Envelope{"authentication_token": token}, http.StatusCreated, nil)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
