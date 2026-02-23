// internal/presentation/dto/requests/create_category_request.go
package requests

// CreateCategoryRequest DTO de request para crear categoría
type CreateCategoryRequest struct {
	Name         string `json:"name"`
	Slug         string `json:"slug,omitempty"`
	Description  string `json:"description,omitempty"`
	ParentID     string `json:"parent_id,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	DisplayOrder int    `json:"display_order,omitempty"`
}
