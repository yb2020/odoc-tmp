package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteWordConfigDAO GORM实现的笔记单词配置DAO
type NoteWordConfigDAO struct {
	*baseDao.GormBaseDAO[model.NoteWordConfig]
	logger logging.Logger
}

// NewNoteWordConfigDAO 创建一个新的笔记单词配置DAO
func NewNoteWordConfigDAO(db *gorm.DB, logger logging.Logger) *NoteWordConfigDAO {
	return &NoteWordConfigDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteWordConfig](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记单词配置
func (d *NoteWordConfigDAO) Create(ctx context.Context, config *model.NoteWordConfig) error {
	return d.GetDB(ctx).Create(config).Error
}

// FindById 根据ID获取笔记单词配置
func (d *NoteWordConfigDAO) FindById(ctx context.Context, id string) (*model.NoteWordConfig, error) {
	var config model.NoteWordConfig
	result := d.GetDB(ctx).Where("id = ?", id).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记单词配置失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &config, nil
}

// GetByNoteID 根据笔记ID获取笔记单词配置
func (d *NoteWordConfigDAO) GetByNoteID(ctx context.Context, noteID string) (*model.NoteWordConfig, error) {
	var config model.NoteWordConfig
	result := d.GetDB(ctx).Where("note_id = ?", noteID).First(&config)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID获取笔记单词配置失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &config, nil
}

// UpdateById 更新笔记单词配置
func (d *NoteWordConfigDAO) UpdateById(ctx context.Context, config *model.NoteWordConfig) error {
	return d.GetDB(ctx).Save(config).Error
}

// DeleteById 删除笔记单词配置
func (d *NoteWordConfigDAO) DeleteById(ctx context.Context, id string) error {
	return d.GetDB(ctx).Delete(&model.NoteWordConfig{}, id).Error
}

// DeleteByNoteID 根据笔记ID删除笔记单词配置
func (d *NoteWordConfigDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.NoteWordConfig{}).Error
}
