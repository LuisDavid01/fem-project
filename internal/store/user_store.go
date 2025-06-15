package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText string
	hash      []byte
}

func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 13)
	if err != nil {
		return err
	}
	p.plainText = plainText
	p.hash = hash
	return nil
}

func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil // Password does not match
		default:
			return false, err //something went wrong in the server

		}
	}
	return true, nil
}

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash password  `json:"-"`
	Email        string    `json:"email"`
	Bio          string    `json:"bio"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{db: db}
}

type UserStore interface {
	CreateUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

func (pg *PostgresUserStore) CreateUser(user *User) error {
	query := `
	INSERT INTO users (username, password_hash, email, bio)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, updated_at`

	err := pg.db.QueryRow(query, user.Username, user.PasswordHash.hash, user.Email, user.Bio).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}
	return nil

}

func (pg *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}
	query := `
	SELECT id, username, password_hash, email, bio, created_at, updated_at
	FROM users
	WHERE username = $1
	`
	err := pg.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash.hash,
		&user.Email,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // User not found
	}
	if err != nil {
		return nil, err // Other error
	}
	return user, nil // User found
}

func (pg *PostgresUserStore) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username = $1, password_hash = $2, email = $3, bio = $4, updated_at = CURRENT_TIMESTAMP
	WHERE id = $5
	returning updated_at
	`
	result, err := pg.db.Exec(query,
		user.Username,
		user.PasswordHash.hash,
		user.Email,
		user.Bio,
		user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // No user found to update
	}
	return nil
}

func (pg *PostgresUserStore) GetUserToken(scope, tokenPlainText string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlainText))
	query := `
	SELECT u.id, u.username, u.password_hash, u.email, u.bio, u.created_at, u.updated_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3
	`
	user := &User{
		PasswordHash: password{},
	}
	err := pg.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash.hash,
		&user.Bio,
		&user.CreatedAt,
		&user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Token not found or expired
	}
	if err != nil {
		return nil, err // Other error
	}
	return user, nil

}
