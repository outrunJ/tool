package dber

import (
	"git.meiqia.com/business_platform/tool"
	"reflect"
	"sync"
	"time"
	"context"
)

type SQL interface {
	QueryContext(ctx context.Context, args ...interface{}) SQLQuery
	Query(args ...interface{}) SQLQuery
}

type dberSQL struct {
	dber DBer
	sql  string
	model func(oi interface{}, fields *[]string, values *[]interface{}) (error)
	closemu sync.RWMutex
}

func (r *dberSQL) QueryContext(ctx context.Context, args ...interface{}) SQLQuery {
	return &dberSQLQuery{
		sql:  r,
		args: args,
		result: &SQLQueryResult{
			closed: false,
		},
		ctx:   ctx,
		model: r.model,
	}
}
func (r *dberSQL) Query(args ...interface{}) SQLQuery {
	return r.QueryContext(context.Background(), args...)
}

type SQLQuery interface {
	query() error
	close() error
	Columns() *[]string
	Each(each func(func(...interface{}) error) error) error
	Model(oi interface{}) error
}

type SQLQueryResult struct {
	rows    Rows
	columns *[]string
	types   *[]reflect.Type
	closed  bool
}

type dberSQLQuery struct {
	sql    *dberSQL
	args   []interface{}
	result *SQLQueryResult
	ctx    context.Context
	model func(oi interface{}, fields *[]string, values *[]interface{}) (error)
}

func (r *dberSQLQuery) query() error {
	if r.result.rows != nil {
		return nil
	}

	r.sql.closemu.RLock()
	defer r.sql.closemu.RUnlock()
	rows, err := r.sql.dber.QueryContext(r.ctx, r.sql.sql, r.args...)
	if err != nil {
		return err
	}
	r.result.rows = rows

	rawColumnTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	columns := make([]string, len(rawColumnTypes))
	types := make([]reflect.Type, len(rawColumnTypes))
	for ind, rawColumnType := range rawColumnTypes {
		columns[ind] = rawColumnType.Name()
		types[ind] = rawColumnType.ScanType()
	}
	r.result.columns = &columns
	r.result.types = &types

	return rows.Err()
}

func (r *dberSQLQuery) close() error {
	rows := r.result.rows
	if rows == nil {
		return nil
	}

	if !r.result.closed {
		r.result.closed = true
		return rows.Close()
	}
	return rows.Err()
}

func (r *dberSQLQuery) Columns() *[]string {
	return r.result.columns
}

func (r *dberSQLQuery) Each(each func(func(...interface{}) error) error) error {
	defer r.close()
	err := r.query()
	if err != nil {
		return err
	}

	rows := r.result.rows
	for rows.Next() {
		rows.ColumnTypes()
		err = each(rows.Scan)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

func (r *dberSQLQuery) Model(oi interface{}) error {
	defer r.close()
	err := r.query()
	if err != nil {
		return err
	}

	isSlice := tool.IsSlice(oi)
	columnName := map[string]string{}
	columnType := map[string]reflect.Type{}
	columns := r.result.columns
	types := r.result.types

	analyzed := false
	analyze := func() {
		if analyzed {
			return
		}
		analyzed = true
		var nameType map[string]reflect.Type
		tmpl := oi
		if isSlice {
			tmpl = tool.SliceTypeNewAddr(oi)
		}
		tool.FieldFillZero(tmpl)
		nameType = tool.FieldTypeMap(tmpl, func(oi interface{}) bool {
			switch oi.(type) {
			case time.Time:
				return true
			default:
				return false
			}
		})
		for name, typ := range nameType {
			column := tool.CamelCaseToUnderlineCase(name)
			columnName[column] = name
			columnType[column] = typ
		}
	}

	makeValues := func() *[]interface{} {
		values := make([]interface{}, len(*columns))
		for ind, column := range *columns {
			typ := columnType[column]
			if typ == nil {
				typ = (*types)[ind]
			}
			values[ind] = WrapNil(tool.StructTypeNewAddr(typ))
		}
		return &values
	}

	eachFn := func(each func(...interface{}) error) error {
		// record
		var record interface{}
		if isSlice {
			record = tool.SliceTypeNewAddr(oi)
			tool.SliceAppend(oi, record)
		} else {
			record = oi
		}
		tool.FieldFillZero(record)

		// model fill
		if r.model != nil {
			l := len(*r.result.columns)
			values := make([]interface{}, l)
			err := r.model(record, r.result.columns, &values)
			if err == nil {
				return each(values...)
			}
		}

		// reflect fill
		analyze()
		values := makeValues()
		err := each(*values...)
		if err != nil {
			return err
		}
		for ind, column := range *columns {
			tool.SetField(record, columnName[column], UnwrapNil((*values)[ind]))
		}
		return nil
	}

	err = r.Each(eachFn)
	if err != nil {
		return err
	}
	return r.result.rows.Err()
}
