package middlewares

import (
	"auth/internal/tokens"
	"auth/internal/users"
	"context"
	"net/http"
	"strings"
)

type tokensStore interface {
	Get(ctx context.Context, token string) (*tokens.Token, error)
}

type usersStore interface {
	Get(ctx context.Context, id int) (*users.User, error)
}

type Middleware struct {
	TokensStore tokensStore
	UsersStore  usersStore
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			plaintext, found := strings.CutPrefix(authHeader, "Bearer ")
			if !found {
				next.ServeHTTP(w, r)
				return
			}

			hash := tokens.Hash(plaintext)

			token, err := m.TokensStore.Get(r.Context(), string(hash))
			if err != nil {
				return
			}

			user, err := m.UsersStore.Get(r.Context(), token.UserID)
			if err != nil {
				return
			}

			ctx := context.WithValue(r.Context(), "user", user)

			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}

func (m *Middleware) RequireAuthentication(next http.Handler) http.Handler {

	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			val := r.Context().Value("user")

			if val == nil {
				http.Error(w, "authentication required", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		},
	)
}
