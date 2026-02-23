// internal/application/queries/list_brands_query.go
package queries

// ListBrandsQuery consulta para listar marcas
type ListBrandsQuery struct{}

func (q ListBrandsQuery) GetName() string {
	return "ListBrandsQuery"
}
