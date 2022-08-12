package log

import (
	"context"
	"log"
)

var DefaultLogger Logger = NewStdLogger(log.Writer())

type Logger interface {
	Log(level Level, keyvals ...interface{}) error
}

type logger struct {
	log    Logger
	prefix []interface{}
	valuer bool
	ctx    context.Context
}

func (c *logger) Log(level Level, keyvals ...interface{}) error {
	kvs := make([]interface{}, 0, len(c.prefix)+len(keyvals))
	kvs = append(kvs, c.prefix...)
	if c.valuer {
		bindValues(c.ctx, kvs)
	}
	kvs = append(kvs, keyvals...)
	return c.log.Log(level, kvs...)
}

func With(l Logger, kv ...interface{}) Logger {
	if c, ok := l.(*logger); ok {
		kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
		kvs = append(kvs, kv...)
		kvs = append(kvs, c.prefix...)
		return &logger{
			log:    c.log,
			prefix: kvs,
			valuer: containsValuer(kvs),
			ctx:    c.ctx,
		}
	}
	return &logger{log: l, prefix: kv, valuer: containsValuer(kv)}
}

func WithContext(ctx context.Context, l Logger) Logger {
	if c, ok := l.(*logger); ok {
		return &logger{
			log:    c.log,
			prefix: c.prefix,
			valuer: c.valuer,
			ctx:    ctx,
		}
	}
	return &logger{log: l, ctx: ctx}
}
