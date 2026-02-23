// internal/presentation/dto/requests/update_category_request.go
package requests

// UpdateCategoryRequest DTO de request para actualizar categoría
type UpdateCategoryRequest struct {
	Name         string `json:"name,omitempty"`
	Slug         string `json:"slug,omitempty"`
	Description  string `json:"description,omitempty"`
	ParentID     string `json:"parent_id,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	IsActive     bool   `json:"is_active"`
	DisplayOrder int    `json:"display_order,omitempty"`
}
