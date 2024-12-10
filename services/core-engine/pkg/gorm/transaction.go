package gorm

import (
	"context"

	"gorm.io/gorm"
)

var key struct{}

func InitTransactionToContext(parentCtx context.Context, writeDB *gorm.DB) (context.Context, *gorm.DB) {
	tx := writeDB.WithContext(parentCtx).Begin()
	parentCtx = context.WithValue(parentCtx, key, tx)
	return parentCtx, tx
}

func GetTransactionFromContext(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(key).(*gorm.DB); ok {
		return tx
	}

	return nil
}
