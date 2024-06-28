package repositories

import (
	"context"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"gorm.io/gorm"
)

type ProductCategory interface {
	FindByName(ctx context.Context, name string) (result entities.ProductCategory, err error)
}

type productCategory struct {
	db *gorm.DB
}

func NewProductCategory(db *gorm.DB) ProductCategory {
	if db == nil {
		panic("db is nil")
	}

	return &productCategory{db}
}

func (r *productCategory) FindByName(ctx context.Context, name string) (result entities.ProductCategory, err error) {
	err = r.db.WithContext(ctx).Where("name = ?", name).First(&result).Error
	if err == gorm.ErrRecordNotFound { // * skip err record not found, checking data empty in service_impl
		err = nil
	}
	return
}
