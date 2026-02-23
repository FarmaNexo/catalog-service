// internal/application/queries/list_products_by_brand_query.go
package queries

// ListProductsByBrandQuery consulta para listar productos por marca
type ListProductsByBrandQuery struct {
	BrandID string `json:"brand_id"`
	Page    int    `json:"page"`
	Limit   int    `json:"limit"`
}

func (q ListProductsByBrandQuery) GetName() string {
	return "ListProductsByBrandQuery"
}
