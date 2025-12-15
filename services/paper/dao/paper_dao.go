package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperDAO 提供论文数据访问功能
type PaperDAO struct {
	*baseDao.GormBaseDAO[model.Paper]
	logger logging.Logger
}

// NewPaperDAO 创建一个新的论文DAO
func NewPaperDAO(db *gorm.DB, logger logging.Logger) *PaperDAO {
	return &PaperDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.Paper](db, logger),
		logger:      logger,
	}
}

// FindById 根据ID获取论文
func (d *PaperDAO) FindById(ctx context.Context, id string) (*model.Paper, error) {
	var paper model.Paper
	result := d.GetDB(ctx).Where("id = ?", id).First(&paper)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &paper, nil
}

// FindByPaperId 根据论文ID获取论文
func (d *PaperDAO) FindByPaperId(ctx context.Context, paperId string) (*model.Paper, error) {
	var paper model.Paper
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).First(&paper)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &paper, nil
}

// FindByOwnerId 根据拥有者ID获取论文列表
func (d *PaperDAO) FindByOwnerId(ctx context.Context, ownerId string) ([]model.Paper, error) {
	var papers []model.Paper
	result := d.GetDB(ctx).Where("owner_id = ?", ownerId).Find(&papers)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文列表失败", "owner_id", ownerId, "error", result.Error.Error())
		return nil, result.Error
	}
	return papers, nil
}

// UpdateById 更新论文
func (d *PaperDAO) UpdateById(ctx context.Context, paper *model.Paper) error {
	return d.Modify(ctx, paper)
}

// DeleteById 删除论文
func (d *PaperDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// List 列出论文
func (d *PaperDAO) List(ctx context.Context, limit, offset int) ([]model.Paper, error) {
	var papers []model.Paper
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&papers)
	if result.Error != nil {
		d.logger.Error("msg", "列出论文失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return papers, nil
}

// Count 获取论文总数
func (d *PaperDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.Paper{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}
