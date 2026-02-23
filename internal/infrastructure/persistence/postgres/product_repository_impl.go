// internal/infrastructure/persistence/postgres/product_repository_impl.go
package postgres

import (
	"context"
	"math"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductRepositoryImpl implementación PostgreSQL del repositorio de productos
type ProductRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewProductRepository(db *gorm.DB, logger *zap.Logger) *ProductRepositoryImpl {
	return &ProductRepositoryImpl{db: db, logger: logger}
}

func (r *ProductRepositoryImpl) Create(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindByIDWithDeleted(ctx context.Context, id string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindBySlug(ctx context.Context, slug string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("slug = ? AND deleted_at IS NULL", slug).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) FindBySKU(ctx context.Context, sku string) (*entities.Product, error) {
	var product entities.Product
	err := r.db.WithContext(ctx).Where("sku = ? AND deleted_at IS NULL", sku).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepositoryImpl) Update(ctx context.Context, product *entities.Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *ProductRepositoryImpl) SoftDelete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).
		Model(&entities.Product{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"deleted_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"is_active":  false,
		}).Error
}

func (r *ProductRepositoryImpl) List(ctx context.Context, page, limit int, sort string) (*repositories.PaginatedResult, error) {
	var products []entities.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Product{}).Where("deleted_at IS NULL")

	query.Count(&total)

	query = r.applySort(query, sort)
	offset := (page - 1) * limit

	err := query.
		Offset(offset).
		Limit(limit).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &repositories.PaginatedResult{
		Products:   products,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}, nil
}

func (r *ProductRepositoryImpl) Search(ctx context.Context, params repositories.ProductSearchParams) (*repositories.PaginatedResult, error) {
	var products []entities.Product
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Product{}).Where("deleted_at IS NULL")

	if params.Query != "" {
		searchPattern := "%" + params.Query + "%"
		query = query.Where("name ILIKE ? OR active_ingredient ILIKE ?", searchPattern, searchPattern)
	}

	if params.CategoryID != "" {
		query = query.Where("category_id = ?", params.CategoryID)
	}

	if params.BrandID != "" {
		query = query.Where("brand_id = ?", params.BrandID)
	}

	if params.RequiresPrescription != nil {
		query = query.Where("requires_prescription = ?", *params.RequiresPrescription)
	}

	if params.IsActive != nil {
		query = query.Where("is_active = ?", *params.IsActive)
	}

	query.Count(&total)

	query = r.applySort(query, params.Sort)
	offset := (params.Page - 1) * params.Limit

	err := query.
		Offset(offset).
		Limit(params.Limit).
		Preload("Category").
		Preload("Brand").
		Preload("Images", func(db *gorm.DB) *gorm.DB {
			return db.Order("display_order ASC")
		}).
		Find(&products).Error
	if err != nil {
		return nil, err
	}

	return &repositories.PaginatedResult{
		Products:   products,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: int(math.Ceil(float64(total) / float64(params.Limit))),
	}, nil
}

func (r *ProductRepositoryImpl) FindByCategoryID(ctx context.Context, categoryID string, page, limit int) (*repositories.PaginatedResult, error) {
	return r.Search(ctx, repositories.ProductSearchParams{
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
		Sort:       "name_asc",
	})
}

func (r *ProductRepositoryImpl) FindByBrandID(ctx context.Context, brandID string, page, limit int) (*repositories.PaginatedResult, error) {
	return r.Search(ctx, repositories.ProductSearchParams{
		BrandID: brandID,
		Page:    page,
		Limit:   limit,
		Sort:    "name_asc",
	})
}

func (r *ProductRepositoryImpl) applySort(query *gorm.DB, sort string) *gorm.DB {
	switch sort {
	case "name_asc":
		return query.Order("name ASC")
	case "name_desc":
		return query.Order("name DESC")
	case "created_at_asc":
		return query.Order("created_at ASC")
	case "created_at_desc":
		return query.Order("created_at DESC")
	default:
		return query.Order("name ASC")
	}
}

// Compile-time interface check
var _ repositories.ProductRepository = (*ProductRepositoryImpl)(nil)
