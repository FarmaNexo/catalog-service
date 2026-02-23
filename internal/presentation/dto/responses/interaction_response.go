// internal/presentation/dto/responses/interaction_response.go
package responses

import "github.com/farmanexo/catalog-service/internal/domain/entities"

// InteractionResponse DTO de respuesta para una interacción medicamentosa
type InteractionResponse struct {
	ID                       string `json:"id"`
	ProductID                string `json:"product_id"`
	ProductName              string `json:"product_name"`
	InteractsWithProductID   string `json:"interacts_with_product_id"`
	InteractsWithProductName string `json:"interacts_with_product_name"`
	Severity                 string `json:"severity"`
	Description              string `json:"description"`
	Recommendation           string `json:"recommendation,omitempty"`
	CreatedAt                string `json:"created_at"`
}

// InteractionListResponse lista de interacciones
type InteractionListResponse struct {
	Interactions []InteractionResponse `json:"interactions"`
	Total        int                   `json:"total"`
}

// ToInteractionResponse convierte una entidad a DTO de respuesta
func ToInteractionResponse(i *entities.DrugInteraction) InteractionResponse {
	resp := InteractionResponse{
		ID:                     i.ID,
		ProductID:              i.ProductID,
		InteractsWithProductID: i.InteractsWithProductID,
		Severity:               i.Severity,
		Description:            i.Description,
		Recommendation:         i.Recommendation,
		CreatedAt:              i.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if i.Product != nil {
		resp.ProductName = i.Product.Name
	}
	if i.InteractsWithProduct != nil {
		resp.InteractsWithProductName = i.InteractsWithProduct.Name
	}

	return resp
}
