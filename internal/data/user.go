package data

import (
	"database/sql"
	"github.com/hayohtee/fumode/internal/validator"
	"time"
)

// User is a struct that holds information about
// a specific user.
type User struct {
	UserID      int64
	Name        string
	Email       string
	Password    password
	Address     sql.NullString
	PhoneNumber sql.NullString
	// Differentiate between types of User (admin, customer)
	Role      string
	CreatedAt time.Time
}

func ValidateUser(v *validator.Validator, user User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, user.Email)

	if user.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *user.Password.plaintext)
	}

	if user.Password.hash == nil {
		panic("missing password hash for user")
	}
}
