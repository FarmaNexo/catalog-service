// internal/domain/entities/frequently_bought_together.go
package entities

import "time"

// FrequentlyBoughtTogether relación de productos comprados juntos
type FrequentlyBoughtTogether struct {
	ID               string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
	ProductID        string    `gorm:"column:product_id;type:uuid;not null" json:"product_id"`
	RelatedProductID string    `gorm:"column:related_product_id;type:uuid;not null" json:"related_product_id"`
	Score            float64   `gorm:"column:score;type:decimal(5,2);not null;default:0" json:"score"`
	CoPurchaseCount  int       `gorm:"column:co_purchase_count;not null;default:0" json:"co_purchase_count"`
	CreatedAt        time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`

	Product        *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	RelatedProduct *Product `gorm:"foreignKey:RelatedProductID" json:"related_product,omitempty"`
}

func (FrequentlyBoughtTogether) TableName() string {
	return `"catalog".frequently_bought_together`
}
