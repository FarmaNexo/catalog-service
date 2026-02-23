// internal/application/handlers/search_products_handler.go
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

type SearchProductsHandler struct {
	productRepo repositories.ProductRepository
	logger      *zap.Logger
}

func NewSearchProductsHandler(productRepo repositories.ProductRepository, logger *zap.Logger) *SearchProductsHandler {
	return &SearchProductsHandler{productRepo: productRepo, logger: logger}
}

func (h *SearchProductsHandler) Handle(ctx context.Context, query queries.SearchProductsQuery) (*common.ApiResponse[responses.ProductListResponse], error) {
	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	isActive := true
	params := repositories.ProductSearchParams{
		Query:                query.Query,
		CategoryID:           query.CategoryID,
		BrandID:              query.BrandID,
		RequiresPrescription: query.RequiresPrescription,
		IsActive:             &isActive,
		Page:                 page,
		Limit:                limit,
		Sort:                 query.Sort,
	}

	result, err := h.productRepo.Search(ctx, params)
	if err != nil {
		h.logger.Error("Error en búsqueda de productos", zap.Error(err))
		return common.InternalServerErrorResponse[responses.ProductListResponse]("Error en búsqueda de productos"), nil
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

var _ mediator.RequestHandler[queries.SearchProductsQuery, responses.ProductListResponse] = (*SearchProductsHandler)(nil)
