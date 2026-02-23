// internal/application/queries/get_product_availability_query.go
package queries

// GetProductAvailabilityQuery consulta para obtener disponibilidad de un producto en farmacias
type GetProductAvailabilityQuery struct {
	ProductID string `json:"product_id"`
}

func (q GetProductAvailabilityQuery) GetName() string {
	return "GetProductAvailabilityQuery"
}
