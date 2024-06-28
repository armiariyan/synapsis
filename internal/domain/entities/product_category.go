package entities

import (
	"time"

	"gorm.io/gorm"
)

type ProductCategory struct {
	ID        uint           `gorm:"column:id" json:"-"`
	UUID      string         `gorm:"column:uuid;type:uuid;uniqueIndex:UQ_product_categories_uuid;not null;default:uuid_generate_v4()" json:"id"`
	Name      string         `gorm:"column:name;type:varchar(150);index:IDX_product_categories_name;not null" json:"name"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;type:timestamp with time zone" json:"deletedAt"`

	// * for relations`
}

func (e *ProductCategory) TableName() string {
	return "product_categories"
}
