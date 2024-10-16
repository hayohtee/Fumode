package main

import (
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/mailer"
	"github.com/hayohtee/fumode/internal/uploader"
	"sync"

	"github.com/hayohtee/fumode/internal/jsonlog"
)

// application holds the dependencies for the handlers, middlewares
// and helpers.
type application struct {
	config       configuration
	logger       *jsonlog.Logger
	wg           sync.WaitGroup
	repositories data.Repositories
	mailer       mailer.Mailer
	s3Uploader   *uploader.S3Uploader
}
