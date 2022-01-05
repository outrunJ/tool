package dber

import (
	"context"
)

type dberContext struct {
	DBer
	ctx          context.Context
	newIDContext func(ctx context.Context) (string, error)
	newTenantIDContext func(ctx context.Context) (string, error)
}

func (d *dberContext) GetContext() context.Context {
	return d.ctx
}
// id
func (d *dberContext) NewID() (string, error) {
	return d.newIDContext(d.ctx)
}
func (d *dberContext) SetNewIDContext(fn func(ctx context.Context) (string, error)) {
	d.newIDContext = fn
}
func (d *dberContext) NewIDContext(ctx context.Context) (string, error) {
	return d.newIDContext(ctx)
}
func (d *dberContext) NewTenantID() (string, error) {
	return d.newTenantIDContext(d.ctx)
}
func (d *dberContext) SetNewTenantIDContext(fn func(ctx context.Context) (string, error)) {
	d.newTenantIDContext = fn
}
func (d *dberContext) NewTenantIDContext(ctx context.Context) (string, error) {
	return d.newTenantIDContext(ctx)
}


func (d *dberContext) Exec(query string, args ...interface{}) (Result, error) {
	return d.ExecContext(d.ctx, query, args...)
}
func (d *dberContext) Query(query string, args ...interface{}) (Rows, error) {
	return d.QueryContext(d.ctx, query, args...)
}
func (d *dberContext) QueryRow(query string, args ...interface{}) Row {
	return d.QueryRowContext(d.ctx, query, args...)
}
func (d *dberContext) Tx(fn func(DBer) error, rb func() error) (err error, dbErr error, rbErr error) {
	return d.TxContext(d.ctx, nil, fn, rb)
}
func (d *dberContext) SQL(sql string) SQL {
	return &sqlContext{
		SQL: d.DBer.SQL(sql),
		ctx: d.ctx,
	}
}
func (d *dberContext) SQLQuery(sql string, args ...interface{}) SQLQuery {
	return d.DBer.SQLQueryContext(d.ctx, sql, args...)
}

type sqlContext struct {
	SQL
	ctx context.Context
}

func (s *sqlContext) Query(args ...interface{}) SQLQuery {
	return s.QueryContext(s.ctx, args...)
}
