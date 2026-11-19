package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"time"
)

const (
	Authenticaton = "AUTHENTICATION"
)

type Token struct {
	ID     int `json:"-"`
	UserID int `json:"-"`

	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func Hash(plainText string) []byte {
	hash := sha256.Sum256([]byte(plainText))
	return hash[:]
}

func GenerateToken(userID int, ttl time.Duration, scope string) *Token {
	token := &Token{
		PlainText: rand.Text(),
		UserID:    userID,
		Expiry:    time.Now().Add(ttl),
		Scope:     scope,
	}
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]

	return token
}
