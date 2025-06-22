package hsserver

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/jacobmiller22/hisight/internal/commands"
	commandsrepository "github.com/jacobmiller22/hisight/internal/commands/repository"
)

func HsServerHttp(args []string) error {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	fset := flag.NewFlagSet("hslog", flag.ContinueOnError)

	var cfg httpConfig
	fset.IntVar(&cfg.http.port, "http-port", 9000, "Port to listen on for HTTP requests")
	fset.StringVar(&cfg.db.dsn, "db-dsn", "db.sqlite", "DSN to database")
	if err := fset.Parse(args); err != nil {
		return nil // NOTE: Hack, ContinueOnError prints usage for us
	}

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
		return err
	}
	return nil
}
