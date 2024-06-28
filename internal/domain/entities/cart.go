package entities

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	ID        uint           `gorm:"column:id" json:"-"`
	UUID      string         `gorm:"column:uuid;type:uuid;uniqueIndex:UQ_carts_uuid;not null;default:uuid_generate_v4()" json:"id"`
	UserID    string         `gorm:"column:user_id;not null" json:"user_id"`
	ProductID string         `gorm:"column:product_id;not null" json:"product_id"`
	Amount    float64        `gorm:"column:amount;not null" json:"amount"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index;type:timestamp with time zone" json:"deletedAt"`

	// * for relations
	Product Product `gorm:"foreignKey:ProductID;references:UUID" json:"product"`
	User    User    `gorm:"foreignKey:UserID;references:UUID" json:"user"`
}

func (e *Cart) TableName() string {
	return "carts"
}
