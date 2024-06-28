package user

import (
	"context"
	"fmt"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"github.com/armiariyan/synapsis/internal/domain/repositories"
	"github.com/armiariyan/synapsis/internal/infrastructure/xendit"
	"github.com/armiariyan/synapsis/internal/pkg/constants"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/armiariyan/synapsis/internal/pkg/utils"
	"gorm.io/gorm"
)

type service struct {
	db             *gorm.DB
	userRepository repositories.User
	cartRepository repositories.Cart
	xenditWrapper  xendit.Wrapper
}

func NewService() *service {
	return &service{}
}

func (s *service) SetDB(db *gorm.DB) *service {
	s.db = db
	return s
}

func (s *service) SetUserRepository(repository repositories.User) *service {
	s.userRepository = repository
	return s
}
func (s *service) SetCartRepository(repository repositories.Cart) *service {
	s.cartRepository = repository
	return s
}
func (s *service) SetXenditWrapper(wrapper xendit.Wrapper) *service {
	s.xenditWrapper = wrapper
	return s
}
func (s *service) Validate() Service {
	if s.db == nil {
		panic("db is nil")
	}
	if s.userRepository == nil {
		panic("userRepository is nil")
	}
	if s.cartRepository == nil {
		panic("cartRepository is nil")
	}
	if s.xenditWrapper == nil {
		panic("xenditWrapper is nil")
	}
	return s
}

func (s *service) Register(ctx context.Context, req RegisterRequest) (res constants.DefaultResponse, err error) {
	if !utils.ValidatePhoneNumberStartWith62(req.PhoneNumber) {
		log.Error(ctx, fmt.Sprintf("user %s phone number dont start with 62", req.Name))
		err = fmt.Errorf("phone number should start with 62")
		return
	}

	if utils.ValidatePasswordsEquals(req.Password, req.PasswordConfirmation) {
		log.Error(ctx, fmt.Sprintf("password and confirmation password not equals during register user %s", req.Name))
		err = fmt.Errorf("password and confirmation not equals")
		return
	}

	// * validate if user already registered
	user, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed find user by email %s", req.Email), err)
		err = fmt.Errorf("something went wrong")
		return
	}

	if user.ID != 0 {
		log.Error(ctx, fmt.Sprintf("user %s already registered", req.Name))
		err = fmt.Errorf("email already registered")
		return
	}

	// * hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Error(ctx, "failed to hash new password", err)
		err = fmt.Errorf("something went wrong")
		return
	}

	entity := entities.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Password:    hashedPassword,
	}

	err = s.userRepository.Create(ctx, entity)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed create user by email %s", req.Email), err)
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

func (s *service) Login(ctx context.Context, req LoginRequest) (res constants.DefaultResponse, err error) {
	// * validate user exist
	userResult, err := s.userRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to find user by email %s", req.Email), err)
		err = fmt.Errorf("something went wrong [0]")
		return
	}

	if userResult.ID == 0 {
		log.Error(ctx, fmt.Sprintf("email %s is unregistered user", req.Email))
		err = fmt.Errorf("invalid credentials")
		return
	}

	// * validate correct password
	if !utils.CheckPasswordHash(userResult.Password, req.Password) {
		log.Error(ctx, fmt.Sprintf("invalid password for email %s", req.Email))
		err = fmt.Errorf("invalid credentials")
		return
	}

	// * use transaction normally there is more than 1 query in login
	tx := s.db.Begin()
	userRepository := repositories.NewUser(tx)

	token, exp, err := utils.JwtSign(utils.JWTClaimsData{
		ID:          userResult.UUID,
		Email:       userResult.Email,
		PhoneNumber: userResult.PhoneNumber,
	})
	if err != nil {
		tx.Rollback()
		log.Error(ctx, "failed to sign token", err)
		err = fmt.Errorf("something went wrong [1]")
		return
	}

	log.Info(ctx, fmt.Sprintf("[STARTING] update user %s token", userResult.Name))

	err = userRepository.UpdateById(ctx, fmt.Sprint(userResult.ID), &entities.User{Token: token})
	if err != nil {
		tx.Rollback()
		log.Error(ctx, fmt.Sprintf("[FAILED] update user %s token", userResult.Name), err)
		err = fmt.Errorf("something went wrong [2]")
		return
	}

	log.Info(ctx, fmt.Sprintf("[FINISHED] update user %s token", userResult.Name))
	tx.Commit()

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Errors:  make([]string, 0),
		Data: LoginResponse{
			Name:        userResult.Name,
			Email:       userResult.Email,
			PhoneNumber: userResult.PhoneNumber,
			AccessCode:  token,
			ExpiredAt:   exp,
		},
	}

	return
}

func (s *service) Checkout(ctx context.Context) (res constants.DefaultResponse, err error) {
	userData, _ := ctx.Value(string("user")).(entities.User)

	// * get user carts
	carts, err := s.cartRepository.FindAll(ctx,
		utils.DBCond{Where: "user_id = ?", WhereArgs: userData.UUID},
		utils.DBCond{Joins: "User"},
		utils.DBCond{Preload: "Product", PreloadArgs: []utils.DBCond{{Joins: "ProductCategory"}}},
	)
	if err != nil {
		log.Error(ctx, fmt.Sprintf("failed to find all carts for user %s during checkout", userData.Name), err)
		err = fmt.Errorf("something went wrong")
		return
	}

	// * hit Xendit to create invoice
	var items []xendit.ItemInvoice
	var totalAmount int
	for i, cart := range carts {
		item := xendit.ItemInvoice{
			Name:     cart.Product.Name,
			Price:    int(cart.Product.Price),
			Quantity: int(cart.Amount),
			Category: cart.Product.ProductCategory.Name,
			URL:      fmt.Sprintf("example-url.com/items/%d", i),
		}

		items = append(items, item)
		totalPerItem := item.Price * item.Quantity

		totalAmount += totalPerItem
	}

	invoiceCode, _ := utils.GenerateRandomString(5)
	invoiceReq := xendit.CreateInvoiceRequest{
		ExternalID: fmt.Sprintf("SYNAPSIS-INV-%s", invoiceCode),
		Amount:     totalAmount,
		Items:      items,
	}

	resp, err := s.xenditWrapper.CreateInvoice(ctx, invoiceReq)
	if err != nil {
		log.Error(ctx, "failed to create invoice to xendit", err, invoiceReq)
		err = fmt.Errorf("something went wrong")
		return
	}

	res = constants.DefaultResponse{
		Status:  constants.STATUS_SUCCESS,
		Message: constants.MESSAGE_SUCCESS,
		Errors:  make([]string, 0),
		Data: CheckoutResponse{
			InvoiceID:  resp.ExternalID,
			InvoiceURL: resp.InvoiceURL,
		},
	}

	return
}
