package adapters

import (
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	produceRetries = 3
	maxBatchSize   = 200 * 1024
)

type MessageProvider struct {
	log logger.Logger
	kcl *kgo.Client
}

func NewMessageProvider(
	log logger.Logger, seedBrokers []string, topic string,
) (MessageProvider, error) {
	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(seedBrokers...),
		kgo.DefaultProduceTopic(topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(produceRetries),
		kgo.ProducerBatchMaxBytes(maxBatchSize),
	)
	if err != nil {
		return MessageProvider{}, err
	}

	return MessageProvider{log, kcl}, nil
}

func (p MessageProvider) Close() {
	p.kcl.Close()
}
