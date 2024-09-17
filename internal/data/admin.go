package data

import (
	"github.com/hayohtee/fumode/internal/validator"
	"time"
)

// Admin is a struct that holds the information about an admin.
type Admin struct {
	AdminID   int64     `json:"admin_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

func ValidateAdmin(v *validator.Validator, admin Admin) {
	v.Check(admin.Name != "", "name", "must be provided")
	v.Check(len(admin.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, admin.Email)

	if admin.Password.plaintext != nil {
		ValidatePasswordPlainText(v, *admin.Password.plaintext)
	}

	if admin.Password.hash == nil {
		panic("missing password hash for user")
	}
}
