// internal/presentation/dto/requests/create_product_request.go
package requests

// CreateProductRequest DTO de request para crear producto
type CreateProductRequest struct {
	Name                 string `json:"name"`
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
}
