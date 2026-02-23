// internal/application/handlers/create_category_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateCategoryHandler struct {
	categoryRepo repositories.CategoryRepository
	cacheService services.CacheService
	logger       *zap.Logger
}

func NewCreateCategoryHandler(categoryRepo repositories.CategoryRepository, cacheService services.CacheService, logger *zap.Logger) *CreateCategoryHandler {
	return &CreateCategoryHandler{categoryRepo: categoryRepo, cacheService: cacheService, logger: logger}
}

func (h *CreateCategoryHandler) Handle(ctx context.Context, cmd commands.CreateCategoryCommand) (*common.ApiResponse[responses.CategoryResponse], error) {
	slug := cmd.Slug
	if slug == "" {
		slug = GenerateSlug(cmd.Name)
	}

	// Verify slug uniqueness
	existing, _ := h.categoryRepo.FindBySlug(ctx, slug)
	if existing != nil {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	// Verify parent exists if provided
	if cmd.ParentID != "" {
		exists, _ := h.categoryRepo.Exists(ctx, cmd.ParentID)
		if !exists {
			return common.NotFoundResponse[responses.CategoryResponse]("Categoría padre no encontrada"), nil
		}
	}

	category := &entities.Category{
		ID:           uuid.New().String(),
		Name:         cmd.Name,
		Slug:         slug,
		Description:  cmd.Description,
		ImageURL:     cmd.ImageURL,
		IsActive:     true,
		DisplayOrder: cmd.DisplayOrder,
	}

	if cmd.ParentID != "" {
		category.ParentID = &cmd.ParentID
	}

	if err := h.categoryRepo.Create(ctx, category); err != nil {
		h.logger.Error("Error creando categoría", zap.Error(err))
		return common.InternalServerErrorResponse[responses.CategoryResponse]("Error creando categoría"), nil
	}

	// Invalidate categories cache
	go func() {
		_ = h.cacheService.Delete(context.Background(), "cache:categories")
	}()

	h.logger.Info("Categoría creada exitosamente",
		zap.String("category_id", category.ID),
		zap.String("name", category.Name),
	)

	return common.CreatedResponse(responses.ToCategoryResponse(category)), nil
}

var _ mediator.RequestHandler[commands.CreateCategoryCommand, responses.CategoryResponse] = (*CreateCategoryHandler)(nil)
