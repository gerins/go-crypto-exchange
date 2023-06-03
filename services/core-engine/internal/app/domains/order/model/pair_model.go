package model

import "time"

type Pair struct {
	ID                int        `json:"id" gorm:"column:id;type:int;primaryKey;autoIncrement"`
	Code              string     `json:"code" gorm:"column:code;type:varchar;size:255"`
	PrimaryCryptoID   int        `json:"primary_crypto_id" gorm:"column:primary_crypto_id;type:int"`
	SecondaryCryptoID int        `json:"secondary_crypto_id" gorm:"column:secondary_crypto_id;type:int"`
	CreatedAt         time.Time  `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeletedAt         *time.Time `json:"deleted_at" gorm:"column:deleted_at;type:datetime"`
}

func (Pair) TableName() string {
	return "pairs"
}
