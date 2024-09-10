package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		ErrorLog:     log.New(app.logger, "", 0),
	}

	// Create shutdown error channel to receive any errors returned by the graceful
	// Shutdown() function.
	shutdownError := make(chan error)

	// Start a background goroutine to listen for OS signals
	go func() {
		quit := make(chan os.Signal, 1)
		// Listen for incoming SIGINT and SIGTERM signals and
		// relay them to the quit channel.
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		// Read the signal from quit channel or block until a signal
		// is received.
		sig := <-quit

		// Log a message to say that the server is shutting down
		app.logger.PrintInfo("shutting down server", map[string]string{
			"signal": sig.String(),
		})

		// Create a context with a 20-second timeout
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownError <- server.Shutdown(ctx)
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  app.config.env,
	})

	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	app.logger.PrintInfo("stopped server", map[string]string{
		"addr": server.Addr,
	})

	return nil
}
