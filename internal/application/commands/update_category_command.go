// internal/application/commands/update_category_command.go
package commands

// UpdateCategoryCommand comando para actualizar una categoría
type UpdateCategoryCommand struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Description  string `json:"description"`
	ParentID     string `json:"parent_id"`
	ImageURL     string `json:"image_url"`
	IsActive     bool   `json:"is_active"`
	DisplayOrder int    `json:"display_order"`
}

func (c UpdateCategoryCommand) GetName() string {
	return "UpdateCategoryCommand"
}
