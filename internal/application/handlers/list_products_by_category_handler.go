// internal/application/handlers/list_products_by_category_handler.go
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

type ListProductsByCategoryHandler struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
	logger       *zap.Logger
}

func NewListProductsByCategoryHandler(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository, logger *zap.Logger) *ListProductsByCategoryHandler {
	return &ListProductsByCategoryHandler{productRepo: productRepo, categoryRepo: categoryRepo, logger: logger}
}

func (h *ListProductsByCategoryHandler) Handle(ctx context.Context, query queries.ListProductsByCategoryQuery) (*common.ApiResponse[responses.ProductListResponse], error) {
	// Verify category exists
	_, err := h.categoryRepo.FindByID(ctx, query.CategoryID)
	if err != nil {
		return common.NotFoundResponse[responses.ProductListResponse]("Categoría no encontrada"), nil
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := h.productRepo.FindByCategoryID(ctx, query.CategoryID, page, limit)
	if err != nil {
		h.logger.Error("Error listando productos por categoría", zap.Error(err))
		return common.InternalServerErrorResponse[responses.ProductListResponse]("Error listando productos"), nil
	}

	productResponses := make([]responses.ProductResponse, len(result.Products))
	for i, p := range result.Products {
		productResponses[i] = responses.ToProductResponse(&p)
	}

	response := responses.ProductListResponse{
		Products:   productResponses,
		Total:      result.Total,
		Page:       result.Page,
		Limit:      result.Limit,
		TotalPages: result.TotalPages,
	}

	return common.OkResponse(response), nil
}

var _ mediator.RequestHandler[queries.ListProductsByCategoryQuery, responses.ProductListResponse] = (*ListProductsByCategoryHandler)(nil)
