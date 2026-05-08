// internal/domain/services/pharmacy_client.go
package services

import "context"

// PharmacyInventoryItem representa un item de inventario de farmacia.
// Los campos District y Address vienen del JOIN con `pharmacy.pharmacies` en
// pharmacy-service (city = distrito por convención DIGEMID Perú; street = dirección).
// DistanceKm se llena solo cuando GetProductAvailability se invoca con geo (HU-014).
// HU-016 — DistrictAvgPrice + IsOverpriced + OverpricePct: viajan calculados
// desde pharmacy-service. catalog-service NO recalcula nada, solo propaga.
type PharmacyInventoryItem struct {
	PharmacyID       string   `json:"pharmacy_id"`
	PharmacySlug     string   `json:"pharmacy_slug"`
	PharmacyName     string   `json:"pharmacy_name"`
	PharmacyDistrict string   `json:"pharmacy_district"`
	PharmacyAddress  string   `json:"pharmacy_address"`
	Stock            int      `json:"stock"`
	Price            float64  `json:"price"`
	IsAvailable      bool     `json:"is_available"`
	DistanceKm       *float64 `json:"distance_km,omitempty"`
	DistrictAvgPrice *float64 `json:"district_avg_price,omitempty"`
	IsOverpriced     bool     `json:"is_overpriced"`
	OverpricePct     *float64 `json:"overprice_pct,omitempty"`
}

// AvailabilityGeo agrupa los parámetros de geolocalización opcionales para
// el endpoint de disponibilidad. Si Lat == 0 && Lng == 0, el cliente
// consume pharmacy-service sin parámetros geo (orden por precio).
type AvailabilityGeo struct {
	Lat      float64
	Lng      float64
	RadiusKm float64
}

func (g AvailabilityGeo) IsActive() bool { return g.Lat != 0 || g.Lng != 0 }

// PharmacyClient interfaz para comunicación con Pharmacy Service
type PharmacyClient interface {
	GetProductAvailability(ctx context.Context, productID string, geo AvailabilityGeo) ([]PharmacyInventoryItem, error)
}
