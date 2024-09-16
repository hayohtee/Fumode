package data

import (
	"errors"
	"time"

	"github.com/hayohtee/fumode/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// Customer is a struct that holds the information about a customer.
type Customer struct {
	CustomerID  int64     `json:"customer_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    password  `json:"-"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
}

// password is a struct which contain the plaintext and hashed
// versions of the password for a user. The plaintext is a pointer
// to a string, so that we are able to distinguish between password
// not present in the struct versus a plaintext password which is empty string "".
type password struct {
	plaintext *string
	hash      []byte
}

// Set method calculates the bcrypt hash of a plaintext password, and
// stores both the hash and the plaintext versions in the struct.
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Matches method checks whether the provided plaintext password matches
// the hashed password stored in the struct, returning true if it matches
// and false otherwise.
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlainText(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateCustomer(v *validator.Validator, customer Customer) {
	v.Check(customer.Name != "", "name", "must be provided")
	v.Check(len(customer.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, customer.Email)

	if customer.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *customer.Password.plaintext)
	}

	if customer.Password.hash == nil {
		panic("missing password hash for user")
	}
}
