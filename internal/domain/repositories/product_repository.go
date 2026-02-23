// internal/domain/repositories/product_repository.go
package repositories

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// ProductSearchParams parámetros de búsqueda avanzada
type ProductSearchParams struct {
	Query                string
	CategoryID           string
	BrandID              string
	RequiresPrescription *bool
	IsActive             *bool
	Page                 int
	Limit                int
	Sort                 string
}

// PaginatedResult resultado paginado
type PaginatedResult struct {
	Products   []entities.Product
	Total      int64
	Page       int
	Limit      int
	TotalPages int
}

// ProductRepository interfaz del repositorio de productos
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	FindByID(ctx context.Context, id string) (*entities.Product, error)
	FindByIDWithDeleted(ctx context.Context, id string) (*entities.Product, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Product, error)
	FindBySKU(ctx context.Context, sku string) (*entities.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	SoftDelete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, sort string) (*PaginatedResult, error)
	Search(ctx context.Context, params ProductSearchParams) (*PaginatedResult, error)
	FindByCategoryID(ctx context.Context, categoryID string, page, limit int) (*PaginatedResult, error)
	FindByBrandID(ctx context.Context, brandID string, page, limit int) (*PaginatedResult, error)
}
