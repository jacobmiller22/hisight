package repository

import (
	"database/sql"
	"errors"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

type Command struct {
	Id              uuid.UUID `json:"id"`
	Aliased         string    `json:"aliased"`
	ExpandedPreview string    `json:"expanded_preview"`
	ExpandedFull    string    `json:"expanded_full"`
	StartTS         string    `json:"start_ts"`
	EndTs           string    `json:"end_ts"`
	PeerIp          string    `json:"peer_ip"`
	TmuxSession     string    `json:"tmux_session"`
	TmuxWindow      string    `json:"tmux_window"`
	TmuxPane        string    `json:"tmux_pane"`
}

type commandStatements struct {
	GetCommandsStatement *sql.Stmt
}

func (s *commandStatements) Close() error {
	errs := make([]error, 1)

	errs = append(errs, s.GetCommandsStatement.Close())

	return errors.Join(errs...)
}

func prepareCommandStatements(db *sql.DB) (*commandStatements, error) {
	getCommandsStatement, err := db.Prepare("SELECT id, aliased, expanded_preview, expanded_full, start_ts, end_ts, peer_ip, tmux_session, tmux_window, tmux_pane FROM COMMANDS LIMIT 10")

	if err != nil {
		return nil, err
	}
	return &commandStatements{GetCommandsStatement: getCommandsStatement}, nil

}

type CommandRepository interface {
	GetCommands() ([]Command, error)
}

func (r *Repo) GetCommands() ([]Command, error) {
	rows, err := r.statements.command.GetCommandsStatement.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var commands []Command

	for rows.Next() {
		var cmd Command
		rows.Scan(&cmd.Id, &cmd.Aliased, &cmd.ExpandedPreview, &cmd.ExpandedFull, &cmd.StartTS, &cmd.EndTs, &cmd.PeerIp, &cmd.TmuxSession, &cmd.TmuxWindow, &cmd.TmuxPane)

		commands = append(commands, cmd)

	}

	return commands, rows.Err()
}
