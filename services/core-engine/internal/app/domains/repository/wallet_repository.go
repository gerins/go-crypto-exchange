package repository

import (
	"context"

	"gorm.io/gorm"

	"core-engine/internal/app/domains/model"
	gormpkg "core-engine/pkg/gorm"
)

type walletRepository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewWalletRepository returns new order Repository.
func NewWalletRepository(readDB *gorm.DB, writeDB *gorm.DB) *walletRepository {
	return &walletRepository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *walletRepository) GetPairDetail(ctx context.Context, code string) (model.Pair, error) {
	var pair model.Pair
	if err := r.readDB.WithContext(ctx).Where("code = ?", code).First(&pair).Error; err != nil {
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *walletRepository) GetPairDetailByID(ctx context.Context, id int) (model.Pair, error) {
	var pair model.Pair
	if err := r.readDB.WithContext(ctx).First(&pair, id).Error; err != nil {
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *walletRepository) Save(ctx context.Context, wallet model.Wallet) error {
	writeDB := r.writeDB
	if tx := gormpkg.GetTransactionFromContext(ctx); tx != nil {
		writeDB = tx
	}

	if err := writeDB.WithContext(ctx).Save(&wallet).First(&wallet).Error; err != nil {
		return err
	}

	return nil
}

func (r *walletRepository) GetUserWallet(ctx context.Context, userID, cryptoID int) (model.Wallet, error) {
	var wallet model.Wallet
	if err := r.readDB.WithContext(ctx).Where("user_id = ? AND crypto_id = ?", userID, cryptoID).First(&wallet).Error; err != nil {
		return model.Wallet{}, err
	}

	return wallet, nil
}

// Increase user wallet
func (r *walletRepository) UpdateUserWallet(ctx context.Context, userID, cryptoID int, amount float64) error {
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
