package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperAccessDAO 提供论文访问权限数据访问功能
type PaperAccessDAO struct {
	*baseDao.GormBaseDAO[model.PaperAccess]
	logger logging.Logger
}

// NewPaperAccessDAO 创建一个新的论文访问权限DAO
func NewPaperAccessDAO(db *gorm.DB, logger logging.Logger) *PaperAccessDAO {
	return &PaperAccessDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperAccess](db, logger),
		logger:      logger,
	}
}

// Create 创建论文访问权限
func (d *PaperAccessDAO) Create(ctx context.Context, access *model.PaperAccess) error {
	return d.GetDB(ctx).Create(access).Error
}

// FindById 根据ID获取论文访问权限
func (d *PaperAccessDAO) FindById(ctx context.Context, id string) (*model.PaperAccess, error) {
	var access model.PaperAccess
	result := d.GetDB(ctx).Where("id = ?", id).First(&access)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文访问权限失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &access, nil
}

// FindByPaperId 根据论文ID获取论文访问权限列表
func (d *PaperAccessDAO) FindByPaperId(ctx context.Context, paperId string) ([]model.PaperAccess, error) {
	var accesses []model.PaperAccess
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Find(&accesses)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文访问权限列表失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return accesses, nil
}

// FindByUserId 根据用户ID获取论文访问权限列表
func (d *PaperAccessDAO) FindByUserId(ctx context.Context, userId string) ([]model.PaperAccess, error) {
	var accesses []model.PaperAccess
	result := d.GetDB(ctx).Where("user_id = ?", userId).Find(&accesses)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文访问权限列表失败", "user_id", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return accesses, nil
}

// FindByPaperIdAndUserId 根据论文ID和用户ID获取论文访问权限
func (d *PaperAccessDAO) FindByPaperIdAndUserId(ctx context.Context, paperId, userId string) (*model.PaperAccess, error) {
	var access model.PaperAccess
	result := d.GetDB(ctx).Where("paper_id = ? AND user_id = ?", paperId, userId).First(&access)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文访问权限失败", "paper_id", paperId, "user_id", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &access, nil
}

// DeleteById 删除论文访问权限
func (d *PaperAccessDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByPaperIdAndUserId 根据论文ID和用户ID删除论文访问权限
func (d *PaperAccessDAO) DeleteByPaperIdAndUserId(ctx context.Context, paperId, userId string) error {
	result := d.GetDB(ctx).Where("paper_id = ? AND user_id = ?", paperId, userId).Delete(&model.PaperAccess{})
	if result.Error != nil {
		d.logger.Error("msg", "删除论文访问权限失败", "paper_id", paperId, "user_id", userId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
