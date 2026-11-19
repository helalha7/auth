package main

import (
	"auth/internal/middlewares"
	"auth/internal/tokens"
	"auth/internal/users"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := DBConnectionPool(cfg.db.dsn)
	if err != nil {
		log.Fatal(err)
	}

	usersStore := users.Store{DB: db}
	tokensStore := tokens.Store{DB: db}
	usersHandler := &users.Handler{
		UsersStore: &usersStore,
	}
	tokensHandler := &tokens.Handler{
		UsersStore:  &usersStore,
		TokensStore: &tokensStore,
	}
	middlewares := &middlewares.Middleware{
		TokensStore: &tokensStore,
		UsersStore:  &usersStore,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/users", usersHandler.CreateUserHandler)
	mux.HandleFunc("POST /v1/tokens/authentication", tokensHandler.CreateAuthenticationTokenHandler)
	mux.Handle(
		"GET /v1/test",
		middlewares.Authenticate(
			middlewares.RequireAuthentication(
				http.HandlerFunc(
					func(w http.ResponseWriter, r *http.Request) {
						value := r.Context().Value("user")
						user := value.(*users.User)
						message := fmt.Sprintf("Hello %s", user.FullName)
						w.Write([]byte(message))
					},
				),
			),
		),
	)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func DBConnectionPool(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
