package dber

import (
	"database/sql"
	"context"
)

func TxOptions2SQLTxOptions(o *TxOptions) *sql.TxOptions {
	if o == nil {
		return nil
	}

	return &sql.TxOptions{
		Isolation: sql.IsolationLevel(o.Isolation),
		ReadOnly:  o.ReadOnly,
	}
}

type sqlRows struct {
	*sql.Rows
}

func (s *sqlRows) ColumnTypes() ([]ColumnType, error) {
	types, err := s.Rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	rst := make([]ColumnType, len(types))
	for ind, typ := range types {
		rst[ind] = typ
	}
	return rst, err
}

type sqlTx struct {
	*sql.Tx
}

func (s *sqlTx) Exec(query string, args ...interface{}) (Result, error) {
	return s.Tx.Exec(query, args...)
}

func (s *sqlTx) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return s.Tx.ExecContext(ctx, query, args...)
}

func (s *sqlTx) Query(query string, args ...interface{}) (Rows, error) {
	rows, err := s.Tx.Query(query, args...)
	return &sqlRows{rows}, err
}

func (s *sqlTx) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := s.Tx.QueryContext(ctx, query, args...)
	return &sqlRows{rows}, err
}

func (s *sqlTx) QueryRow(query string, args ...interface{}) Row {
	return s.Tx.QueryRow(query, args...)
}

func (s *sqlTx) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return s.Tx.QueryRowContext(ctx, query, args...)
}

type sqlDB struct {
	*sql.DB
}

func NewSqlDBer(driverName string, dataSourceName string, maxIdle int, maxOpen int) (DBer, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}
	if maxIdle != 0 {
		db.SetMaxIdleConns(maxIdle)
	}
	if maxOpen != 0 {
		db.SetMaxOpenConns(maxOpen)
	}

	sdb := &sqlDB{db}
	return NewDBer(sdb), nil
}

func (s *sqlDB) Begin() (Tx, error) {
	tx, err := s.DB.Begin()
	return &sqlTx{tx}, err
}

func (s *sqlDB) BeginTx(ctx context.Context, opts *TxOptions) (Tx, error) {
	tx, err := s.DB.BeginTx(ctx, TxOptions2SQLTxOptions(opts))
	return &sqlTx{tx}, err
}

func (s *sqlDB) Exec(query string, args ...interface{}) (Result, error) {
	return s.DB.Exec(query, args...)
}

func (s *sqlDB) ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error) {
	return s.DB.ExecContext(ctx, query, args ...)
}

func (s *sqlDB) Query(query string, args ...interface{}) (Rows, error) {
	rows, err := s.DB.Query(query, args...)
	return &sqlRows{rows}, err
}

func (s *sqlDB) QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error) {
	rows, err := s.DB.QueryContext(ctx, query, args...)
	return &sqlRows{rows}, err
}

func (s *sqlDB) QueryRow(query string, args ...interface{}) Row {
	return s.DB.QueryRow(query, args...)
}

func (s *sqlDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) Row {
	return s.DB.QueryRowContext(ctx, query, args...)
}
