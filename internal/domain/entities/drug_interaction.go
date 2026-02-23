// internal/domain/entities/drug_interaction.go
package entities

import "time"

// DrugInteraction representa una interacción medicamentosa entre dos productos
type DrugInteraction struct {
	ID                     string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProductID              string    `gorm:"column:product_id;type:uuid;not null" json:"product_id"`
	InteractsWithProductID string    `gorm:"column:interacts_with_product_id;type:uuid;not null" json:"interacts_with_product_id"`
	Severity               string    `gorm:"column:severity;type:varchar(20);not null" json:"severity"`
	Description            string    `gorm:"column:description;type:text;not null" json:"description"`
	Recommendation         string    `gorm:"column:recommendation;type:text" json:"recommendation"`
	CreatedAt              time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Product              *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	InteractsWithProduct *Product `gorm:"foreignKey:InteractsWithProductID" json:"interacts_with_product,omitempty"`
}

func (DrugInteraction) TableName() string {
	return `"catalog".drug_interactions`
}

// ValidSeverities valores permitidos de severidad
var ValidSeverities = []string{"leve", "moderada", "grave"}

// IsValidSeverity verifica si la severidad es válida
func IsValidSeverity(severity string) bool {
	for _, s := range ValidSeverities {
		if s == severity {
			return true
		}
	}
	return false
}
