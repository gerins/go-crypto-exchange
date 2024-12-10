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

func (r *orderRepository) GetPairDetail(ctx context.Context, code string) (model.Pair, error) {
	var pair model.Pair
	if err := r.readDB.WithContext(ctx).Where("code = ?", code).First(&pair).Error; err != nil {
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *orderRepository) GetPairDetailByID(ctx context.Context, id int) (model.Pair, error) {
	var pair model.Pair
	if err := r.readDB.WithContext(ctx).First(&pair, id).Error; err != nil {
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *orderRepository) GetUserWallet(ctx context.Context, userID, cryptoID int) (model.Wallet, error) {
	var wallet model.Wallet
	if err := r.readDB.WithContext(ctx).Where("user_id = ? AND crypto_id = ?", userID, cryptoID).First(&wallet).Error; err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil
}

// Increase user wallet
func (r *orderRepository) UpdateUserWallet(ctx context.Context, userID, cryptoID int, amount float64) error {
	writeDB := r.writeDB
	if tx := gormpkg.GetTransactionFromContext(ctx); tx != nil {
		writeDB = tx
	}

	rawQuery := `UPDATE wallet SET quantity = quantity + ? WHERE user_id = ? AND crypto_id = ?`
	if err := writeDB.WithContext(ctx).Exec(rawQuery, amount, userID, cryptoID).Error; err != nil {
		return err
	}

	return nil
}
