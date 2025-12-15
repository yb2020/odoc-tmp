package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteShapeDAO GORM实现的笔记形状DAO
type NoteShapeDAO struct {
	*baseDao.GormBaseDAO[model.NoteShape]
	logger logging.Logger
}

// NewNoteShapeDAO 创建一个新的笔记形状DAO
func NewNoteShapeDAO(db *gorm.DB, logger logging.Logger) *NoteShapeDAO {
	return &NoteShapeDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteShape](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记形状
func (d *NoteShapeDAO) Create(ctx context.Context, shape *model.NoteShape) error {
	return d.GetDB(ctx).Create(shape).Error
}

// FindById 根据ID获取笔记形状
func (d *NoteShapeDAO) FindById(ctx context.Context, id string) (*model.NoteShape, error) {
	var shape model.NoteShape
	result := d.GetDB(ctx).Where("id = ?", id).First(&shape)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记形状失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &shape, nil
}

// GetByUUID 根据UUID获取笔记形状
func (d *NoteShapeDAO) GetByUUID(ctx context.Context, uuid string) (*model.NoteShape, error) {
	var shape model.NoteShape
	result := d.GetDB(ctx).Where("uuid = ?", uuid).First(&shape)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据UUID获取笔记形状失败", "uuid", uuid, "error", result.Error.Error())
		return nil, result.Error
	}
	return &shape, nil
}

// GetByNoteID 根据笔记ID获取笔记形状列表
func (d *NoteShapeDAO) GetByNoteID(ctx context.Context, noteID string) ([]model.NoteShape, error) {
	var shapes []model.NoteShape
	result := d.GetDB(ctx).Where("note_id = ? AND is_deleted = false", noteID).Order("page_number").Find(&shapes)
	if result.Error != nil {
		d.logger.Error("msg", "根据笔记ID获取笔记形状列表失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return shapes, nil
}

// GetByNoteIDAndPage 根据笔记ID和页码获取笔记形状列表
func (d *NoteShapeDAO) GetByNoteIDAndPage(ctx context.Context, noteID string, pageNumber int) ([]model.NoteShape, error) {
	var shapes []model.NoteShape
	result := d.GetDB(ctx).Where("note_id = ? AND page_number = ?", noteID, pageNumber).Find(&shapes)
	if result.Error != nil {
		d.logger.Error("msg", "根据笔记ID和页码获取笔记形状列表失败", "noteID", noteID, "pageNumber", pageNumber, "error", result.Error.Error())
		return nil, result.Error
	}
	return shapes, nil
}

// UpdateById 更新笔记形状
func (d *NoteShapeDAO) UpdateById(ctx context.Context, shape *model.NoteShape) error {
	return d.GetDB(ctx).Save(shape).Error
}

// DeleteByUUID 根据UUID删除笔记形状
func (d *NoteShapeDAO) DeleteByUUID(ctx context.Context, uuid string) error {
	return d.GetDB(ctx).Where("uuid = ?", uuid).Delete(&model.NoteShape{}).Error
}

// DeleteByNoteID 根据笔记ID删除所有笔记形状
func (d *NoteShapeDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.NoteShape{}).Error
}
