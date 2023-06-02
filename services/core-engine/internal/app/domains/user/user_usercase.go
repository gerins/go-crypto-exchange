package user

import (
	"context"

	"github.com/go-playground/validator/v10"
)

type usecase struct {
	userRepository Repository
	validator      *validator.Validate
}

// NewUsecase returns new user usecase.
func NewUsecase(validator *validator.Validate, userRepository Repository) *usecase {
	return &usecase{
		validator:      validator,
		userRepository: userRepository,
	}
}

func (u *usecase) Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error) {
	return LoginResponse{}, nil
}

func (u *usecase) Register(ctx context.Context, registerReq RegisterRequest) (User, error) {
	return User{}, nil
}
