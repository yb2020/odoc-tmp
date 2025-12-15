package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteLatestReadDAO GORM实现的笔记最近阅读DAO
type NoteLatestReadDAO struct {
	*baseDao.GormBaseDAO[model.NoteLatestRead]
	logger logging.Logger
}

// NewNoteLatestReadDAO 创建一个新的笔记最近阅读DAO
func NewNoteLatestReadDAO(db *gorm.DB, logger logging.Logger) *NoteLatestReadDAO {
	return &NoteLatestReadDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteLatestRead](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记最近阅读记录
func (d *NoteLatestReadDAO) Create(ctx context.Context, latestRead *model.NoteLatestRead) error {
	return d.GetDB(ctx).Create(latestRead).Error
}

// FindById 根据ID获取笔记最近阅读记录
func (d *NoteLatestReadDAO) FindById(ctx context.Context, id string) (*model.NoteLatestRead, error) {
	var latestRead model.NoteLatestRead
	result := d.GetDB(ctx).Where("id = ?", id).First(&latestRead)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记最近阅读记录失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &latestRead, nil
}

// GetByNoteID 根据笔记ID获取笔记最近阅读记录
func (d *NoteLatestReadDAO) GetByNoteID(ctx context.Context, noteID string) (*model.NoteLatestRead, error) {
	var latestRead model.NoteLatestRead
	result := d.GetDB(ctx).Where("note_id = ?", noteID).First(&latestRead)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID获取笔记最近阅读记录失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &latestRead, nil
}

// UpdateById 更新笔记最近阅读记录
func (d *NoteLatestReadDAO) UpdateById(ctx context.Context, latestRead *model.NoteLatestRead) error {
	return d.GetDB(ctx).Save(latestRead).Error
}

// DeleteById 删除笔记最近阅读记录
func (d *NoteLatestReadDAO) DeleteById(ctx context.Context, id string) error {
	return d.GetDB(ctx).Delete(&model.NoteLatestRead{}, id).Error
}

// DeleteByNoteID 根据笔记ID删除笔记最近阅读记录
func (d *NoteLatestReadDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.NoteLatestRead{}).Error
}
