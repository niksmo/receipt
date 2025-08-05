package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/niksmo/receipt/internal/core/domain"
	"github.com/niksmo/receipt/internal/core/port"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kerr"
	"github.com/twmb/franz-go/pkg/kgo"
)

var _ port.EventProducer = (*KafkaProducer)(nil)

const (
	produceRetries    = 3
	maxBatchSize      = 200 * 1024
	partitions        = 3
	replicationFactor = 3
)

var minInsyncReplicas = "2"

type KafkaProducer struct {
	log   logger.Logger
	kcl   *kgo.Client
	topic string
}

func NewKafkaProducer(
	ctx context.Context, log logger.Logger, brokers []string, topic string,
) (*KafkaProducer, error) {
	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(brokers...),
		kgo.DefaultProduceTopic(topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(produceRetries),
		kgo.ProducerBatchMaxBytes(maxBatchSize),
	)
	if err != nil {
		return nil, err
	}

	p := KafkaProducer{log, kcl, topic}
	if err := p.initTopic(ctx); err != nil {
		return nil, err
	}

	return &p, nil
}

func (p *KafkaProducer) ProduceEvent(
	ctx context.Context, rct domain.Receipt,
) error {
	const op = "KafkaProducer.ProduceEvent"

	kr, err := p.createRecord(rct)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	p.kcl.Produce(ctx, &kr, p.promise(rct))
	return nil
}

func (p *KafkaProducer) createRecord(rct domain.Receipt) (kgo.Record, error) {
	const op = "KafkaProducer.createRecord"

	v, err := json.Marshal(rct)
	if err != nil {
		return kgo.Record{}, fmt.Errorf("%s: %w", op, err)
	}

	return kgo.Record{Value: v}, nil
}

func (p *KafkaProducer) promise(rct domain.Receipt) func(*kgo.Record, error) {
	const op = "KafkaProducer.promise"
	log := p.log.WithOp(op)

	return func(r *kgo.Record, err error) {
		if err != nil {
			log.Error().Err(
				err).Str("receiptUUID", rct.UUID).Msg(
				"failed to produce record")
		}
	}
}

func (p *KafkaProducer) Close() {
	const op = "KafkaProducer.Close"
	log := p.log.WithOp(op)

	log.Info().Msg("closing producer")
	p.kcl.Close()
	log.Info().Msg("producer is closed")
}

func (p *KafkaProducer) initTopic(ctx context.Context) error {
	const op = "KafkaProducer.initTopic"
	log := p.log.WithOp(op)

	_, err := kadm.NewClient(p.kcl).CreateTopic(
		ctx,
		partitions,
		replicationFactor,
		map[string]*string{"min.insync.replicas": &minInsyncReplicas},
		p.topic,
	)
	if err != nil {
		if errors.Is(err, kerr.TopicAlreadyExists) {
			log.Info().Str("topic", p.topic).Msg("topic already exists")
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info().Str("topic", p.topic).Msg("topic created")
	return nil
}
