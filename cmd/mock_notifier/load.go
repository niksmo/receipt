package main

import (
	"errors"
	"fmt"

	"github.com/niksmo/receipt/config"
	"github.com/niksmo/receipt/pkg/env"
)

type AppConfig struct {
	LogLevel       string
	HTTPServerAddr string
	RateLimit      int
}

func LoadConfig() AppConfig {
	logLevel := config.LoadLogLevel("NOTIFIER_LOG_LEVEL", defaultLogLevel)

	addr, err := config.LoadHTTPServerAddr("NOTIFIER_HTTP_ADDR", defaultHTTPServerAddr)
	if err != nil {
		panic(err)
	}

	rateLimit := loadRateLimit()

	return AppConfig{logLevel, addr, rateLimit}
}

func (c AppConfig) Print() {
	fmt.Printf(
		`Configuration:
LogLevel:          %q
HTTPServerAddress: %q
RateLimit:        % d

`,
		c.LogLevel,
		c.HTTPServerAddr,
		c.RateLimit,
	)
}

func loadRateLimit() int {
	rateLimit, err := env.Int("NOTIFIER_RATE_LIMIT", nil)
	if errors.Is(err, env.ErrNotSet) {
		rateLimit = defaultRateLimit
	}
	return rateLimit
}

func PrintAppTitle() {
	fmt.Printf(`
+-----------------------------+
|🌈MOCK NOTIFIER APPLICATION🚀|
+-----------------------------+
`)
}
