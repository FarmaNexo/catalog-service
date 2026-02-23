// internal/presentation/dto/responses/availability_response.go
package responses

// PharmacyAvailability disponibilidad de un producto en una farmacia
type PharmacyAvailability struct {
	PharmacyID   string  `json:"pharmacy_id"`
	PharmacyName string  `json:"pharmacy_name"`
	Stock        int     `json:"stock"`
	Price        float64 `json:"price"`
	IsAvailable  bool    `json:"is_available"`
}

// AvailabilityResponse respuesta de disponibilidad de un producto
type AvailabilityResponse struct {
	ProductID   string                 `json:"product_id"`
	ProductName string                 `json:"product_name"`
	Pharmacies  []PharmacyAvailability `json:"pharmacies"`
	Total       int                    `json:"total_pharmacies"`
}
