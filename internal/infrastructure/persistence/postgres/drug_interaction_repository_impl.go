// internal/infrastructure/persistence/postgres/drug_interaction_repository_impl.go
package postgres

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DrugInteractionRepositoryImpl implementación PostgreSQL del repositorio de interacciones
type DrugInteractionRepositoryImpl struct {
	db     *gorm.DB
	logger *zap.Logger
}

func NewDrugInteractionRepository(db *gorm.DB, logger *zap.Logger) *DrugInteractionRepositoryImpl {
	return &DrugInteractionRepositoryImpl{db: db, logger: logger}
}

func (r *DrugInteractionRepositoryImpl) Create(ctx context.Context, interaction *entities.DrugInteraction) error {
	return r.db.WithContext(ctx).Create(interaction).Error
}

func (r *DrugInteractionRepositoryImpl) FindByProductID(ctx context.Context, productID string) ([]entities.DrugInteraction, error) {
	var interactions []entities.DrugInteraction
	err := r.db.WithContext(ctx).
		Where("product_id = ? OR interacts_with_product_id = ?", productID, productID).
		Preload("Product").
		Preload("InteractsWithProduct").
		Order("severity DESC, created_at ASC").
		Find(&interactions).Error
	if err != nil {
		return nil, err
	}
	return interactions, nil
}

func (r *DrugInteractionRepositoryImpl) FindByID(ctx context.Context, id string) (*entities.DrugInteraction, error) {
	var interaction entities.DrugInteraction
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		Preload("Product").
		Preload("InteractsWithProduct").
		First(&interaction).Error
	if err != nil {
		return nil, err
	}
	return &interaction, nil
}

func (r *DrugInteractionRepositoryImpl) Exists(ctx context.Context, productID, interactsWithProductID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&entities.DrugInteraction{}).
		Where("(product_id = ? AND interacts_with_product_id = ?) OR (product_id = ? AND interacts_with_product_id = ?)",
			productID, interactsWithProductID, interactsWithProductID, productID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Compile-time interface check
var _ repositories.DrugInteractionRepository = (*DrugInteractionRepositoryImpl)(nil)
