// internal/presentation/dto/responses/fbt_response.go
package responses

import "github.com/farmanexo/catalog-service/internal/domain/entities"

// FBTItemResponse un producto frecuentemente comprado junto
type FBTItemResponse struct {
	ProductID       string  `json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductSlug     string  `json:"product_slug"`
	BrandName       string  `json:"brand_name,omitempty"`
	ImageURL        string  `json:"image_url,omitempty"`
	Score           float64 `json:"score"`
	CoPurchaseCount int     `json:"co_purchase_count"`
}

// FBTListResponse lista de productos frecuentemente comprados juntos
type FBTListResponse struct {
	ProductID string            `json:"product_id"`
	Items     []FBTItemResponse `json:"items"`
	Total     int               `json:"total"`
}

// ToFBTItemResponse convierte una entidad FBT a DTO
func ToFBTItemResponse(fbt *entities.FrequentlyBoughtTogether) FBTItemResponse {
	resp := FBTItemResponse{
		ProductID:       fbt.RelatedProductID,
		Score:           fbt.Score,
		CoPurchaseCount: fbt.CoPurchaseCount,
	}

	if fbt.RelatedProduct != nil {
		resp.ProductName = fbt.RelatedProduct.Name
		resp.ProductSlug = fbt.RelatedProduct.Slug
		if fbt.RelatedProduct.Brand != nil {
			resp.BrandName = fbt.RelatedProduct.Brand.Name
		}
		if len(fbt.RelatedProduct.Images) > 0 {
			resp.ImageURL = fbt.RelatedProduct.Images[0].ImageURL
		}
	}

	return resp
}
