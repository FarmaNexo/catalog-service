// internal/application/postprocessors/log_audit_postprocessor.go
package postprocessors

import (
	"context"
	"time"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/application/queries"
	"github.com/farmanexo/catalog-service/pkg/mediator"
	"go.uber.org/zap"
)

// LogAuditPostProcessor registra eventos de auditoría
type LogAuditPostProcessor struct {
	logger *zap.Logger
}

func NewLogAuditPostProcessor(logger *zap.Logger) *LogAuditPostProcessor {
	return &LogAuditPostProcessor{logger: logger}
}

func (p *LogAuditPostProcessor) Process(ctx context.Context, request interface{}, response interface{}) error {
	userID := p.getUserIDFromContext(ctx)
	correlationID := mediator.GetCorrelationID(ctx)
	isSuccess := p.checkSuccess(response)

	switch request.(type) {
	case commands.CreateProductCommand, *commands.CreateProductCommand:
		p.logAudit("PRODUCT_CREATED", userID, correlationID, isSuccess)
	case commands.UpdateProductCommand, *commands.UpdateProductCommand:
		p.logAudit("PRODUCT_UPDATED", userID, correlationID, isSuccess)
	case commands.DeleteProductCommand, *commands.DeleteProductCommand:
		p.logAudit("PRODUCT_DELETED", userID, correlationID, isSuccess)
	case commands.UploadProductImagesCommand, *commands.UploadProductImagesCommand:
		p.logAudit("PRODUCT_IMAGES_UPLOADED", userID, correlationID, isSuccess)
	case commands.CreateCategoryCommand, *commands.CreateCategoryCommand:
		p.logAudit("CATEGORY_CREATED", userID, correlationID, isSuccess)
	case commands.UpdateCategoryCommand, *commands.UpdateCategoryCommand:
		p.logAudit("CATEGORY_UPDATED", userID, correlationID, isSuccess)
	case commands.CreateBrandCommand, *commands.CreateBrandCommand:
		p.logAudit("BRAND_CREATED", userID, correlationID, isSuccess)
	case commands.UpdateBrandCommand, *commands.UpdateBrandCommand:
		p.logAudit("BRAND_UPDATED", userID, correlationID, isSuccess)
	case queries.SearchProductsQuery, *queries.SearchProductsQuery:
		p.logAudit("PRODUCTS_SEARCHED", userID, correlationID, isSuccess)
	default:
		p.logger.Debug("Post-processor: comando sin auditoría configurada")
	}

	return nil
}

func (p *LogAuditPostProcessor) logAudit(eventType, userID, correlationID string, success bool) {
	p.logger.Info("AUDIT",
		zap.String("event_type", eventType),
		zap.Bool("success", success),
		zap.String("correlation_id", correlationID),
		zap.String("user_id", userID),
		zap.Time("timestamp", time.Now()),
	)
}

func (p *LogAuditPostProcessor) checkSuccess(response interface{}) bool {
	if resp, ok := response.(interface{ IsValid() bool }); ok {
		return resp.IsValid()
	}
	return false
}

func (p *LogAuditPostProcessor) getUserIDFromContext(ctx context.Context) string {
	userID, _ := mediator.GetUserID(ctx)
	if userID == "" {
		return "ANONYMOUS"
	}
	return userID
}
