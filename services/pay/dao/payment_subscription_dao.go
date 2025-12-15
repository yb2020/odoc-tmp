package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pay/model"
	"gorm.io/gorm"
)

// PaymentSubscriptionDAO GORM实现的支付订阅DAO
type PaymentSubscriptionDAO struct {
	*baseDao.GormBaseDAO[model.PaymentSubscription]
	db     *gorm.DB
	logger logging.Logger
}

// NewPaymentSubscriptionDAO 创建一个新的支付订阅DAO
func NewPaymentSubscriptionDAO(db *gorm.DB, logger logging.Logger) *PaymentSubscriptionDAO {
	return &PaymentSubscriptionDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaymentSubscription](db, logger),
		db:          db,
		logger:      logger,
	}
}

// GetByProviderSubscriptionId 根据支付提供商订阅ID获取支付记录
func (d *PaymentSubscriptionDAO) GetByProviderSubscriptionId(ctx context.Context, providerSubscriptionId string) (*model.PaymentSubscription, error) {
	var paymentSubscription model.PaymentSubscription
	result := d.GetDB(ctx).Where("provider_subscription_id = ? and is_deleted = false", providerSubscriptionId).First(&paymentSubscription)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据支付提供商订阅获取支付记录失败", "providerSubscriptionId", providerSubscriptionId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &paymentSubscription, nil
}
