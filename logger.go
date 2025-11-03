package eslgo

import (
	"log"
)

type Logger interface {
	Debugf(format string, args ...any)
	Infof(format string, args ...any)
	Warnf(format string, args ...any)
	Errorf(format string, args ...any)
}

type NilLogger struct{}
type NormalLogger struct{}

func (l NormalLogger) Debugf(format string, args ...any) {
	log.Print("DEBUG: ")
	log.Printf(format, args...)
}
func (l NormalLogger) Infof(format string, args ...any) {
	log.Print("INFO: ")
	log.Printf(format, args...)
}
func (l NormalLogger) Warnf(format string, args ...any) {
	log.Print("WARN: ")
	log.Printf(format, args...)
}
func (l NormalLogger) Errorf(format string, args ...any) {
	log.Print("ERROR: ")
	log.Printf(format, args...)
}

func (l NilLogger) Debugf(string, ...any) {}
func (l NilLogger) Infof(string, ...any)  {}
func (l NilLogger) Warnf(string, ...any)  {}
func (l NilLogger) Errorf(string, ...any) {}
