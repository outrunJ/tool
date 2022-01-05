package dber

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

type NullBool struct {
	sql.NullBool
	Bool  *bool
	Valid bool
}

func (ns *NullBool) Scan(value interface{}) error {
	err := ns.NullBool.Scan(value)
	if err != nil {
		return err
	}
	*ns.Bool = ns.NullBool.Bool
	ns.Valid = ns.NullBool.Valid
	return nil
}
func (ns NullBool) Value() (driver.Value, error) {
	return ns.NullBool.Value()
}

type NullString struct {
	sql.NullString
	String *string
	Valid  bool
}

func (ns *NullString) Scan(value interface{}) error {
	err := ns.NullString.Scan(value)
	if err != nil {
		return err
	}
	*ns.String = ns.NullString.String
	ns.Valid = ns.NullString.Valid
	return nil
}
func (ns NullString) Value() (driver.Value, error) {
	return ns.NullString.Value()
}

type NullInt64 struct {
	sql.NullInt64
	Int64 *int64
	Valid bool
}

func (ns *NullInt64) Scan(value interface{}) error {
	err := ns.NullInt64.Scan(value)
	if err != nil {
		return err
	}
	*ns.Int64 = ns.NullInt64.Int64
	ns.Valid = ns.NullInt64.Valid
	return nil
}
func (ns NullInt64) Value() (driver.Value, error) {
	return ns.NullInt64.Value()
}

type NullInt struct {
	sql.NullInt64
	Int   *int
	Valid bool
}

func (ns *NullInt) Scan(value interface{}) error {
	err := ns.NullInt64.Scan(value)
	if err != nil {
		return err
	}
	*ns.Int = int(ns.NullInt64.Int64)
	ns.Valid = ns.NullInt64.Valid
	return nil
}
func (ns NullInt) Value() (driver.Value, error) {
	v, err := ns.NullInt64.Value()
	if err != nil {
		return v, err
	}
	return driver.Value(int(v.(int64))), nil
}

type NullInt32 struct {
	sql.NullInt64
	Int32 *int32
	Valid bool
}

func (ns *NullInt32) Scan(value interface{}) error {
	err := ns.NullInt64.Scan(value)
	if err != nil {
		return err
	}
	*ns.Int32 = int32(ns.NullInt64.Int64)
	ns.Valid = ns.NullInt64.Valid
	return nil
}
func (ns NullInt32) Value() (driver.Value, error) {
	v, err := ns.NullInt64.Value()
	if err != nil {
		return v, err
	}
	return driver.Value(int(v.(int32))), nil
}

type NullFloat64 struct {
	sql.NullFloat64
	Float64 *float64
	Valid   bool
}

func (ns *NullFloat64) Scan(value interface{}) error {
	err := ns.NullFloat64.Scan(value)
	if err != nil {
		return err
	}
	*ns.Float64 = ns.NullFloat64.Float64
	ns.Valid = ns.NullFloat64.Valid
	return nil
}
func (ns NullFloat64) Value() (driver.Value, error) {
	return ns.NullFloat64.Value()
}

type NullFloat32 struct {
	sql.NullFloat64
	Float32 *float32
	Valid   bool
}

func (ns *NullFloat32) Scan(value interface{}) error {
	err := ns.NullFloat64.Scan(value)
	if err != nil {
		return err
	}
	*ns.Float32 = float32(ns.NullFloat64.Float64)
	ns.Valid = ns.NullFloat64.Valid
	return nil
}
func (ns NullFloat32) Value() (driver.Value, error) {
	v, err := ns.NullFloat64.Value()
	if err != nil {
		return v, err
	}
	return driver.Value(float32(v.(float64))), nil
}

type NullTime struct {
	Time  *time.Time
	Valid bool
}

func (n *NullTime) Scan(value interface{}) error {
	*n.Time, n.Valid = value.(time.Time)
	return nil
}
func (n NullTime) Value() (driver.Value, error) {
	if !n.Valid {
		return nil, nil
	}
	return n.Time, nil
}

func WrapNil(oi interface{}) interface{} {
	switch oi.(type) {
	case string:
		o := oi.(string)
		return &NullString{
			NullString: sql.NullString{},
			String:     &o,
		}
	case *string:
		o := oi.(*string)
		return &NullString{
			NullString: sql.NullString{},
			String:     o,
		}
	case bool:
		o := oi.(bool)
		return &NullBool{
			NullBool: sql.NullBool{},
			Bool:     &o,
		}
	case *bool:
		o := oi.(*bool)
		return &NullBool{
			NullBool: sql.NullBool{},
			Bool:     o,
		}
	case int:
		o := oi.(int)
		return &NullInt{
			NullInt64: sql.NullInt64{},
			Int:       &o,
		}
	case *int:
		o := oi.(*int)
		return &NullInt{
			NullInt64: sql.NullInt64{},
			Int:       o,
		}
	case int32:
		o := oi.(int32)
		return &NullInt32{
			NullInt64: sql.NullInt64{},
			Int32:     &o,
		}
	case *int32:
		o := oi.(*int32)
		return &NullInt32{
			NullInt64: sql.NullInt64{},
			Int32:     o,
		}
	case int64:
		o := oi.(int64)
		o1 := int64(o)
		return &NullInt64{
			NullInt64: sql.NullInt64{},
			Int64:     &o1,
		}
	case *int64:
		o := oi.(*int64)
		o1 := int64(*o)
		return &NullInt64{
			NullInt64: sql.NullInt64{},
			Int64:     &o1,
		}
	case float32:
		o := oi.(float32)
		return &NullFloat32{
			NullFloat64: sql.NullFloat64{},
			Float32:     &o,
		}
	case *float32:
		o := oi.(*float32)
		return &NullFloat32{
			NullFloat64: sql.NullFloat64{},
			Float32:     o,
		}
	case float64:
		o := oi.(float64)
		o1 := float64(o)
		return &NullFloat64{
			NullFloat64: sql.NullFloat64{},
			Float64:     &o1,
		}
	case *float64:
		o := oi.(*float64)
		o1 := float64(*o)
		return &NullFloat64{
			NullFloat64: sql.NullFloat64{},
			Float64:     &o1,
		}
	case time.Time:
		o := oi.(time.Time)
		return &NullTime{
			Time: &o,
		}
	case *time.Time:
		o := oi.(*time.Time)
		return &NullTime{
			Time: o,
		}
	default:
		return oi
	}
}

func UnwrapNil(oi interface{}) interface{} {
	switch oi.(type) {
	case *NullString:
		v := oi.(*NullString)
		return v.String
	case *NullBool:
		v := oi.(*NullBool)
		return v.Bool
	case *NullInt:
		v := oi.(*NullInt)
		return v.Int
	case *NullInt32:
		v := oi.(*NullInt32)
		return v.Int32
	case *NullInt64:
		v := oi.(*NullInt64)
		return v.Int64
	case *NullFloat32:
		v := oi.(*NullFloat32)
		return v.Float32
	case *NullFloat64:
		v := oi.(*NullFloat64)
		return v.Float64
	case *NullTime:
		v := oi.(*NullTime)
		return v.Time
	default:
		return oi
	}
}

func UnwrapNilList(ois *[]interface{}) *[]interface{} {
	retList := make([]interface{}, len(*ois))
	for ind, oi := range *ois {
		retList[ind] = UnwrapNil(oi)
	}
	return &retList
}
