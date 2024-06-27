package user

import (
	"context"
	"fmt"

	"github.com/armiariyan/synapsis/internal/domain/entities"
	"github.com/armiariyan/synapsis/internal/domain/repositories"
	"github.com/armiariyan/synapsis/internal/pkg/constants"
	"github.com/armiariyan/synapsis/internal/pkg/log"
	"github.com/armiariyan/synapsis/internal/pkg/utils"
	"gorm.io/gorm"
)

type service struct {
	db             *gorm.DB
	userRepository repositories.User
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

func (s *service) Validate() Service {
	if s.db == nil {
		panic("db is nil")
	}
	if s.userRepository == nil {
		panic("userRepository is nil")
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
	//TODO armia check if record not found
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

// func (s *service) FindAll(ctx context.Context, req FindAllRequest) (res FindAllResponse, err error) {
// 	datas, err := s.userRepository.FindAll(ctx)
// 	if err != nil {
// 		log.Error(ctx, "failed to find all users", err)
// 		return
// 	}
// 	res = FindAllResponse{
// 		DefaultResponse: constants.DefaultResponse{
// 			Status:  constants.STATUS_SUCCESS,
// 			Message: constants.MESSAGE_SUCCESS,
// 			Data:    datas,
// 			Errors:  make([]string, 0),
// 		},
// 	}
// 	return
// }
