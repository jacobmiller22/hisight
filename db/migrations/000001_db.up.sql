CREATE TABLE commands (
    id TEXT PRIMARY KEY NOT NULL,
    aliased TEXT NOT NULL,
    expanded_preview TEXT NOT NULL,
    expanded_full TEXT NOT NULL,
    start_ts TEXT NOT NULL,
    end_ts TEXT NOT NULL,
    peer_ip TEXT NOT NULL,
    tmux_session TEXT NOT NULL,
    tmux_window TEXT NOT NULL,
    tmux_pane TEXT NOT NULL
);
