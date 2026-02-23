// internal/domain/entities/brand.go
package entities

import (
	"time"
)

// Brand representa una marca/laboratorio farmacéutico
type Brand struct {
	ID          string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;type:varchar(255);not null;uniqueIndex" json:"name"`
	Slug        string    `gorm:"column:slug;type:varchar(255);uniqueIndex" json:"slug"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	LogoURL     string    `gorm:"column:logo_url;type:varchar(500)" json:"logo_url,omitempty"`
	Website     string    `gorm:"column:website;type:varchar(500)" json:"website,omitempty"`
	Country     string    `gorm:"column:country;type:varchar(100)" json:"country,omitempty"`
	IsActive    bool      `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// TableName especifica la tabla en el esquema catalog
func (Brand) TableName() string {
	return `"catalog".brands`
}
