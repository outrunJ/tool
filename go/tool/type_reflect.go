package tool

import (
	"reflect"
	"time"
)

// poly
func isPtr(i interface{}) bool {
	switch i.(type) {
	case reflect.Value:
		return i.(reflect.Value).Kind() == reflect.Ptr
	case reflect.Type:
		return i.(reflect.Type).Kind() == reflect.Ptr
	default:
		return reflect.TypeOf(i).Kind() == reflect.Ptr
	}
}

// type
func isTypeSame(type1 reflect.Type, type2 reflect.Type) bool {
	return type1.AssignableTo(type2) && type2.AssignableTo(type1)
}
func isTypeTime(typ reflect.Type) bool {
	return typ.Kind() == reflect.Struct && typ.Name() == "Time"
}

func typeElem(typ reflect.Type) reflect.Type {
	for isPtr(typ) {
		typ = typ.Elem()
	}
	return typ
}

func valueZeroValue(value reflect.Value) reflect.Value {
	typ := value.Type()
	return typeZeroValue(typ)
}

func typeZeroValue(typ reflect.Type) reflect.Value {
	kind := typ.Kind()
	switch kind {
	case reflect.String:
		return reflect.ValueOf("")
	}

	// can set time.Time
	return reflect.Zero(typ)
}
func typeZeroValueAddr(typ reflect.Type) reflect.Value {
	return reflect.New(typ)
}

// value
func isValueZero(value reflect.Value) bool {
	if !value.IsValid() {
		return true
	}

	valueKind := value.Kind()
	switch valueKind {
	case
		reflect.Ptr,
		reflect.Chan,
		reflect.Func,
		reflect.Interface,
		reflect.Map,
		reflect.Slice:
		return value.IsNil()
	}

	if !value.CanInterface() {
		return true
	}

	valueInterfac := value.Interface()
	switch valueInterfac.(type) {
	case string:
		v := valueInterfac.(string)
		return v == ""
	case *string:
		v := valueInterfac.(*string)
		return *v == ""
	case time.Time:
		v := valueInterfac.(time.Time)
		return v == time.Time{}
		//case *time.Time:
		//	v := valueInterfac.(*time.Time)
		//	return *v == time.Time{}
	}

	return false
}

func valueSetZero(value reflect.Value) {
	if isValueZero(value) {
		return
	}

	if value.CanSet() {
		value.Set(valueZeroValue(value))
	}
}

func valueElem(value reflect.Value) reflect.Value {
	for isPtr(value) {
		value = value.Elem()
	}
	return value
}
func valueAddr(value reflect.Value) reflect.Value {
	if isPtr(value) {
		return value
	} else if value.CanAddr() {
		return value.Addr()
	} else {
		return value
	}
}
func valueInterface(value reflect.Value) interface{} {
	if value.IsValid() && value.CanInterface() {
		return value.Interface()
	} else {
		return nil
	}
}
func isValueSame(value1 reflect.Value, value2 reflect.Value) bool {
	return reflect.DeepEqual(valueInterface(value1), valueInterface(value2))
}

