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

// Repositories is a container that holds all the database repositories for this project.
type Repositories struct {
	Users UserRepository
}

// NewRepositories returns a Repositories which contains all initialized repositories for
// interacting with the database.
// for the project.
func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Users: UserRepository{DB: db},
	}
}
