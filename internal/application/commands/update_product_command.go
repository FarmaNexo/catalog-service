// internal/application/commands/update_product_command.go
package commands

// UpdateProductCommand comando para actualizar un producto
type UpdateProductCommand struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Slug                 string `json:"slug"`
	Description          string `json:"description"`
	ActiveIngredient     string `json:"active_ingredient"`
	Presentation         string `json:"presentation"`
	Concentration        string `json:"concentration"`
	Form                 string `json:"form"`
	RegistryNumber       string `json:"registry_number"`
	Manufacturer         string `json:"manufacturer"`
	SourceProductCode    *int   `json:"source_product_code"`
	RequiresPrescription bool   `json:"requires_prescription"`
	CategoryID           string `json:"category_id"`
	BrandID              string `json:"brand_id"`
	SKU                  string `json:"sku"`
	Barcode              string `json:"barcode"`
	IsActive             bool   `json:"is_active"`
}

func (c UpdateProductCommand) GetName() string {
	return "UpdateProductCommand"
}
