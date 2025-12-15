package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// UserDocClassifyDAO GORM实现的用户文档分类DAO
type UserDocClassifyDAO struct {
	*baseDao.GormBaseDAO[model.UserDocClassify]
	logger logging.Logger
}

// NewUserDocClassifyDAO 创建一个新的用户文档分类DAO
func NewUserDocClassifyDAO(db *gorm.DB, logger logging.Logger) *UserDocClassifyDAO {
	return &UserDocClassifyDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.UserDocClassify](db, logger),
		logger:      logger,
	}
}

// GetUserDocClassifiesByUserID 根据用户ID获取用户文档分类
// 只返回未删除的分类
func (d *UserDocClassifyDAO) GetUserDocClassifiesByUserID(ctx context.Context, userID string) ([]model.UserDocClassify, error) {
	var classifies []model.UserDocClassify
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userID).Find(&classifies)
	if result.Error != nil {
		d.logger.Error("msg", "获取用户文档分类失败", "userID", userID, "error", result.Error.Error())
		return nil, result.Error
	}
	return classifies, nil
}
