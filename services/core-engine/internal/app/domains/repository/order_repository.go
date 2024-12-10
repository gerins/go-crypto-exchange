package repository

import (
	"context"

	"gorm.io/gorm"

	"core-engine/internal/app/domains/model"
	gormpkg "core-engine/pkg/gorm"
)

type orderRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewOrderRepository returns new order Repository.
func NewOrderRepository(readDB *gorm.DB, writeDB *gorm.DB) *orderRepository {
	return &orderRepository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *orderRepository) SaveOrder(ctx context.Context, order model.Order) (model.Order, error) {
	writeDB := r.writeDB
	if tx := gormpkg.GetTransactionFromContext(ctx); tx != nil {
		writeDB = tx
	}

	if err := writeDB.WithContext(ctx).Save(&order).Error; err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) GetOrder(ctx context.Context, id int) (model.Order, error) {
	var order model.Order
	if err := r.writeDB.WithContext(ctx).First(&order, id).Error; err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *orderRepository) SaveMatchOrder(ctx context.Context, matchOrder model.MatchOrder) error {
	writeDB := r.writeDB
	if tx := gormpkg.GetTransactionFromContext(ctx); tx != nil {
		writeDB = tx
	}

	if err := writeDB.WithContext(ctx).Save(&matchOrder).Error; err != nil {
		return err
	}

	return nil
}
