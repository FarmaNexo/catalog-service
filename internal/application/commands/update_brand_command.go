// internal/application/commands/update_brand_command.go
package commands

// UpdateBrandCommand comando para actualizar una marca
type UpdateBrandCommand struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	Website     string `json:"website"`
	Country     string `json:"country"`
	IsActive    bool   `json:"is_active"`
}

func (c UpdateBrandCommand) GetName() string {
	return "UpdateBrandCommand"
}
