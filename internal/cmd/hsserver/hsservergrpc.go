package hsserver

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/jacobmiller22/hisight/internal/commands"
	pb "github.com/jacobmiller22/hisight/internal/commands/proto"
	"github.com/jacobmiller22/hisight/internal/commands/repository"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

func HsServerGrpc(args []string) error {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	fset := flag.NewFlagSet("hslog", flag.ContinueOnError)

	var cfg config
	fset.IntVar(&cfg.http.port, "http-port", 9001, "Port to listen on for HTTP requests")
	fset.StringVar(&cfg.db.dsn, "db-dsn", ":memory:", "DSN to database")
	if err := fset.Parse(args); err != nil {
		return nil // NOTE: Hack, ContinueOnError prints usage for us

	}

	db, err := sql.Open("sqlite3", cfg.db.dsn)
	if err != nil {
		log.Fatalf("Error opening db: %v\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS commands (id INTEGER PRIMARY KEY, command_text TEXT NOT NULL, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		log.Fatalf("Error creating commands table: %v\n", err)
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", cfg.http.port))
	if err != nil {
		log.Fatalf("Failed trying to listen on port %d: %v", cfg.http.port, err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	cmdRepo := repository.New(db)
	cmdSvc := commands.CommandService{
		Repo:   cmdRepo,
		Logger: logger,
	}
	pb.RegisterCommandServiceServer(grpcServer, cmdSvc)

	log.Printf("Starting server on port %d\n", cfg.http.port)

	grpcServer.Serve(listener)
	return nil
}
