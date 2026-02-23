// internal/domain/entities/product_image.go
package entities

import (
	"time"
)

// ProductImage representa una imagen de un producto
type ProductImage struct {
	ID           string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProductID    string    `gorm:"column:product_id;type:uuid;not null" json:"product_id"`
	ImageURL     string    `gorm:"column:image_url;type:varchar(500);not null" json:"image_url"`
	IsPrimary    bool      `gorm:"column:is_primary;not null;default:false" json:"is_primary"`
	DisplayOrder int       `gorm:"column:display_order;not null;default:0" json:"display_order"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

// TableName especifica la tabla en el esquema catalog
func (ProductImage) TableName() string {
	return `"catalog".product_images`
}
