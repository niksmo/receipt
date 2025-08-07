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
	produceRetries = 3
	maxBatchSize   = 200 * 1024
)

type KafkaProducer struct {
	log   logger.Logger
	kcl   *kgo.Client
	topic string
}

func NewKafkaProducer(
	log logger.Logger, seedBrokers []string, topic string,
) *KafkaProducer {
	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(seedBrokers...),
		kgo.DefaultProduceTopic(topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(produceRetries),
		kgo.ProducerBatchMaxBytes(maxBatchSize),
	)
	if err != nil {
		panic(err) // developer mistake
	}

	return &KafkaProducer{log, kcl, topic}
}

func (p *KafkaProducer) ProduceEvent(
	ctx context.Context, rct domain.Receipt,
) error {
	const op = "KafkaProducer.ProduceEvent"

	kr, err := p.createRecord(rct)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = p.kcl.ProduceSync(ctx, &kr).FirstErr()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (p *KafkaProducer) Close() {
	const op = "KafkaProducer.Close"
	log := p.log.WithOp(op)

	log.Info().Msg("closing producer")
	p.kcl.Close()
	log.Info().Msg("producer is closed")
}

func (p *KafkaProducer) InitTopic(
	ctx context.Context, partitions int, repFactor int, onFall func(err error),
) {
	const op = "KafkaProducer.InitTopic"
	log := p.log.WithOp(op)

	log.Info().Str("topic", p.topic).Msg("initializing topic...")
	_, err := kadm.NewClient(p.kcl).CreateTopic(
		ctx, int32(partitions), int16(repFactor), nil, p.topic,
	)
	if err != nil {
		if errors.Is(err, kerr.TopicAlreadyExists) {
			log.Info().Str("topic", p.topic).Msg("topic already exists")
			return
		}
		onFall(fmt.Errorf("%s: %w", op, err))
		return
	}

	log.Info().Str("topic", p.topic).Msg("topic created")
}

func (p *KafkaProducer) createRecord(rct domain.Receipt) (kgo.Record, error) {
	const op = "KafkaProducer.createRecord"

	v, err := json.Marshal(rct)
	if err != nil {
		return kgo.Record{}, fmt.Errorf("%s: %w", op, err)
	}

	return kgo.Record{Topic: p.topic, Value: v}, nil
}
