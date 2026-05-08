// internal/application/queries/get_product_by_source_query.go
package queries

// GetProductBySourceQuery busca un producto por su clave natural DIGEMID.
// Lo usa pharmacy-service (HTTP) para resolver eventos INVENTORY_DISCOVERED
// a un product_id local.
type GetProductBySourceQuery struct {
	SourceProductCode int    `json:"source_product_code"`
	Concentration     string `json:"concentration"`
}

func (q GetProductBySourceQuery) GetName() string {
	return "GetProductBySourceQuery"
}
