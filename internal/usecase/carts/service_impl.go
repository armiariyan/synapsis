package carts

import (
	"context"
	"fmt"
	"math"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"github.com/armiariyan/synapsis/internal/domain/repositories"
	"github.com/armiariyan/synapsis/internal/pkg/constants"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/armiariyan/synapsis/internal/pkg/utils"
	"gorm.io/gorm"
)

type service struct {
	db                        *gorm.DB
	productCategoryRepository repositories.ProductCategory
	productRepository         repositories.Product
	cartRepository            repositories.Cart
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *gorm.DB) *service {
	s.db = db
	return s
}

func (s *service) SetProductCategoryRepository(repository repositories.ProductCategory) *service {
	s.productCategoryRepository = repository
	return s
}

func (s *service) SetProductRepository(repository repositories.Product) *service {
	s.productRepository = repository
	return s
}

func (s *service) SetCartRepository(repository repositories.Cart) *service {
	s.cartRepository = repository
	return s
}

func (s *service) Validate() Service {
	if s.db == nil {
		panic("db is nil")
	}
	if s.productCategoryRepository == nil {
		panic("productCategoryRepository is nil")
	}
	if s.productRepository == nil {
		panic("productRepository is nil")
	}
	if s.cartRepository == nil {
		panic("cartRepository is nil")
	}
	return s
}

func (s *service) FindAll(ctx context.Context, req FindAllRequest) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(string("user")).(entities.User)

	carts, count, err := s.cartRepository.FindAllAndCount(ctx, req.PaginationRequest,
		utils.DBCond{Where: "user_id = ?", WhereArgs: userData.UUID},
		utils.DBCond{Joins: "User"},
		// * we need ProductCategory in relation Product
		utils.DBCond{Preload: "Product", PreloadArgs: []utils.DBCond{{Joins: "ProductCategory"}}},
	)

	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to find all and count carts for user %s", userData.Name), err)
		err = fmt.Errorf("something went wrong")
		return
	}

	var results []FindAllResponse
	for _, cart := range carts {
		results = append(results, FindAllResponse{
			ID:       cart.UUID,
			Name:     cart.Product.Name,
			Price:    cart.Product.Price,
			Amount:   cart.Amount,
			Category: cart.Product.ProductCategory.Name,
		})
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data: constants.PaginationResponseData{
			Results: results,
			PaginationData: constants.PaginationData{
				Page:        req.Page,
				Limit:       req.Limit,
				TotalPages:  uint(math.Ceil(float64(count) / float64(req.Limit))),
				TotalItems:  uint(count),
				HasNext:     req.Page < uint(math.Ceil(float64(count)/float64(req.Limit))),
				HasPrevious: req.Page > 1,
			},
		},
		Errors: make([]string, 0),
	}
	return
}

func (s *service) DeleteProduct(ctx context.Context, uuid string) (res constants.DefaultResponse, err error) {
	err = s.cartRepository.DeleteByUUID(ctx, uuid)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed delete data from cart"), err)
		err = fmt.Errorf("something went wrong")
		return
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Data:    struct{}{},
		Errors:  make([]string, 0),
	}

	return
}
