// internal/domain/repositories/product_image_repository.go
package repositories

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// ProductImageRepository interfaz del repositorio de imágenes de productos
type ProductImageRepository interface {
	Create(ctx context.Context, image *entities.ProductImage) error
	FindByProductID(ctx context.Context, productID string) ([]entities.ProductImage, error)
	ClearPrimaryByProductID(ctx context.Context, productID string) error
	Delete(ctx context.Context, id string) error
	CountByProductID(ctx context.Context, productID string) (int64, error)
}
