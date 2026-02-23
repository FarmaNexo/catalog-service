// internal/application/queries/list_fbt_query.go
package queries

// ListFBTQuery consulta para listar productos frecuentemente comprados juntos
type ListFBTQuery struct {
	ProductID string `json:"product_id"`
	Limit     int    `json:"limit"`
}

func (q ListFBTQuery) GetName() string {
	return "ListFBTQuery"
}
