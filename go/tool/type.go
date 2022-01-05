package tool

import (
	"reflect"
)

func IsZero(oi interface{}) bool {
	if oi == nil {
		return true
	}

	switch oi.(type) {
	case *string:
		return *oi.(*string) == ""
	case string:
		return oi.(string) == ""
	case *int:
		return *oi.(*int) == 0
	case *int32:
		return *oi.(*int32) == 0
	case *int64:
		return *oi.(*int64) == 0
	case int:
		return oi.(int) == 0
	case int32:
		return oi.(int32) == 0
	case int64:
		return oi.(int64) == 0
	case *float32:
		return *oi.(*float32) == 0
	case *float64:
		return *oi.(*float64) == 0
	case float32:
		return oi.(float32) == 0
	case float64:
		return oi.(float64) == 0
	default:
		return isValueZero(reflect.ValueOf(oi))
	}
}

var IsNil = IsZero

// pointer
func Elem(oi interface{}) interface{} {
	return valueInterface(valueElem(reflect.ValueOf(oi)))
}
func Addr(oi interface{}) interface{} {
	oValue := reflect.ValueOf(oi)
	if oValue.CanAddr() {
		return valueInterface(oValue.Addr())
	} else {
		return oi
	}
}
