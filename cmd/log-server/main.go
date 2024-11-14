package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	pb "github.com/jacobmiller22/hisight/pkg/history"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
)

type LogServer struct {
	pb.UnimplementedHistoryServer
	DB *sql.DB
}

func NewServer(db *sql.DB) *LogServer {
	return &LogServer{
		DB: db,
	}
}

func (s *LogServer) LogCommand(ctx context.Context, command *pb.Command) (*pb.Ack, error) {
	log.Printf("Received log command %v\n", command.Aliased)

	// s.DB.Exec("INSERT INTO commands (command_text) VALUES (?);", command.Aliased)
	stmt, _ := s.DB.Prepare("INSERT INTO commands (command_text) VALUES (?)")

	stmt.Exec(command.Aliased)

	return &pb.Ack{}, nil
}

func main() {
	args := os.Args

	portStr := args[1]
	port, err := strconv.ParseUint(portStr, 10, 16)

	db, err := sql.Open("sqlite3", "./database.db")
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
	pb.RegisterHistoryServer(grpcServer, NewServer(db))

	log.Printf("Starting server on port %d\n", port)

	grpcServer.Serve(lis)
}
