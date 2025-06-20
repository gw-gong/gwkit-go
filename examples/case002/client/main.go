package main

import (
	"context"

	"github.com/gw-gong/gwkit-go/log"
	gwkit_common "github.com/gw-gong/gwkit-go/utils/common"
)

func main() {
	syncFn, err := log.InitGlobalLogger(log.NewDefaultLoggerConfig())
	gwkit_common.ExitOnErr(context.Background(), err)
	defer syncFn()

	testClient, err := NewTestClient("127.0.0.1:8500", "test_service", "test", "")
	gwkit_common.ExitOnErr(context.Background(), err)

	response, err := testClient.TestFunc(context.Background(), "test")
	gwkit_common.ExitOnErr(context.Background(), err)

	log.Info("rpc 调用成功", log.Str("response", response))
}
