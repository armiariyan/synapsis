package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID          uint           `gorm:"column:id" json:"-"`
	UUID        string         `gorm:"column:uuid;type:uuid;uniqueIndex:UQ_users_uuid;not null;default:uuid_generate_v4()" json:"id"`
	Name        string         `gorm:"column:name;type:varchar(150);index:IDX_users_name;not null" json:"name"`
	Email       string         `gorm:"column:email;type:varchar(64);uniqueIndex:UQ_users_email;not null" json:"email"`
	PhoneNumber string         `gorm:"column:phone_number;type:varchar(15);default:null" json:"phoneNumber"`
	Password    string         `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Token       string         `gorm:"column:token;type:text;default:null" json:"-"`
	CreatedAt   time.Time      `gorm:"column:created_at;type:timestamp with time zone;default:current_timestamp" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;type:timestamp with time zone;default:current_timestamp" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index;type:timestamp with time zone" json:"deletedAt"`

	// * for relations
}

func (e *User) TableName() string {
	return "users"
}
