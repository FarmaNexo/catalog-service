// internal/presentation/dto/requests/update_product_request.go
package requests

// UpdateProductRequest DTO de request para actualizar producto
type UpdateProductRequest struct {
	Name                 string `json:"name,omitempty"`
	Slug                 string `json:"slug,omitempty"`
	Description          string `json:"description,omitempty"`
	ActiveIngredient     string `json:"active_ingredient,omitempty"`
	Presentation         string `json:"presentation,omitempty"`
	Concentration        string `json:"concentration,omitempty"`
	RequiresPrescription bool   `json:"requires_prescription"`
	CategoryID           string `json:"category_id,omitempty"`
	BrandID              string `json:"brand_id,omitempty"`
	SKU                  string `json:"sku,omitempty"`
	Barcode              string `json:"barcode,omitempty"`
	IsActive             bool   `json:"is_active"`
}
