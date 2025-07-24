package service

import (
	"fmt"

	"github.com/niksmo/receipt/internal/schema"
	"github.com/niksmo/receipt/pkg/logger"
)

type Sender interface {
	Send(to string, msg []byte) error
}

type EmailNotifier struct {
	log logger.Logger
	s   Sender
}

func NewEmailProvider(log logger.Logger, sender Sender) EmailNotifier {
	return EmailNotifier{log, sender}
}

func (n EmailNotifier) Send(receipt schema.Receipt) error {
	const op = "EmailNotifier.Send"
	// make msg
	msg := []byte("hello wolrd")

	err := n.s.Send(receipt.BuyerEmail, msg)
	if err != nil {
		return fmt.Errorf("%s: failed to send email: %w", op, err)
	}
	return nil
}

func (n EmailNotifier) createMessage(receipt schema.Receipt) []byte {
	return nil
}
