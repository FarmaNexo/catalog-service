// internal/domain/services/pharmacy_client.go
package services

import "context"

// PharmacyInventoryItem representa un item de inventario de farmacia
type PharmacyInventoryItem struct {
	PharmacyID   string  `json:"pharmacy_id"`
	PharmacyName string  `json:"pharmacy_name"`
	Stock        int     `json:"stock"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
}

// PharmacyClient interfaz para comunicación con Pharmacy Service
type PharmacyClient interface {
	GetProductAvailability(ctx context.Context, productID string) ([]PharmacyInventoryItem, error)
}
