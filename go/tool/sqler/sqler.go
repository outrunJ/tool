package sqler

import (
	"fmt"
	"git.meiqia.com/business_platform/tool"
	"sort"
	"strings"
	"time"
)

type Sqle struct {
	model          interface{}
	table          string
	nameSlice      *[]string
	nameSliceAll   *[]string
	nameTagAll     map[string]*[]string
	nameValue      map[string]interface{}
	nameValueAll   map[string]interface{}
	nameColumn     map[string]string
	nameColumnAll  map[string]string
	columnValue    map[string]interface{}
	columnValueAll map[string]interface{}
}

func tableNameAmend(name string) string {
	if strings.HasPrefix(name, "db_") {
		return name[3:]
	}
	return name
}

func (s *Sqle) initNameColumn() {
	s.nameColumn = map[string]string{}
	s.nameColumnAll = map[string]string{}
	for name := range s.nameValueAll {
		column := (*s.nameTagAll[name])[0]
		if column == "" {
			column = tool.CamelCaseToUnderlineCase(name)
		}

		s.nameColumnAll[name] = column
		if s.nameValue[name] != nil {
			s.nameColumn[name] = column
		}
	}
}

func (s *Sqle) initColumnValue() {
	s.columnValueAll = map[string]interface{}{}
	for name, value := range s.nameValueAll {
		s.columnValueAll[s.nameColumnAll[name]] = value
	}
	s.columnValue = *tool.MapCloneDeleteZero_KStringVInterface(&s.columnValueAll)
}

func (s *Sqle) initModel() {
	flattenSkip := func(oi interface{}) bool {
		switch oi.(type) {
		case time.Time:
			return true
		default:
			return false
		}
	}
	nameValueAll, _, nameTagAll := tool.FieldMap(s.model, flattenSkip, "dbName")
	table := tool.StructName(s.model)
	if table != "" {
		table = tool.Pluralize(tool.CamelCaseToUnderlineCase(table))
	}

	s.table = tableNameAmend(table)
	s.nameTagAll = nameTagAll
	s.nameValue = *tool.MapCloneDeleteZero_KStringVInterface(&nameValueAll)
	s.nameValueAll = nameValueAll
	s.nameSlice = tool.MapKeys_KStringVInterface(&s.nameValue)
	s.nameSliceAll = tool.MapKeys_KStringVInterface(&s.nameValueAll)
	sort.Strings(*s.nameSlice)
	sort.Strings(*s.nameSliceAll)

	s.initNameColumn()
	s.initColumnValue()
}

func Sqler(oi interface{}) *Sqle {
	sqle := &Sqle{
		model: oi,
	}
	sqle.initModel()
	return sqle
}

// table
func (s *Sqle) Table() string {
	return (*Quote("`", s.table))[0]
}

// name
func (s *Sqle) names(all bool, names ...string) *[]string {
	var nameSlice *[]string
	if all {
		nameSlice = s.nameSliceAll
	} else {
		nameSlice = s.nameSlice
	}
	if len(names) == 0 {
		return nameSlice
	}

	return tool.SliceIntersection_TString(nameSlice, &names)
}
func (s *Sqle) namesEx(all bool, names ...string) *[]string {
	return tool.SliceWithout_TString(
		s.names(all),
		*s.names(all, names...)...,
	)
}

// name column
func (s *Sqle) columns(all bool, names ...string) *[]string {
	var nameColumn map[string]string
	if all {
		nameColumn = s.nameColumnAll
	} else {
		nameColumn = s.nameColumn
	}

	// sort
	if len(names) == 0 {
		names = *tool.MapKeys_KStringVString(&nameColumn)
	}
	tool.SliceSort_TString(&names)

	return Quote("`", *tool.Slice2Slice_TStringRTInterface(tool.MapValues_KStringVString(&nameColumn, names...))...)
}

func (s *Sqle) columnsEx(all bool, names ...string) *[]string {
	return tool.SliceWithout_TString(
		s.columns(all),
		*s.columns(all, names...)...)
}

func (s *Sqle) Columns(names ...string) *[]string {
	return s.columns(false, names...)
}

func (s *Sqle) ColumnsEx(names ...string) *[]string {
	return s.columnsEx(false, names...)
}

