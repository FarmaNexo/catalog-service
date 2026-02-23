// internal/application/handlers/update_category_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

type UpdateCategoryHandler struct {
	categoryRepo repositories.CategoryRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewUpdateCategoryHandler(categoryRepo repositories.CategoryRepository, cacheService services.CacheService, logger *zap.Logger) *UpdateCategoryHandler {
	return &UpdateCategoryHandler{categoryRepo: categoryRepo, cacheService: cacheService, logger: logger}
}

func (h *UpdateCategoryHandler) Handle(ctx context.Context, cmd commands.UpdateCategoryCommand) (*common.ApiResponse[responses.CategoryResponse], error) {
	category, err := h.categoryRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return common.NotFoundResponse[responses.CategoryResponse]("Categoría no encontrada"), nil
	}

	// Verify slug uniqueness if changed
	if cmd.Slug != "" && cmd.Slug != category.Slug {
		existing, _ := h.categoryRepo.FindBySlug(ctx, cmd.Slug)
		if existing != nil && existing.ID != category.ID {
			return common.ConflictResponse[responses.CategoryResponse]("BUS_006", "El slug ya existe"), nil
		}
		category.Slug = cmd.Slug
	}

	if cmd.Name != "" {
		category.Name = cmd.Name
	}
	if cmd.Description != "" {
		category.Description = cmd.Description
	}
	if cmd.ImageURL != "" {
		category.ImageURL = cmd.ImageURL
	}
	if cmd.ParentID != "" {
		exists, _ := h.categoryRepo.Exists(ctx, cmd.ParentID)
		if !exists {
			return common.NotFoundResponse[responses.CategoryResponse]("Categoría padre no encontrada"), nil
		}
		category.ParentID = &cmd.ParentID
	}
	category.IsActive = cmd.IsActive
	category.DisplayOrder = cmd.DisplayOrder

	if err := h.categoryRepo.Update(ctx, category); err != nil {
		h.logger.Error("Error actualizando categoría", zap.Error(err))
		return common.InternalServerErrorResponse[responses.CategoryResponse]("Error actualizando categoría"), nil
	}

	// Invalidate categories cache
	go func() {
		_ = h.cacheService.Delete(context.Background(), "cache:categories")
	}()

	h.logger.Info("Categoría actualizada exitosamente",
		zap.String("category_id", category.ID),
	)

	return common.OkResponse(responses.ToCategoryResponse(category)), nil
}

var _ mediator.RequestHandler[commands.UpdateCategoryCommand, responses.CategoryResponse] = (*UpdateCategoryHandler)(nil)
