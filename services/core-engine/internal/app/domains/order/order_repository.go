package order

import (
	"context"

	"github.com/gerins/log"
	"gorm.io/gorm"

	"core-engine/internal/app/domains/order/model"
)

type repository struct {
	readDB  *gorm.DB
	writeDB *gorm.DB
}

// NewRepository returns new order Repository.
func NewRepository(readDB *gorm.DB, writeDB *gorm.DB) model.Repository {
	return &repository{
		readDB:  readDB,
		writeDB: writeDB,
	}
}

func (r *repository) SaveOrder(ctx context.Context, order model.Order) (model.Order, error) {
	defer log.Context(ctx).RecordDuration("save order detail to database").Stop()

	if err := r.writeDB.Save(&order).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Order{}, err
	}

	return order, nil
}

func (r *repository) GetOrder(ctx context.Context, id int) (model.Order, error) {
	defer log.Context(ctx).RecordDuration("Get order detail").Stop()

	var order model.Order
	if err := r.writeDB.First(&order, id).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Order{}, err
	}

	return order, nil
}

func (r *repository) SaveMatchOrder(ctx context.Context, order model.MatchOrder) error {
	defer log.Context(ctx).RecordDuration("save match order to database").Stop()

	if err := r.writeDB.Save(&order).Error; err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	return nil
}

func (r *repository) GetPairDetail(ctx context.Context, code string) (model.Pair, error) {
	defer log.Context(ctx).RecordDuration("get pair detail").Stop()

	var pair model.Pair
	if err := r.readDB.Where("code = ?", code).First(&pair).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *repository) GetPairDetailByID(ctx context.Context, id int) (model.Pair, error) {
	defer log.Context(ctx).RecordDuration("get pair detail").Stop()

	var pair model.Pair
	if err := r.readDB.First(&pair, id).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Pair{}, err
	}

	return pair, nil
}

func (r *repository) GetUserWallet(ctx context.Context, userID, cryptoID int) (model.Wallet, error) {
	defer log.Context(ctx).RecordDuration("get user wallet").Stop()

	var wallet model.Wallet
	if err := r.readDB.Where("user_id = ? AND crypto_id = ?", userID, cryptoID).First(&wallet).Error; err != nil {
		log.Context(ctx).Error(err)
		return model.Wallet{}, err
	}

	return wallet, nil
}

// Increase user wallet
func (r *repository) UpdateUserWallet(ctx context.Context, userID, cryptoID int, amount float64) error {
	defer log.Context(ctx).RecordDuration("update user wallet").Stop()

	rawQuery := `UPDATE wallet SET quantity = quantity + ? WHERE user_id = ? AND crypto_id = ?`

	if err := r.writeDB.Exec(rawQuery, amount, userID, cryptoID).Error; err != nil {
		log.Context(ctx).Error(err)
		return err
	}

	return nil
}
