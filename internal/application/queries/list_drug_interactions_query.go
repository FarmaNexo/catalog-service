// internal/application/queries/list_drug_interactions_query.go
package queries

// ListDrugInteractionsQuery consulta para listar interacciones de un producto
type ListDrugInteractionsQuery struct {
	ProductID string `json:"product_id"`
}

func (q ListDrugInteractionsQuery) GetName() string {
	return "ListDrugInteractionsQuery"
}
