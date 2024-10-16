package main

import (
	"context"
	"flag"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/hayohtee/fumode/internal/data"
	"github.com/hayohtee/fumode/internal/mailer"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"

	"github.com/hayohtee/fumode/internal/jsonlog"
	"github.com/joho/godotenv"
)

func main() {
	var cfg configuration
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	err := godotenv.Load(".env")
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("FUMODE_DB_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConn, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConn, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection idle time")

	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 587, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	flag.Parse()

	db, err := openDB(cfg)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer db.Close()
	logger.PrintInfo("database connection pool established", nil)

	client, err := mailer.NewMailClient(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password)
	if err != nil {
		logger.PrintFatal(err, nil)
	}
	defer client.Close()

	awsConfig, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(os.Getenv("AWS_ACCESS_KEY"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
		),
	)
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	app := application{
		config:       cfg,
		logger:       logger,
		repositories: data.NewRepositories(db),
		mailer:       mailer.New(client, cfg.smtp.sender),
		s3Client:     s3.NewFromConfig(awsConfig),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}
}
