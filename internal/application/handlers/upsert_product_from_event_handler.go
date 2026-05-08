// internal/application/handlers/upsert_product_from_event_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/domain/events"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"go.uber.org/zap"
)

// UpsertProductFromEventHandler procesa un PRODUCT_DISCOVERED proveniente
// del scraper. NO usa mediator (no hay request HTTP detrás): se invoca
// directamente desde el SQS consumer.
type UpsertProductFromEventHandler struct {
	productRepo repositories.ProductRepository
	logger      *zap.Logger
}

func NewUpsertProductFromEventHandler(productRepo repositories.ProductRepository, logger *zap.Logger) *UpsertProductFromEventHandler {
	return &UpsertProductFromEventHandler{productRepo: productRepo, logger: logger}
}

// Handle hace UPSERT por (source_product_code, concentration). Idempotente.
func (h *UpsertProductFromEventHandler) Handle(ctx context.Context, data events.ProductDiscoveredData) (string, error) {
	id, err := h.productRepo.UpsertBySource(ctx, repositories.ProductUpsertParams{
		SourceProductCode:    data.SourceProductCode,
		CanonicalName:        data.CanonicalName,
		ActiveIngredient:     data.ActiveIngredient,
		Concentration:        data.Concentration,
		Form:                 data.Form,
		SourceFormCode:       data.SourceFormCode,
		Presentation:         data.Presentation,
		RegistryNumber:       data.RegistryNumber,
		Manufacturer:         data.Manufacturer,
		Holder:               data.Holder,
		RequiresPrescription: data.RequiresPrescription,
	})
	if err != nil {
		h.logger.Warn("UPSERT producto desde evento falló",
			zap.Int("source_product_code", data.SourceProductCode),
			zap.String("concentration", data.Concentration),
			zap.Error(err),
		)
		return "", err
	}
	h.logger.Info("Producto upserteado desde scraper",
		zap.String("id", id),
		zap.Int("source_product_code", data.SourceProductCode),
		zap.String("name", data.CanonicalName),
	)
	return id, nil
}
