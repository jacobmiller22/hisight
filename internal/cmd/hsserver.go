package cmd

import (
	"context"
	"fmt"

	"github.com/jacobmiller22/hisight/internal/config"
	"github.com/jacobmiller22/hisight/internal/logkeys"

	"github.com/jacobmiller22/gossentials/clog"
)

const hsServerUsage string = "Usage:\n\thsserver [http | grpc] [args]"

func HsServer(ctx context.Context, args []string) error {
	l := clog.FromContext(ctx)
	cfg := config.LoadConfig(args)

	l.Debug(logkeys.CommandStart, logkeys.Command, "HSSERVER", logkeys.Config, cfg)

	return fmt.Errorf(hsServerUsage)
}
