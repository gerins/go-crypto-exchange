package user

import (
	"context"
	"time"

	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"core-engine/config"
	serverError "core-engine/pkg/error"
)

type usecase struct {
	securityConfig config.Security
	userRepository Repository
	validator      *validator.Validate
}

// NewUsecase returns new user usecase.
func NewUsecase(securityConfig config.Security, validator *validator.Validate, userRepository Repository) *usecase {
	return &usecase{
		securityConfig: securityConfig,
		validator:      validator,
		userRepository: userRepository,
	}
}

func (u *usecase) Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error) {
	// Find user detail in repository
	userDetail, err := u.userRepository.FindUserByEmail(ctx, loginReq.Email)
	if err != nil {
		return LoginResponse{}, err
	}

	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(loginReq.Password)); err != nil {
		log.Context(ctx).Error(err)
		return LoginResponse{}, serverError.ErrInvalidUsernameOrPassword(err)
	}

	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userDetail.ID
	claims["email"] = userDetail.Email
	claims["exp"] = time.Now().Add(u.securityConfig.Jwt.Duration).Unix()

	// Generate the token string
	tokenString, err := token.SignedString([]byte(u.securityConfig.Jwt.Key))
	if err != nil {
		log.Context(ctx).Errorf("failed generate token string, %v", err)
		return LoginResponse{}, ErrInvalidPassword
	}

	loginResponse := LoginResponse{
		Email: userDetail.Email,
		Token: tokenString,
	}

	return loginResponse, nil
}

func (u *usecase) Register(ctx context.Context, registerReq RegisterRequest) (User, error) {
	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Context(ctx).Error(err)
		return User{}, err
	}

	// Create new model for inserting to database
	newUser := User{
		FullName:    registerReq.FullName,
		Email:       registerReq.Email,
		PhoneNumber: registerReq.PhoneNumber,
		Password:    string(hashedPassword),
		Status:      true, // Active
	}

	newUser, err = u.userRepository.RegisterNewUser(ctx, newUser)
	if err != nil {
		return User{}, err
	}

	return newUser, nil
}
