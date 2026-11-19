package users

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`

	FullName     string `json:"full_name"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	PasswordHash string `json:"-"`
}

func (u *User) SetPassword(plaintext string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plaintext), 12)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashed)

	return nil
}

func (u *User) PasswordMatches(plaintext string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(plaintext))
	if err != nil {
		return err
	}

	return nil
}
