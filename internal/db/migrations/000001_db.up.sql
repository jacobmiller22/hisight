CREATE TABLE commands (
    id TEXT PRIMARY KEY NOT NULL,
    aliased TEXT NOT NULL,
    expanded TEXT NOT NULL,
    ts TEXT NOT NULL,
    shell TEXT NOT NULL
);
