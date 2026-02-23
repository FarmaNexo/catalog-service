// internal/application/commands/create_product_command.go
package commands

// CreateProductCommand comando para crear un producto
type CreateProductCommand struct {
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	Description          string `json:"description"`
	ActiveIngredient     string `json:"active_ingredient"`
	Presentation         string `json:"presentation"`
	Concentration        string `json:"concentration"`
	RequiresPrescription bool   `json:"requires_prescription"`
	CategoryID           string `json:"category_id"`
	BrandID              string `json:"brand_id"`
	SKU                  string `json:"sku"`
	Barcode              string `json:"barcode"`
}

func (c CreateProductCommand) GetName() string {
	return "CreateProductCommand"
}
