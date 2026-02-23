// internal/application/handlers/list_drug_interactions_handler.go
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

type ListDrugInteractionsHandler struct {
	interactionRepo repositories.DrugInteractionRepository
	productRepo     repositories.ProductRepository
	cacheService    services.CacheService
	logger          *zap.Logger
}

func NewListDrugInteractionsHandler(
	interactionRepo repositories.DrugInteractionRepository,
	productRepo repositories.ProductRepository,
	cacheService services.CacheService,
	logger *zap.Logger,
) *ListDrugInteractionsHandler {
	return &ListDrugInteractionsHandler{
		interactionRepo: interactionRepo,
		productRepo:     productRepo,
		cacheService:    cacheService,
		logger:          logger,
	}
}

func (h *ListDrugInteractionsHandler) Handle(ctx context.Context, query queries.ListDrugInteractionsQuery) (*common.ApiResponse[responses.InteractionListResponse], error) {
	if query.ProductID == "" {
		return common.BadRequestResponse[responses.InteractionListResponse]("VAL_001", "El ID del producto es requerido"), nil
	}

	// Verify product exists
	_, err := h.productRepo.FindByID(ctx, query.ProductID)
	if err != nil {
		return common.NotFoundResponse[responses.InteractionListResponse]("Producto no encontrado"), nil
	}

	// Try cache first
	cacheKey := "cache:product:" + query.ProductID + ":interactions"
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var listResp responses.InteractionListResponse
		if err := json.Unmarshal([]byte(cached), &listResp); err == nil {
			h.logger.Debug("Interacciones obtenidas de caché", zap.String("product_id", query.ProductID))
			return common.OkResponse(listResp), nil
		}
	}

	interactions, err := h.interactionRepo.FindByProductID(ctx, query.ProductID)
	if err != nil {
		h.logger.Error("Error listando interacciones", zap.String("product_id", query.ProductID), zap.Error(err))
		return common.InternalServerErrorResponse[responses.InteractionListResponse]("Error listando interacciones"), nil
	}

	interactionResponses := make([]responses.InteractionResponse, len(interactions))
	for i, interaction := range interactions {
		interactionResponses[i] = responses.ToInteractionResponse(&interaction)
	}

	listResp := responses.InteractionListResponse{
		Interactions: interactionResponses,
		Total:        len(interactionResponses),
	}

	// Cache for 24 hours
	if data, err := json.Marshal(listResp); err == nil {
		go func() {
			_ = h.cacheService.Set(context.Background(), cacheKey, string(data), 24*time.Hour)
		}()
	}

	return common.OkResponse(listResp), nil
}

var _ mediator.RequestHandler[queries.ListDrugInteractionsQuery, responses.InteractionListResponse] = (*ListDrugInteractionsHandler)(nil)
