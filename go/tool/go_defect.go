package tool

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime/debug"
)

// ok err
func Err(_ interface{}, err error) error {
	return err
}
func OkErrToErr(oi interface{}, ok bool, err error) (interface{}, error) {
	if err != nil {
		return oi, err
	}

	if !ok {
		return oi, errors.New("result not ok")
	} else {
		return oi, nil
	}
}
func OkToErr(oi interface{}, ok bool) (interface{}, error) {
	return OkErrToErr(oi, ok, nil)
}
func PanicErr(err error) {
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
}

func PanicErrWith(o interface{}, err error) interface{} {
	if err != nil {
		debug.PrintStack()
		panic(err)
	}
	return o
}

// debug
func DebugBP(oiList ...interface{}) {
	fmt.Printf("%v", oiList)
	return
}

// func
func Identity(oi interface{}) interface{} {
	return oi
}
