package tool

import (
	"fmt"
	"testing"
	"time"
)

func TestSetField(t *testing.T) {
	type S struct {
		Name string
	}
	s := &S{Name: "a"}
	SetField(s, "Same", "b")
	fmt.Println(s)
}

func TestSetFieldZero(t *testing.T) {
	type S struct {
		Name string
		Time *time.Time
	}

	now := time.Now()
	s := &S{
		Name: "a",
		Time: &now,
	}
	SetFieldZero(s, "Name")
	SetFieldZero(s, "Time")
	fmt.Println(s)
}
