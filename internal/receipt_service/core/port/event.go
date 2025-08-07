package port

import (
	"context"

	"github.com/niksmo/receipt/internal/receipt_service/core/domain"
)

type EventSaver interface {
	SaveEvent(context.Context, domain.Receipt) error
}

type EventProducer interface {
	ProduceEvent(context.Context, domain.Receipt) error
}

type EventProcessor interface {
	ProcessEvent(context.Context, []domain.Receipt)
}
