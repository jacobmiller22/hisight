package repository

import (
	"database/sql"
	"errors"
)

type statements struct {
	command commandStatements
}

func (s *statements) Close() error {
	errs := make([]error, 1)

	errs = append(errs, s.command.Close())

	return errors.Join(errs...)
}

type Repository interface {
	CommandRepository
}

type Repo struct {
	DB         *sql.DB
	statements statements
}

func NewRepo(dsn string) (*Repo, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)

	return &Repo{
		DB:         db,
		statements: statements{},
	}, nil
}

func (r *Repo) Close() {
	r.Close()
}
