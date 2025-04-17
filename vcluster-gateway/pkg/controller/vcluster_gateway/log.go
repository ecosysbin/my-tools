package vcluster_gateway

import (
	"io"
	"strings"

	"gitlab.datacanvas.com/AlayaNeW/OSM/gokit/log"

	"github.com/go-logr/logr"
	vclusterlogger "github.com/loft-sh/log"
	"github.com/loft-sh/log/survey"
	"github.com/sirupsen/logrus"
)

var _ vclusterlogger.Logger = &VClusterLogger{}

type VClusterLogger struct {
	logger *log.Logger
}

func NewVClusterLogger(logger *log.Logger) *VClusterLogger {
	return &VClusterLogger{
		logger: logger,
	}
}

func (l *VClusterLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *VClusterLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *VClusterLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *VClusterLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *VClusterLogger) Done(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *VClusterLogger) Donef(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *VClusterLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *VClusterLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *VClusterLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *VClusterLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *VClusterLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *VClusterLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *VClusterLogger) Print(level logrus.Level, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		l.Info(args...)
	case logrus.DebugLevel:
		l.Debug(args...)
	case logrus.WarnLevel:
		l.Warn(args...)
	case logrus.ErrorLevel:
		l.Error(args...)
	case logrus.FatalLevel:
		l.Fatal(args...)
	case logrus.PanicLevel:
		l.Fatal(args...)
	case logrus.TraceLevel:
		l.Debug(args...)
	}
}

func (l *VClusterLogger) Printf(level logrus.Level, format string, args ...interface{}) {
	switch level {
	case logrus.InfoLevel:
		l.Info(args...)
	case logrus.DebugLevel:
		l.Debug(args...)
	case logrus.WarnLevel:
		l.Warn(args...)
	case logrus.ErrorLevel:
		l.Error(args...)
	case logrus.FatalLevel:
		l.Fatal(args...)
	case logrus.PanicLevel:
		l.Fatal(args...)
	case logrus.TraceLevel:
		l.Debug(args...)
	}
}

func (l *VClusterLogger) SetLevel(level logrus.Level) {
	switch strings.ToLower(level.String()) {
	case "panic":
		l.logger.Level = "debug"
	case "fatal":
		l.logger.Level = "debug"
	case "error":
		l.logger.Level = "error"
	case "warn", "warning":
		l.logger.Level = "warn"
	case "info":
		l.logger.Level = "info"
	case "debug":
		l.logger.Level = "debug"
	case "trace":
		l.logger.Level = "debug"
	}
}

func (l *VClusterLogger) GetLevel() logrus.Level {
	return logrus.DebugLevel
}

func (l *VClusterLogger) LogrLogSink() logr.LogSink {
	return nil
}

func (l *VClusterLogger) Question(params *survey.QuestionOptions) (string, error) {
	return "", nil
}

func (l *VClusterLogger) ErrorStreamOnly() vclusterlogger.Logger {
	return l
}

func (l *VClusterLogger) Writer(level logrus.Level, raw bool) io.WriteCloser {
	return &NopCloser{io.Discard}
}

func (l *VClusterLogger) WriteString(level logrus.Level, message string) {
}

type NopCloser struct {
	io.Writer
}

func (NopCloser) Close() error { return nil }
