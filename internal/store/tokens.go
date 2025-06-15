package store

import (
	"database/sql"
	"time"

	"github.com/LuisDavid01/femProject/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{db: db}
}

type TokenStore interface {
	Insert(tokens *tokens.Token) error
	CreateNewToken(userId int64, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userID int, scope string) error
}

func (t *PostgresTokenStore) CreateNewToken(userID int64, ttl time.Duration, Scope string) (*tokens.Token, error) {

	token, err := tokens.GenerateToken(userID, ttl, Scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(token)

	return token, err
}

func (t *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope) 
	VALUES ($1, $2, $3, $4)
	`
	_, err := t.db.Exec(query, token.Hash, token.UserID, token.Expiry, token.Scope)
	return err
}

func (t *PostgresTokenStore) DeleteAllTokensForUser(userID int, scope string) error {
	query := `DELETE FROM tokens WHERE user_id = $1 AND scope = $2`
	_, err := t.db.Exec(query, userID, scope)
	if err != nil {
		return err
	}
	return nil
}
