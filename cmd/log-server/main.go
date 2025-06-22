package main

import (
	"database/sql"
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

func main() {
	// ctx := context.Background()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	port := 9001

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Error opening db: %v\n", err)
	}
	defer db.Close()

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS commands (id INTEGER PRIMARY KEY, command_text TEXT NOT NULL, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		log.Fatalf("Error creating commands table: %v\n", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		log.Fatalf("Failed trying to listen on port %d: %v", port, err)
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	cmdRepo := repository.New(db)
	cmdSvc := commands.CommandService{
		Repo:   cmdRepo,
		Logger: logger,
	}
	pb.RegisterCommandServiceServer(grpcServer, cmdSvc)

	log.Printf("Starting server on port %d\n", port)

	grpcServer.Serve(lis)
}
