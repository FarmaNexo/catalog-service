// internal/application/handlers/delete_product_handler.go
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

type DeleteProductHandler struct {
	productRepo    repositories.ProductRepository
	eventPublisher services.EventPublisher
	cacheService   services.CacheService
	logger         *zap.Logger
}

func NewDeleteProductHandler(
	productRepo repositories.ProductRepository,
	eventPublisher services.EventPublisher,
	cacheService services.CacheService,
	logger *zap.Logger,
) *DeleteProductHandler {
	return &DeleteProductHandler{
		productRepo:    productRepo,
		eventPublisher: eventPublisher,
		cacheService:   cacheService,
		logger:         logger,
	}
}

func (h *DeleteProductHandler) Handle(ctx context.Context, cmd commands.DeleteProductCommand) (*common.ApiResponse[responses.EmptyResponse], error) {
	product, err := h.productRepo.FindByID(ctx, cmd.ID)
	if err != nil {
		return common.NotFoundResponse[responses.EmptyResponse]("Producto no encontrado"), nil
	}

	if err := h.productRepo.SoftDelete(ctx, product.ID); err != nil {
		h.logger.Error("Error eliminando producto", zap.Error(err))
		return common.InternalServerErrorResponse[responses.EmptyResponse]("Error eliminando producto"), nil
	}

	// Invalidate cache
	go func() {
		_ = h.cacheService.Delete(context.Background(), "cache:product:"+product.ID)
		_ = h.cacheService.DeleteByPattern(context.Background(), "cache:search:*")
	}()

	// Fire-and-forget event
	go func() {
		event := events.NewCatalogEvent(events.EventProductDeleted, product.ID)
		if err := h.eventPublisher.Publish(context.Background(), event); err != nil {
			h.logger.Error("Error publicando evento PRODUCT_DELETED", zap.Error(err))
		}
	}()

	h.logger.Info("Producto eliminado (soft delete)",
		zap.String("product_id", product.ID),
	)

	resp := common.NewApiResponse[responses.EmptyResponse]()
	resp.SetHttpStatus(200)
	resp.SetData(responses.EmptyResponse{})
	resp.AddMessageWithType("SUCCESS_004", "Producto eliminado exitosamente", "SUCCESS")
	return resp, nil
}

var _ mediator.RequestHandler[commands.DeleteProductCommand, responses.EmptyResponse] = (*DeleteProductHandler)(nil)
