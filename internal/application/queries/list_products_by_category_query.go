// internal/application/queries/list_products_by_category_query.go
package queries

// ListProductsByCategoryQuery consulta para listar productos por categoría
type ListProductsByCategoryQuery struct {
	CategoryID string `json:"category_id"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

func (q ListProductsByCategoryQuery) GetName() string {
	return "ListProductsByCategoryQuery"
}
