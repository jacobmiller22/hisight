-- name: CommandById :one
SELECT id, aliased, expanded_preview, expanded_full, start_ts, end_ts, peer_ip, tmux_session, tmux_window, tmux_pane
FROM commands
WHERE id = ?
LIMIT 1;

-- name: ListCommands :many
SELECT id, aliased, expanded_preview, expanded_full, start_ts, end_ts, peer_ip, tmux_session, tmux_window, tmux_pane
FROM commands;

-- name: InsertCommand :exec
INSERT INTO commands
(id, aliased, expanded_preview, expanded_full, start_ts, end_ts, peer_ip, tmux_session, tmux_window, tmux_pane)
VALUES (
?, ?, ?, ?, ?, ?, ?, ?, ?, ?
)
RETURNING id, aliased, expanded_preview, expanded_full, start_ts, end_ts, peer_ip, tmux_session, tmux_window, tmux_pane;

-- name: SearchCommandByAliased :many
SELECT id, aliased, expanded_preview, expanded_full, start_ts, end_ts, peer_ip, tmux_session, tmux_window, tmux_pane
FROM commands
WHERE aliased = ?;
