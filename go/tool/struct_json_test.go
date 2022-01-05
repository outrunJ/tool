package tool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJSON2String(t *testing.T) {
	assert := assert.New(t)
	s, err := JSON2String(map[string]int{"a": 1})
	assert.NoError(err)
	assert.Equal("{\"a\":1}", s)
}
