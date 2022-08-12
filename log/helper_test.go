package log

import (
	"testing"
	"time"
)

func TestHelper(t *testing.T) {
	logging := With(DefaultLogger, "ts", Timestamp(time.RFC3339))
	log := NewHelper(logging)

	log.Log(LevelDebug, "msg", "test debug")
	log.Debug("test debug")
	log.Debugf("test %s", "debug")
	log.Debugw("log", "test debug")

	log.Warn("test warn")
	log.Warnf("test %s", "warn")
	log.Warnw("log", "test warn")
}