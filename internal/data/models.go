package data

import (
	"database/sql"
	"errors"
)

var (
	// ErrRecordNotFound is a custom error that is returned when looking
	// for a specific record that is not in the database.
	ErrRecordNotFound = errors.New("record not found")

	// ErrEditConflict is a custom error that is returned when two or more
	// users try to access the same record concurrently.
	ErrEditConflict = errors.New("edit conflict")

	// ErrDuplicateEmail is a custom error that is returned when there
	// is a duplicate email in the database.
	ErrDuplicateEmail = errors.New("duplicate email")
)

// Models is a container that holds all the database models for this project.
type Models struct {
	Customers CustomerModel
	Admins    AdminModel
}

// NewModels returns a Model which contains all initialized database models
// for the project.
func NewModels(db *sql.DB) Models {
	return Models{
		Customers: CustomerModel{DB: db},
		Admins:    AdminModel{DB: db},
	}
}
