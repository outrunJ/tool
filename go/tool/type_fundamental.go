package tool

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// os
func PathCurrent() (string, error) {
	path, err := exec.LookPath(os.Args[0])
	return string(path[0:strings.LastIndex(path, "/")]), err
}

// time
func TimeString(t *time.Time) string {
	return t.String()[:19]
}

// int
func ParseUint(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 32)
	return uint(id), err
}

// byte
func IsByteUpper(b byte) bool {
	if b >= 'A' && b <= 'Z' {
		return true
	} else {
		return false
	}
}

func ByteToLower(b byte) byte {
	if IsByteUpper(b) {
		return b + ASCII_ALPHA_LENTH
	} else {
		return b
	}
}

// string
var ASCII_ALPHA_LENTH byte = 'a' - 'A'

type StringAble interface {
	String() string
}

func Stringify(o interface{}) string {
	if s, ok := o.(StringAble); ok {
		return s.String()
	} else {
		return ""
	}
}

func CamelCaseToUnderlineCase(s string) string {
	sLen := len(s)
	retS := make([]byte, 2*sLen, 2*sLen)
	j := 1
	var pushRetS = func(b byte) {
		retS[j] = ByteToLower(b)
		j += 1
	}

	var isPreUpper = func(ind int) bool {
		if ind < 1 {
			return true
		} else {
			return IsByteUpper(s[ind-1])
		}
	}

	var isNxtUpper = func(ind int) bool {
		if ind > sLen-2 {
			return true
		} else {
			return IsByteUpper(s[ind+1])
		}
	}

	if sLen < 2 {
		return strings.ToLower(s)
	}

	retS[0] = ByteToLower(s[0])
	for i := 1; i < sLen; i++ {
		if IsByteUpper(s[i]) {
			if !isNxtUpper(i) || !isPreUpper(i) {
				pushRetS('_')
			}
		}

		pushRetS(s[i])
	}

	return string(retS[:j])
}

func Pluralize(s string) string {
	length := len(s)
	switch s[length-1] {
	case 'y':
		return s[:length-1] + "ies"
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		return s
	default:
		return s + "s"
	}
}
func StringSplit(s string, sep string) []string {
	rst := make([]string, 0)
	sepLen := len(sep)
	for ind := strings.Index(s, sep); ind != -1; ind = strings.Index(s, sep) {
		rst = append(rst, s[0:ind])
		s = s[ind+sepLen:]
	}


	return append(rst,s)
}
