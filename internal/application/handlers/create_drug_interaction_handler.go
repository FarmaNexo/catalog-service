// internal/application/handlers/create_drug_interaction_handler.go
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

type CreateDrugInteractionHandler struct {
	interactionRepo repositories.DrugInteractionRepository
	productRepo     repositories.ProductRepository
	cacheService    services.CacheService
	logger          *zap.Logger
}

func NewCreateDrugInteractionHandler(
	interactionRepo repositories.DrugInteractionRepository,
	productRepo repositories.ProductRepository,
	cacheService services.CacheService,
	logger *zap.Logger,
) *CreateDrugInteractionHandler {
	return &CreateDrugInteractionHandler{
		interactionRepo: interactionRepo,
		productRepo:     productRepo,
		cacheService:    cacheService,
		logger:          logger,
	}
}

func (h *CreateDrugInteractionHandler) Handle(ctx context.Context, cmd commands.CreateDrugInteractionCommand) (*common.ApiResponse[responses.InteractionResponse], error) {
	// Validate severity
	if !entities.IsValidSeverity(cmd.Severity) {
		return common.BadRequestResponse[responses.InteractionResponse]("VAL_001", "Severidad inválida. Valores permitidos: leve, moderada, grave"), nil
	}

	if cmd.ProductID == "" || cmd.InteractsWithProductID == "" {
		return common.BadRequestResponse[responses.InteractionResponse]("VAL_006", "product_id e interacts_with_product_id son requeridos"), nil
	}

	if cmd.ProductID == cmd.InteractsWithProductID {
		return common.BadRequestResponse[responses.InteractionResponse]("VAL_001", "Un producto no puede interactuar consigo mismo"), nil
	}

	if cmd.Description == "" {
		return common.BadRequestResponse[responses.InteractionResponse]("VAL_006", "La descripción es requerida"), nil
	}

	// Verify both products exist
	_, err := h.productRepo.FindByID(ctx, cmd.ProductID)
	if err != nil {
		return common.NotFoundResponse[responses.InteractionResponse]("Producto origen no encontrado"), nil
	}

	_, err = h.productRepo.FindByID(ctx, cmd.InteractsWithProductID)
	if err != nil {
		return common.NotFoundResponse[responses.InteractionResponse]("Producto destino no encontrado"), nil
	}

	// Check if interaction already exists
	exists, _ := h.interactionRepo.Exists(ctx, cmd.ProductID, cmd.InteractsWithProductID)
	if exists {
		return common.ConflictResponse[responses.InteractionResponse]("BUS_008", "La interacción ya existe entre estos productos"), nil
	}

	interaction := &entities.DrugInteraction{
		ID:                     uuid.New().String(),
		ProductID:              cmd.ProductID,
		InteractsWithProductID: cmd.InteractsWithProductID,
		Severity:               cmd.Severity,
		Description:            cmd.Description,
		Recommendation:         cmd.Recommendation,
	}

	if err := h.interactionRepo.Create(ctx, interaction); err != nil {
		h.logger.Error("Error creando interacción", zap.Error(err))
		return common.InternalServerErrorResponse[responses.InteractionResponse]("Error creando interacción"), nil
	}

	// Invalidate cache for both products
	go func() {
		_ = h.cacheService.Delete(context.Background(), "cache:product:"+cmd.ProductID+":interactions")
		_ = h.cacheService.Delete(context.Background(), "cache:product:"+cmd.InteractsWithProductID+":interactions")
	}()

	// Reload with relations
	created, err := h.interactionRepo.FindByID(ctx, interaction.ID)
	if err != nil {
		return common.CreatedResponse(responses.ToInteractionResponse(interaction)), nil
	}

	h.logger.Info("Interacción creada exitosamente",
		zap.String("interaction_id", interaction.ID),
		zap.String("product_id", cmd.ProductID),
		zap.String("interacts_with", cmd.InteractsWithProductID),
		zap.String("severity", cmd.Severity),
	)

	return common.CreatedResponse(responses.ToInteractionResponse(created)), nil
}

var _ mediator.RequestHandler[commands.CreateDrugInteractionCommand, responses.InteractionResponse] = (*CreateDrugInteractionHandler)(nil)
