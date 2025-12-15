package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// PaperNoteAccessDAO 提供论文笔记访问记录数据访问功能

// PaperNoteAccessDAO GORM实现的论文笔记访问记录DAO
type PaperNoteAccessDAO struct {
	*baseDao.GormBaseDAO[model.PaperNoteAccess]
	logger logging.Logger
}

// NewPaperNoteAccessDAO 创建一个新的论文笔记访问记录DAO
func NewPaperNoteAccessDAO(db *gorm.DB, logger logging.Logger) *PaperNoteAccessDAO {
	return &PaperNoteAccessDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperNoteAccess](db, logger),
		logger:      logger,
	}
}

// Create 创建论文笔记访问记录
func (d *PaperNoteAccessDAO) Create(ctx context.Context, access *model.PaperNoteAccess) error {
	return d.GetDB(ctx).Create(access).Error
}

// FindById 根据ID查找论文笔记访问记录
func (d *PaperNoteAccessDAO) FindById(ctx context.Context, id string) (*model.PaperNoteAccess, error) {
	var access model.PaperNoteAccess
	result := d.GetDB(ctx).Where("id = ?", id).First(&access)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据ID查找论文笔记访问记录失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &access, nil
}

// GetByNoteID 根据笔记ID获取论文笔记访问记录
func (d *PaperNoteAccessDAO) GetByNoteID(ctx context.Context, noteID string) (*model.PaperNoteAccess, error) {
	var access model.PaperNoteAccess
	result := d.GetDB(ctx).Where("note_id = ? and is_deleted = false", noteID).First(&access)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID获取论文笔记访问记录失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &access, nil
}

// UpdateById 更新论文笔记访问记录
func (d *PaperNoteAccessDAO) UpdateById(ctx context.Context, access *model.PaperNoteAccess) error {
	return d.Modify(ctx, access)
}

// DeleteById 删除论文笔记访问记录
func (d *PaperNoteAccessDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByNoteID 根据笔记ID删除论文笔记访问记录
func (d *PaperNoteAccessDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.PaperNoteAccess{}).Error
}

// List 列出论文笔记访问记录
func (d *PaperNoteAccessDAO) List(ctx context.Context, limit, offset int) ([]model.PaperNoteAccess, error) {
	var accesses []model.PaperNoteAccess
	result := d.GetDB(ctx).Limit(limit).Offset(offset).Find(&accesses)
	if result.Error != nil {
		d.logger.Error("msg", "列出论文笔记访问记录失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return accesses, nil
}

// Count 获取论文笔记访问记录总数
func (d *PaperNoteAccessDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PaperNoteAccess{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文笔记访问记录总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}
