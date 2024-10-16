package main

// configuration holds all the configuration settings for the app.
type configuration struct {
	// The network port the server is listening on.
	port int
	// The current operating environment for the application
	// (development, staging, production, etc...)
	env string

	// Configurations for database.
	db struct {
		// The data source name.
		dsn string
		// The maximum number of open connections.
		maxOpenConn int
		// The maximum number of idle connections.
		maxIdleConn int
		// The time duration for idle connections in string ("5s" - 5 seconds, "3m" - 3 minutes).
		maxIdleTime string
	}

	// Configurations for rate limiter.
	limiter struct {
		// Request per second
		rps float64
		// Maximum burst
		burst int
		// Check if rate limiter should be enabled
		enabled bool
	}

	// Configurations for SMTP
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}
