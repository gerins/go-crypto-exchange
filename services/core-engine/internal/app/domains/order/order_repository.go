package order

import (
	"context"

	"gorm.io/gorm"

	"core-engine/internal/app/domains/order/model"
	gormpkg "core-engine/pkg/gorm"
)

type repository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewRepository returns new order Repository.
func NewRepository(readDB *gorm.DB, writeDB *gorm.DB) *repository {
	return &repository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *repository) SaveOrder(ctx context.Context, order model.Order) (model.Order, error) {
	writeDB := r.writeDB
	if tx := gormpkg.GetTransactionFromContext(ctx); tx != nil {
		writeDB = tx
	}

	if err := writeDB.WithContext(ctx).Save(&order).Error; err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *repository) GetOrder(ctx context.Context, id int) (model.Order, error) {
	var order model.Order
	if err := r.writeDB.WithContext(ctx).First(&order, id).Error; err != nil {
		return model.Order{}, err
	}

	return order, nil
}

func (r *repository) SaveMatchOrder(ctx context.Context, matchOrder model.MatchOrder) error {
	writeDB := r.writeDB
	if tx := gormpkg.GetTransactionFromContext(ctx); tx != nil {
		writeDB = tx
	}

	if err := writeDB.WithContext(ctx).Save(&matchOrder).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) GetPairDetail(ctx context.Context, code string) (model.Pair, error) {
	var pair model.Pair
	if err := r.readDB.WithContext(ctx).Where("code = ?", code).First(&pair).Error; err != nil {
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *repository) GetPairDetailByID(ctx context.Context, id int) (model.Pair, error) {
	var pair model.Pair
	if err := r.readDB.WithContext(ctx).First(&pair, id).Error; err != nil {
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *repository) GetUserWallet(ctx context.Context, userID, cryptoID int) (model.Wallet, error) {
	var wallet model.Wallet
	if err := r.readDB.WithContext(ctx).Where("user_id = ? AND crypto_id = ?", userID, cryptoID).First(&wallet).Error; err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil
}

// Increase user wallet
func (r *repository) UpdateUserWallet(ctx context.Context, userID, cryptoID int, amount float64) error {
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
