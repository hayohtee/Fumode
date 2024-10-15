package main

import (
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/mailer"
	"sync"

	"github.com/hayohtee/fumode/internal/jsonlog"
)

// application holds the dependencies for the handlers, middlewares
// and helpers.
type application struct {
	config       config
	logger       *jsonlog.Logger
	wg           sync.WaitGroup
	repositories data.Repositories
	mailer       mailer.Mailer
}
