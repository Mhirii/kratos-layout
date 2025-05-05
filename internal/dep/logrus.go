package dep

import (
	"fmt"
	"os"

	"layout/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	Logger *logrus.Logger
}

func NewLogrusLogger(c *conf.Bootstrap) *LogrusLogger {
	var formatter logrus.Formatter
	txtFmt := &logrus.TextFormatter{
		ForceQuote:       true,
		QuoteEmptyFields: true,
		DisableColors:    false,
		DisableTimestamp: false,
		//FullTimestamp:    true,
		//TimestampFormat: "",
	}
	jsonFmt := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "time",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "msg",
			logrus.FieldKeyFunc:  "func",
			logrus.FieldKeyFile:  "file",
		},
		DisableTimestamp: false,
	}

	switch c.Metadata.Env {
	case 0, 1:
		formatter = txtFmt
	case 2, 3:
		formatter = jsonFmt
	default:
		formatter = txtFmt
	}
	l := logrus.New()
	l.SetFormatter(formatter)
	if c.Log.Filepath != "" {
		file, err := os.OpenFile(c.Log.Filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err == nil {
			l.SetOutput(file)
		} else {
			l.Info("Failed to Logger to file, using default stderr")
		}
	}

	return &LogrusLogger{Logger: l}
}
func (l *LogrusLogger) Log(level log.Level, keyvals ...interface{}) error {
	if len(keyvals) == 0 || len(keyvals)%2 != 0 {
		l.Logger.Warn(fmt.Sprint("Keyvalues must appear in pairs: ", keyvals))
		return nil
	}

	var msg string
	var fields logrus.Fields = make(logrus.Fields)

	for i := 0; i < len(keyvals); i += 2 {
		if keyvals[i] == "msg" {
			msg = fmt.Sprint(keyvals[i+1])
		} else {
			fields[fmt.Sprint(keyvals[i])] = fmt.Sprint(keyvals[i+1])
		}
	}
	fields[fmt.Sprint("logger")] = "logrus"

	lg := l.Logger.WithFields(fields)
	switch level {
	case log.LevelDebug:
		lg.Debug(msg)
	case log.LevelInfo:
		lg.Info(msg)
	case log.LevelWarn:
		lg.Warn(msg)
	case log.LevelError:
		lg.Error(msg)
	case log.LevelFatal:
		lg.Fatal(msg)
	}
	return nil
}
