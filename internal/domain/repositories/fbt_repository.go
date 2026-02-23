// internal/domain/repositories/fbt_repository.go
package repositories

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// FBTRepository interfaz del repositorio de productos frecuentemente comprados juntos
type FBTRepository interface {
	FindByProductID(ctx context.Context, productID string, limit int) ([]entities.FrequentlyBoughtTogether, error)
}
