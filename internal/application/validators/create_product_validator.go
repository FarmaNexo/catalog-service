// internal/application/validators/create_product_validator.go
package validators

import (
	"context"
	"fmt"
	"strings"

	"github.com/farmanexo/catalog-service/internal/application/commands"
	"github.com/farmanexo/catalog-service/internal/presentation/dto/responses"
	"github.com/farmanexo/catalog-service/pkg/mediator"
)

// CreateProductValidator valida el comando de creación de producto
type CreateProductValidator struct{}

func NewCreateProductValidator() *CreateProductValidator {
	return &CreateProductValidator{}
}

func (v *CreateProductValidator) Validate(ctx context.Context, cmd commands.CreateProductCommand) error {
	var errors []string

	if strings.TrimSpace(cmd.Name) == "" {
		errors = append(errors, "El nombre del producto es requerido")
	}
	if len(cmd.Name) > 500 {
		errors = append(errors, "El nombre no puede exceder 500 caracteres")
	}

	if len(errors) > 0 {
		return fmt.Errorf("errores de validación: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Compile-time interface check
var _ mediator.Validator[commands.CreateProductCommand, responses.ProductResponse] = (*CreateProductValidator)(nil)
