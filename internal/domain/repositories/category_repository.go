// internal/domain/repositories/category_repository.go
package repositories

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// CategoryRepository interfaz del repositorio de categorías
type CategoryRepository interface {
	Create(ctx context.Context, category *entities.Category) error
	FindByID(ctx context.Context, id string) (*entities.Category, error)
	FindBySlug(ctx context.Context, slug string) (*entities.Category, error)
	FindAll(ctx context.Context) ([]entities.Category, error)
	Update(ctx context.Context, category *entities.Category) error
	Exists(ctx context.Context, id string) (bool, error)
}
