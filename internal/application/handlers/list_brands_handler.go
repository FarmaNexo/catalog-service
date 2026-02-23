// internal/application/handlers/list_brands_handler.go
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

type ListBrandsHandler struct {
	brandRepo repositories.BrandRepository
	logger    *zap.Logger
}

func NewListBrandsHandler(brandRepo repositories.BrandRepository, logger *zap.Logger) *ListBrandsHandler {
	return &ListBrandsHandler{brandRepo: brandRepo, logger: logger}
}

func (h *ListBrandsHandler) Handle(ctx context.Context, query queries.ListBrandsQuery) (*common.ApiResponse[responses.BrandListResponse], error) {
	brands, err := h.brandRepo.FindAll(ctx)
	if err != nil {
		h.logger.Error("Error listando marcas", zap.Error(err))
		return common.InternalServerErrorResponse[responses.BrandListResponse]("Error listando marcas"), nil
	}

	brandResponses := make([]responses.BrandResponse, len(brands))
	for i, b := range brands {
		brandResponses[i] = responses.ToBrandResponse(&b)
	}

	response := responses.BrandListResponse{
		Brands: brandResponses,
	}

	return common.OkResponse(response), nil
}

var _ mediator.RequestHandler[queries.ListBrandsQuery, responses.BrandListResponse] = (*ListBrandsHandler)(nil)
