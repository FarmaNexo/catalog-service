// internal/application/commands/delete_product_command.go
package commands

// DeleteProductCommand comando para soft-delete un producto
type DeleteProductCommand struct {
	ID string `json:"id"`
}

func (c DeleteProductCommand) GetName() string {
	return "DeleteProductCommand"
}
