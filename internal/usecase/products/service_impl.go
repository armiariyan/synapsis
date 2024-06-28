package products

import (
	"context"
	"fmt"
	"math"
	"strings"

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

func (s *service) FindByCategory(ctx context.Context, req FindByCategoryRequest) (res constants.DefaultResponse, err error) {
	// * find category name
	category, err := s.productCategoryRepository.FindByName(ctx, strings.ToLower(req.Category))
	if err != nil {
		log.Error(ctx, "failed find product category by name", err)
		err = fmt.Errorf("something went wrong")
		return
	}

	if category.ID == 0 {
		log.Error(ctx, "invalid category name", req.Category)
		err = fmt.Errorf("invalid category name")
		return
	}

	// * find and count products
	products, count, err := s.productRepository.FindAllAndCount(ctx, req.PaginationRequest,
		utils.DBCond{Where: "category_id = ?", WhereArgs: category.UUID},
	)
	if err != nil {
		log.Error(ctx, "failed to find all and count products", err)
		err = fmt.Errorf("something went wrong")
		return
	}

	var results []FindByCategoryResponse
	for _, product := range products {
		results = append(results, FindByCategoryResponse{
			ID:           product.UUID,
			Name:         product.Name,
			Price:        product.Price,
			CategoryName: category.Name,
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

func (s *service) AddToCart(ctx context.Context, req AddToCartRequest) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(string("user")).(entities.User)

	// * validate product exist
	product, err := s.productRepository.FindByUUID(ctx, req.ProductID)
	if err != nil {
		log.Error(ctx, "failed find product by uuid", err)
		err = fmt.Errorf("something went wrong")
		return
	}

	if product.ID == 0 {
		log.Error(ctx, "invalid product id")
		err = fmt.Errorf("invalid product id")
		return
	}

	entity := entities.Cart{
		UserID:    userData.UUID,
		ProductID: product.UUID,
		Amount:    req.Amount,
	}

	// * insert into cart
	err = s.cartRepository.Create(ctx, entity)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed create cart for user %s", userData.Name), err, req)
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
