package dao

import (
	"context"
	"time"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	pb "github.com/yb2020/odoc/proto/gen/go/membership"
	"github.com/yb2020/odoc/services/membership/model"
	"gorm.io/gorm"
)

// CreditPaymentRecordDAO GORM实现的Cedit支付记录DAO
type CreditPaymentRecordDAO struct {
	*baseDao.GormBaseDAO[model.CreditPaymentRecord]
	logger logging.Logger
}

// NewCreditPaymentRecordDAO 创建一个新的Cedit支付记录DAO
func NewCreditPaymentRecordDAO(db *gorm.DB, logger logging.Logger) *CreditPaymentRecordDAO {
	return &CreditPaymentRecordDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.CreditPaymentRecord](db, logger),
		logger:      logger,
	}
}

// GetConfirmExpiredList 获取确认积分支付过期的列表
func (d *CreditPaymentRecordDAO) GetConfirmExpiredList(ctx context.Context, size int) ([]model.CreditPaymentRecord, error) {
	var entities []model.CreditPaymentRecord
	result := d.GetDB(ctx).Where("is_deleted = false and status = ? and confirm_expired_at <= ?", int32(pb.CreditPaymentStatus_CREDIT_PAYMENT_STATUS_PAY_AWAITING_CONFIRMATION), time.Now()).Order("confirm_expired_at desc").Limit(size).Find(&entities)
	if result.Error != nil {
		d.logger.Error("msg", "获取待确认支付的过期记录失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return entities, nil
}
