// internal/application/queries/get_product_availability_query.go
package queries

// GetProductAvailabilityQuery consulta para obtener disponibilidad de un producto en farmacias.
//
// Geolocalización (HU-014, opcional):
//   - Latitude/Longitude: si ambos son != 0, la respuesta incluye distance_km
//     por farmacia y se ordena por cercanía.
//   - RadiusKm: si > 0 y hay lat/lng, filtra farmacias dentro de ese radio.
type GetProductAvailabilityQuery struct {
	ProductID string  `json:"product_id"`
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
	RadiusKm  float64 `json:"radius_km,omitempty"`
}

func (q GetProductAvailabilityQuery) GetName() string {
	return "GetProductAvailabilityQuery"
}
