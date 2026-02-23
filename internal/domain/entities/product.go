// internal/domain/entities/product.go
package entities

import (
	"time"
)

// Product representa un producto farmacéutico
type Product struct {
	ID                   string     `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name                 string     `gorm:"column:name;type:varchar(500);not null" json:"name"`
	Slug                 string     `gorm:"column:slug;type:varchar(500);uniqueIndex" json:"slug"`
	Description          string     `gorm:"column:description;type:text" json:"description"`
	ActiveIngredient     string     `gorm:"column:active_ingredient;type:varchar(500)" json:"active_ingredient"`
	Presentation         string     `gorm:"column:presentation;type:varchar(255)" json:"presentation"`
	Concentration        string     `gorm:"column:concentration;type:varchar(100)" json:"concentration"`
	RequiresPrescription bool       `gorm:"column:requires_prescription;not null;default:false" json:"requires_prescription"`
	CategoryID           *string    `gorm:"column:category_id;type:uuid" json:"category_id"`
	BrandID              *string    `gorm:"column:brand_id;type:uuid" json:"brand_id"`
	SKU                  string     `gorm:"column:sku;type:varchar(100);uniqueIndex" json:"sku"`
	Barcode              string     `gorm:"column:barcode;type:varchar(100)" json:"barcode"`
	IsActive             bool       `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt            time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt            time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt            *time.Time `gorm:"column:deleted_at" json:"deleted_at,omitempty"`

	// Relaciones
	Category *Category      `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Brand    *Brand         `gorm:"foreignKey:BrandID" json:"brand,omitempty"`
	Images   []ProductImage `gorm:"foreignKey:ProductID" json:"images,omitempty"`
}

// TableName especifica la tabla en el esquema catalog
func (Product) TableName() string {
	return `"catalog".products`
}

// IsDeleted verifica si el producto está eliminado (soft delete)
func (p *Product) IsDeleted() bool {
	return p.DeletedAt != nil
}

// SoftDelete marca el producto como eliminado
func (p *Product) SoftDelete() {
	now := time.Now()
	p.DeletedAt = &now
	p.IsActive = false
}
