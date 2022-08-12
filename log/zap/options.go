package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Option func(options *options)

type options struct {
	encoder zapcore.Encoder
	syncer  zapcore.WriteSyncer
	level   zapcore.LevelEnabler
	zaps    []zap.Option
}

func WithEncoder(encoder zapcore.Encoder) Option {
	return func(o *options) {
		o.encoder = encoder
	}
}

func WithWriteSyncer(syncer zapcore.WriteSyncer) Option {
	return func(o *options) {
		o.syncer = syncer
	}
}

func WithLevelEnabler(level zapcore.LevelEnabler) Option {
	return func(o *options) {
		o.level = level
	}
}

func WithZapOptions(opts ...zap.Option) Option {
	return func(o *options) {
		o.zaps = opts
	}
}
