// internal/application/validators/create_category_validator.go
package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/pkg/mediator"
)

// CreateCategoryValidator valida el comando de creación de categoría
type CreateCategoryValidator struct{}

func NewCreateCategoryValidator() *CreateCategoryValidator {
	return &CreateCategoryValidator{}
}

func (v *CreateCategoryValidator) Validate(ctx context.Context, cmd commands.CreateCategoryCommand) error {
	var errors []string

	if strings.TrimSpace(cmd.Name) == "" {
		errors = append(errors, "El nombre de la categoría es requerido")
	}
	if len(cmd.Name) > 255 {
		errors = append(errors, "El nombre no puede exceder 255 caracteres")
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Compile-time interface check
var _ mediator.Validator[commands.CreateCategoryCommand, responses.CategoryResponse] = (*CreateCategoryValidator)(nil)
