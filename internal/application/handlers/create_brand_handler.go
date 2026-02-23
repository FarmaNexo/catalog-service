// internal/application/handlers/create_brand_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateBrandHandler struct {
	brandRepo repositories.BrandRepository
	logger    *zap.Logger
}

func NewCreateBrandHandler(brandRepo repositories.BrandRepository, logger *zap.Logger) *CreateBrandHandler {
	return &CreateBrandHandler{brandRepo: brandRepo, logger: logger}
}

func (h *CreateBrandHandler) Handle(ctx context.Context, cmd commands.CreateBrandCommand) (*common.ApiResponse[responses.BrandResponse], error) {
	slug := cmd.Slug
	if slug == "" {
		slug = GenerateSlug(cmd.Name)
	}

	// Verify slug uniqueness
	existing, _ := h.brandRepo.FindBySlug(ctx, slug)
	if existing != nil {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	brand := &entities.Brand{
		ID:          uuid.New().String(),
		Name:        cmd.Name,
		Slug:        slug,
		Description: cmd.Description,
		LogoURL:     cmd.LogoURL,
		Website:     cmd.Website,
		Country:     cmd.Country,
		IsActive:    true,
	}

	if err := h.brandRepo.Create(ctx, brand); err != nil {
		h.logger.Error("Error creando marca", zap.Error(err))
		return common.InternalServerErrorResponse[responses.BrandResponse]("Error creando marca"), nil
	}

	h.logger.Info("Marca creada exitosamente",
		zap.String("brand_id", brand.ID),
		zap.String("name", brand.Name),
	)

	return common.CreatedResponse(responses.ToBrandResponse(brand)), nil
}

var _ mediator.RequestHandler[commands.CreateBrandCommand, responses.BrandResponse] = (*CreateBrandHandler)(nil)
