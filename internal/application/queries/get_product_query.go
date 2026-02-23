// internal/application/queries/get_product_query.go
package queries

// GetProductQuery consulta para obtener un producto por ID
type GetProductQuery struct {
	ID      string `json:"id"`
	IsAdmin bool   `json:"is_admin"`
}

func (q GetProductQuery) GetName() string {
	return "GetProductQuery"
}
