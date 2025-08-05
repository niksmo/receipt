package adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/niksmo/receipt/pkg/logger"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	produceRetries    = 3
	maxBatchSize      = 200 * 1024
	partitions        = 3
	replicationFactor = 3
)

var minInsyncReplicas = "2"

type MessageProducer struct {
	log   logger.Logger
	kcl   *kgo.Client
	topic string
}

func NewMessageProducer(
	ctx context.Context, log logger.Logger, seedBrokers []string, topic string,
) (MessageProducer, error) {
	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(seedBrokers...),
		kgo.DefaultProduceTopic(topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(produceRetries),
		kgo.ProducerBatchMaxBytes(maxBatchSize),
	)
	if err != nil {
		return MessageProducer{}, err
	}

	p := MessageProducer{log, kcl, topic}
	if err := p.initTopic(ctx); err != nil {
		return MessageProducer{}, err
	}

	return p, nil
}

func (p MessageProducer) ProduceReceipt() {}

func (p MessageProducer) Close() {
	const op = "MessageProducer.Close"
	log := p.log.WithOp(op)

	log.Info().Msg("closing producer")
	p.kcl.Close()
	log.Info().Msg("producer is closed")
}

func (p MessageProducer) initTopic(ctx context.Context) error {
	const op = "MessageProvider.initTopic"

	_, err := kadm.NewClient(p.kcl).CreateTopic(
		ctx,
		partitions,
		replicationFactor,
		map[string]*string{"min.insync.replicas": &minInsyncReplicas},
		p.topic,
	)
	if err != nil {
		if errors.Is(err, kerr.TopicAlreadyExists) {
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
