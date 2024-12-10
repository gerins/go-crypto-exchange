package usecase

import (
	"context"
	"time"

	"github.com/gerins/log"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"core-engine/config"
	"core-engine/internal/app/domains/dto"
	"core-engine/internal/app/domains/model"
	serverError "core-engine/pkg/error"
)

type userUsecase struct {
	validator        *validator.Validate
	securityConfig   config.Security
	userRepository   model.UserRepository
	walletRepository model.WalletRepository
}

// NewUserUsecase returns new user userUsecase.
func NewUserUsecase(
	validator *validator.Validate,
	securityConfig config.Security,
	userRepository model.UserRepository,
	walletRepository model.WalletRepository,
) *userUsecase {
	return &userUsecase{
		validator:        validator,
		securityConfig:   securityConfig,
		userRepository:   userRepository,
		walletRepository: walletRepository,
	}
}

func (u *userUsecase) Login(ctx context.Context, loginReq dto.LoginRequest) (dto.LoginResponse, error) {
	// Find user detail in repository
	userDetail, err := u.userRepository.FindUserByEmail(ctx, loginReq.Email)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	// Comparing the password with the hash
	if err = bcrypt.CompareHashAndPassword([]byte(userDetail.Password), []byte(loginReq.Password)); err != nil {
		log.Context(ctx).Error(err)
		return dto.LoginResponse{}, serverError.ErrInvalidUsernameOrPassword(err)
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
		return dto.LoginResponse{}, serverError.ErrInvalidUsernameOrPassword(err)
	}

	loginResponse := dto.LoginResponse{
		Email: userDetail.Email,
		Token: tokenString,
	}

	return loginResponse, nil
}

func (u *userUsecase) Register(ctx context.Context, registerReq dto.RegisterRequest) (model.User, error) {
	// Hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerReq.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Context(ctx).Error(err)
		return model.User{}, err
	}

	// Create new model for inserting to database
	newUser := model.User{
		FullName:    registerReq.FullName,
		Email:       registerReq.Email,
		PhoneNumber: registerReq.PhoneNumber,
		Password:    string(hashedPassword),
		Status:      true, // Active
	}

	newUser, err = u.userRepository.RegisterNewUser(ctx, newUser)
	if err != nil {
		return model.User{}, err
	}

	// Inject initial balance for testing purpose
	if err := u.injectInitialBalance(ctx, newUser.ID); err != nil {
		return model.User{}, err
	}

	return newUser, nil
}

func (u *userUsecase) injectInitialBalance(ctx context.Context, userID int) error {
	pair, err := u.walletRepository.GetPairDetail(ctx, "DOGEIDRT")
	if err != nil {
		return err
	}

	primaryWallet := model.Wallet{
		UserID:   userID,
		CryptoID: pair.PrimaryCryptoID,
		Quantity: 1000000000,
	}

	secondaryWallet := model.Wallet{
		UserID:   userID,
		CryptoID: pair.SecondaryCryptoID,
		Quantity: 1000000000,
	}

	if err := u.walletRepository.Save(ctx, primaryWallet); err != nil {
		return err
	}
	if err := u.walletRepository.Save(ctx, secondaryWallet); err != nil {
		return err
	}

	return nil
}
