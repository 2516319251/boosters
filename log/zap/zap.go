package zap

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(opts ...Option) *zap.Logger {
	o := options{
		encoder: encoder(),
		syncer:  syncer(),
		level:   zapcore.DebugLevel,
	}

	for _, opt := range opts {
		opt(&o)
	}

	core := zapcore.NewCore(o.encoder, o.syncer, o.level)
	return zap.New(core, o.zaps...)
}

func encoder() zapcore.Encoder {
	enc := zap.NewProductionEncoderConfig()
	enc.EncodeTime = zapcore.ISO8601TimeEncoder
	return zapcore.NewJSONEncoder(enc)
}

func syncer() zapcore.WriteSyncer {
	return zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./app.log",
		MaxSize:    100,
		MaxBackups: 1024,
		MaxAge:     30,
		Compress:   false,
	})
}
