package main

import (
	"errors"
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/validator"
	"net/http"
)

func (app *application) registerAdminHandler(w http.ResponseWriter, r *http.Request) {
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

	admin := data.Admin{
		Name:  input.Name,
		Email: input.Email,
	}

	err = admin.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateAdmin(v, admin); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Admins.Insert(&admin)
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
	// TODO: will update later for admin
	//app.background(func() {
	//
	//})

	err = app.writeJSON(w, http.StatusCreated, envelope{"admin": admin}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