func (s *Sqle) AllColumns(names ...string) *[]string {
	return s.columns(true, names...)
}

func (s *Sqle) AllColumnsEx(names ...string) *[]string {
	return s.columnsEx(true, names...)
}

// name value
func (s *Sqle) values(all bool, names ...string) *[]interface{} {
	names = *s.names(all, names...)
	values := make([]interface{}, len(names))
	var value interface{}
	for ind, name := range names {
		value = s.nameValueAll[name]
		switch value.(type) {
		case *time.Time:
			v := value.(*time.Time)
			value = v.String()
		}
		values[ind] = value
	}
	return &values
}
func (s *Sqle) valuesEx(all bool, names ...string) *[]interface{} {
	names = *s.namesEx(all, names...)
	return s.values(true, names...)
}
func (s *Sqle) Values(names ...string) *[]interface{} {
	return s.values(false, names...)
}
func (s *Sqle) ValuesEx(names ...string) *[]interface{} {
	return s.valuesEx(false, names...)
}

// clause
func (s *Sqle) ClauseColumns(names ...string) string {
	return strings.Join(*s.Columns(names...), " , ")
}

func (s *Sqle) ClauseColumnsEx(names ...string) string {
	return strings.Join(*s.ColumnsEx(names...), " , ")
}

func (s *Sqle) ClauseAllColumns(names ...string) string {
	return strings.Join(*s.AllColumns(names...), " , ")
}

func (s *Sqle) ClauseAllColumnsEx(names ...string) string {
	return strings.Join(*s.AllColumnsEx(names...), " , ")
}

func (s *Sqle) ClausePatterns(names ...string) string {
	columns := s.Columns(names...)
	return Pattern(len(*columns))
}

func (s *Sqle) ClausePatternsEx(names ...string) string {
	columns := s.ColumnsEx(names...)
	return Pattern(len(*columns))
}

// clause where
func (s *Sqle) clauseEqualAnd(columns ...string) string {
	equalColumns, _ := tool.SliceMap_TStringRTString(&columns, func(s string) (string, error) {
		return s + " = ?", nil
	})
	return strings.Join(*equalColumns, " and ")
}
func (s *Sqle) clauseEqualComma(columns ...string) string {
	equalColumns, _ := tool.SliceMap_TStringRTString(&columns, func(s string) (string, error) {
		return s + " = ?", nil
	})
	return strings.Join(*equalColumns, " , ")
}

func (s *Sqle) ClauseEqualAnd(names ...string) string {
	return s.clauseEqualAnd(*s.Columns(names...)...)
}

func (s *Sqle) ClauseEqualAndEx(names ...string) string {
	return s.clauseEqualAnd(*s.ColumnsEx(names...)...)
}

func (s *Sqle) ClauseEqualComma(names ...string) string {
	return s.clauseEqualComma(*s.Columns(names...)...)
}

func (s *Sqle) ClauseEqualCommaEx(names ...string) string {
	return s.clauseEqualComma(*s.ColumnsEx(names...)...)
}

// sql
func (s *Sqle) SqlSelect() string {
	return fmt.Sprintf(" SELECT %v FROM %v ", s.ClauseColumns(), s.Table())
}

func (s *Sqle) SqlSelectWhereEqualAnd() string {
	return fmt.Sprintf(" SELECT %v FROM %v WHERE %v ", s.ClauseColumns(), s.Table(), s.clauseEqualAnd())
}

func (s *Sqle) SqlInsert() string {
	return fmt.Sprintf(" INSERT INTO %v (%v) VALUES (%v) ", s.Table(), s.ClauseColumns(), s.ClausePatterns())
}

func (s *Sqle) SqlDelete() string {
	return fmt.Sprintf(" DELETE FROM %v WHERE %v ", s.Table(), s.ClauseEqualAnd())
}

func (s *Sqle) SqlUpdate() string {
	return fmt.Sprintf(
		" UPDATE %v (%v) SET %v ",
		s.Table(),
		s.Columns(),
		s.ClauseEqualComma())
}

func (s *Sqle) SqlUpdateWhereID() string {
	return s.SqlUpdate() + " WHERE id = ? "
}
