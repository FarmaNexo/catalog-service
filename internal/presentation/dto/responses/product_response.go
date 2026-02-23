// internal/presentation/dto/responses/product_response.go
package responses

import (
	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// ProductResponse DTO de respuesta de producto
type ProductResponse struct {
	ID                   string          `json:"id"`
	Name                 string          `json:"name"`
	Slug                 string          `json:"slug"`
	Description          string          `json:"description"`
	ActiveIngredient     string          `json:"active_ingredient"`
	Presentation         string          `json:"presentation"`
	Concentration        string          `json:"concentration"`
	RequiresPrescription bool            `json:"requires_prescription"`
	CategoryID           string          `json:"category_id,omitempty"`
	CategoryName         string          `json:"category_name,omitempty"`
	BrandID              string          `json:"brand_id,omitempty"`
	BrandName            string          `json:"brand_name,omitempty"`
	SKU                  string          `json:"sku,omitempty"`
	Barcode              string          `json:"barcode,omitempty"`
	IsActive             bool            `json:"is_active"`
	Images               []ImageResponse `json:"images,omitempty"`
	CreatedAt            string          `json:"created_at"`
	UpdatedAt            string          `json:"updated_at"`
}

// ProductListResponse DTO de respuesta de lista paginada de productos
type ProductListResponse struct {
	Products   []ProductResponse `json:"products"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
	TotalPages int               `json:"total_pages"`
}

// ToProductResponse convierte una entidad Product a ProductResponse
func ToProductResponse(p *entities.Product) ProductResponse {
	resp := ProductResponse{
		ID:                   p.ID,
		Name:                 p.Name,
		Slug:                 p.Slug,
		Description:          p.Description,
		ActiveIngredient:     p.ActiveIngredient,
		Presentation:         p.Presentation,
		Concentration:        p.Concentration,
		RequiresPrescription: p.RequiresPrescription,
		SKU:                  p.SKU,
		Barcode:              p.Barcode,
		IsActive:             p.IsActive,
		CreatedAt:            p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:            p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if p.CategoryID != nil {
		resp.CategoryID = *p.CategoryID
	}
	if p.Category != nil {
		resp.CategoryName = p.Category.Name
	}
	if p.BrandID != nil {
		resp.BrandID = *p.BrandID
	}
	if p.Brand != nil {
		resp.BrandName = p.Brand.Name
	}

	if p.Images != nil {
		resp.Images = make([]ImageResponse, len(p.Images))
		for i, img := range p.Images {
			resp.Images[i] = ImageResponse{
				ID:           img.ID,
				ImageURL:     img.ImageURL,
				IsPrimary:    img.IsPrimary,
				DisplayOrder: img.DisplayOrder,
			}
		}
	}

	return resp
}