// struct
func isStruct(value reflect.Value) bool {
	return value.Kind() == reflect.Struct
}
func isTypeStruct(typ reflect.Type) (bool) {
	if typ.Kind() == reflect.Struct {
		return true
	} else {
		return false
	}
}
func structNewAddr(typ reflect.Type) reflect.Value {
	return reflect.New(typ)
}
func fieldEach(value reflect.Value, each func(reflect.Value, reflect.StructField) error, flatten bool, notFlatten func(reflect.Value) bool) error {
	var err error
	value = valueElem(value)
	typ := value.Type()
	num := value.NumField()
	for i := 0; i < num; i++ {
		fieldValue := value.Field(i)
		fieldStruct := typ.Field(i)
		fieldValueElem := valueElem(fieldValue)
		fieldTypeElem := typeElem(fieldStruct.Type)
		if flatten && isTypeStruct(fieldTypeElem) && (notFlatten == nil || !notFlatten(fieldValueElem)) {
			if !isValueZero(fieldValue) {
				err = fieldEach(fieldValue, each, flatten, notFlatten)
				if err != nil {
					return err
				}
			}
		} else {
			err = each(fieldValue, fieldStruct)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func fieldTypeMap(value reflect.Value, flattenSkip func(reflect.Value) bool) map[string]reflect.Type {
	typeMap := map[string]reflect.Type{}
	fieldEach(value, func(v reflect.Value, s reflect.StructField) error {
		typeMap[s.Name] = v.Type()
		return nil
	}, true, flattenSkip)
	return typeMap
}
func fieldMap(value reflect.Value, flattenSkip func(reflect.Value) bool, tagNameSlice ...string) (map[string]interface{}, map[string]reflect.Type, map[string]*[]string) {
	tagLength := len(tagNameSlice)

	valueMap := map[string]interface{}{}
	typeMap := map[string]reflect.Type{}
	tagMap := map[string]*[]string{}

	fieldEach(value, func(v reflect.Value, s reflect.StructField) error {
		valueMap[s.Name] = valueInterface(v)
		typeMap[s.Name] = v.Type()
		tagSlice := make([]string, tagLength)
		tagMap[s.Name] = &tagSlice

		for ind, tagName := range tagNameSlice {
			if tag, ok := s.Tag.Lookup(tagName); ok {
				tagSlice[ind] = tag
			}
		}
		return nil
	}, true, flattenSkip)
	return valueMap, typeMap, tagMap
}
func getField(value reflect.Value, fieldName string) reflect.Value {
	value = valueElem(value)
	retValue := value.FieldByName(fieldName)
	if retValue.IsValid() {
		return valueAddr(retValue)
	} else {
		return retValue
	}
}
func setField(value reflect.Value, fieldName string, fieldValue reflect.Value) {
	value = valueElem(value)
	field := value.FieldByName(fieldName)

	hasSet := false
	set := func(fieldValue reflect.Value) {
		if isTypeSame(field.Type(), fieldValue.Type()) {
			field.Set(fieldValue)
			hasSet = true
		}
	}

	if field.CanSet() {
		if fieldValue.IsValid() {
			set(fieldValue)
			if !hasSet {
				set(valueElem(fieldValue))
			}
			if !hasSet {
				set(valueAddr(fieldValue))
			}
		} else {
			valueSetZero(field)
		}
	}
}
func asField(value reflect.Value, typ reflect.Type, fieldName string) reflect.Value {
	retValue := reflect.New(typ)
	setField(retValue, fieldName, value)
	return retValue
}
func isFieldZero(value reflect.Value, fieldName string) bool {
	fieldValue := valueElem(getField(value, fieldName))
	return isValueZero(fieldValue)
}

// slice
func isSlice(i reflect.Value) bool {
	if i.Kind() == reflect.Slice {
		return true
	} else {
		return false
	}
}

func isSliceType(t reflect.Type) bool {
	if t.Kind() == reflect.Slice {
		return true
	} else {
		return false
	}
}

func sliceType(t reflect.Type) reflect.Type {
	if isSliceType(t) {
		return t.Elem()
	} else {
		return t
	}
}

func sliceNewAddrOfType(typ reflect.Type) reflect.Value {
	sliceValue := reflect.MakeSlice(reflect.SliceOf(typ), 0, 0)
	sliceValueAddr := reflect.New(sliceValue.Type())
	return sliceValueAddr
}
func sliceLen(sliceValue reflect.Value) int {
	if !isSlice(sliceValue) {
		return 0
	}
	return sliceValue.Len()
}
func sliceAppend(sliceValue reflect.Value, itemValue reflect.Value) {
	sliceValue = valueElem(sliceValue)
	sliceValue.Set(reflect.Append(sliceValue, itemValue))
}
func sliceForEach(sliceValue reflect.Value, fn func(reflect.Value) error) error {
	length := sliceLen(sliceValue)
	for i := 0; i < length; i++ {
		err := fn(sliceValue.Index(i))
		if err != nil {
			return err
		}
	}
	return nil
}
func sliceMap(sliceValue reflect.Value, typ reflect.Type, fn func(reflect.Value) (reflect.Value, error)) (reflect.Value, error) {
	length := sliceLen(sliceValue)
	retSliceValue := reflect.MakeSlice(reflect.SliceOf(typ), length, length)
	for i := 0; i < length; i++ {
		retValue, err := fn(sliceValue.Index(i))
		if err != nil {
			return retValue, err
		}
		retSliceValue.Index(i).Set(retValue)
	}
	return retSliceValue, nil
}
func sliceIndexOf(sliceValue reflect.Value, value reflect.Value) int {
	length := sliceLen(sliceValue)
	for i := 0; i < length; i++ {
		if isValueSame(sliceValue.Index(i), value) {
			return i
		}
	}
	return -1
}
