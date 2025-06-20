package common

import (
	"context"
	"os"

	"github.com/gw-gong/gwkit-go/log"
)

func ExitOnErr(ctx context.Context, err error) {
	if err != nil {
		log.Errorc(ctx, "ExitOnErr", log.Err(err))
		os.Exit(1)
	}
}
