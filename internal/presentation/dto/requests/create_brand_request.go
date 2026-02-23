// internal/presentation/dto/requests/create_brand_request.go
package requests

// CreateBrandRequest DTO de request para crear marca
type CreateBrandRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug,omitempty"`
	Description string `json:"description,omitempty"`
	LogoURL     string `json:"logo_url,omitempty"`
	Website     string `json:"website,omitempty"`
	Country     string `json:"country,omitempty"`
}
