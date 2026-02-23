// internal/application/queries/search_products_query.go
package queries

// SearchProductsQuery consulta para búsqueda avanzada de productos
type SearchProductsQuery struct {
	Query                string `json:"query"`
	CategoryID           string `json:"category_id"`
	BrandID              string `json:"brand_id"`
	RequiresPrescription *bool  `json:"requires_prescription"`
	Page                 int    `json:"page"`
	Limit                int    `json:"limit"`
	Sort                 string `json:"sort"`
}

func (q SearchProductsQuery) GetName() string {
	return "SearchProductsQuery"
}
