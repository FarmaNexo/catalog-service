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
	ActiveIngredient     string // Búsqueda exacta por DCI (Denominación Común Internacional)
	ExcludeID            string // Excluir un producto específico (útil para alternativas)
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

// ProductUpsertParams — payload para UpsertBySource. Refleja los campos del
// PRODUCT_DISCOVERED del scraper. La clave UPSERT es
// (source_product_code, concentration), definida en migration 000004.
type ProductUpsertParams struct {
	SourceProductCode    int
	CanonicalName        string
	ActiveIngredient     string
	Concentration        string
	Form                 string
	SourceFormCode       string
	Presentation         string
	RegistryNumber       string
	Manufacturer         string
	Holder               string
	RequiresPrescription bool
}

// ProductRepository interfaz del repositorio de productos
type ProductRepository interface {
	Create(ctx context.Context, product *entities.Product) error
	FindByID(ctx context.Context, id string) (*entities.Product, error)
	FindByIDWithDeleted(ctx context.Context, id string) (*entities.Product, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Product, error)
	FindBySKU(ctx context.Context, sku string) (*entities.Product, error)
	FindByBarcode(ctx context.Context, barcode string) (*entities.Product, error)
	// FindBySourceCode busca un producto por su clave natural DIGEMID
	// (source_product_code + concentration). Lo usa pharmacy-service via
	// HTTP para resolver inventory events.
	FindBySourceCode(ctx context.Context, sourceProductCode int, concentration string) (*entities.Product, error)
	Update(ctx context.Context, product *entities.Product) error
	// UpsertBySource hace INSERT ... ON CONFLICT (source_product_code, concentration)
	// DO UPDATE. Retorna el id (UUID) del producto resultante. Usado por el
	// consumer SQS de PRODUCT_DISCOVERED.
	UpsertBySource(ctx context.Context, params ProductUpsertParams) (string, error)
	SoftDelete(ctx context.Context, id string) error
	List(ctx context.Context, page, limit int, sort string) (*PaginatedResult, error)
	Search(ctx context.Context, params ProductSearchParams) (*PaginatedResult, error)
	FindByCategoryID(ctx context.Context, categoryID string, page, limit int) (*PaginatedResult, error)
	FindByBrandID(ctx context.Context, brandID string, page, limit int) (*PaginatedResult, error)
}
