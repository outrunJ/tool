package tool

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncrypt(t *testing.T) {
	var key = []byte("00000000")

	assert := assert.New(t)
	str := "abc"

	bits, err := Encrypt([]byte(str), key)
	assert.NoError(err)

	rstBits, err := Decrypt(bits, key)
	assert.NoError(err)

	assert.Equal(str, string(rstBits))
}

func TestSign(t *testing.T) {
	assert := assert.New(t)
	s := Sign([]byte(""), []byte(""))
	assert.Equal("2jmj7l5rSw0yVb/vlWAYkK/YBwk=", s)
}
