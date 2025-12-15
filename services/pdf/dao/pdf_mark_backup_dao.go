package dao

import (
	"context"
	"errors"
	"fmt"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/model"
	"gorm.io/gorm"
)

// PdfMarkBackupDAO GORM实现的PDF标记备份DAO
type PdfMarkBackupDAO struct {
	*baseDao.GormBaseDAO[model.PdfMarkBackup]
	logger logging.Logger
}

// NewPdfMarkBackupDAO 创建一个新的PDF标记备份DAO
func NewPdfMarkBackupDAO(db *gorm.DB, logger logging.Logger) *PdfMarkBackupDAO {
	return &PdfMarkBackupDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfMarkBackup](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF标记备份
func (d *PdfMarkBackupDAO) Create(ctx context.Context, backup *model.PdfMarkBackup) error {
	return d.GetDB(ctx).Create(backup).Error
}

// GetPdfMarkBackupByID 根据ID获取PDF标记备份
func (d *PdfMarkBackupDAO) GetPdfMarkBackupByID(ctx context.Context, id string) (*model.PdfMarkBackup, error) {
	var backup model.PdfMarkBackup
	result := d.GetDB(ctx).Where("id = ?", id).First(&backup)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记备份失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &backup, nil
}

// GetPdfMarkBackupsByPaperID 根据论文ID获取PDF标记备份列表
func (d *PdfMarkBackupDAO) GetPdfMarkBackupsByPaperID(ctx context.Context, paperID string) ([]model.PdfMarkBackup, error) {
	var backups []model.PdfMarkBackup
	result := d.GetDB(ctx).Where("paper_id = ?", paperID).Find(&backups)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记备份列表失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return backups, nil
}

// GetPdfMarkBackupsByNoteID 根据笔记ID获取PDF标记备份列表
func (d *PdfMarkBackupDAO) GetPdfMarkBackupsByNoteID(ctx context.Context, noteID string) ([]model.PdfMarkBackup, error) {
	var backups []model.PdfMarkBackup
	result := d.GetDB(ctx).Where("note_id = ?", noteID).Find(&backups)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记备份列表失败", "note_id", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return backups, nil
}

// UpdatePdfMarkBackup 更新PDF标记备份
func (d *PdfMarkBackupDAO) UpdatePdfMarkBackup(ctx context.Context, backup *model.PdfMarkBackup) error {
	return d.Modify(ctx, backup)
}

// DeletePdfMarkBackup 删除PDF标记备份
func (d *PdfMarkBackupDAO) DeletePdfMarkBackup(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPdfMarkBackups 列出PDF标记备份
func (d *PdfMarkBackupDAO) ListPdfMarkBackups(ctx context.Context, limit, offset int) ([]model.PdfMarkBackup, error) {
	var backups []model.PdfMarkBackup
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&backups)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF标记备份失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return backups, nil
}

// CountPdfMarkBackups 获取PDF标记备份总数
func (d *PdfMarkBackupDAO) CountPdfMarkBackups(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfMarkBackup{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记备份总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}
