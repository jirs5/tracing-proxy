package logger

import (
	"bytes"
	"fmt"
	"github.com/jirs5/tracing-proxy/config"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
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

	if entry.HasCaller() {
		b.WriteString(fmt.Sprintf(" [%s:%v] ", entry.Caller.File[strings.LastIndex(entry.Caller.File, "/")+1:], entry.Caller.Line))
	} else {
		b.WriteString(" [:-(] ")
	}

	b.WriteString(entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
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
		})
	}

	return l, nil
}
