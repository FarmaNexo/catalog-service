// internal/infrastructure/persistence/postgres/brand_repository_impl.go
package postgres

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// BrandRepositoryImpl implementación PostgreSQL del repositorio de marcas
type BrandRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewBrandRepository(db *gorm.DB, logger *zap.Logger) *BrandRepositoryImpl {
	return &BrandRepositoryImpl{db: db, logger: logger}
}

func (r *BrandRepositoryImpl) Create(ctx context.Context, brand *entities.Brand) error {
	return r.db.WithContext(ctx).Create(brand).Error
}

func (r *BrandRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Brand, error) {
	var brand entities.Brand
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *BrandRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entities.Brand, error) {
	var brand entities.Brand
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&brand).Error
	if err != nil {
		return nil, err
	}
	return &brand, nil
}

func (r *BrandRepositoryImpl) FindAll(ctx context.Context) ([]entities.Brand, error) {
	var brands []entities.Brand
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&brands).Error
	if err != nil {
		return nil, err
	}
	return brands, nil
}

func (r *BrandRepositoryImpl) Update(ctx context.Context, brand *entities.Brand) error {
	return r.db.WithContext(ctx).Save(brand).Error
}

func (r *BrandRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Brand{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Compile-time interface check
var _ repositories.BrandRepository = (*BrandRepositoryImpl)(nil)
