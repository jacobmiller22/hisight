package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/jacobmiller22/hisight/internal/commands/proto"
	"google.golang.org/grpc"
)

func main() {

	ctx := context.Background()

	args := os.Args

	fp, err := os.OpenFile("/Users/jacobmiller22/text.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	defer fp.Close()

	_, err = fp.WriteString("i was called")
	if err != nil {
		log.Fatalf("Error 2: %v", err)
	}

	// var opts []grpc.DialOption
	// opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	//
	// port := 3000
	//
	// conn, err := grpc.Dial(fmt.Sprintf("127.0.0.1:%d", port), opts...)
	target := fmt.Sprintf("127.0.0.1:%d", 3000)

	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Error while dialiing grpc server: %v", err)
	}

	defer conn.Close()

	cmdSvc := pb.NewCommandServiceClient(conn)

	_, err = cmdSvc.LogCommand(ctx, &pb.Command{
		Aliased:         args[1],
		ExpandedPreview: args[2],
		ExpandedFull:    args[3],
	})
	if err != nil {
		log.Fatalf("Error received from LogCommand: %v", err)
	}
}
