package logger

import "github.com/kataras/golog"

type DatabaseLogger struct {
	*golog.Logger
}

func (l *DatabaseLogger) Panic(v ...interface{}) {
	l.Fatal(v...)
}

func (l *DatabaseLogger) Panicf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}

func NewDatabaseLogger(logger *golog.Logger) *DatabaseLogger {
	return &DatabaseLogger{logger}
}
