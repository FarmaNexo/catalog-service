// internal/application/commands/create_drug_interaction_command.go
package commands

// CreateDrugInteractionCommand comando para crear una interacción medicamentosa
type CreateDrugInteractionCommand struct {
	ProductID              string `json:"product_id"`
	InteractsWithProductID string `json:"interacts_with_product_id"`
	Severity               string `json:"severity"`
	Description            string `json:"description"`
	Recommendation         string `json:"recommendation"`
}

func (c CreateDrugInteractionCommand) GetName() string {
	return "CreateDrugInteractionCommand"
}
