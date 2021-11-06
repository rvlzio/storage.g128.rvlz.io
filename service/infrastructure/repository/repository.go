package repository

import (
	"errors"
	"github.com/jackc/pgx/v4"
)

type DBContext interface {
	Execute(string, ...interface{}) error
	QueryOne(string, ...interface{}) pgx.Row
	Prepare(string, string) error
}

var (
	NotFoundErr   = errors.New("not_found_error")
	RepositoryErr = errors.New("repository_error")
)
