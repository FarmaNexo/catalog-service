// internal/application/handlers/list_categories_handler.go
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

type ListCategoriesHandler struct {
	categoryRepo repositories.CategoryRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewListCategoriesHandler(categoryRepo repositories.CategoryRepository, cacheService services.CacheService, logger *zap.Logger) *ListCategoriesHandler {
	return &ListCategoriesHandler{categoryRepo: categoryRepo, cacheService: cacheService, logger: logger}
}

func (h *ListCategoriesHandler) Handle(ctx context.Context, query queries.ListCategoriesQuery) (*common.ApiResponse[responses.CategoryListResponse], error) {
	// Try cache
	cacheKey := "cache:categories"
	if cached, err := h.cacheService.Get(ctx, cacheKey); err == nil && cached != "" {
		var listResp responses.CategoryListResponse
		if err := json.Unmarshal([]byte(cached), &listResp); err == nil {
			h.logger.Debug("Categorías obtenidas de caché")
			return common.OkResponse(listResp), nil
		}
	}

	categories, err := h.categoryRepo.FindAll(ctx)
	if err != nil {
		h.logger.Error("Error listando categorías", zap.Error(err))
		return common.InternalServerErrorResponse[responses.CategoryListResponse]("Error listando categorías"), nil
	}

	categoryResponses := make([]responses.CategoryResponse, len(categories))
	for i, c := range categories {
		categoryResponses[i] = responses.ToCategoryResponse(&c)
	}

	response := responses.CategoryListResponse{
		Categories: categoryResponses,
	}

	// Cache for 6 hours
	if data, err := json.Marshal(response); err == nil {
		go func() {
			_ = h.cacheService.Set(context.Background(), cacheKey, string(data), 6*time.Hour)
		}()
	}

	return common.OkResponse(response), nil
}

var _ mediator.RequestHandler[queries.ListCategoriesQuery, responses.CategoryListResponse] = (*ListCategoriesHandler)(nil)
