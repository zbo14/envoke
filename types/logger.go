package types

import log "github.com/inconshreveable/log15"

type Logger log.Logger

func NewLogger(module string) Logger {
	return log.New("module", module)
}
