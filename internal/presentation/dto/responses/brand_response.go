// internal/presentation/dto/responses/brand_response.go
package responses

import (
	"github.com/farmanexo/catalog-service/internal/domain/entities"
)

// BrandResponse DTO de respuesta de marca
type BrandResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url,omitempty"`
	Website     string `json:"website,omitempty"`
	Country     string `json:"country,omitempty"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// BrandListResponse DTO de respuesta de lista de marcas
type BrandListResponse struct {
	Brands []BrandResponse `json:"brands"`
}

// ToBrandResponse convierte una entidad Brand a BrandResponse
func ToBrandResponse(b *entities.Brand) BrandResponse {
	return BrandResponse{
		ID:          b.ID,
		Name:        b.Name,
		Slug:        b.Slug,
		Description: b.Description,
		LogoURL:     b.LogoURL,
		Website:     b.Website,
		Country:     b.Country,
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
