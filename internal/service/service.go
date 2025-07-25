package service

import (
	"context"
	"fmt"

	"github.com/niksmo/receipt/internal/schema"
)

type Sender interface {
	Send(ctx context.Context, to string, sub string, payload []byte) error
}

type EmailNotifier struct {
	s Sender
}

func NewEmailNotifier(sender Sender) EmailNotifier {
	return EmailNotifier{sender}
}

func (n EmailNotifier) SendReceipt(ctx context.Context, receipt schema.Receipt) error {
	const op = "EmailNotifier.SendReceipt"
	payload := n.renderReciept(receipt)

	err := n.s.Send(ctx, receipt.CustomerEmail, "Receipt", payload)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (n EmailNotifier) renderReciept(receipt schema.Receipt) []byte {
	return renderReciept(receipt)
}
