package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

// AdminModel is a type which wraps around a sql.DB connection pool
// and provide methods for creating and managing admins to and from
// the database.
type AdminModel struct {
	DB *sql.DB
}

// Insert a new admin record to the database.
func (a AdminModel) Insert(admin *Admin) error {
	query := `
		INSERT INTO admin(name, email, password)
		VALUES ($1, $2, $3)
		RETURNING admin_id, created_at, role`

	args := []any{admin.Name, admin.Email, admin.Password.hash}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, args...).Scan(
		&admin.AdminID,
		&admin.CreatedAt,
		&admin.Role,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), `duplicate key value violates unique constraint "admin_email_key"`):
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

// GetByEmail retrieve the Admin details from the database based
// on the admin email address.
func (a AdminModel) GetByEmail(email string) (*Admin, error) {
	query := `
		SELECT admin_id, name, email, password, role, created_at
		FROM admin
		WHERE email = $1`

	var admin Admin

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := a.DB.QueryRowContext(ctx, query, email).Scan(
		&admin.AdminID,
		&admin.Name,
		&admin.Email,
		&admin.Password.hash,
		&admin.Role,
		&admin.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &admin, nil
}
