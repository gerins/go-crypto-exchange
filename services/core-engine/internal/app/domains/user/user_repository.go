package user

import (
	"context"

	"gorm.io/gorm"
)

type repository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewRepository returns new user Repository.
func NewRepository(readDB *gorm.DB, writeDB *gorm.DB) *repository {
	return &repository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *repository) FindUserByEmail(ctx context.Context, email string) (User, error) {
	var user User
	if err := r.readDB.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}

func (r *repository) RegisterNewUser(ctx context.Context, user User) (User, error) {
	if err := r.writeDB.WithContext(ctx).Save(&user).Error; err != nil {
		return User{}, err
	}

	return user, nil
}
