// internal/application/handlers/get_product_handler.go
package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/farmanexo/catalog-service/internal/application/queries"
	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

type GetProductHandler struct {
	productRepo  repositories.ProductRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewGetProductHandler(productRepo repositories.ProductRepository, cacheService services.CacheService, logger *zap.Logger) *GetProductHandler {
	return &GetProductHandler{productRepo: productRepo, cacheService: cacheService, logger: logger}
}

func (h *GetProductHandler) Handle(ctx context.Context, query queries.GetProductQuery) (*common.ApiResponse[responses.ProductResponse], error) {
	// Try cache first
	cacheKey := "cache:product:" + query.ID
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var productResp responses.ProductResponse
		if err := json.Unmarshal([]byte(cached), &productResp); err == nil {
			h.logger.Debug("Producto obtenido de caché", zap.String("product_id", query.ID))
			return common.OkResponse(productResp), nil
		}
	}

	var product *entities.Product
	var err error

	if query.IsAdmin {
		product, err = h.productRepo.FindByIDWithDeleted(ctx, query.ID)
	} else {
		product, err = h.productRepo.FindByID(ctx, query.ID)
	}

	if err != nil {
		h.logger.Warn("Producto no encontrado", zap.String("product_id", query.ID), zap.Error(err))
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado"), nil
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

var _ mediator.RequestHandler[queries.GetProductQuery, responses.ProductResponse] = (*GetProductHandler)(nil)
