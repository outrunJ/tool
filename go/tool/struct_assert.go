package tool

import (
	"github.com/stretchr/testify/assert"
	"log"
)

func AssertErrWith(assert *assert.Assertions) func(interface{}, error) interface{} {
	return func(o interface{}, err error) interface{} {
		if err != nil {
			log.Printf("assert err: %v", err)
			log.Println()
		}
		assert.Nil(err)
		return o
	}
}
func AssertOkWith(assert *assert.Assertions) func(interface{}, bool) interface{} {
	return func(o interface{}, ok bool) interface{} {
		if ok != true {
			log.Printf("assert ok false: %v", o)
			log.Println()
		}
		assert.True(ok)
		return o
	}
}
func AssertOkErrWith(assert *assert.Assertions) func(interface{}, bool, error) interface{} {
	return func(o interface{}, ok bool, err error) interface{} {
		AssertOkWith(assert)(o, ok)
		return AssertErrWith(assert)(o, err)
	}
}
