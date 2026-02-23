// internal/presentation/dto/requests/update_brand_request.go
package requests

// UpdateBrandRequest DTO de request para actualizar marca
type UpdateBrandRequest struct {
	Name        string `json:"name,omitempty"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
	Website     string `json:"website,omitempty"`
	Country     string `json:"country,omitempty"`
	IsActive    bool   `json:"is_active"`
}
