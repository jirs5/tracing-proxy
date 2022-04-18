package logger

import (
	"bytes"
	"fmt"
	"github.com/jirs5/tracing-proxy/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
	"runtime"
	"time"
)

// Custom Log Formatter for OpsRamp
type OpsRampLogFormat struct {
	TimestampFormat string
}

func (f *OpsRampLogFormat) Format(entry *logrus.Entry) ([]byte, error) {
	b := &bytes.Buffer{}

	if entry.Buffer != nil {
		b = entry.Buffer
	}

	if f.TimestampFormat == "" {
		f.TimestampFormat = time.RFC3339
	}

	b.WriteString(fmt.Sprintf("%s [%s]", entry.Time.Format(f.TimestampFormat), entry.Level.String()))
	b.WriteByte(' ')

	if entry.HasCaller() {
		_, filepath := opsRampCallerPrettyfier(entry.Caller)
		b.WriteString(fmt.Sprintf("[%v]", filepath))
	} else {
		b.WriteString(" [:-(] ")
	}
	b.WriteByte(' ')

	if len(entry.Data) > 0 {
		for key, value := range entry.Data {
			f.appendKeyValue(b, key, value)
		}
		b.WriteByte(' ')
	}

	b.WriteString(entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *OpsRampLogFormat) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	b.WriteString(fmt.Sprintf("[%s", key))
	b.WriteByte('=')
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}
	b.WriteString(fmt.Sprintf("%q]", stringVal))
}

func opsRampCallerPrettyfier(frame *runtime.Frame) (function string, file string) {
	return frame.Function, fmt.Sprintf("%s:%d", path.Base(frame.File), frame.Line)
}

func GetLoggerImplementation(c config.Config) (*logrus.Logger, error) {

	l := logrus.New()

	// set log level
	logLevel, err := c.GetLoggingLevel()
	if err != nil {
		return nil, err
	}

	logrusLevel, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, err
	}
	l.SetLevel(logrusLevel)

	l.SetReportCaller(true)

	logrusConfig, err := c.GetLogrusConfig()
	if err != nil {
		return nil, err
	}

	switch logrusConfig.LogOutput {
	case "stdout":
		l.SetOutput(os.Stdout)
	case "stderr":
		l.SetOutput(os.Stderr)
	case "file":
		l.SetOutput(&lumberjack.Logger{
			Filename:   logrusConfig.File.FileName,
			MaxSize:    logrusConfig.File.MaxSize,
			MaxBackups: logrusConfig.File.MaxBackups,
			Compress:   logrusConfig.File.Compress,
		})
	}

	switch logrusConfig.LogFormatter {
	case "opsramp":
		l.SetFormatter(&OpsRampLogFormat{})
	case "logfmt":
		l.SetFormatter(&logrus.TextFormatter{
			DisableColors:          true,
			ForceQuote:             true,
			FullTimestamp:          true,
			DisableLevelTruncation: true,
			QuoteEmptyFields:       true,
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyFile:  "file",
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
			CallerPrettyfier: opsRampCallerPrettyfier,
		})
	case "json":
		l.SetFormatter(&logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{
				logrus.FieldKeyFile:  "file",
				logrus.FieldKeyTime:  "timestamp",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
				logrus.FieldKeyFunc:  "caller",
			},
			CallerPrettyfier: opsRampCallerPrettyfier,
		})
	}

	return l, nil
}
