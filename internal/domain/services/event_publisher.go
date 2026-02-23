// internal/domain/services/event_publisher.go
package services

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/events"
)

// EventPublisher interfaz para publicar eventos del catálogo
type EventPublisher interface {
	Publish(ctx context.Context, event events.CatalogEvent) error
}
