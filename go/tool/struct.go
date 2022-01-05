package tool

import (
	"bytes"
	"container/list"
	"encoding/gob"
	"reflect"
)

type TypeStruct struct {
	Typ          reflect.Type
	Val          reflect.Value
	Name         string
	Instance     interface{}
	InstanceAddr interface{}
}

func GenerateTypeStruct(oi interface{}) *TypeStruct {
	oAddrValue := reflect.ValueOf(oi)
	oValue := valueElem(oAddrValue)
	oType := oValue.Type()

	return &TypeStruct{
		Instance:     oValue.Interface(),
		InstanceAddr: oi,
		Val:          oValue,
		Typ:          oType,
		Name:         oType.Name(),
	}
}

func Serialize(o interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := gob.NewEncoder(buf)
	err := encoder.Encode(o)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Deserialize(b []byte, o interface{}) error {
	buf := &bytes.Buffer{}
	buf.Write(b)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(o)
}

type FnType func() error
type FnTypeList struct {
	l *list.List
}

func NewFnTypeList() *FnTypeList {
	return &FnTypeList{l: list.New()}
}

func (f *FnTypeList) PushFront(fn FnType) {
	l := f.l
	l.PushFront(fn)
}

func (f *FnTypeList) PushBack(fn FnType) {
	l := f.l
	l.PushBack(fn)
}

func (f *FnTypeList) Exec() []error {
	l := f.l
	var ret = make([]error, l.Len())

	ind := 0
	for ele := l.Front(); ele != nil; ele = ele.Next() {
		if fn := ele.Value.(FnType); fn != nil {
			ret[ind] = fn()
		}
		ind++
	}
	return ret
}
