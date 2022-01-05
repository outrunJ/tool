package dber

import "reflect"

type ColumnType interface {
	Name() string
	ScanType() reflect.Type
}

type Row interface {
	Scan(...interface{}) (error)
}

type Rows interface {
	ColumnTypes() ([]ColumnType, error)
	Err() error
	Close() error
	Next() bool
	Scan(...interface{}) error
}

type Result interface {
	RowsAffected() (int64, error)
	LastInsertId() (int64, error)
}

type TxOptions struct {
	Isolation int
	ReadOnly  bool
}
