package tool

import (
	"github.com/json-iterator/go"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func String2JSON(s string, o interface{}) error {
	err := JSON.Unmarshal([]byte(s), o)
	return err
}
func JSON2ByteSlice(j interface{}) ([]byte, error) {
	b, err := JSON.Marshal(j)
	if err != nil {
		return nil, err
	}
	return b, nil
}
func JSON2String(j interface{}) (string, error) {
	bs, err := JSON2ByteSlice(j)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
