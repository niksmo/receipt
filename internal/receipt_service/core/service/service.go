package service

import (
	"context"
	"fmt"
	"time"

	"github.com/niksmo/receipt/internal/receipt_service/core/domain"
	"github.com/niksmo/receipt/internal/receipt_service/core/port"
	"github.com/niksmo/receipt/pkg/logger"
)

const produceTimeout = 3 * time.Second

var _ port.EventSaver = (*Service)(nil)
var _ port.EventProcessor = (*Service)(nil)

type Service struct {
	log  logger.Logger
	evtP port.EventProducer
}

func NewService(log logger.Logger, evtP port.EventProducer) *Service {
	return &Service{log, evtP}
}

func (s *Service) SaveEvent(ctx context.Context, rct domain.Receipt) error {
	const op = "Service.SaveEvent"
	err := s.evtP.ProduceEvent(ctx, rct)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Service) ProcessEvent(ctx context.Context, rcts []domain.Receipt) {
	const op = "Service.ProcessEvent"
	log := s.log.WithOp(op)
	log.Debug().Msg("processed")
}
