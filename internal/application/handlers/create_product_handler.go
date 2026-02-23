// internal/application/handlers/create_product_handler.go
package handlers

import (
	"context"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/domain/entities"
	"github.com/farmanexo/catalog-service/internal/domain/events"
	"github.com/farmanexo/catalog-service/internal/domain/repositories"
	"github.com/farmanexo/catalog-service/internal/domain/services"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/internal/shared/common"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateProductHandler struct {
	productRepo    repositories.ProductRepository
	categoryRepo   repositories.CategoryRepository
	brandRepo      repositories.BrandRepository
	eventPublisher services.EventPublisher
	cacheService   services.CacheService
	logger         *zap.Logger
}

func NewCreateProductHandler(
	productRepo repositories.ProductRepository,
	categoryRepo repositories.CategoryRepository,
	brandRepo repositories.BrandRepository,
	eventPublisher services.EventPublisher,
	cacheService services.CacheService,
	logger *zap.Logger,
) *CreateProductHandler {
	return &CreateProductHandler{
		productRepo:    productRepo,
		categoryRepo:   categoryRepo,
		brandRepo:      brandRepo,
		eventPublisher: eventPublisher,
		cacheService:   cacheService,
		logger:         logger,
	}
}

func (h *CreateProductHandler) Handle(ctx context.Context, cmd commands.CreateProductCommand) (*common.ApiResponse[responses.ProductResponse], error) {
	slug := cmd.Slug
	if slug == "" {
		slug = GenerateSlug(cmd.Name)
	}

	// Verify slug uniqueness
	existing, _ := h.productRepo.FindBySlug(ctx, slug)
	if existing != nil {
		slug = slug + "-" + uuid.New().String()[:8]
	}

	// Verify SKU uniqueness if provided
	if cmd.SKU != "" {
		existingSKU, _ := h.productRepo.FindBySKU(ctx, cmd.SKU)
		if existingSKU != nil {
			return common.ConflictResponse[responses.ProductResponse]("BUS_007", "El SKU ya existe"), nil
		}
	}

	// Verify category exists if provided
	if cmd.CategoryID != "" {
		exists, _ := h.categoryRepo.Exists(ctx, cmd.CategoryID)
		if !exists {
			return common.NotFoundResponse[responses.ProductResponse]("Categoría no encontrada"), nil
		}
	}

	// Verify brand exists if provided
	if cmd.BrandID != "" {
		exists, _ := h.brandRepo.Exists(ctx, cmd.BrandID)
		if !exists {
			return common.NotFoundResponse[responses.ProductResponse]("Marca no encontrada"), nil
		}
	}

	product := &entities.Product{
		ID:                   uuid.New().String(),
		Name:                 cmd.Name,
		Slug:                 slug,
		Description:          cmd.Description,
		ActiveIngredient:     cmd.ActiveIngredient,
		Presentation:         cmd.Presentation,
		Concentration:        cmd.Concentration,
		RequiresPrescription: cmd.RequiresPrescription,
		SKU:                  cmd.SKU,
		Barcode:              cmd.Barcode,
		IsActive:             true,
	}

	if cmd.CategoryID != "" {
		product.CategoryID = &cmd.CategoryID
	}
	if cmd.BrandID != "" {
		product.BrandID = &cmd.BrandID
	}

	if err := h.productRepo.Create(ctx, product); err != nil {
		h.logger.Error("Error creando producto", zap.Error(err))
		return common.InternalServerErrorResponse[responses.ProductResponse]("Error creando producto"), nil
	}

	// Invalidate search cache
	go func() {
		_ = h.cacheService.DeleteByPattern(context.Background(), "cache:search:*")
	}()

	// Fire-and-forget event
	go func() {
		event := events.NewCatalogEvent(events.EventProductCreated, product.ID)
		if err := h.eventPublisher.Publish(context.Background(), event); err != nil {
			h.logger.Error("Error publicando evento PRODUCT_CREATED", zap.Error(err))
		}
	}()

	// Reload with relations
	created, err := h.productRepo.FindByID(ctx, product.ID)
	if err != nil {
		return common.CreatedResponse(responses.ToProductResponse(product)), nil
	}

	h.logger.Info("Producto creado exitosamente",
		zap.String("product_id", product.ID),
		zap.String("name", product.Name),
	)

	return common.CreatedResponse(responses.ToProductResponse(created)), nil
}

var _ mediator.RequestHandler[commands.CreateProductCommand, responses.ProductResponse] = (*CreateProductHandler)(nil)
