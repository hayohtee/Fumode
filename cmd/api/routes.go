package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /v1/customers", app.registerCustomerHandler)
	mux.HandleFunc("POST /v1/customers/login", app.loginUserHandler)

	mux.HandleFunc("POST /v1/admins", app.registerAdminHandler)
	mux.HandleFunc("POST /v1/admins/login", app.loginUserHandler)

	return mux
}
