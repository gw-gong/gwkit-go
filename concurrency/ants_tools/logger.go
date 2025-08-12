package ants_tools

import (
	"context"

	"github.com/gw-gong/gwkit-go/log"
)

/*

Implement the Logger interface of github.com/panjf2000/ants/v2

type Logger interface {
	// Printf must have the same semantics as log.Printf.
	Printf(format string, args ...any)
}

*/

type SugarLogger struct {
	Ctx context.Context
}

func (l *SugarLogger) Printf(format string, args ...interface{}) {
	if l == nil || l.Ctx == nil {
		log.Infof(format, args...)
		return
	}
	log.Infofc(l.Ctx, format, args...)
}
