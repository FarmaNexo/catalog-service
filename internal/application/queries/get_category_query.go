// internal/application/queries/get_category_query.go
package queries

// GetCategoryQuery consulta para obtener una categoría por ID
type GetCategoryQuery struct {
	ID string `json:"id"`
}

func (q GetCategoryQuery) GetName() string {
	return "GetCategoryQuery"
}
