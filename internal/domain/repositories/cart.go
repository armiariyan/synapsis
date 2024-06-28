package repositories

import (
	"context"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"github.com/armiariyan/synapsis/internal/pkg/constants"
	"github.com/armiariyan/synapsis/internal/pkg/utils"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Cart interface {
	Create(ctx context.Context, entity entities.Cart) (err error)
	UpdateById(ctx context.Context, id string, entity *entities.Cart) (err error)
	FindByUUID(ctx context.Context, uuid string) (result entities.Cart, err error)
	FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.Cart, count int64, err error)
	DeleteByUUID(ctx context.Context, uuid string) (err error)
}

type cart struct {
	db *gorm.DB
}

func NewCart(db *gorm.DB) Cart {
	if db == nil {
		panic("db is nil")
	}

	return &cart{db}
}

func (r *cart) Create(ctx context.Context, entity entities.Cart) (err error) {
	err = r.db.WithContext(ctx).Create(&entity).Error
	return
}

func (r *cart) UpdateById(ctx context.Context, id string, entity *entities.Cart) (err error) {
	tx := r.db.WithContext(ctx).Where("id = ?", id).Updates(&entity)
	err = tx.Error
	if err == nil && tx.RowsAffected < 1 {
		err = gorm.ErrRecordNotFound
	}
	return
}

func (r *cart) FindByUUID(ctx context.Context, uuid string) (result entities.Cart, err error) {
	err = r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&result).Error
	return
}

func (r *cart) FindAllAndCount(ctx context.Context, pagination constants.PaginationRequest, conds ...utils.DBCond) (result []entities.Cart, count int64, err error) {
	limit := pagination.Limit
	offset := (pagination.Page - 1) * pagination.Limit
	eg, egCtx := errgroup.WithContext(ctx)
	eg.Go(func() (egErr error) {
		queryPayload := r.db.WithContext(egCtx).Limit(int(limit)).Offset(int(offset))
		return utils.CompileConds(queryPayload, conds...).Find(&result).Error
	})
	eg.Go(func() (egErr error) {
		countPayload := r.db.WithContext(egCtx).Model(&entities.Cart{})
		return utils.CompileConds(countPayload, conds...).Count(&count).Error
	})
	err = eg.Wait()
	return
}

func (r *cart) DeleteByUUID(ctx context.Context, uuid string) (err error) {
	err = r.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&entities.Cart{}).Error
	return

}
