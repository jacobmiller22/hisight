package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	pb "github.com/jacobmiller22/hisight/internal/commands/protocol/pb"
	"github.com/jacobmiller22/hisight/internal/config"
	"github.com/jacobmiller22/hisight/internal/logkeys"

	"github.com/jacobmiller22/gossentials/clog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrUnknownLogLevel error = errors.New("unknown log level")

const hsLogUsage string = "Usage:\n\thslog [flags] [command]"

func HsLog(ctx context.Context, args []string) error {

	l := clog.FromContext(ctx)
	cfg := config.LoadConfig(args)

	l.Debug(logkeys.CommandStart, logkeys.Command, "HSLOG", logkeys.Config, cfg)

	if len(args) < 1 {
		return fmt.Errorf(hsLogUsage)
	}

	conn, err := grpc.NewClient(
		cfg.Server.GRPC.Host,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		l.Error("GRPC_CHANNEL_CREATION_ERROR", "err", err)
		return err
	}

	defer conn.Close()

	cmdSvc := pb.NewCommandServiceClient(conn)

	_, err = cmdSvc.LogCommand(ctx, &pb.Command{
		Shell:   "",
		Command: strings.Join(args, " "),
	})
	if err != nil {
		log.Fatalf("Error received from LogCommand: %v", err)
	}

	return nil
}
