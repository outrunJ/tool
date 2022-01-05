package tool

import (
	"strconv"
	"fmt"
	"github.com/cheekybits/genny/generic"
)

type Typ generic.Type

// type convert
func string2string(v string) (string) {
	return v
}
func String2String(v string) (string) {
	return string2string(v)
}
func bool2bool(v bool) (bool) {
	return v
}
func Bool2Bool(v bool) (bool) {
	return bool2bool(v)
}
func int2int(v int) (int) {
	return v
}
func Int2Int(v int) (int) {
	return int2int(v)
}
func interface2interface(v interface{}) (interface{}) {
	return v
}
func Interface2Interface(v interface{}) (interface{}) {
	return interface2interface(v)
}

func string2bool(v string) (bool) {
	ret, _ := strconv.ParseBool(v)
	return ret
}
func String2Bool(v string) (bool) {
	return string2bool(v)
}
func bool2string(v bool) (string) {
	return fmt.Sprintf("%v", v)
}
func Bool2String(v bool) (string) {
	return bool2string(v)
}
func string2int(v string) (int) {
	ret, _ := strconv.ParseInt(v, 10, 0)
	return int(ret)
}
func String2Int(v string) (int) {
	return string2int(v)
}
func int2string(v int) (string) {
	return fmt.Sprintf("%v", v)
}
func Int2String(v int) (string) {
	return int2string(v)
}
func string2interface(v string) (interface{}) {
	return v
}
func String2Interface(v string) (interface{}) {
	return string2interface(v)
}
func interface2string(v interface{}) (string) {
	return fmt.Sprintf("%v", v)
}
func Interface2String(v interface{}) (string) {
	return interface2string(v)
}

func bool2int(v bool) (int) {
	if v {
		return 1
	} else {
		return 0
	}
}
func Bool2Int(v bool) (int) {
	return bool2int(v)
}
func int2bool(v int) (bool) {
	if v == 0 {
		return false
	} else {
		return true
	}
}
func Int2Bool(v int) (bool) {
	return int2bool(v)
}
func bool2interface(v bool) (interface{}) {
	return v
}
func Bool2Interface(v bool) (interface{}) {
	return bool2interface(v)
}
func interface2bool(v interface{}) (bool) {
	return string2bool(interface2string(v))
}
func Interface2Bool(v interface{}) (bool) {
	return interface2bool(v)
}
func int2interface(v int) (interface{}) {
	return v
}
func Int2Interface(v int) (interface{}) {
	return int2interface(v)
}
func interface2int(v interface{}) (int) {
	return string2int(interface2string(v))
}
func Interface2Int(v interface{}) (int) {
	return interface2int(v)
}
