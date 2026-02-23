// internal/application/handlers/update_product_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/events"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

type UpdateProductHandler struct {
	productRepo    repositories.ProductRepository
	categoryRepo   repositories.CategoryRepository
	brandRepo      repositories.BrandRepository
	eventPublisher services.EventPublisher
	cacheService   services.CacheService
	logger         *zap.Logger
}

func NewUpdateProductHandler(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	brandRepo repositories.BrandRepository,
	eventPublisher services.EventPublisher,
	cacheService services.CacheService,
	logger *zap.Logger,
) *UpdateProductHandler {
	return &UpdateProductHandler{
		productRepo:    productRepo,
		categoryRepo:   categoryRepo,
		brandRepo:      brandRepo,
		eventPublisher: eventPublisher,
		cacheService:   cacheService,
		logger:         logger,
	}
}

func (h *UpdateProductHandler) Handle(ctx context.Context, cmd commands.UpdateProductCommand) (*common.ApiResponse[responses.ProductResponse], error) {
	product, err := h.productRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return common.NotFoundResponse[responses.ProductResponse]("Producto no encontrado"), nil
	}

	// Verify slug uniqueness if changed
	if cmd.Slug != "" && cmd.Slug != product.Slug {
		existing, _ := h.productRepo.FindBySlug(ctx, cmd.Slug)
		if existing != nil && existing.ID != product.ID {
			return common.ConflictResponse[responses.ProductResponse]("BUS_006", "El slug ya existe"), nil
		}
		product.Slug = cmd.Slug
	}

	// Verify SKU uniqueness if changed
	if cmd.SKU != "" && cmd.SKU != product.SKU {
		existing, _ := h.productRepo.FindBySKU(ctx, cmd.SKU)
		if existing != nil && existing.ID != product.ID {
			return common.ConflictResponse[responses.ProductResponse]("BUS_007", "El SKU ya existe"), nil
		}
		product.SKU = cmd.SKU
	}

	// Verify category exists if changed
	if cmd.CategoryID != "" {
		exists, _ := h.categoryRepo.Exists(ctx, cmd.CategoryID)
		if !exists {
			return common.NotFoundResponse[responses.ProductResponse]("Categoría no encontrada"), nil
		}
		product.CategoryID = &cmd.CategoryID
	}

	// Verify brand exists if changed
	if cmd.BrandID != "" {
		exists, _ := h.brandRepo.Exists(ctx, cmd.BrandID)
		if !exists {
			return common.NotFoundResponse[responses.ProductResponse]("Marca no encontrada"), nil
		}
		product.BrandID = &cmd.BrandID
	}

	if cmd.Name != "" {
		product.Name = cmd.Name
	}
	if cmd.Description != "" {
		product.Description = cmd.Description
	}
	if cmd.ActiveIngredient != "" {
		product.ActiveIngredient = cmd.ActiveIngredient
	}
	if cmd.Presentation != "" {
		product.Presentation = cmd.Presentation
	}
	if cmd.Concentration != "" {
		product.Concentration = cmd.Concentration
	}
	product.RequiresPrescription = cmd.RequiresPrescription
	product.IsActive = cmd.IsActive
	if cmd.Barcode != "" {
		product.Barcode = cmd.Barcode
	}

	if err := h.productRepo.Update(ctx, product); err != nil {
		h.logger.Error("Error actualizando producto", zap.Error(err))
		return common.InternalServerErrorResponse[responses.ProductResponse]("Error actualizando producto"), nil
	}

	// Invalidate cache
	go func() {
		_ = h.cacheService.Delete(context.Background(), "cache:product:"+product.ID)
		_ = h.cacheService.DeleteByPattern(context.Background(), "cache:search:*")
	}()

	// Fire-and-forget event
	go func() {
		event := events.NewCatalogEvent(events.EventProductUpdated, product.ID)
		if err := h.eventPublisher.Publish(context.Background(), event); err != nil {
			h.logger.Error("Error publicando evento PRODUCT_UPDATED", zap.Error(err))
		}
	}()

	updated, _ := h.productRepo.FindByID(ctx, product.ID)
	if updated != nil {
		product = updated
	}

	h.logger.Info("Producto actualizado exitosamente",
		zap.String("product_id", product.ID),
	)

	return common.OkResponse(responses.ToProductResponse(product)), nil
}

var _ mediator.RequestHandler[commands.UpdateProductCommand, responses.ProductResponse] = (*UpdateProductHandler)(nil)
