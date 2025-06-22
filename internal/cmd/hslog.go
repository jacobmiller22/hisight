package cmd

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/jacobmiller22/hisight/internal/commands/proto"
	"google.golang.org/grpc"
)

func HsLog(args []string) error {

	fmt.Println("hslog!")
	ctx := context.Background()

	target := fmt.Sprintf("127.0.0.1:%d", 3000)

	conn, err := grpc.NewClient(
		target,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

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

	return nil
}
