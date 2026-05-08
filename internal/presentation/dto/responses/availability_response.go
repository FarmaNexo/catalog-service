// internal/presentation/dto/responses/availability_response.go
package responses

// PharmacyAvailability disponibilidad de un producto en una farmacia.
// Cumple los datos mínimos del Excel HU-013: nombre, distrito, dirección, precio.
// HU-014: distance_km solo presente si la query incluyó lat/lng.
// HU-016: district_avg_price + is_overpriced + overprice_pct calculados por
// pharmacy-service y propagados sin modificar.
type PharmacyAvailability struct {
	PharmacyID       string   `json:"pharmacy_id"`
	PharmacySlug     string   `json:"pharmacy_slug,omitempty"`
	PharmacyName     string   `json:"pharmacy_name"`
	PharmacyDistrict string   `json:"pharmacy_district,omitempty"`
	PharmacyAddress  string   `json:"pharmacy_address,omitempty"`
	Stock            int      `json:"stock"`
	Price            float64  `json:"price"`
	IsAvailable      bool     `json:"is_available"`
	DistanceKm       *float64 `json:"distance_km,omitempty"`
	DistrictAvgPrice *float64 `json:"district_avg_price,omitempty"`
	IsOverpriced     bool     `json:"is_overpriced"`
	OverpricePct     *float64 `json:"overprice_pct,omitempty"`
}

// AvailabilityResponse respuesta de disponibilidad de un producto
type AvailabilityResponse struct {
	ProductID   string                 `json:"product_id"`
	ProductName string                 `json:"product_name"`
	Pharmacies  []PharmacyAvailability `json:"pharmacies"`
	Total       int                    `json:"total_pharmacies"`
}
