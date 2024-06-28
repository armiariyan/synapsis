package entities

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint           `gorm:"column:id" json:"-"`
	UUID       string         `gorm:"column:uuid;type:uuid;uniqueIndex:UQ_products_uuid;not null;default:uuid_generate_v4()" json:"id"`
	Name       string         `gorm:"column:name;type:varchar(150);index:IDX_products_name;not null" json:"name"`
	Price      float64        `gorm:"column:price;type:decimal;not null" json:"price"`
	CategoryID string         `gorm:"column:category_id;not null" json:"category_id"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"createdAt"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp" json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at;index;type:timestamp with time zone" json:"deletedAt"`

	// * for relations
	ProductCategory ProductCategory `gorm:"foreignkey:CategoryID;references:UUID" json:"productCategory"`
}

func (e *Product) TableName() string {
	return "products"
}
