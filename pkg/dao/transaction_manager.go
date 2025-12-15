// pkg/dao/transaction.go
package dao

import (
	"context"

	"gorm.io/gorm"
)

// TxKey 事务上下文键
type TxKey string

const (
	// TransactionContextKey 用于在上下文中存储事务对象的键
	TransactionContextKey TxKey = "transaction"
)

// TransactionManager 事务管理器
type TransactionManager struct {
	db *gorm.DB
}

// NewTransactionManager 创建事务管理器
func NewTransactionManager(db *gorm.DB) *TransactionManager {
	return &TransactionManager{
		db: db,
	}
}

// GetDB 获取数据库连接，如果上下文中有事务则返回事务对象
func (tm *TransactionManager) GetDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(TransactionContextKey).(*gorm.DB); ok {
		return tx
	}
	return tm.db.WithContext(ctx)
}

// ExecuteInTransaction 在事务中执行函数
func (tm *TransactionManager) ExecuteInTransaction(ctx context.Context, fn func(txCtx context.Context) error) error {
	// 如果上下文中已经有事务，直接使用该事务
	if _, ok := ctx.Value(TransactionContextKey).(*gorm.DB); ok {
		return fn(ctx)
	}

	// 开始新事务
	tx := tm.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// 创建包含事务的新上下文
	txCtx := context.WithValue(ctx, TransactionContextKey, tx)

	// 确保事务最终会提交或回滚
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // 重新抛出 panic
		}
	}()

	// 执行事务函数
	if err := fn(txCtx); err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
