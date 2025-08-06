package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/niksmo/receipt/internal/core/domain"
	"github.com/niksmo/receipt/internal/core/port"
	"github.com/niksmo/receipt/pkg/logger"
	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaConsumer struct {
	log    logger.Logger
	kcl    *kgo.Client
	ep     port.EventProcessor
	nRecs  int
	nBytes int
}

func NewKafkaConsumer(
	log logger.Logger,
	seedBrokers []string,
	topic string,
	ep port.EventProcessor,
) *KafkaConsumer {
	kcl, err := kgo.NewClient(
		kgo.DisableAutoCommit(),
		kgo.ConsumeTopics(topic),
	)
	if err != nil {
		panic(err)
	}
	return &KafkaConsumer{log: log, kcl: kcl, ep: ep}
}

func (c *KafkaConsumer) Run(ctx context.Context) {
	const op = "KafkaConsumer.Run"
	log := c.log.WithOp(op)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			c.logCapacity()
		default:
			fetches, err := c.pollFetches(ctx)
			if err != nil {
				if errors.Is(context.DeadlineExceeded, err) {
					log.Info().Msg("poll fetches interrupted")
					return
				}
				log.Error().Err(err).Msg("failed to poll fetches")
				continue
			}

			rcts := c.getReceipts(fetches)
			c.handleReceipts(ctx, rcts)
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

func (c *KafkaConsumer) logCapacity() {
	const op = "KafkaConsumer.logCapacity"
	log := c.log.WithOp(op)

	log.Debug().Msgf(
		"consume %.2f MiB/s, %.2f records/s",
		float64(c.nBytes)/1024*1024, float64(c.nRecs)/1000,
	)
	c.nRecs = 0
	c.nBytes = 0
}

func (c *KafkaConsumer) pollFetches(ctx context.Context) (kgo.Fetches, error) {
	const op = "KafkaConsumer.pollFetches"

	fetches := c.kcl.PollFetches(ctx)
	err := fetches.Err0()
	if errors.Is(err, context.DeadlineExceeded) {
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
	return fetches, nil
}

func (c *KafkaConsumer) getReceipts(fetches kgo.Fetches) []domain.Receipt {
	const op = "KafkaConsumer.getReceipts"
	log := c.log.WithOp(op)

	var rcts []domain.Receipt
	fetches.EachRecord(func(rec *kgo.Record) {
		c.nRecs++
		c.nBytes += len(rec.Value)

		var rct domain.Receipt
		err := json.Unmarshal(rec.Value, &rct)
		if err != nil {
			log.Error().Err(err).Msg("failed to unmarshal record value")
			return
		}
		rcts = append(rcts, rct)
	})

	return rcts
}

func (c *KafkaConsumer) handleReceipts(
	ctx context.Context, rcts []domain.Receipt,
) {
	const op = "KafkaConsumer.handleReceipts"
	log := c.log.WithOp(op)

	if len(rcts) == 0 {
		log.Info().Msg("no receipts to handle")
		return
	}

	c.ep.ProcessEvent(ctx, rcts)

	if err := c.kcl.CommitUncommittedOffsets(ctx); err != nil {
		log.Error().Err(err).Msg("failed to commit offsets")
	}
}
