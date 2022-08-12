package log

import (
	"log"
	"testing"
)

func TestNewStdLogger(t *testing.T) {
	l := NewStdLogger(log.Writer())

	l.Log(LevelDebug, "msg", "test debug")
	l.Log(LevelInfo, "msg", "test info")
	l.Log(LevelWarn, "msg", "test warn")
	l.Log(LevelError, "msg", "test error")
	l.Log(LevelFatal, "msg", "test fatal")
}
