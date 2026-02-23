// internal/application/preprocessors/sanitize_input_preprocessor.go
package preprocessors

import (
	"context"
	"strings"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"go.uber.org/zap"
)

// SanitizeInputPreProcessor limpia y normaliza los inputs
type SanitizeInputPreProcessor struct {
	logger *zap.Logger
}

func NewSanitizeInputPreProcessor(logger *zap.Logger) *SanitizeInputPreProcessor {
	return &SanitizeInputPreProcessor{logger: logger}
}

func (p *SanitizeInputPreProcessor) Process(ctx context.Context, request interface{}) error {
	switch cmd := request.(type) {
	case *commands.CreateProductCommand:
		cmd.Name = strings.TrimSpace(cmd.Name)
		cmd.Description = strings.TrimSpace(cmd.Description)
		cmd.ActiveIngredient = strings.TrimSpace(cmd.ActiveIngredient)
		cmd.Presentation = strings.TrimSpace(cmd.Presentation)
		cmd.Concentration = strings.TrimSpace(cmd.Concentration)
		cmd.SKU = strings.TrimSpace(cmd.SKU)
		cmd.Barcode = strings.TrimSpace(cmd.Barcode)
		p.logger.Debug("Input sanitizado", zap.String("command", "CreateProductCommand"))

	case *commands.UpdateProductCommand:
		cmd.Name = strings.TrimSpace(cmd.Name)
		cmd.Description = strings.TrimSpace(cmd.Description)
		cmd.ActiveIngredient = strings.TrimSpace(cmd.ActiveIngredient)
		cmd.Presentation = strings.TrimSpace(cmd.Presentation)
		cmd.Concentration = strings.TrimSpace(cmd.Concentration)
		cmd.SKU = strings.TrimSpace(cmd.SKU)
		cmd.Barcode = strings.TrimSpace(cmd.Barcode)
		p.logger.Debug("Input sanitizado", zap.String("command", "UpdateProductCommand"))

	case *commands.CreateCategoryCommand:
		cmd.Name = strings.TrimSpace(cmd.Name)
		cmd.Description = strings.TrimSpace(cmd.Description)
		p.logger.Debug("Input sanitizado", zap.String("command", "CreateCategoryCommand"))

	case *commands.UpdateCategoryCommand:
		cmd.Name = strings.TrimSpace(cmd.Name)
		cmd.Description = strings.TrimSpace(cmd.Description)
		p.logger.Debug("Input sanitizado", zap.String("command", "UpdateCategoryCommand"))

	case *commands.CreateBrandCommand:
		cmd.Name = strings.TrimSpace(cmd.Name)
		cmd.Description = strings.TrimSpace(cmd.Description)
		cmd.Website = strings.TrimSpace(cmd.Website)
		cmd.Country = strings.TrimSpace(cmd.Country)
		p.logger.Debug("Input sanitizado", zap.String("command", "CreateBrandCommand"))

	case *commands.UpdateBrandCommand:
		cmd.Name = strings.TrimSpace(cmd.Name)
		cmd.Description = strings.TrimSpace(cmd.Description)
		cmd.Website = strings.TrimSpace(cmd.Website)
		cmd.Country = strings.TrimSpace(cmd.Country)
		p.logger.Debug("Input sanitizado", zap.String("command", "UpdateBrandCommand"))
	}

	return nil
}
