// internal/application/queries/list_categories_query.go
package queries

// ListCategoriesQuery consulta para listar categorías
type ListCategoriesQuery struct{}

func (q ListCategoriesQuery) GetName() string {
	return "ListCategoriesQuery"
}
