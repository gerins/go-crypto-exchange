//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
package user

import (
	"context"
	"time"
)

type User struct {
	ID          int        `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	FullName    string     `json:"full_name" gorm:"column:full_name;type:varchar;size:255"`
	Email       string     `json:"email" gorm:"column:email;type:varchar;size:255"`
	PhoneNumber string     `json:"phone_number" gorm:"column:phone_number;type:varchar;size:255"`
	Password    string     `json:"password" gorm:"column:password;type:varchar;size:255"`
	Status      bool       `json:"status" gorm:"column:status;type:tinyint"`
	CreatedAt   time.Time  `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime"`
}

func (User) TableName() string {
	return "users"
}

//counterfeiter:generate -o ./mock . Usecase
type Usecase interface {
	Login(ctx context.Context, loginReq LoginRequest) (LoginResponse, error)
	Register(ctx context.Context, registerReq RegisterRequest) (User, error)
}

//counterfeiter:generate -o ./mock . Repository
type Repository interface {
	FindUserByEmail(ctx context.Context, email string) (User, error)
	RegisterNewUser(ctx context.Context, user User) (User, error)
}
