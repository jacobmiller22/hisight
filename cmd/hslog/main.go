package main

import (
	"errors"
	"fmt"

	"github.com/jacobmiller22/hisight/internal/cmd"
	"github.com/jacobmiller22/hisight/internal/cmd/hsserver"

	"os"
)

const usage string = "hs [hook | log | server] [args]"

func main() {

	hslog := cmd.
		NewCliNode(os.Args[0], nil).
		WithChildren([]cmd.CliNode{
			*cmd.NewCliNode("hook", cmd.HsHook),
			*cmd.NewCliNode("log", cmd.HsLog),
			*cmd.NewCliNode("server", cmd.HsServer).
				WithChildren(
					[]cmd.CliNode{
						*cmd.NewCliNode("http", hsserver.HsServerHttp),
						*cmd.NewCliNode("grpc", hsserver.HsServerGrpc),
					},
				),
		},
		)

	target, targetArgs := hslog.Search(os.Args)

	if target == nil {
		fmt.Printf("target not found.\n\n%s\n", usage)
		os.Exit(0)
	}

	if err := target.Entrypoint(targetArgs); err != nil {
		if errors.Is(err, cmd.ErrNoEntrypoint) {
			fmt.Println(usage)
		} else {
			fmt.Println(err)
		}
		os.Exit(0)
	}
}
