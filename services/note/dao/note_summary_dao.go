package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// NoteSummaryDAO GORM实现的笔记摘要DAO
type NoteSummaryDAO struct {
	*baseDao.GormBaseDAO[model.NoteSummary]
	logger logging.Logger
}

// NewNoteSummaryDAO 创建一个新的笔记摘要DAO
func NewNoteSummaryDAO(db *gorm.DB, logger logging.Logger) *NoteSummaryDAO {
	return &NoteSummaryDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.NoteSummary](db, logger),
		logger:      logger,
	}
}

// Create 创建笔记摘要
func (d *NoteSummaryDAO) Create(ctx context.Context, summary *model.NoteSummary) error {
	return d.GetDB(ctx).Create(summary).Error
}

// FindById 根据ID获取笔记摘要
func (d *NoteSummaryDAO) FindById(ctx context.Context, id string) (*model.NoteSummary, error) {
	var summary model.NoteSummary
	result := d.GetDB(ctx).Where("id = ? and is_deleted = false", id).First(&summary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取笔记摘要失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &summary, nil
}

// GetByNoteID 根据笔记ID获取笔记摘要
func (d *NoteSummaryDAO) GetByNoteID(ctx context.Context, noteID string) (*model.NoteSummary, error) {
	var summary model.NoteSummary
	result := d.GetDB(ctx).Where("note_id = ? and is_deleted = false", noteID).First(&summary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据笔记ID获取笔记摘要失败", "noteID", noteID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &summary, nil
}

// GetByUserId 根据用户ID获取笔记摘要列表
func (d *NoteSummaryDAO) GetByUserId(ctx context.Context, userId string) ([]model.NoteSummary, error) {
	var summaries []model.NoteSummary
	result := d.GetDB(ctx).Where("user_id = ? and is_deleted = false", userId).Find(&summaries)
	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID获取笔记摘要列表失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return summaries, nil
}

// UpdateById 更新笔记摘要
func (d *NoteSummaryDAO) UpdateById(ctx context.Context, summary *model.NoteSummary) error {
	return d.GetDB(ctx).Save(summary).Error
}

// DeleteById 删除笔记摘要
func (d *NoteSummaryDAO) DeleteById(ctx context.Context, id string) error {
	return d.GetDB(ctx).Delete(&model.NoteSummary{}, id).Error
}

// DeleteByNoteID 根据笔记ID删除笔记摘要
func (d *NoteSummaryDAO) DeleteByNoteID(ctx context.Context, noteID string) error {
	return d.GetDB(ctx).Where("note_id = ?", noteID).Delete(&model.NoteSummary{}).Error
}
