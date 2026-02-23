// internal/application/commands/create_brand_command.go
package commands

// CreateBrandCommand comando para crear una marca
type CreateBrandCommand struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	Website     string `json:"website"`
	Country     string `json:"country"`
}

func (c CreateBrandCommand) GetName() string {
	return "CreateBrandCommand"
}
