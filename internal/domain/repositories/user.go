package repositories

import (
	"context"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"gorm.io/gorm"
)

type User interface {
	Create(ctx context.Context, entity entities.User) (err error)
	FindByEmail(ctx context.Context, email string) (result entities.User, err error)
	UpdateById(ctx context.Context, id string, entity *entities.User) (err error)
}

type user struct {
	db *gorm.DB
}

func NewUser(db *gorm.DB) User {
	if db == nil {
		panic("db is nil")
	}

	return &user{db}
}

func (r *user) Create(ctx context.Context, entity entities.User) (err error) {
	err = r.db.WithContext(ctx).Create(&entity).Error
	return
}

func (r *user) FindByEmail(ctx context.Context, email string) (result entities.User, err error) {
	err = r.db.WithContext(ctx).Where(&entities.User{Email: email}).First(&result).Error
	if err == gorm.ErrRecordNotFound { // * skip err record not found, checking data empty in service_impl
		err = nil
	}
	return
}

func (r *user) UpdateById(ctx context.Context, id string, entity *entities.User) (err error) {
	tx := r.db.WithContext(ctx).Where("id = ?", id).Updates(&entity)
	err = tx.Error
	if err == nil && tx.RowsAffected < 1 {
		err = gorm.ErrRecordNotFound
	}
	return
}
