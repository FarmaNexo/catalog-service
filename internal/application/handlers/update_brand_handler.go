// internal/application/handlers/update_brand_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

type UpdateBrandHandler struct {
	brandRepo repositories.BrandRepository
	logger    *zap.Logger
}

func NewUpdateBrandHandler(brandRepo repositories.BrandRepository, logger *zap.Logger) *UpdateBrandHandler {
	return &UpdateBrandHandler{brandRepo: brandRepo, logger: logger}
}

func (h *UpdateBrandHandler) Handle(ctx context.Context, cmd commands.UpdateBrandCommand) (*common.ApiResponse[responses.BrandResponse], error) {
	brand, err := h.brandRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return common.NotFoundResponse[responses.BrandResponse]("Marca no encontrada"), nil
	}

	// Verify slug uniqueness if changed
	if cmd.Slug != "" && cmd.Slug != brand.Slug {
		existing, _ := h.brandRepo.FindBySlug(ctx, cmd.Slug)
		if existing != nil && existing.ID != brand.ID {
			return common.ConflictResponse[responses.BrandResponse]("BUS_006", "El slug ya existe"), nil
		}
		brand.Slug = cmd.Slug
	}

	if cmd.Name != "" {
		brand.Name = cmd.Name
	}
	if cmd.Description != "" {
		brand.Description = cmd.Description
	}
	if cmd.LogoURL != "" {
		brand.LogoURL = cmd.LogoURL
	}
	if cmd.Website != "" {
		brand.Website = cmd.Website
	}
	if cmd.Country != "" {
		brand.Country = cmd.Country
	}
	brand.IsActive = cmd.IsActive

	if err := h.brandRepo.Update(ctx, brand); err != nil {
		h.logger.Error("Error actualizando marca", zap.Error(err))
		return common.InternalServerErrorResponse[responses.BrandResponse]("Error actualizando marca"), nil
	}

	h.logger.Info("Marca actualizada exitosamente",
		zap.String("brand_id", brand.ID),
	)

	return common.OkResponse(responses.ToBrandResponse(brand)), nil
}

var _ mediator.RequestHandler[commands.UpdateBrandCommand, responses.BrandResponse] = (*UpdateBrandHandler)(nil)
