package main

import (
	"errors"
	"fmt"
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/validator"
	"net/http"
)

func (app *application) registerCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	customer := data.Customer{
		Name:  input.Name,
		Email: input.Email,
	}

	err = customer.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateCustomer(v, customer); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Customers.Insert(&customer)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Launch a goroutine to send welcome email
	app.background(func() {
		err = app.mailer.Send(customer.Email, "user_welcome.tmpl", customer)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	})

	err = app.writeJSON(w, http.StatusCreated, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) loginCustomerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlainText(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	customer, err := app.models.Customers.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := generateJWT(customer.CustomerID, customer.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	err = app.writeJSON(w, http.StatusOK, envelope{"customer": customer}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
