package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteExportHistoryDAO GORM实现的笔记导出历史DAO
type NoteExportHistoryDAO struct {
	*baseDao.GormBaseDAO[model.NoteExportHistory]
	logger logging.Logger
}

// NewNoteExportHistoryDAO 创建一个新的笔记导出历史DAO
func NewNoteExportHistoryDAO(db *gorm.DB, logger logging.Logger) *NoteExportHistoryDAO {
	return &NoteExportHistoryDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteExportHistory](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记导出历史
func (d *NoteExportHistoryDAO) Create(ctx context.Context, history *model.NoteExportHistory) error {
	return d.GetDB(ctx).Create(history).Error
}

// FindById 根据ID获取笔记导出历史
func (d *NoteExportHistoryDAO) FindById(ctx context.Context, id string) (*model.NoteExportHistory, error) {
	var history model.NoteExportHistory
	result := d.GetDB(ctx).Where("id = ?", id).First(&history)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记导出历史失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &history, nil
}

// GetByNoteID 根据笔记ID获取笔记导出历史列表
func (d *NoteExportHistoryDAO) GetByNoteID(ctx context.Context, noteID string) ([]model.NoteExportHistory, error) {
	var histories []model.NoteExportHistory
	result := d.GetDB(ctx).Where("note_id = ?", noteID).Find(&histories)
	if result.Error != nil {
		d.logger.Error("msg", "根据笔记ID获取笔记导出历史列表失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return histories, nil
}

// GetLatestByNoteID 根据笔记ID获取最新的笔记导出历史
func (d *NoteExportHistoryDAO) GetLatestByNoteID(ctx context.Context, noteID string) (*model.NoteExportHistory, error) {
	var history model.NoteExportHistory
	result := d.GetDB(ctx).Where("note_id = ?", noteID).Order("version DESC").First(&history)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID获取最新的笔记导出历史失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &history, nil
}

// GetByNoteIDAndVersion 根据笔记ID和版本获取笔记导出历史
func (d *NoteExportHistoryDAO) GetByNoteIDAndVersion(ctx context.Context, noteID string, version int64) (*model.NoteExportHistory, error) {
	var history model.NoteExportHistory
	result := d.GetDB(ctx).Where("note_id = ? AND version = ?", noteID, version).First(&history)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID和版本获取笔记导出历史失败", "noteID", noteID, "version", version, "error", result.Error.Error())
		return nil, result.Error
	}
	return &history, nil
}

// UpdateById 更新笔记导出历史
func (d *NoteExportHistoryDAO) UpdateById(ctx context.Context, history *model.NoteExportHistory) error {
	return d.GetDB(ctx).Save(history).Error
}

// DeleteById 删除笔记导出历史
func (d *NoteExportHistoryDAO) DeleteById(ctx context.Context, id string) error {
	return d.GetDB(ctx).Delete(&model.NoteExportHistory{}, id).Error
}

// DeleteByNoteID 根据笔记ID删除所有笔记导出历史
func (d *NoteExportHistoryDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.NoteExportHistory{}).Error
}
