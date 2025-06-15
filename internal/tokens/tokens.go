package tokens

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"time"
)

const (
	ScopeAuth = "authentication"
)

type Token struct {
	PlainText string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

func GenerateToken(userID int64, ttl time.Duration, Scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  Scope,
	}

	emptyByte := make([]byte, 32)
	_, err := rand.Read(emptyByte)
	if err != nil {
		return nil, err
	}

	token.PlainText = base64.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(emptyByte)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil

}
