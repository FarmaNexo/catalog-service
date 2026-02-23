// internal/infrastructure/persistence/postgres/fbt_repository_impl.go
package postgres

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// FBTRepositoryImpl implementación PostgreSQL del repositorio FBT
type FBTRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewFBTRepository(db *gorm.DB, logger *zap.Logger) *FBTRepositoryImpl {
	return &FBTRepositoryImpl{db: db, logger: logger}
}

func (r *FBTRepositoryImpl) FindByProductID(ctx context.Context, productID string, limit int) ([]entities.FrequentlyBoughtTogether, error) {
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	var results []entities.FrequentlyBoughtTogether
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Preload("RelatedProduct").
		Preload("RelatedProduct.Brand").
		Preload("RelatedProduct.Images", func(db *gorm.DB) *gorm.DB {
			return db.Where("is_primary = ?", true).Limit(1)
		}).
		Order("score DESC").
		Limit(limit).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// Compile-time interface check
var _ repositories.FBTRepository = (*FBTRepositoryImpl)(nil)
