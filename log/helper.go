package log

import (
	"context"
	"fmt"
	"os"
)

type Option func(helper *Helper)

func WithMessageKey(key string) Option {
	return func(helper *Helper) {
		helper.msgKey = key
	}
}

const DefaultMsgKey = "msg"

type Helper struct {
	msgKey string
	logger Logger
}

func NewHelper(logger Logger, opts ...Option) *Helper {
	options := &Helper{
		msgKey: DefaultMsgKey,
		logger: logger,
	}

	for _, o := range opts {
		o(options)
	}

	return options
}

func (helper *Helper) WithContext(ctx context.Context) *Helper {
	return &Helper{
		msgKey: helper.msgKey,
		logger: WithContext(ctx, helper.logger),
	}
}

func (helper *Helper) Log(level Level, keyvals ...interface{}) {
	_ = helper.logger.Log(level, keyvals...)
}

func (helper *Helper) Debug(a ...interface{}) {
	helper.Log(LevelDebug, helper.msgKey, fmt.Sprint(a...))
}

func (helper *Helper) Debugf(format string, a ...interface{}) {
	helper.Log(LevelDebug, helper.msgKey, fmt.Sprintf(format, a...))
}

func (helper *Helper) Debugw(keyvals ...interface{}) {
	helper.Log(LevelDebug, keyvals...)
}

func (helper *Helper) Info(a ...interface{}) {
	helper.Log(LevelInfo, helper.msgKey, fmt.Sprint(a...))
}

func (helper *Helper) Infof(format string, a ...interface{}) {
	helper.Log(LevelInfo, helper.msgKey, fmt.Sprintf(format, a...))
}

func (helper *Helper) Infow(keyvals ...interface{}) {
	helper.Log(LevelInfo, keyvals...)
}

func (helper *Helper) Warn(a ...interface{}) {
	helper.Log(LevelWarn, helper.msgKey, fmt.Sprint(a...))
}

func (helper *Helper) Warnf(format string, a ...interface{}) {
	helper.Log(LevelWarn, helper.msgKey, fmt.Sprintf(format, a...))
}

func (helper *Helper) Warnw(keyvals ...interface{}) {
	helper.Log(LevelWarn, keyvals...)
}

func (helper *Helper) Error(a ...interface{}) {
	helper.Log(LevelError, helper.msgKey, fmt.Sprint(a...))
}

func (helper *Helper) Errorf(format string, a ...interface{}) {
	helper.Log(LevelError, helper.msgKey, fmt.Sprintf(format, a...))
}

func (helper *Helper) Errorw(keyvals ...interface{}) {
	helper.Log(LevelError, keyvals...)
}

func (helper *Helper) Fatal(a ...interface{}) {
	helper.Log(LevelFatal, helper.msgKey, fmt.Sprint(a...))
	os.Exit(1)
}

func (helper *Helper) Fatalf(format string, a ...interface{}) {
	helper.Log(LevelFatal, helper.msgKey, fmt.Sprintf(format, a...))
	os.Exit(1)
}

func (helper *Helper) Fatalw(keyvals ...interface{}) {
	helper.Log(LevelFatal, keyvals...)
	os.Exit(1)
}
