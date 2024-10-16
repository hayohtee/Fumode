package main

import (
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/mailer"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/s3"
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
	s3Client     *s3.Client
}
