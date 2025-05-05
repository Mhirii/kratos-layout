package dep

import (
	"strings"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(c *conf.Bootstrap) log.Logger {
	switch strings.ToUpper(c.Log.Logger) {
	case conf.Log_ZAP.String():
		level := zap.NewAtomicLevelAt(zap.DebugLevel)
		return NewZapLogger(zapcore.EncoderConfig{
			MessageKey:     "msg",
			LevelKey:       "level",
			TimeKey:        "ts",
			NameKey:        "name",
			CallerKey:      "caller",
			FunctionKey:    "fn",
			StacktraceKey:  "stack",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
			EncodeName:     zapcore.FullNameEncoder,
			LineEnding:     zapcore.DefaultLineEnding,
		}, level)
	case conf.Log_LOGRUS.String():
		return NewLogrusLogger(c)
	default:
		return NewLogrusLogger(c)
	}
}
