package port

import (
	"context"

	"github.com/niksmo/receipt/internal/mock_notifier/core/domain"
)

type MessagePrinter interface {
	PrintMessage(ctx context.Context, msg domain.Message) (domain.MessageID, error)
}
