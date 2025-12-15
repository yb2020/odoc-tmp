package dao

import (
	"context"
	"errors"

	"gorm.io/gorm"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/user/model"
)

// UserDAO PostgreSQL 实现的用户 DAO
type UserDAO struct {
	*baseDao.GormBaseDAO[model.User]
	logger logging.Logger
}

// NewUserDAO 创建一个新的用户 DAO
func NewUserDAO(db *gorm.DB, logger logging.Logger) *UserDAO {
	return &UserDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.User](db, logger),
		logger:      logger,
	}
}

// GetUserByEmail 根据电子邮件获取用户
func (d *UserDAO) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	result := d.GetDB(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		d.logger.Error("msg", "根据电子邮件获取用户失败", "email", email, "error", result.Error.Error())
		return nil, result.Error
	}
	return &user, nil
}

// FindByGoogleID 根据 Google Open ID 查找用户
func (d *UserDAO) FindByGoogleID(ctx context.Context, googleID string) (*model.User, error) {
	var user model.User
	err := d.GetDB(ctx).Where("google_open_id = ?", googleID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 返回 nil, nil 表示未找到，由 service 层决定是创建还是返回业务错误
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
