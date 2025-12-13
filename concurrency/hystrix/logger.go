package hystrix

import (
	"context"
	"fmt"

	"github.com/gw-gong/gwkit-go/log"
)

/*

Implement the Logger interface of github.com/afex/hystrix-go/hystrix

type logger interface {
    Printf(format string, items ...interface{})
}

*/

type SugarLogger struct {
	Ctx context.Context
}

func (l *SugarLogger) Printf(format string, items ...interface{}) {
	format = fmt.Sprintf("[hystrix] %s", format)
	if l == nil || l.Ctx == nil {
		log.Infof(format, items...)
		return
	}
	log.Infofc(l.Ctx, format, items...)
}
