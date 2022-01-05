package dber

import (
	"context"
	"fmt"
)

type SQLer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Exec(query string, args ...interface{}) (Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	Query(query string, args ...interface{}) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	QueryRow(query string, args ...interface{}) Row
}

type DB interface {
	SQLer
	BeginTx(ctx context.Context, opts *TxOptions) (Tx, error)
	Begin() (Tx, error)
}
type Tx interface {
	SQLer
	Rollback() error
	Commit() error
}

type DBer interface {
	SQLer
	Tx(fn func(DBer) error, rb func() error) (err error, dbErr error, rbErr error)
	TxContext(ctx context.Context, opts *TxOptions, fn func(DBer) error, rb func() error) (err error, dbErr error, rbErr error)
	SQL(query string) SQL
	SQLQuery(sql string, args ...interface{}) SQLQuery
	SQLQueryContext(ctx context.Context, sql string, args ...interface{}) SQLQuery
	// diy
	Context(ctx context.Context) DBer
	GetContext() context.Context
	SetNewID(func() (string, error))
	NewID() (string, error)
	SetNewIDContext(func(ctx context.Context) (string, error))
	NewIDContext(ctx context.Context) (string, error)
	SetNewTenantID(func() (string, error))
	NewTenantID() (string, error)
	SetNewTenantIDContext(func(ctx context.Context) (string, error))
	NewTenantIDContext(ctx context.Context) (string, error)
	SetModel(func(oi interface{}, fields *[]string, values *[]interface{}) (error))
}

type dber struct {
	SQLer
	db                 DB
	txMark             TxMark
	newID              func() (string, error)
	newIDContext       func(ctx context.Context) (string, error)
	newTenantID        func() (string, error)
	newTenantIDContext func(ctx context.Context) (string, error)
	model func(oi interface{}, fields *[]string, values *[]interface{}) (error)
}

func NewDBer(db DB) DBer {
	return &dber{
		SQLer:  db,
		db:     db,
		txMark: TX_MARK_NONE,
		newID: func() (string, error) {
			return "", fmt.Errorf("please set NewID func")
		},
		newIDContext: func(ctx context.Context) (string, error) {
			return "", fmt.Errorf("please set NewIDContext func")
		},
		newTenantID: func() (string, error) {
			return "", fmt.Errorf("please set NewTenantID func")
		},
		newTenantIDContext: func(ctx context.Context) (string, error) {
			return "", fmt.Errorf("please set NewTenantIDContext func")
		},
		model: func(oi interface{}, fields *[]string, values *[]interface{}) (error) {
			return fmt.Errorf("no model func")
		},
	}
}

// no goroutine id , no thread local control. tx outer db executing block can not avoid
//func (d *dber) Exec(query string, args ...interface{}) (sql.Result, error) {
//	if d.txMark
//
//}

func (d *dber) clone() *dber {
	return &dber{
		SQLer:        d.SQLer,
		db:           d.db,
		txMark:       d.txMark,
		newID:        d.newID,
		newIDContext: d.newIDContext,
		model:        d.model,
	}
}
func (d *dber) Context(ctx context.Context) DBer {
	return &dberContext{
		DBer:         d,
		ctx:          ctx,
		newIDContext: d.newIDContext,
	}
}
func (d *dber) GetContext() context.Context {
	return nil
}

// id
func (d *dber) SetNewID(fn func() (string, error)) {
	d.newID = fn
}
func (d *dber) NewID() (string, error) {
	return d.newID()
}
func (d *dber) SetNewIDContext(fn func(ctx context.Context) (string, error)) {
	d.newIDContext = fn
}
func (d *dber) NewIDContext(ctx context.Context) (string, error) {
	return d.newIDContext(ctx)
}
func (d *dber) SetNewTenantID(fn func() (string, error)) {
	d.newTenantID = fn
}
func (d *dber) NewTenantID() (string, error) {
	return d.newTenantID()
}
func (d *dber) SetNewTenantIDContext(fn func(ctx context.Context) (string, error)) {
	d.newTenantIDContext = fn
}
func (d *dber) NewTenantIDContext(ctx context.Context) (string, error) {
	return d.newTenantIDContext(ctx)
}
func (d *dber) SetModel(fn func(oi interface{}, fields *[]string, values *[]interface{}) (error)) {
	d.model = fn
}

func (d *dber) SQL(sql string) SQL {
	return &dberSQL{
		dber:  d,
		sql:   sql,
		model: d.model,
	}
}

func (d *dber) SQLQueryContext(ctx context.Context, sql string, args ...interface{}) SQLQuery {

	return &dberSQLQuery{
		sql: &dberSQL{
			dber:  d,
			sql:   sql,
			model: d.model,
		},
		args: args,
		result: &SQLQueryResult{
			closed: false,
		},
		ctx:   ctx,
		model: d.model,
	}
}

func (d *dber) SQLQuery(sql string, args ...interface{}) SQLQuery {
	return d.SQLQueryContext(context.Background(), sql, args...)
}
