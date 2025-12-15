package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/model"
	"gorm.io/gorm"
)

// CreditDAO GORM实现的会员积分DAO
type CreditDAO struct {
	*baseDao.GormBaseDAO[model.Credit]
	logger logging.Logger
}

// NewCreditDAO 创建一个新的会员积分DAO
func NewCreditDAO(db *gorm.DB, logger logging.Logger) *CreditDAO {
	return &CreditDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.Credit](db, logger),
		logger:      logger,
	}
}

// GetByUserId 根据用户ID获取用户会员积分账号
func (d *CreditDAO) GetByUserId(ctx context.Context, userId string) (*model.Credit, error) {
	var membershipCredit model.Credit
	// 使用 GetDB 从 context 获取 DB/事务，确保在事务中正确执行
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userId).First(&membershipCredit)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户会员积分账号失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &membershipCredit, nil
}

// GetByMembershipId 根据会员ID获取用户会员积分账号
func (d *CreditDAO) GetByMembershipId(ctx context.Context, membershipId string) (*model.Credit, error) {
	var membershipCredit model.Credit
	// 使用 GetDB 从 context 获取 DB/事务，确保在事务中正确执行
	result := d.GetDB(ctx).Where("membership_id = ? and is_deleted = false", membershipId).First(&membershipCredit)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户会员积分账号失败", "membershipId", membershipId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &membershipCredit, nil
}
