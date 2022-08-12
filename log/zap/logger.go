package zap

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/2516319251/boosters/log"
)

var _ log.Logger = (*Logger)(nil)

type Logger struct {
	log *zap.Logger
}

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.log.Warn(fmt.Sprint("Key values must appear in pairs: ", keyvals))
		return nil
	}

	var msg string
	var data []zap.Field
	for i := 0; i < len(keyvals); i += 2 {

		if keyvals[i] == log.DefaultMsgKey {
			msg = fmt.Sprintf("%v", keyvals[i+1])
			continue
		}

		data = append(data, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.LevelDebug:
		l.log.Debug(msg, data...)
	case log.LevelInfo:
		l.log.Info(msg, data...)
	case log.LevelWarn:
		l.log.Warn(msg, data...)
	case log.LevelError:
		l.log.Error(msg, data...)
	case log.LevelFatal:
		l.log.Fatal(msg, data...)
	}

	return nil
}

func (l *Logger) Sync() error {
	return l.log.Sync()
}
