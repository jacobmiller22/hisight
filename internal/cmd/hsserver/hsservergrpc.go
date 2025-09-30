package hsserver

import (
	"context"
	"fmt"
	"net"

	"github.com/jacobmiller22/hisight/internal/commands"
	"github.com/jacobmiller22/hisight/internal/commands/protocol/pb"
	repository "github.com/jacobmiller22/hisight/internal/commands/protocol/sql"
	"github.com/jacobmiller22/hisight/internal/config"
	"github.com/jacobmiller22/hisight/internal/logkeys"

	"github.com/jacobmiller22/gossentials/clog"

	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func HsServerGrpc(ctx context.Context, args []string) error {

	l := clog.FromContext(ctx)
	cfg := config.LoadConfig(args)

	l.Debug(logkeys.CommandStart, logkeys.Command, "HSSERVER_GRPC", logkeys.Config, cfg)

	db, err := openDb(cfg.DB.DSN)
	if err != nil {
		l.Info(logkeys.DbConnectError, "dsn", cfg.DB.DSN, logkeys.Error, err)
	}
	defer db.Close()

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.Server.GRPC.Port))
	if err != nil {
		l.Info("SERVER_INIT_LISTEN_ERROR", "port", cfg.Server.GRPC.Port, "protocol", "grpc", "err", err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	cmdRepo := repository.New(db)
	cmdSvc := &commands.CommandService{
		Repo:   cmdRepo,
		Logger: l,
	}
	cmdSvr := &pb.GrpcCommandServiceServer{
		Cmd: cmdSvc,
	}
	pb.RegisterCommandServiceServer(grpcServer, cmdSvr)

	l.Info("SERVER_INIT_SERVE_START", "port", cfg.Server.GRPC.Port, "protocol", "grpc")
	if err := grpcServer.Serve(listener); err != nil {
		l.Info("SERVER_INIT_SERVE_ERROR", "port", cfg.Server.GRPC.Port, "protocol", "grpc", "err", err)
	}
	return nil
}
