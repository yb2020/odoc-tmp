package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pay/model"
	"gorm.io/gorm"
)

// PaymentRecordDAO GORM实现的支付记录DAO
type PaymentRecordDAO struct {
	*baseDao.GormBaseDAO[model.PaymentRecord]
	db     *gorm.DB
	logger logging.Logger
}

// NewPaymentRecordDAO 创建一个新的支付记录DAO
func NewPaymentRecordDAO(db *gorm.DB, logger logging.Logger) *PaymentRecordDAO {
	return &PaymentRecordDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaymentRecord](db, logger),
		db:          db,
		logger:      logger,
	}
}

// GetByOrderID 根据订单ID获取支付记录
func (d *PaymentRecordDAO) GetByOrderID(ctx context.Context, orderID string) (*model.PaymentRecord, error) {
	var paymentRecord model.PaymentRecord
	result := d.GetDB(ctx).Where("order_id = ? and is_deleted = false", orderID).First(&paymentRecord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据订单ID获取支付记录失败", "orderID", orderID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &paymentRecord, nil
}

// GetByProviderTxId 根据支付提供商交易ID获取支付记录
func (d *PaymentRecordDAO) GetByProviderTxId(ctx context.Context, providerTxId string) (*model.PaymentRecord, error) {
	var paymentRecord model.PaymentRecord
	result := d.GetDB(ctx).Where("provider_tx_id = ? and is_deleted = false", providerTxId).First(&paymentRecord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据支付提供商交易ID获取支付记录失败", "providerTxId", providerTxId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &paymentRecord, nil
}

// GetBySubscriptionId 根据订阅ID获取支付记录
// 注意：一个订阅ID可能关联多个支付记录（初始支付和续费）。此方法返回找到的第一个记录。
func (d *PaymentRecordDAO) GetBySubscriptionIdAndInvoiceId(ctx context.Context, subscriptionId string, invoiceId string) (*model.PaymentRecord, error) {
	var paymentRecord model.PaymentRecord
	result := d.GetDB(ctx).Where("provider_subscription_id = ? AND invoice_id = ? AND is_deleted = false", subscriptionId, invoiceId).First(&paymentRecord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据订阅ID获取支付记录失败", "subscriptionId", subscriptionId, "invoiceId", invoiceId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &paymentRecord, nil
}

// UpdateFields 更新支付记录的特定字段
// 'fields'是数据库列名到其新值的映射
// 此方法适用于部分更新，包括将字段设置为零值
func (d *PaymentRecordDAO) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	if id == "" {
		return gorm.ErrMissingWhereClause
	}
	result := d.GetDB(ctx).Model(&model.PaymentRecord{}).Where("id = ?", id).Updates(fields)
	if result.Error != nil {
		d.logger.Error("msg", "更新支付记录字段失败", "id", id, "error", result.Error.Error())
		return result.Error
	}
	return nil
}

// GetByUserIdAndStatus 根据用户ID和支付状态获取支付记录列表
func (d *PaymentRecordDAO) GetByUserIdAndStatus(ctx context.Context, userId string, status string) ([]model.PaymentRecord, error) {
	var records []model.PaymentRecord
	result := d.GetDB(ctx).Where("user_id = ? AND status = ? AND is_deleted = false", userId, status).
		Order("id DESC").Find(&records)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID和状态获取支付记录失败", "userId", userId, "status", status, "error", result.Error.Error())
		return nil, result.Error
	}
	return records, nil
}

// GetByChannelAndStatus 根据支付渠道和状态获取支付记录列表
func (d *PaymentRecordDAO) GetByChannelAndStatus(ctx context.Context, channel string, status string) ([]model.PaymentRecord, error) {
	var records []model.PaymentRecord
	result := d.GetDB(ctx).Where("channel = ? AND status = ? AND is_deleted = false", channel, status).
		Order("id DESC").Find(&records)
	if result.Error != nil {
		d.logger.Error("msg", "根据支付渠道和状态获取支付记录失败", "channel", channel, "status", status, "error", result.Error.Error())
		return nil, result.Error
	}
	return records, nil
}
