package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jacobmiller22/gossentials/clog"
	"github.com/jacobmiller22/hisight/cli"
	"github.com/jacobmiller22/hisight/internal/cmd"
	"github.com/jacobmiller22/hisight/internal/cmd/hsserver"
	"github.com/jacobmiller22/hisight/internal/config"

	"os"
)

const usage string = "hs [hook | log | server] [args]"

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(os.Args)

	l, err := cfg.Logger()
	if err != nil {
		fmt.Printf("Error creating logger: %v\n", err)
		os.Exit(1)
	}

	ctx = clog.WithContext(ctx, l)

	if err := config.InitDefaultConfig(); err != nil {
		l.InfoContext(ctx, "could not create default config")
	}

	namespace := os.Args[0]

	hslog := cli.New(namespace, nil,
		cli.WithChildren([]cli.Node{
			*cli.New("hook", cmd.HsHook),
			*cli.New("log", cmd.HsLog),
			*cli.New("server", cmd.HsServer,
				cli.WithChildren([]cli.Node{
					*cli.New("http", hsserver.HsServerHttp),
					*cli.New("grpc", hsserver.HsServerGrpc),
				},
				),
			),
		},
		),
	)

	target, targetArgs := hslog.Search(os.Args)

	if target == nil {
		fmt.Printf("target not found.\n\n%s\n", usage)
		os.Exit(0)
	}

	if err := target.Entrypoint(ctx, targetArgs); err != nil {
		if errors.Is(err, cli.ErrNoEntrypoint) {
			fmt.Println(usage)
		} else {
			fmt.Println(err)
		}
		os.Exit(0)
	}
}
