package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteWordDAO GORM实现的笔记单词DAO
type NoteWordDAO struct {
	*baseDao.GormBaseDAO[model.NoteWord]
	logger logging.Logger
}

// NewNoteWordDAO 创建一个新的笔记单词DAO
func NewNoteWordDAO(db *gorm.DB, logger logging.Logger) *NoteWordDAO {
	return &NoteWordDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteWord](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记单词
func (d *NoteWordDAO) Create(ctx context.Context, noteWord *model.NoteWord) error {
	return d.GetDB(ctx).Create(noteWord).Error
}

// FindById 根据ID获取笔记单词
func (d *NoteWordDAO) FindById(ctx context.Context, id string) (*model.NoteWord, error) {
	var noteWord model.NoteWord
	result := d.GetDB(ctx).Where("id = ? AND is_deleted = false", id).First(&noteWord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记单词失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &noteWord, nil
}

// GetByNoteID 根据笔记ID获取笔记单词列表
func (d *NoteWordDAO) GetByNoteID(ctx context.Context, noteID string) ([]model.NoteWord, error) {
	var noteWords []model.NoteWord
	result := d.GetDB(ctx).Where("note_id = ? AND is_deleted = false", noteID).Find(&noteWords)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词列表失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return noteWords, nil
}

// UpdateById 更新笔记单词
func (d *NoteWordDAO) UpdateById(ctx context.Context, noteWord *model.NoteWord) error {
	return d.GetDB(ctx).Save(noteWord).Error
}

// DeleteByNoteID 根据笔记ID删除所有笔记单词
func (d *NoteWordDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ? AND is_deleted = false", noteID).Delete(&model.NoteWord{}).Error
}

// List 列出笔记单词
func (d *NoteWordDAO) List(ctx context.Context, limit, offset int) ([]model.NoteWord, error) {
	var noteWords []model.NoteWord
	result := d.GetDB(ctx).Limit(limit).Offset(offset).Find(&noteWords)
	if result.Error != nil {
		d.logger.Error("msg", "列出笔记单词失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return noteWords, nil
}

// Count 获取笔记单词总数
func (d *NoteWordDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.NoteWord{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetByNoteIDWithMinLoadedId 根据笔记ID获取笔记单词列表，并且只返回ID小于minLoadedId的记录，按ID降序排序，限制返回数量
func (d *NoteWordDAO) GetByNoteIDWithMinLoadedId(ctx context.Context, noteID string, minLoadedId string, limit int) ([]model.NoteWord, error) {
	var noteWords []model.NoteWord
	result := d.GetDB(ctx).Where("note_id = ? AND id < ? AND is_deleted = false", noteID, minLoadedId).Order("id DESC").Limit(limit).Find(&noteWords)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词列表失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return noteWords, nil
}

// GetListByNoteID 根据笔记ID获取笔记单词列表，并且只返回ID小于minLoadedId的记录，按ID降序排序，限制返回数量
func (d *NoteWordDAO) GetListByNoteID(ctx context.Context, noteID string, limit, offset int) ([]model.NoteWord, error) {
	var noteWords []model.NoteWord
	result := d.GetDB(ctx).Where("note_id = ? AND is_deleted = false", noteID).Order("id DESC").Limit(limit).Offset(offset).Find(&noteWords)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词列表失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return noteWords, nil
}

// GetListByNoteIds 根据笔记Ids获取笔记单词列表
func (d *NoteWordDAO) GetListByNoteIds(ctx context.Context, noteIds []string, limit, offset int) ([]model.NoteWord, error) {
	var noteWords []model.NoteWord
	result := d.GetDB(ctx).Where("note_id IN ? AND is_deleted = false", noteIds).Order("id DESC").Limit(limit).Offset(offset).Find(&noteWords)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词列表失败", "noteIds", noteIds, "error", result.Error.Error())
		return nil, result.Error
	}
	return noteWords, nil
}

// GetCountByNoteId 获取笔记单词总数By NoteId
func (d *NoteWordDAO) GetCountByNoteId(ctx context.Context, noteId string) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.NoteWord{}).Where("note_id = ? AND is_deleted = false", noteId).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetCountByNoteIds 获取笔记单词总数By NoteIds
func (d *NoteWordDAO) GetCountByNoteIds(ctx context.Context, noteIds []string) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.NoteWord{}).Where("note_id IN ? AND is_deleted = false", noteIds).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取笔记单词总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// FindById 根据ID获取笔记单词
func (d *NoteWordDAO) GetByNoteIdAndWord(ctx context.Context, noteId string, word string) (*model.NoteWord, error) {
	var noteWord model.NoteWord
	result := d.GetDB(ctx).Where("note_id = ? AND word = ? AND is_deleted = false", noteId, word).First(&noteWord)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记单词失败", "noteId", noteId, "word", word, "error", result.Error.Error())
		return nil, result.Error
	}
	return &noteWord, nil
}
