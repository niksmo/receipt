package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

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

type KafkaProducerOpts struct {
	SeedBrokers       []string
	Topic             string
	Partitions        int
	ReplicationFactor int
	MinInsyncReplicas int
}

type KafkaProducer struct {
	log  logger.Logger
	kcl  *kgo.Client
	opts KafkaProducerOpts
}

func NewKafkaProducer(
	ctx context.Context, log logger.Logger, opts KafkaProducerOpts,
) (*KafkaProducer, error) {
	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(opts.SeedBrokers...),
		kgo.DefaultProduceTopic(opts.Topic),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.RecordRetries(produceRetries),
		kgo.ProducerBatchMaxBytes(maxBatchSize),
	)
	if err != nil {
		panic(err)
	}

	p := KafkaProducer{log, kcl, opts}
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

func (p *KafkaProducer) initTopic(ctx context.Context) error {
	const op = "KafkaProducer.initTopic"
	log := p.log.WithOp(op)

	minInsyncReplicas := strconv.Itoa(p.opts.MinInsyncReplicas)

	log.Info().Str("topic", p.opts.Topic).Msg("topic initialization")
	_, err := kadm.NewClient(p.kcl).CreateTopic(
		ctx,
		int32(p.opts.Partitions),
		int16(p.opts.ReplicationFactor),
		map[string]*string{"min.insync.replicas": &minInsyncReplicas},
		p.opts.Topic,
	)
	if err != nil {
		if errors.Is(err, kerr.TopicAlreadyExists) {
			log.Info().Str("topic", p.opts.Topic).Msg("topic already exists")
			return nil
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info().Str("topic", p.opts.Topic).Msg("topic created")
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
