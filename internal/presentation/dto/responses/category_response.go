// internal/presentation/dto/responses/category_response.go
package responses

import (
	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// CategoryResponse DTO de respuesta de categoría
type CategoryResponse struct {
	ID           string             `json:"id"`
	Name         string             `json:"name"`
	Slug         string             `json:"slug"`
	Description  string             `json:"description"`
	ParentID     string             `json:"parent_id,omitempty"`
	ImageURL     string             `json:"image_url,omitempty"`
	IsActive     bool               `json:"is_active"`
	DisplayOrder int                `json:"display_order"`
	Children     []CategoryResponse `json:"children,omitempty"`
	CreatedAt    string             `json:"created_at"`
	UpdatedAt    string             `json:"updated_at"`
}

// CategoryListResponse DTO de respuesta de lista de categorías
type CategoryListResponse struct {
	Categories []CategoryResponse `json:"categories"`
}

// ToCategoryResponse convierte una entidad Category a CategoryResponse
func ToCategoryResponse(c *entities.Category) CategoryResponse {
	resp := CategoryResponse{
		ID:           c.ID,
		Name:         c.Name,
		Slug:         c.Slug,
		Description:  c.Description,
		ImageURL:     c.ImageURL,
		IsActive:     c.IsActive,
		DisplayOrder: c.DisplayOrder,
		CreatedAt:    c.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    c.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if c.ParentID != nil {
		resp.ParentID = *c.ParentID
	}

	if c.Children != nil {
		resp.Children = make([]CategoryResponse, len(c.Children))
		for i, child := range c.Children {
			resp.Children[i] = ToCategoryResponse(&child)
		}
	}

	return resp
}
