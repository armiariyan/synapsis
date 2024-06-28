package repositories

import (
	"context"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"github.com/armiariyan/synapsis/internal/pkg/constants"
	"github.com/armiariyan/synapsis/internal/pkg/utils"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Product interface {
	FindByUUID(ctx context.Context, uuid string) (result entities.Product, err error)
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.Product, count int64, err error)
}

type product struct {
	db *gorm.DB
}

func NewProduct(db *gorm.DB) Product {
	if db == nil {
		panic("db is nil")
	}

	return &product{db}
}

func (r *product) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.Product, count int64, err error) {
	limit := pagination.Limit
	offset := (pagination.Page - 1) * pagination.Limit
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		queryPayload := r.db.WithContext(egCtx).Limit(int(limit)).Offset(int(offset))
		return utils.CompileConds(queryPayload, conds...).Find(&result).Error
	})
	eg.Go(func() (egErr error) {
		countPayload := r.db.WithContext(egCtx).Model(&entities.Product{})
		return utils.CompileConds(countPayload, conds...).Count(&count).Error
	})
	err = eg.Wait()
	return
}

func (r *product) FindByUUID(ctx context.Context, uuid string) (result entities.Product, err error) {
	err = r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&result).Error
	return
}
