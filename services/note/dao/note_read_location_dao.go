package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteReadLocationDAO GORM实现的笔记阅读位置DAO
type NoteReadLocationDAO struct {
	*baseDao.GormBaseDAO[model.NoteReadLocation]
	logger logging.Logger
}

// NewNoteReadLocationDAO 创建一个新的笔记阅读位置DAO
func NewNoteReadLocationDAO(db *gorm.DB, logger logging.Logger) *NoteReadLocationDAO {
	return &NoteReadLocationDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteReadLocation](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记阅读位置
func (d *NoteReadLocationDAO) Create(ctx context.Context, location *model.NoteReadLocation) error {
	return d.GetDB(ctx).Create(location).Error
}

// FindById 根据ID获取笔记阅读位置
func (d *NoteReadLocationDAO) FindById(ctx context.Context, id string) (*model.NoteReadLocation, error) {
	var location model.NoteReadLocation
	result := d.GetDB(ctx).Where("id = ?", id).First(&location)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记阅读位置失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &location, nil
}

// GetByNoteID 根据笔记ID获取笔记阅读位置
func (d *NoteReadLocationDAO) GetByNoteID(ctx context.Context, noteId string) (*model.NoteReadLocation, error) {
	var location model.NoteReadLocation
	result := d.GetDB(ctx).Where("note_id = ? AND is_deleted = false", noteId).First(&location)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID获取笔记阅读位置失败", "noteId", noteId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &location, nil
}

// UpdateById 更新笔记阅读位置
func (d *NoteReadLocationDAO) UpdateById(ctx context.Context, location *model.NoteReadLocation) error {
	return d.GetDB(ctx).Save(location).Error
}

// DeleteById 删除笔记阅读位置
func (d *NoteReadLocationDAO) DeleteById(ctx context.Context, id string) error {
	return d.GetDB(ctx).Delete(&model.NoteReadLocation{}, id).Error
}

// DeleteByNoteID 根据笔记ID删除笔记阅读位置
func (d *NoteReadLocationDAO) DeleteByNoteId(ctx context.Context, noteId string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteId).Delete(&model.NoteReadLocation{}).Error
}

// GetByUserIdAndNoteId 根据用户ID和笔记ID获取笔记阅读位置
func (d *NoteReadLocationDAO) GetByUserIdAndNoteId(ctx context.Context, userId string, noteId string) (*model.NoteReadLocation, error) {
	var location model.NoteReadLocation
	result := d.GetDB(ctx).Where("creator_id = ? AND note_id = ? AND is_deleted = false", userId, noteId).First(&location)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据用户ID和笔记ID获取笔记阅读位置失败", "userId", userId, "noteId", noteId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &location, nil
}
