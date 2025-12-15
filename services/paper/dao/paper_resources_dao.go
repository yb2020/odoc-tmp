package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperResourcesDAO 提供论文资源数据访问功能
type PaperResourcesDAO struct {
	*baseDao.GormBaseDAO[model.PaperResources]
	logger logging.Logger
}

// NewPaperResourcesDAO 创建一个新的论文资源DAO
func NewPaperResourcesDAO(db *gorm.DB, logger logging.Logger) *PaperResourcesDAO {
	return &PaperResourcesDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperResources](db, logger),
		logger:      logger,
	}
}

// Create 创建论文资源
func (d *PaperResourcesDAO) Create(ctx context.Context, resources *model.PaperResources) error {
	return d.GetDB(ctx).Create(resources).Error
}

// FindById 根据ID获取论文资源
func (d *PaperResourcesDAO) FindById(ctx context.Context, id string) (*model.PaperResources, error) {
	var resources model.PaperResources
	result := d.GetDB(ctx).Where("id = ?", id).First(&resources)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文资源失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &resources, nil
}

// FindByPaperId 根据论文ID获取论文资源列表
func (d *PaperResourcesDAO) FindByPaperId(ctx context.Context, paperId string) ([]model.PaperResources, error) {
	var resources []model.PaperResources
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Find(&resources)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文资源列表失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return resources, nil
}

// FindByPaperTitle 根据论文标题获取论文资源列表
func (d *PaperResourcesDAO) FindByPaperTitle(ctx context.Context, paperTitle string) ([]model.PaperResources, error) {
	var resources []model.PaperResources
	result := d.GetDB(ctx).Where("paper_title LIKE ?", "%"+paperTitle+"%").Find(&resources)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文资源列表失败", "paper_title", paperTitle, "error", result.Error.Error())
		return nil, result.Error
	}
	return resources, nil
}

// FindByResourceTitle 根据资源标题获取论文资源列表
func (d *PaperResourcesDAO) FindByResourceTitle(ctx context.Context, resourceTitle string) ([]model.PaperResources, error) {
	var resources []model.PaperResources
	result := d.GetDB(ctx).Where("resource_title LIKE ?", "%"+resourceTitle+"%").Find(&resources)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文资源列表失败", "resource_title", resourceTitle, "error", result.Error.Error())
		return nil, result.Error
	}
	return resources, nil
}

// UpdateById 更新论文资源
func (d *PaperResourcesDAO) UpdateById(ctx context.Context, resources *model.PaperResources) error {
	return d.Modify(ctx, resources)
}

// DeleteById 删除论文资源
func (d *PaperResourcesDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByPaperId 根据论文ID删除论文资源
func (d *PaperResourcesDAO) DeleteByPaperId(ctx context.Context, paperId string) error {
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Delete(&model.PaperResources{})
	if result.Error != nil {
		d.logger.Error("msg", "删除论文资源失败", "paper_id", paperId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
