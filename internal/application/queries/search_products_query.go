// internal/application/queries/search_products_query.go
package queries

// SearchProductsQuery consulta para búsqueda avanzada de productos
type SearchProductsQuery struct {
	Query                string `json:"query"`
	CategoryID           string `json:"category_id"`
	BrandID              string `json:"brand_id"`
	ActiveIngredient     string `json:"active_ingredient,omitempty"` // DCI — match exacto case-insensitive
	ExcludeID            string `json:"exclude_id,omitempty"`        // Excluye un producto del resultado (HU-015 alternativas)
	RequiresPrescription *bool  `json:"requires_prescription"`
	Page                 int    `json:"page"`
	Limit                int    `json:"limit"`
	Sort                 string `json:"sort"`
}

func (q SearchProductsQuery) GetName() string {
	return "SearchProductsQuery"
}
