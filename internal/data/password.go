package data

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

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
