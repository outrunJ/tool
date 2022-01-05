package tool

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestStringSplit(t *testing.T) {
	assert := assert.New(t)
	rst := StringSplit(",1,",",")
	assert.Equal([]string{"","1",""}, rst)

	rst = StringSplit(",1,3",",")
	assert.Equal([]string{"","1","3"}, rst)

	rst = StringSplit("1,2",",")
	assert.Equal([]string{"1","2"}, rst)

	rst = StringSplit("",",")
	assert.Equal([]string{""}, rst)

	rst = StringSplit(",,",",")
	assert.Equal([]string{"","",""}, rst)

	rst = StringSplit("1,,2,,",",,")
	assert.Equal([]string{"1","2",""}, rst)
}