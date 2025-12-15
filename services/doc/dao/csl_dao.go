package dao

import (
	"context"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
	"gorm.io/gorm"
)

// CslDAO GORM实现的引用样式DAO
type CslDAO struct {
	*baseDao.GormBaseDAO[model.Csl]
	logger logging.Logger
}

// NewCslDAO 创建一个新的引用样式DAO
func NewCslDAO(db *gorm.DB, logger logging.Logger) *CslDAO {
	return &CslDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.Csl](db, logger),
		logger:      logger,
	}
}

// GetDefaultCsl 获取默认的引用样式
func (d *CslDAO) GetDefaultCsl(ctx context.Context) ([]model.Csl, error) {
	var csls []model.Csl
	result := d.GetDB(ctx).Where("is_deleted = false AND is_default = true").Find(&csls)
	if result.Error != nil {
		d.logger.Error("获取默认引用样式失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return csls, nil
}

// SelectByUserId 根据用户ID获取用户的引用样式列表
func (d *CslDAO) SelectByUserId(ctx context.Context, userId string) ([]model.Csl, error) {
	var csls []model.Csl
	result := d.GetDB(ctx).
		Table("t_csl c").
		Joins("inner join t_user_csl_relation r on c.id = r.csl_id").
		Where("c.is_deleted = false AND r.is_deleted = false AND r.user_id = ?", userId).
		Order("r.sort").
		Find(&csls)

	if result.Error != nil {
		d.logger.Error("根据用户ID获取引用样式列表失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return csls, nil
}
