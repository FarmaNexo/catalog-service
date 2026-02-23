// internal/domain/repositories/drug_interaction_repository.go
package repositories

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// DrugInteractionRepository interfaz del repositorio de interacciones medicamentosas
type DrugInteractionRepository interface {
	Create(ctx context.Context, interaction *entities.DrugInteraction) error
	FindByProductID(ctx context.Context, productID string) ([]entities.DrugInteraction, error)
	FindByID(ctx context.Context, id string) (*entities.DrugInteraction, error)
	Exists(ctx context.Context, productID, interactsWithProductID string) (bool, error)
}
