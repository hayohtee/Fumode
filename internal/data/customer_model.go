package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"
)

var (
	// ErrDuplicateEmail is a custom error that is returned when there
	// is a duplicate email in the database.
	ErrDuplicateEmail = errors.New("duplicate email")
)

// CustomerModel is a type which wraps around a sql.DB connection pool
// and provide methods for creating and managing customers to and from
// the database.
type CustomerModel struct {
	DB *sql.DB
}

// Insert a new customer record to the database.
func (c CustomerModel) Insert(customer *Customer) error {
	query := `
		INSERT INTO customer(name, email, password, address, phone_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING customer_id, created_at, role`

	args := []any{customer.Name, customer.Email, customer.Password.hash, customer.Address, customer.PhoneNumber}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, args...).Scan(
		&customer.CustomerID,
		&customer.CreatedAt,
		&customer.Role,
	)

	if err != nil {
		switch {
		case strings.Contains(err.Error(), `duplicate key value violates unique constraint "customer_email_key"`):
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// GetByEmail retrieve the Customer details from the database based
// on the customer's email address.
func (c CustomerModel) GetByEmail(email string) (*Customer, error) {
	query := `
		SELECT customer_id, name, email, password, address, phone_number, created_at, role
		FROM customer
		WHERE email = $1`

	var customer Customer

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := c.DB.QueryRowContext(ctx, query, email).Scan(
		&customer.CustomerID,
		&customer.Name,
		&customer.Email,
		&customer.Password.hash,
		&customer.Address,
		&customer.PhoneNumber,
		&customer.CreatedAt,
		&customer.Role,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &customer, nil
}
