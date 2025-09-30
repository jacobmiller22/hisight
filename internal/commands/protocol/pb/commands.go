package pb

import (
	"context"

	"github.com/jacobmiller22/hisight/internal/commands"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcCommandServiceServer struct {
	UnimplementedCommandServiceServer
	Cmd *commands.CommandService
	// mustEmbedUnimplementedCommandServiceServer()
}

func (svr *GrpcCommandServiceServer) LogCommand(ctx context.Context, c *Command) (*emptypb.Empty, error) {

	cmd := commands.Command{
		Shell:   c.Shell,
		Command: c.Command,
	}

	if err := svr.Cmd.CreateCommand(ctx, &cmd); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
