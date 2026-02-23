// internal/application/handlers/get_category_handler.go
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

type GetCategoryHandler struct {
	categoryRepo repositories.CategoryRepository
	logger       *zap.Logger
}

func NewGetCategoryHandler(categoryRepo repositories.CategoryRepository, logger *zap.Logger) *GetCategoryHandler {
	return &GetCategoryHandler{categoryRepo: categoryRepo, logger: logger}
}

func (h *GetCategoryHandler) Handle(ctx context.Context, query queries.GetCategoryQuery) (*common.ApiResponse[responses.CategoryResponse], error) {
	category, err := h.categoryRepo.FindByID(ctx, query.ID)
	if err != nil {
		return common.NotFoundResponse[responses.CategoryResponse]("Categoría no encontrada"), nil
	}

	return common.OkResponse(responses.ToCategoryResponse(category)), nil
}

var _ mediator.RequestHandler[queries.GetCategoryQuery, responses.CategoryResponse] = (*GetCategoryHandler)(nil)
