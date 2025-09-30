package hsserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jacobmiller22/hisight/internal/commands"
	"github.com/jacobmiller22/hisight/internal/commands/protocol/json"
	"github.com/jacobmiller22/hisight/internal/commands/protocol/sql"
	"github.com/jacobmiller22/hisight/internal/config"
	"github.com/jacobmiller22/hisight/internal/logkeys"

	"github.com/jacobmiller22/gossentials/clog"
)

func HsServerHttp(ctx context.Context, args []string) error {
	l := clog.FromContext(ctx)
	cfg := config.LoadConfig(args)

	l.Debug(logkeys.CommandStart, logkeys.Command, "HSSERVER_HTTP", logkeys.Config, cfg)

	db, err := openDb(cfg.DB.DSN)
	if err != nil {
		l.Info(logkeys.DbConnectError, "dsn", cfg.DB.DSN, logkeys.Error, err)
	}
	defer db.Close()

	cmdRepo := sql.New(db)

	cmdSvc := &commands.CommandService{Repo: cmdRepo, Logger: l}

	cmdRoutes := json.CommandRoutes{CmdSvc: cmdSvc}

	mux := http.NewServeMux()

	mux.Handle("GET /commands", cmdRoutes.GetCommandsHandler())
	mux.Handle("POST /commands", cmdRoutes.CreateCommandHandler())
	mux.Handle("GET /commands/{commandId}", cmdRoutes.GetCommandHandler())

	addr := fmt.Sprintf(":%d", cfg.Server.HTTP.Port)
	l.Info("SERVER_INIT_START", "port", cfg.Server.HTTP.Port, "protocol", "http")
	if err := http.ListenAndServe(addr, mux); err != nil {
		l.Error("SERVER_INIT_ERROR", "protocol", "http", "err", err)
		return err
	}
	return nil
}
