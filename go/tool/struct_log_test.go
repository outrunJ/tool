package tool

import (
	"testing"
	"time"
)

func TestSequenceLogger(t *testing.T) {
	sl := NewSequenceLogger(NewLogger())
	sl.Info("info")
	sl.Warn("warn")
	//sl.Error("error")

	time.Sleep(5 * time.Millisecond)
}
