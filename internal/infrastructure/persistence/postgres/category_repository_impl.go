// internal/infrastructure/persistence/postgres/category_repository_impl.go
package postgres

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CategoryRepositoryImpl implementación PostgreSQL del repositorio de categorías
type CategoryRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewCategoryRepository(db *gorm.DB, logger *zap.Logger) *CategoryRepositoryImpl {
	return &CategoryRepositoryImpl{db: db, logger: logger}
}

func (r *CategoryRepositoryImpl) Create(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *CategoryRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Children", "is_active = ?", true).
		First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entities.Category, error) {
	var category entities.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepositoryImpl) FindAll(ctx context.Context) ([]entities.Category, error) {
	var categories []entities.Category
	err := r.db.WithContext(ctx).
		Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Preload("Children", "is_active = ?", true).
		Find(&categories).Error
	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepositoryImpl) Update(ctx context.Context, category *entities.Category) error {
	return r.db.WithContext(ctx).Save(category).Error
}

func (r *CategoryRepositoryImpl) Exists(ctx context.Context, id string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.Category{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// Compile-time interface check
var _ repositories.CategoryRepository = (*CategoryRepositoryImpl)(nil)
