package service

import (
	"context"

	"github.com/niksmo/receipt/internal/mock_notifier/core/domain"
	"github.com/niksmo/receipt/internal/mock_notifier/core/port"
	"github.com/niksmo/receipt/pkg/logger"
)

var _ port.MessagePrinter = (*Service)(nil)

type Service struct {
	log logger.Logger
}

func NewService(log logger.Logger) Service {
	return Service{log}
}

func (s Service) PrintMessage(
	ctx context.Context, msg domain.Message,
) (domain.MessageID, error) {
	const op = "Service.PrintMessage"
	log := s.log.WithOp(op)

	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	log.Debug().Str(
		"to", msg.ToEmail).Str("subject", msg.Subject).Msg("send message")

	return domain.NewMessageID(), nil
}
