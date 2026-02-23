// internal/domain/entities/category.go
package entities

import (
	"time"
)

// Category representa una categoría de productos farmacéuticos
type Category struct {
	ID           string     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name         string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Slug         string     `gorm:"column:slug;type:varchar(255);uniqueIndex" json:"slug"`
	Description  string     `gorm:"column:description;type:text" json:"description"`
	ParentID     *string    `gorm:"column:parent_id;type:uuid" json:"parent_id,omitempty"`
	ImageURL     string     `gorm:"column:image_url;type:varchar(500)" json:"image_url,omitempty"`
	IsActive     bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	DisplayOrder int        `gorm:"column:display_order;not null;default:0" json:"display_order"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	// Relaciones
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// TableName especifica la tabla en el esquema catalog
func (Category) TableName() string {
	return `"catalog".categories`
}
