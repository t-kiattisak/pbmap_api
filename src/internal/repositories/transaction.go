package repositories

import (
	"context"

	"gorm.io/gorm"
)

type TransactionManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

type txKey struct{}

type gormTransactionManager struct {
	db *gorm.DB
}

func NewTransactionManager(db *gorm.DB) TransactionManager {
	return &gormTransactionManager{db}
}

func (m *gormTransactionManager) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		txCtx := context.WithValue(ctx, txKey{}, tx)
		return fn(txCtx)
	})
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}
	return defaultDB
}
