// internal/domain/repositories/brand_repository.go
package repositories

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// BrandRepository interfaz del repositorio de marcas
type BrandRepository interface {
	Create(ctx context.Context, brand *entities.Brand) error
	FindByID(ctx context.Context, id string) (*entities.Brand, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Brand, error)
	FindAll(ctx context.Context) ([]entities.Brand, error)
	Update(ctx context.Context, brand *entities.Brand) error
	Exists(ctx context.Context, id string) (bool, error)
}
