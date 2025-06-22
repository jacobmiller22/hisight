package commands

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/jacobmiller22/hisight/internal/commands/proto"
	"github.com/jacobmiller22/hisight/internal/commands/repository"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

var ErrInvalidOption error = errors.New("invalid option")
var ErrBadId error = errors.New("bad identifier")

/*
Routes
*/
type CommandRoutes struct {
	CmdSvc *CommandService
}

func (rte *CommandRoutes) GetCommandsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cmds, err := rte.CmdSvc.Repo.ListCommands(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err = json.NewEncoder(w).Encode(cmds); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (rte *CommandRoutes) GetCommandHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("commandId")

		if id == "" {
			http.Error(w, fmt.Sprintf("%v: empty id", ErrBadId), http.StatusBadRequest)
			return
		}

		cmd, err := rte.CmdSvc.Repo.CommandById(r.Context(), id)

		if json.NewEncoder(w).Encode(cmd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (rte *CommandRoutes) CreateCommandHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var cmd *repository.InsertCommandParams
		if err := json.NewDecoder(r.Body).Decode(cmd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := rte.CmdSvc.CreateCommand(r.Context(), cmd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

/*
Services
*/

type CommandService struct {
	Repo   *repository.Queries
	Logger *slog.Logger
	proto.UnimplementedCommandServiceServer
}

func (s *CommandService) CreateCommand(ctx context.Context, c *repository.InsertCommandParams) error {
	if c.ID != "" {
		return fmt.Errorf("%w: nonempty id", ErrBadId)
	}

	c.ID = uuid.NewString()

	return s.Repo.InsertCommand(ctx, *c)
}

func (s CommandService) LogCommand(ctx context.Context, c *proto.Command) (*emptypb.Empty, error) {
	p := &repository.InsertCommandParams{
		Aliased:         c.Aliased,
		ExpandedPreview: c.ExpandedPreview,
		ExpandedFull:    c.ExpandedFull,
		StartTs:         c.StartTs.String(),
		EndTs:           c.EndTs.String(),
		PeerIp:          c.PeerInfo.Ip,
		TmuxSession:     c.TmuxIndo.Session,
		TmuxWindow:      c.TmuxIndo.Window,
		TmuxPane:        c.TmuxIndo.Pane,
	}

	if err := s.CreateCommand(ctx, p); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
