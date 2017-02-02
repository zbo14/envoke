package util

import log "github.com/inconshreveable/log15"

type Logger struct {
	log.Logger
}

func NewLogger(module string) Logger {
	return Logger{log.New("module", module)}
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.Error(Sprintf(format, args...))
}
