package dao

import (
	"context"
	"errors"
	"time"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/membership/model"
	"gorm.io/gorm"
)

// UserMembershipDAO GORM实现的用户会员DAO
type UserMembershipDAO struct {
	*baseDao.GormBaseDAO[model.UserMembership]
	logger logging.Logger
}

// NewUserMembershipDAO 创建一个新的用户会员DAO
func NewUserMembershipDAO(db *gorm.DB, logger logging.Logger) *UserMembershipDAO {
	return &UserMembershipDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserMembership](db, logger),
		logger:      logger,
	}
}

// GetByUserId 根据用户ID获取用户会员
func (d *UserMembershipDAO) GetByUserId(ctx context.Context, userId string) (*model.UserMembership, error) {
	var userMembership model.UserMembership
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userId).First(&userMembership)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取用户会员失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &userMembership, nil
}

// GetExpiredList 获取过期的会员列表
func (d *UserMembershipDAO) GetExpiredList(ctx context.Context, size int) ([]model.UserMembership, error) {
	var userMemberships []model.UserMembership
	result := d.GetDB(ctx).Where("is_deleted = false and end_at <= ?", time.Now()).Order("end_at desc").Limit(size).Find(&userMemberships)
	if result.Error != nil {
		d.logger.Error("msg", "获取过期的会员列表失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return userMemberships, nil
}

// GetExpiredListByMemberType 获取过期的会员列表
func (d *UserMembershipDAO) GetExpiredListByMemberType(ctx context.Context, memberType int, size int) ([]model.UserMembership, error) {
	var userMemberships []model.UserMembership
	result := d.GetDB(ctx).Where("is_deleted = false and end_at <= ? and type = ?", time.Now(), memberType).Order("end_at desc").Limit(size).Find(&userMemberships)
	if result.Error != nil {
		d.logger.Error("msg", "获取过期的会员列表失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return userMemberships, nil
}
