// internal/application/handlers/list_fbt_handler.go
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

type ListFBTHandler struct {
	fbtRepo      repositories.FBTRepository
	productRepo  repositories.ProductRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewListFBTHandler(
	fbtRepo repositories.FBTRepository,
	productRepo repositories.ProductRepository,
	cacheService services.CacheService,
	logger *zap.Logger,
) *ListFBTHandler {
	return &ListFBTHandler{
		fbtRepo:      fbtRepo,
		productRepo:  productRepo,
		cacheService: cacheService,
		logger:       logger,
	}
}

func (h *ListFBTHandler) Handle(ctx context.Context, query queries.ListFBTQuery) (*common.ApiResponse[responses.FBTListResponse], error) {
	if query.ProductID == "" {
		return common.BadRequestResponse[responses.FBTListResponse]("VAL_001", "El ID del producto es requerido"), nil
	}

	// Verify product exists
	_, err := h.productRepo.FindByID(ctx, query.ProductID)
	if err != nil {
		return common.NotFoundResponse[responses.FBTListResponse]("Producto no encontrado"), nil
	}

	// Try cache first
	cacheKey := "cache:product:" + query.ProductID + ":fbt"
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var listResp responses.FBTListResponse
		if err := json.Unmarshal([]byte(cached), &listResp); err == nil {
			h.logger.Debug("FBT obtenido de caché", zap.String("product_id", query.ProductID))
			return common.OkResponse(listResp), nil
		}
	}

	limit := query.Limit
	if limit <= 0 || limit > 20 {
		limit = 10
	}

	fbts, err := h.fbtRepo.FindByProductID(ctx, query.ProductID, limit)
	if err != nil {
		h.logger.Error("Error listando FBT", zap.String("product_id", query.ProductID), zap.Error(err))
		return common.InternalServerErrorResponse[responses.FBTListResponse]("Error listando productos relacionados"), nil
	}

	items := make([]responses.FBTItemResponse, len(fbts))
	for i, fbt := range fbts {
		items[i] = responses.ToFBTItemResponse(&fbt)
	}

	listResp := responses.FBTListResponse{
		ProductID: query.ProductID,
		Items:     items,
		Total:     len(items),
	}

	// Cache for 6 hours
	if data, err := json.Marshal(listResp); err == nil {
		go func() {
			_ = h.cacheService.Set(context.Background(), cacheKey, string(data), 6*time.Hour)
		}()
	}

	return common.OkResponse(listResp), nil
}

var _ mediator.RequestHandler[queries.ListFBTQuery, responses.FBTListResponse] = (*ListFBTHandler)(nil)
