package service

import (
	"context"
	"fmt"

	"github.com/niksmo/receipt/internal/mail-sender/domain"
	"github.com/niksmo/receipt/pkg/logger"
)

type MailProvider interface {
	SendMail(context.Context, domain.Receipt) error
}

type ReceiverService struct {
	log          logger.Logger
	mailProvider MailProvider
}

func NewReceiverService(log logger.Logger, mailProvider MailProvider) ReceiverService {
	return ReceiverService{log, mailProvider}
}

func (s ReceiverService) SendReceipt(
	ctx context.Context, receipt domain.Receipt,
) error {
	const op = "ReceiverService.SendReceipt"

	if ctx.Err() != nil {
		return ctx.Err()
	}

	err := s.mailProvider.SendMail(ctx, receipt)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
