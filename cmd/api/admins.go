package main

import (
	"errors"
	"fmt"
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

func (app *application) loginAdminHandler(w http.ResponseWriter, r *http.Request) {
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

	admin, err := app.models.Admins.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.errorResponse(w, r, http.StatusNotFound, "the provided email address could not be found")
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := admin.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.unauthorizedResponse(w, r, "invalid credentials. Please check your email and password")
		return
	}

	token, err := generateJWT(admin.AdminID, admin.Role)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))

	err = app.writeJSON(w, http.StatusOK, envelope{"admin": admin}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
