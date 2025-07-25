package service

import (
	"context"
	"fmt"

	"github.com/niksmo/receipt/internal/schema"
	"github.com/niksmo/receipt/pkg/logger"
)

type Sender interface {
	Send(ctx context.Context, to string, payload []byte) error
}

type EmailNotifier struct {
	log logger.Logger
	s   Sender
}

func NewEmailNotifier(log logger.Logger, sender Sender) EmailNotifier {
	return EmailNotifier{log, sender}
}

func (n EmailNotifier) SendReceipt(ctx context.Context, receipt schema.Receipt) error {
	const op = "EmailNotifier.SendReceipt"
	payload := n.renderReciept(receipt)

	err := n.s.Send(ctx, receipt.BuyerEmail, payload)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (n EmailNotifier) renderReciept(receipt schema.Receipt) []byte {
	return renderReciept(receipt)
}
