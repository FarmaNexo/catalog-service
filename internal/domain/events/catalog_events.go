// internal/domain/events/catalog_events.go
package events

import "time"

// CatalogEvent representa un evento del catálogo
type CatalogEvent struct {
	EventType string            `json:"event_type"`
	ProductID string            `json:"product_id,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Metadata  map[string]string `json:"metadata"`
}

// Tipos de eventos del catálogo
const (
	EventProductCreated = "PRODUCT_CREATED"
	EventProductUpdated = "PRODUCT_UPDATED"
	EventProductDeleted = "PRODUCT_DELETED"
)

// NewCatalogEvent crea un nuevo evento del catálogo
func NewCatalogEvent(eventType, productID string) CatalogEvent {
	return CatalogEvent{
		EventType: eventType,
		ProductID: productID,
		Timestamp: time.Now(),
		Metadata: map[string]string{
			"source":  "catalog-service",
			"version": "1.0",
		},
	}
}
