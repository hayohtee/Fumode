package data

import (
	"time"

	"github.com/hayohtee/fumode/internal/validator"
)

// Customer is a struct that holds the information about a customer.
type Customer struct {
	CustomerID  int64     `json:"customer_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Password    password  `json:"-"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
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
