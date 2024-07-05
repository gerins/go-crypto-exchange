package user

import (
	"context"

	"github.com/gerins/log"
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
	defer log.Context(ctx).RecordDuration("find user by email").Stop()

	var user User
	if err := r.readDB.Where("email = ?", email).First(&user).Error; err != nil {
		log.Context(ctx).Error(err)
		return User{}, err
	}

	return user, nil
}

func (r *repository) RegisterNewUser(ctx context.Context, user User) (User, error) {
	defer log.Context(ctx).RecordDuration("register new user").Stop()

	if err := r.writeDB.Save(&user).Error; err != nil {
		log.Context(ctx).Error(err)
		return User{}, err
	}

	return user, nil
}
