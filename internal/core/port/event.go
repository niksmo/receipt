package port

import (
	"context"

	"github.com/niksmo/receipt/internal/core/domain"
)

type EventSaver interface {
	SaveEvent(context.Context, domain.Receipt) error
}
