// internal/application/handlers/get_product_by_barcode_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/farmanexo/catalog-service/internal/application/queries"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

type GetProductByBarcodeHandler struct {
	productRepo  repositories.ProductRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewGetProductByBarcodeHandler(productRepo repositories.ProductRepository, cacheService services.CacheService, logger *zap.Logger) *GetProductByBarcodeHandler {
	return &GetProductByBarcodeHandler{productRepo: productRepo, cacheService: cacheService, logger: logger}
}

func (h *GetProductByBarcodeHandler) Handle(ctx context.Context, query queries.GetProductByBarcodeQuery) (*common.ApiResponse[responses.ProductResponse], error) {
	if query.Barcode == "" {
		return common.BadRequestResponse[responses.ProductResponse]("VAL_001", "El código de barras es requerido"), nil
	}

	// Try cache first
	cacheKey := "cache:product:barcode:" + query.Barcode
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var productResp responses.ProductResponse
		if err := json.Unmarshal([]byte(cached), &productResp); err == nil {
			h.logger.Debug("Producto obtenido de caché por barcode", zap.String("barcode", query.Barcode))
			return common.OkResponse(productResp), nil
		}
	}

	product, err := h.productRepo.FindByBarcode(ctx, query.Barcode)
	if err != nil {
		h.logger.Warn("Producto no encontrado por barcode", zap.String("barcode", query.Barcode), zap.Error(err))
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado con el código de barras proporcionado"), nil
	}

	productResp := responses.ToProductResponse(product)

	// Cache for 1 hour
	if data, err := json.Marshal(productResp); err == nil {
		go func() {
			_ = h.cacheService.Set(context.Background(), cacheKey, string(data), 1*time.Hour)
		}()
	}

	return common.OkResponse(productResp), nil
}

var _ mediator.RequestHandler[queries.GetProductByBarcodeQuery, responses.ProductResponse] = (*GetProductByBarcodeHandler)(nil)
