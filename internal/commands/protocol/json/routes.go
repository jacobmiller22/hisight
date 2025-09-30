package json

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jacobmiller22/hisight/internal/commands"
)

var ErrInvalidId error = errors.New("invalid id")

type Command struct {
	Shell   string `json:"shell"`
	Command string `json:"command"`
}

type CommandRoutes struct {
	CmdSvc *commands.CommandService
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
			http.Error(w, fmt.Sprintf("%v: empty id", ErrInvalidId), http.StatusBadRequest)
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
		var jsonCmd Command
		if err := json.NewDecoder(r.Body).Decode(&jsonCmd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		cmd := &commands.Command{
			Shell:   jsonCmd.Shell,
			Command: jsonCmd.Command,
		}

		if err := rte.CmdSvc.CreateCommand(r.Context(), cmd); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
