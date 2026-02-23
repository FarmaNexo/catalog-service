// internal/application/commands/create_category_command.go
package commands

// CreateCategoryCommand comando para crear una categoría
type CreateCategoryCommand struct {
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	ParentID     string `json:"parent_id"`
	ImageURL     string `json:"image_url"`
	DisplayOrder int    `json:"display_order"`
}

func (c CreateCategoryCommand) GetName() string {
	return "CreateCategoryCommand"
}
