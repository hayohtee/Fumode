package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// UserRepository is a type which wraps around a sql.DB connection pool
// and provide methods for creating and managing users to and from
// the database.
type UserRepository struct {
	DB *sql.DB
}

// Insert a user record to the database.
func (u UserRepository) Insert(user *User) error {
	query := `
		INSERT INTO users(name, email, password, address, phone_number, role)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING user_id, created_at`

	args := []any{user.Name, user.Email, user.Password.hash, user.Address, user.PhoneNumber, user.Role}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.UserID,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), `duplicate key value violates unique constraint "users_email_key"`):
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// GetByID retrieve a specific User from the database given the
// userID.
func (u UserRepository) GetByID(userID int64) (User, error) {
	query := `
		SELECT user_id, name, email, password, address, phone_number, role, created_at
		FROM users
		WHERE user_id = $1`

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.DB.QueryRowContext(ctx, query, userID).Scan(
		&user.UserID,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Address,
		&user.PhoneNumber,
		&user.Role,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return User{}, ErrRecordNotFound
		default:
			return User{}, err
		}
	}
	return user, nil
}
