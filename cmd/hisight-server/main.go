package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jacobmiller22/hisight/internal/commands"
	commandsrepository "github.com/jacobmiller22/hisight/internal/commands/repository"
)

type config struct {
	http struct {
		port int
	}
	db struct {
		dsn string
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	var cfg config
	flag.IntVar(&cfg.http.port, "http-port", 8000, "Port to listen on for HTTP requests")
	flag.StringVar(&cfg.db.dsn, "db-dsn", "db.sqlite", "DSN to database")
	flag.Parse()

	cmdRepo := &commandsrepository.Queries{}

	cmdSvc := &commands.CommandService{Repo: cmdRepo, Logger: logger}

	cmdRoutes := commands.CommandRoutes{CmdSvc: cmdSvc}

	mux := http.NewServeMux()

	mux.Handle("GET /commands", cmdRoutes.GetCommandsHandler())
	mux.Handle("POST /commands", cmdRoutes.CreateCommandHandler())
	mux.Handle("GET /commands/{commandId}", cmdRoutes.GetCommandHandler())

	addr := fmt.Sprintf(":%d", cfg.http.port)
	logger.Info("SERVER_INIT_START", "port", cfg.http.port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		logger.Error("SERVER_INIT_ERROR", "err", err)
	}

}
