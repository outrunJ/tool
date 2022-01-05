package tool

import (
	"github.com/pkg/errors"
	"reflect"
	"sort"
)

func IsSlice(oi interface{}) bool {
	return isSlice(valueElem(reflect.ValueOf(oi)))
}

func SliceTypeNewAddr(oi interface{}) interface{} {
	typSlice := typeElem(reflect.TypeOf(oi))
	typ := typeElem(sliceType(typSlice))
	ret := valueInterface(reflect.New(typ))
	return ret
}

func SliceNewAddr(slicei interface{}) interface{} {
	return valueInterface(sliceNewAddrOfType(reflect.TypeOf(slicei)))
}
func SliceLen(slicei interface{}) int {
	return sliceLen(valueElem(reflect.ValueOf(slicei)))
}
func SliceAppend(slicei interface{}, itemi interface{}) {
	sliceAppend(reflect.ValueOf(slicei), reflect.ValueOf(itemi))
}
func SliceForEach(slicei interface{}, fn func(interface{}) error) error {
	return sliceForEach(valueElem(reflect.ValueOf(slicei)), func(itemValue reflect.Value) error {
		return fn(valueInterface(itemValue))
	})
}
func SliceMap(slicei interface{}, typei interface{}, fn func(interface{}) (interface{}, error)) (interface{}, error) {
	sliceValue := valueElem(reflect.ValueOf(slicei))
	typeType := reflect.TypeOf(typei)

	ret, err := sliceMap(sliceValue, typeType, func(itemValue reflect.Value) (reflect.Value, error) {
		iValue, err := fn(valueInterface(itemValue))
		if err != nil {
			return reflect.Zero(typeType), err
		}
		return reflect.ValueOf(iValue), nil
	})
	if err != nil {
		return nil, err
	}
	return valueInterface(ret), nil
}

func SliceFeedField(slicei interface{}, typ interface{}, fieldName string) interface{} {
	slice := slicei.([]interface{})
	retSlice := make([]interface{}, len(slice))
	tType := reflect.TypeOf(typ)

	var value reflect.Value
	for ind, val := range slice {
		value = reflect.New(tType)
		value.FieldByName(fieldName).Set(reflect.ValueOf(val))
		retSlice[ind] = valueInterface(value)
	}
	return &retSlice
}
func SlicePluckField(slicei interface{}, typ interface{}, fieldName string) (interface{}, error) {
	sliceValue := valueElem(reflect.ValueOf(slicei))
	if !isSlice(sliceValue) {
		return nil, errors.Errorf("not slice: %v", slicei)
	}

	typType := reflect.TypeOf(typ)
	typValue := reflect.ValueOf(typ)
	typIsPtr := isPtr(typValue)

	length := sliceLen(sliceValue)
	retSliceValue := reflect.MakeSlice(reflect.SliceOf(typType), length, length)
	if length != 0 {
		var value reflect.Value

		for i := 0; i < length; i++ {
			value = valueElem(sliceValue.Index(i)).FieldByName(fieldName)
			if typIsPtr {
				value = value.Addr()
			}
			retSliceValue.Index(i).Set(value)
		}
	}
	retSlice := valueInterface(retSliceValue)
	return retSlice, nil
}

// poly
func SliceSort_TString(so *[]string) {
	sort.Strings(*so)
}

func SliceSort_TBool(so *[]bool) {
	length := len(*so)
	if length < 2 {
		return
	}

	i, j := 0, length-1
	for i < j && !(*so)[i] {
		i++
	}
	if i == j {
		return
	}

	for i < j {
		for i < j && (*so)[j] {
			j--
		}
		(*so)[i] = (*so)[j]
		for i < j && !(*so)[i] {
			i++
		}
		(*so)[j] = (*so)[i]
	}
	(*so)[i] = true
}

func SliceSort_TInt(so *[]int) {
	sort.Ints(*so)
}

func qsort_with(ss *[]string, is *[]interface{}, low int, high int) {
	assign := func(from int, to int) {
		(*ss)[to] = (*ss)[from]
		(*is)[to] = (*is)[from]
	}

	pivot := (*ss)[low]
	i, j := low, high
	for i < j {
		for i < j && (*ss)[j] >= pivot {
			j--
		}
		assign(j, i)
		for i < j && (*ss)[i] >= pivot {
			i++
		}
		assign(i, j)
	}
	(*ss)[i] = pivot
	(*is)[i] = low

	qsort_with(ss, is, low, i-1)
	qsort_with(ss, is, i+1, high)
}

func SliceSort_TInterface(so *[]interface{}) {
	length := len(*so)
	ss := Slice2Slice_TInterfaceRTString(so)
	qsort_with(ss, so, 0, length)
}
