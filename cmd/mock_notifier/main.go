package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/niksmo/receipt/internal/mock_notifier/adapter"
	"github.com/niksmo/receipt/internal/mock_notifier/core/service"
	"github.com/niksmo/receipt/pkg/env"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/niksmo/receipt/pkg/sig"
)

const (
	defaultLogLevel       = "info"
	defaultHTTPServerAddr = ":8080"
	defaultRateLimit      = 1000 // RPS
)

func main() {
	printAppTitle()
	sigCtx, stop := sig.NotifyContext()
	defer stop()

	log := logger.New(loadLogLevel())

	service := service.NewService(log)

	mux := http.NewServeMux()
	adapter.RegisterSendMailHandler(log, mux, service, loadRateLimit())

	httpServer := httpserver.New(log, loadHTTPAddr(), mux)
	go httpServer.Run(stop)

	<-sigCtx.Done()
	httpServer.Close()
}

func loadLogLevel() string {
	lvl, err := env.String("NOTIFIER_LOG_LEVEL", nil)
	if errors.Is(err, env.ErrNotSet) {
		lvl = defaultLogLevel
	}
	return lvl
}

func loadHTTPAddr() string {
	addr, err := env.String(
		"NOTIFIER_HTTP_ADDR",
		func(v string) error {
			_, err := net.ResolveTCPAddr("tcp", v)
			return err
		},
	)

	if err != nil {
		if errors.Is(err, env.ErrNotSet) {
			addr = defaultHTTPServerAddr
		} else {
			panic(err)
		}
	}
	return addr
}

func loadRateLimit() int {
	rateLimit, err := env.Int("NOTIFIER_RATE_LIMIT", nil)
	if errors.Is(err, env.ErrNotSet) {
		rateLimit = defaultRateLimit
	}
	return rateLimit
}

func printAppTitle() {
	fmt.Printf(`
+-----------------------------+
|🏓MOCK NOTIFIER APPLICATION🚀|
+-----------------------------+
`)
}
