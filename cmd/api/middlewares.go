package main

import (
	"errors"
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	AdminRole    = "admin"
	CustomerRole = "customer"
)

// recoverPanic is a middleware that recover panics in any http.Handler
// and send 500 Internal Server error response to the client.
func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// rateLimit is a middleware that uses token bucket rate limit implementation
// to limit the number of requests to the endpoints.
func (app *application) rateLimit(next http.Handler) http.Handler {
	// client is a struct for holding the rate limiter and last seen
	// for each client.
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	// Launch background goroutine which remove old entries from the clients
	// map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			// Extract client's IP address from the request
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			mu.Lock()
			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			// Update the last seen for the client.
			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()

		}
		next.ServeHTTP(w, r)
	})
}

// authorize is a middleware that authorize the user. It checks for Authorization Header in
// the request and validates it using the secret jwt key it then extract the role from the
// user claims and see if it matches the provided role.
func (app *application) authorize(role string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		data := strings.Split(authorizationHeader, " ")

		if len(data) != 2 {
			w.WriteHeader(http.StatusExpectationFailed)
			return
		}

		payload, err := validateJWT(data[0])
		if err != nil {
			switch {
			case errors.Is(err, errInvalidToken):
				app.unauthorizedResponse(w, r, "invalid or expired token, please authenticate again")
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		if role == AdminRole {
			if !strings.EqualFold(payload.Role, AdminRole) {
				app.forbiddenResponse(w, r, "you do not have permission to access this resource, admin role required")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
