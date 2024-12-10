package repository

import (
	"context"

	"gorm.io/gorm"

	"core-engine/internal/app/domains/model"
)

type userRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewUserRepository returns new user Repository.
func NewUserRepository(readDB *gorm.DB, writeDB *gorm.DB) *userRepository {
	return &userRepository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	if err := r.readDB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *userRepository) RegisterNewUser(ctx context.Context, user model.User) (model.User, error) {
	if err := r.writeDB.WithContext(ctx).Save(&user).Error; err != nil {
		return model.User{}, err
	}

	return user, nil
}
