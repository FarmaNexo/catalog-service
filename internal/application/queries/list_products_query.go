// internal/application/queries/list_products_query.go
package queries

// ListProductsQuery consulta para listar productos con paginación
type ListProductsQuery struct {
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Sort  string `json:"sort"`
}

func (q ListProductsQuery) GetName() string {
	return "ListProductsQuery"
}
