// internal/application/handlers/list_products_by_brand_handler.go
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

type ListProductsByBrandHandler struct {
	productRepo repositories.ProductRepository
	brandRepo   repositories.BrandRepository
	logger      *zap.Logger
}

func NewListProductsByBrandHandler(productRepo repositories.ProductRepository, brandRepo repositories.BrandRepository, logger *zap.Logger) *ListProductsByBrandHandler {
	return &ListProductsByBrandHandler{productRepo: productRepo, brandRepo: brandRepo, logger: logger}
}

func (h *ListProductsByBrandHandler) Handle(ctx context.Context, query queries.ListProductsByBrandQuery) (*common.ApiResponse[responses.ProductListResponse], error) {
	// Verify brand exists
	_, err := h.brandRepo.FindByID(ctx, query.BrandID)
	if err != nil {
		return common.NotFoundResponse[responses.ProductListResponse]("Marca no encontrada"), nil
	}

	page := query.Page
	if page < 1 {
		page = 1
	}
	limit := query.Limit
	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := h.productRepo.FindByBrandID(ctx, query.BrandID, page, limit)
	if err != nil {
		h.logger.Error("Error listando productos por marca", zap.Error(err))
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

var _ mediator.RequestHandler[queries.ListProductsByBrandQuery, responses.ProductListResponse] = (*ListProductsByBrandHandler)(nil)
