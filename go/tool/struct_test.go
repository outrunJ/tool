package tool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFnTypeList(t *testing.T) {
	l := NewFnTypeList()
	l.PushBack(nil)
	l.Exec()
}

func TestSerialize(t *testing.T) {
	assert := assert.New(t)

	type User struct {
		Name string
		Cash int
	}
	u := &User{
		Name: "a",
		Cash: 1,
	}
	b, err := Serialize(u)
	assert.NoError(err)

	u1 := &User{}
	err = Deserialize(b, u1)
	assert.NoError(err)

	assert.Equal(u, u1)
}
