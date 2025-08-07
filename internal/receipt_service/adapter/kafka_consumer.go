package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/niksmo/receipt/internal/receipt_service/core/domain"
	"github.com/niksmo/receipt/internal/receipt_service/core/port"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	fetchMinBytes = 20 * 1024
	fetchMaxWait  = 2 * time.Second
)

type KafkaConsumer struct {
	log    logger.Logger
	kcl    *kgo.Client
	ep     port.EventProcessor
	nRecs  atomic.Int64
	nBytes atomic.Int64
}

func NewKafkaConsumer(
	log logger.Logger,
	seedBrokers []string, topic string, group string, ep port.EventProcessor,
) *KafkaConsumer {
	kcl, err := kgo.NewClient(
		kgo.SeedBrokers(seedBrokers...),
		kgo.DisableAutoCommit(),
		kgo.ConsumeTopics(topic),
		kgo.FetchMinBytes(fetchMinBytes),
		kgo.FetchMaxWait(fetchMaxWait),
		kgo.ConsumerGroup(group),
	)
	if err != nil {
		panic(err) // developer mistake
	}
	return &KafkaConsumer{log: log, kcl: kcl, ep: ep}
}

func (c *KafkaConsumer) Run(ctx context.Context) {
	const op = "KafkaConsumer.Run"
	log := c.log.WithOp(op)

	log.Info().Msg("kafka consumer is running")

	go c.capacity(ctx)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			c.consume(ctx)
		}
	}
}

func (c *KafkaConsumer) Close() {
	const op = "KafkaConsumer.Close"
	log := c.log.WithOp(op)

	log.Info().Msg("closing consumer")
	c.kcl.Close()
	log.Info().Msg("consumer is closed")
}

func (c *KafkaConsumer) consume(ctx context.Context) {
	const op = "KafkaConsumer.consume"
	log := c.log.WithOp(op)

	fetches, err := c.pollFetches(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			log.Info().Msg("interrupted")
			return
		}
		log.Error().Err(err).Msg("failed to poll fetches")
		return
	}

	rcts := c.retrieveReceipts(fetches)
	if !c.isReceipts(rcts) {
		return
	}
	c.handleReceipts(ctx, rcts)
	c.commitOffset(ctx)
}

func (c *KafkaConsumer) capacity(ctx context.Context) {
	const op = "KafkaConsumer.capacity"
	log := c.log.WithOp(op)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			nBytes := fmt.Sprintf("%.2f", float64(c.nBytes.Swap(0)))
			nRecs := c.nRecs.Swap(0)
			log.Debug().Str("KiB/s", nBytes).Int64("recs/s", nRecs).Send()
		}
	}
}

func (c *KafkaConsumer) pollFetches(ctx context.Context) (kgo.Fetches, error) {
	const op = "KafkaConsumer.pollFetches"
	log := c.log.WithOp(op)

	log.Debug().Msg("start polling")
	fetches := c.kcl.PollFetches(ctx)
	log.Debug().Msg("complete polling")
	err := fetches.Err0()
	if errors.Is(err, context.Canceled) {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var errs []error
	fetches.EachError(func(t string, p int32, err error) {
		if err != nil {
			fetchErr := fmt.Errorf("topic %q partition %d: %w", t, p, err)
			errs = append(errs, fetchErr)
		}
	})

	if len(errs) != 0 {
		return nil, fmt.Errorf("%s: %w", op, errors.Join(errs...))
	}

	recs := int64(fetches.NumRecords())
	c.nRecs.Add(recs)
	log.Debug().Int64("nRecord", recs).Send()
	return fetches, nil
}

func (c *KafkaConsumer) retrieveReceipts(
	fetches kgo.Fetches,
) []domain.Receipt {
	const op = "KafkaConsumer.retrieveReceipts"
	log := c.log.WithOp(op)

	var rcts []domain.Receipt
	fetches.EachRecord(func(rec *kgo.Record) {
		c.nBytes.Add(int64(len(rec.Value)))
		rct, err := c.unmarshalReceipt(rec.Value)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal record value")
			return
		}
		rcts = append(rcts, rct)
	})
	log.Debug().Int("nReceipts", len(rcts)).Send()

	return rcts
}

func (c *KafkaConsumer) unmarshalReceipt(b []byte) (domain.Receipt, error) {
	const op = "KafkaConsumer.unmarshalReceipt"

	var rct domain.Receipt
	err := json.Unmarshal(b, &rct)
	if err != nil {
		return domain.Receipt{}, fmt.Errorf("%s: %w", op, err)
	}
	return rct, nil
}

func (c *KafkaConsumer) isReceipts(rcts []domain.Receipt) bool {
	return len(rcts) != 0
}

func (c *KafkaConsumer) handleReceipts(
	ctx context.Context, rcts []domain.Receipt,
) {
	c.ep.ProcessEvent(ctx, rcts)
}

func (c *KafkaConsumer) commitOffset(ctx context.Context) {
	const op = "KafkaConsumer.commitOffset"
	log := c.log.WithOp(op)

	if err := c.kcl.CommitUncommittedOffsets(ctx); err != nil {
		log.Error().Err(err).Msg("failed to commit offsets")
	}
	log.Debug().Msg("successfuly committed")
}
