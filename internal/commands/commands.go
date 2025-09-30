package commands

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jacobmiller22/hisight/internal/commands/protocol/sql"
)

type CommandService struct {
	Repo   *sql.Queries
	Logger *slog.Logger
}

func (s *CommandService) CreateCommand(ctx context.Context, c *Command) error {

	s.Logger.DebugContext(ctx, "COMMAND_CREATE")

	v := sql.InsertCommandParams{
		ID:       uuid.NewString(),
		Aliased:  "",
		Expanded: c.Command,
		Ts:       time.Now().Format(time.RFC3339),
		Shell:    c.Shell,
	}

	return s.Repo.InsertCommand(ctx, v)
}
