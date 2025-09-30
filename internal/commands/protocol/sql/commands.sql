-- name: CommandById :one
SELECT id, aliased, expanded, ts, shell
FROM commands
WHERE id = ?
LIMIT 1;

-- name: ListCommands :many
SELECT id, aliased, expanded, ts, shell
FROM commands;

-- name: InsertCommand :exec
INSERT INTO commands
(id, aliased, expanded, ts, shell)
VALUES (
?, ?, ?, ?, ?
)
RETURNING id, aliased, expanded, ts, shell;

-- name: SearchCommandByAliased :many
SELECT id, aliased, expanded, ts, shell
FROM commands
WHERE aliased = ?;
