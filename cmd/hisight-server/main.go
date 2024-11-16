package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/jacobmiller22/hisight/internal/repository"
)

type config struct {
	http_port int
	db_dsn    string
}

type hsApp struct {
	config config
	// db     database.Database
	repo repository.Repository
	log  *slog.Logger
}

func main() {
	var cfg config
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	flag.IntVar(&cfg.http_port, "http-port", 8000, "Port to listen on for HTTP requests")
	flag.StringVar(&cfg.db_dsn, "db-dsn", "db.sqlite", "DSN to database")

	flag.Parse()

	repo, err := repository.NewRepo(cfg.db_dsn)

	if err != nil {
		log.Error("Failed to create repository", "err", err)
	}

	app := &hsApp{
		config: cfg,
		repo:   repo,
		log:    log,
	}

	app.serveHttp()
}
