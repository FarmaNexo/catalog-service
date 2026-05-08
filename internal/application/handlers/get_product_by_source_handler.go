// internal/application/handlers/get_product_by_source_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/queries"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

// GetProductBySourceHandler resuelve un producto por su clave natural
// DIGEMID (source_product_code, concentration). Sin caché: la llamada
// interna desde pharmacy-service es de baja frecuencia (1 vez por evento
// INVENTORY_DISCOVERED) y queremos ver datos siempre frescos.
type GetProductBySourceHandler struct {
	productRepo repositories.ProductRepository
	logger      *zap.Logger
}

func NewGetProductBySourceHandler(productRepo repositories.ProductRepository, logger *zap.Logger) *GetProductBySourceHandler {
	return &GetProductBySourceHandler{productRepo: productRepo, logger: logger}
}

func (h *GetProductBySourceHandler) Handle(ctx context.Context, query queries.GetProductBySourceQuery) (*common.ApiResponse[responses.ProductResponse], error) {
	if query.SourceProductCode == 0 || query.Concentration == "" {
		return common.BadRequestResponse[responses.ProductResponse]("VAL_001", "source_product_code y concentration son requeridos"), nil
	}

	product, err := h.productRepo.FindBySourceCode(ctx, query.SourceProductCode, query.Concentration)
	if err != nil || product == nil {
		h.logger.Debug("Producto no encontrado por source_code+concentration",
			zap.Int("source_product_code", query.SourceProductCode),
			zap.String("concentration", query.Concentration),
		)
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado"), nil
	}

	if !product.IsActive {
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado"), nil
	}

	return common.OkResponse(responses.ToProductResponse(product)), nil
}

var _ mediator.RequestHandler[queries.GetProductBySourceQuery, responses.ProductResponse] = (*GetProductBySourceHandler)(nil)
