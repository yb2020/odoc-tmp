package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteDrawEntityDAO GORM实现的笔记绘制实体DAO
type NoteDrawEntityDAO struct {
	*baseDao.GormBaseDAO[model.NoteDrawEntity]
	logger logging.Logger
}

// NewNoteDrawEntityDAO 创建一个新的笔记绘制实体DAO
func NewNoteDrawEntityDAO(db *gorm.DB, logger logging.Logger) *NoteDrawEntityDAO {
	return &NoteDrawEntityDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteDrawEntity](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记绘制实体
func (d *NoteDrawEntityDAO) Create(ctx context.Context, drawEntity *model.NoteDrawEntity) error {
	return d.GetDB(ctx).Create(drawEntity).Error
}

// FindById 根据ID获取笔记绘制实体
func (d *NoteDrawEntityDAO) FindById(ctx context.Context, id string) (*model.NoteDrawEntity, error) {
	var drawEntity model.NoteDrawEntity
	result := d.GetDB(ctx).Where("id = ?", id).First(&drawEntity)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记绘制实体失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &drawEntity, nil
}

// GetByNoteID 根据笔记ID获取笔记绘制实体列表
func (d *NoteDrawEntityDAO) GetByNoteID(ctx context.Context, noteID string) ([]model.NoteDrawEntity, error) {
	var drawEntities []model.NoteDrawEntity
	result := d.GetDB(ctx).Where("note_id = ?", noteID).Find(&drawEntities)
	if result.Error != nil {
		d.logger.Error("msg", "根据笔记ID获取笔记绘制实体列表失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return drawEntities, nil
}

// GetByNoteIDAndPage 根据笔记ID和页码获取笔记绘制实体列表
func (d *NoteDrawEntityDAO) GetByNoteIDAndPage(ctx context.Context, noteID string, pageNumber int) ([]model.NoteDrawEntity, error) {
	var drawEntities []model.NoteDrawEntity
	result := d.GetDB(ctx).Where("note_id = ? AND page_number = ?", noteID, pageNumber).Find(&drawEntities)
	if result.Error != nil {
		d.logger.Error("msg", "根据笔记ID和页码获取笔记绘制实体列表失败", "noteID", noteID, "pageNumber", pageNumber, "error", result.Error.Error())
		return nil, result.Error
	}
	return drawEntities, nil
}

// UpdateById 更新笔记绘制实体
func (d *NoteDrawEntityDAO) UpdateById(ctx context.Context, drawEntity *model.NoteDrawEntity) error {
	return d.GetDB(ctx).Save(drawEntity).Error
}

// DeleteById 删除笔记绘制实体
func (d *NoteDrawEntityDAO) DeleteById(ctx context.Context, id string) error {
	return d.GetDB(ctx).Delete(&model.NoteDrawEntity{}, id).Error
}

// DeleteByNoteID 根据笔记ID删除所有笔记绘制实体
func (d *NoteDrawEntityDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.NoteDrawEntity{}).Error
}
