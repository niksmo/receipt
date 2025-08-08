package main

import (
	"context"
	"fmt"

	"github.com/niksmo/receipt/pkg/logger"
)

func OnInitTopicFall(log logger.Logger, stop context.CancelFunc) func(error) {
	return func(err error) {
		log.Error().Err(err).Msg("failed to init broker topic")
		stop()
	}
}

func PrintAppTitle() {
	fmt.Printf(`
+-----------------------+
|🧾RECEIPT APPLICATION🚀|
+-----------------------+
`)
}
