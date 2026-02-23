// internal/presentation/dto/requests/interaction_request.go
package requests

// CreateInteractionRequest DTO para crear una interacción medicamentosa
type CreateInteractionRequest struct {
	ProductID              string `json:"product_id"`
	InteractsWithProductID string `json:"interacts_with_product_id"`
	Severity               string `json:"severity"`
	Description            string `json:"description"`
	Recommendation         string `json:"recommendation,omitempty"`
}
