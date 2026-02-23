// internal/presentation/dto/requests/search_products_request.go
package requests

// SearchProductsRequest DTO de request para búsqueda avanzada
type SearchProductsRequest struct {
	Query                string `json:"query,omitempty"`
	CategoryID           string `json:"category_id,omitempty"`
	BrandID              string `json:"brand_id,omitempty"`
	RequiresPrescription *bool  `json:"requires_prescription,omitempty"`
	Page                 int    `json:"page,omitempty"`
	Limit                int    `json:"limit,omitempty"`
	Sort                 string `json:"sort,omitempty"`
}
