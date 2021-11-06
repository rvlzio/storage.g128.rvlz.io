package infrastructure

import (
	"context"
	"github.com/jackc/pgx/v4"
)

type Context struct {
	transaction pgx.Tx
}

func (ctx *Context) SaveChanges() {
	ctx.transaction.Commit(context.Background())
}

func (ctx *Context) Prepare(name, sql string) error {
	_, err := ctx.transaction.Prepare(context.Background(), name, sql)
	return err
}

func (ctx *Context) Execute(name string, args ...interface{}) error {
	_, err := ctx.transaction.Exec(context.Background(), name, args...)
	return err
}

func (ctx *Context) QueryOne(sql string, args ...interface{}) pgx.Row {
	return ctx.transaction.QueryRow(context.Background(), sql, args...)
}

func NewContext(transaction pgx.Tx) *Context {
	return &Context{transaction: transaction}
}
