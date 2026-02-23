// internal/infrastructure/persistence/postgres/product_image_repository_impl.go
package postgres

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductImageRepositoryImpl implementación PostgreSQL del repositorio de imágenes
type ProductImageRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProductImageRepository(db *gorm.DB, logger *zap.Logger) *ProductImageRepositoryImpl {
	return &ProductImageRepositoryImpl{db: db, logger: logger}
}

func (r *ProductImageRepositoryImpl) Create(ctx context.Context, image *entities.ProductImage) error {
	return r.db.WithContext(ctx).Create(image).Error
}

func (r *ProductImageRepositoryImpl) FindByProductID(ctx context.Context, productID string) ([]entities.ProductImage, error) {
	var images []entities.ProductImage
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&images).Error
	if err != nil {
		return nil, err
	}
	return images, nil
}

func (r *ProductImageRepositoryImpl) ClearPrimaryByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).
		Model(&entities.ProductImage{}).
		Where("product_id = ?", productID).
		Update("is_primary", false).Error
}

func (r *ProductImageRepositoryImpl) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entities.ProductImage{}).Error
}

func (r *ProductImageRepositoryImpl) CountByProductID(ctx context.Context, productID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.ProductImage{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

// Compile-time interface check
var _ repositories.ProductImageRepository = (*ProductImageRepositoryImpl)(nil)
