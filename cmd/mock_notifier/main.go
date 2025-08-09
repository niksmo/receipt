package main

import (
	"net/http"

	"github.com/niksmo/receipt/internal/mock_notifier/adapter"
	"github.com/niksmo/receipt/internal/mock_notifier/core/service"
	"github.com/niksmo/receipt/pkg/httpserver"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/niksmo/receipt/pkg/middleware"
	"github.com/niksmo/receipt/pkg/sig"
)

func main() {
	PrintAppTitle()

	cfg := LoadConfig()
	cfg.Print()

	sigCtx, stop := sig.NotifyContext()
	defer stop()

	log := logger.New(cfg.LogLevel)

	service := service.NewService(log)

	mux := http.NewServeMux()
	adapter.RegisterSendMailHandler(log, mux, service, cfg.RateLimit)

	handler := middleware.LogResposeStatus(log, middleware.AcceptJSON(mux))
	httpServer := httpserver.New(log, cfg.HTTPServerAddr, handler)
	go httpServer.Run(stop)

	<-sigCtx.Done()
	httpServer.Close()
}
