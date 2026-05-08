// internal/application/queries/get_product_by_slug_query.go
package queries

// GetProductBySlugQuery consulta para obtener un producto por slug (SEO-friendly URL).
type GetProductBySlugQuery struct {
	Slug    string `json:"slug"`
	IsAdmin bool   `json:"is_admin"`
}

func (q GetProductBySlugQuery) GetName() string {
	return "GetProductBySlugQuery"
}
