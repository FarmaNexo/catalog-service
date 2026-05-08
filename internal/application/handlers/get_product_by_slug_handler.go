// internal/application/handlers/get_product_by_slug_handler.go
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

type GetProductBySlugHandler struct {
	productRepo  repositories.ProductRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewGetProductBySlugHandler(productRepo repositories.ProductRepository, cacheService services.CacheService, logger *zap.Logger) *GetProductBySlugHandler {
	return &GetProductBySlugHandler{productRepo: productRepo, cacheService: cacheService, logger: logger}
}

func (h *GetProductBySlugHandler) Handle(ctx context.Context, query queries.GetProductBySlugQuery) (*common.ApiResponse[responses.ProductResponse], error) {
	cacheKey := "cache:product:slug:" + query.Slug
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var productResp responses.ProductResponse
		if err := json.Unmarshal([]byte(cached), &productResp); err == nil {
			h.logger.Debug("Producto obtenido de caché por slug", zap.String("slug", query.Slug))
			return common.OkResponse(productResp), nil
		}
	}

	product, err := h.productRepo.FindBySlug(ctx, query.Slug)
	if err != nil || product == nil {
		h.logger.Warn("Producto no encontrado por slug", zap.String("slug", query.Slug), zap.Error(err))
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado"), nil
	}

	if !query.IsAdmin && !product.IsActive {
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado"), nil
	}

	productResp := responses.ToProductResponse(product)

	if data, err := json.Marshal(productResp); err == nil {
		go func() {
			_ = h.cacheService.Set(context.Background(), cacheKey, string(data), 1*time.Hour)
		}()
	}

	return common.OkResponse(productResp), nil
}

var _ mediator.RequestHandler[queries.GetProductBySlugQuery, responses.ProductResponse] = (*GetProductBySlugHandler)(nil)
