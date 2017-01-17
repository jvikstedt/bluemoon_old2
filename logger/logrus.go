package logger

import (
	"io"

	log "github.com/Sirupsen/logrus"
)

const (
	PanicLevel Level = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

type Level uint8

type LogrusLogger struct {
}

func NewLogrusLogger(out io.Writer, level Level) *LogrusLogger {
	log.SetOutput(out)
	log.SetLevel(log.Level(level))
	return &LogrusLogger{}
}

func (ll *LogrusLogger) Panicln(args ...interface{}) {
	log.Panicln(args)
}
func (ll *LogrusLogger) Fatalln(args ...interface{}) {
	log.Fatalln(args)
}
func (ll *LogrusLogger) Errorln(args ...interface{}) {
	log.Errorln(args)
}
func (ll *LogrusLogger) Warnln(args ...interface{}) {
	log.Warnln(args)
}
func (ll *LogrusLogger) Infoln(args ...interface{}) {
	log.Infoln(args)
}
func (ll *LogrusLogger) Debugln(args ...interface{}) {
	log.Debugln(args)
}
