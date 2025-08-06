package adapter

import (
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	log logger.Logger
	kcl *kgo.Client
}

func NewKafkaConsumer(log logger.Logger) *KafkaConsumer {
	kcl, err := kgo.NewClient()
	if err != nil {
		panic(err)
	}
	return &KafkaConsumer{log, kcl}
}
