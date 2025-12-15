package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/model"
	"gorm.io/gorm"
)

// PaperNoteDAO 提供论文笔记数据访问功能

// PaperNoteDAO GORM实现的论文笔记DAO
type PaperNoteDAO struct {
	*baseDao.GormBaseDAO[model.PaperNote]
	logger logging.Logger
}

// NewPaperNoteDAO 创建一个新的论文笔记DAO
func NewPaperNoteDAO(db *gorm.DB, logger logging.Logger) *PaperNoteDAO {
	return &PaperNoteDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperNote](db, logger),
		logger:      logger,
	}
}

// CreatePaperNote 创建笔记
func (d *PaperNoteDAO) Create(ctx context.Context, note *model.PaperNote) error {
	return d.GetDB(ctx).Create(note).Error
}

// GetPaperNoteByID 根据ID获取笔记
func (d *PaperNoteDAO) FindById(ctx context.Context, id string) (*model.PaperNote, error) {
	var note model.PaperNote
	result := d.GetDB(ctx).Where("id = ?", id).First(&note)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文笔记失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &note, nil
}

// GetPaperNoteByPaperID 根据论文ID获取笔记
func (d *PaperNoteDAO) GetByPaperID(ctx context.Context, paperID string) (*model.PaperNote, error) {
	var note model.PaperNote
	result := d.GetDB(ctx).Where("paper_id = ?", paperID).First(&note)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文笔记失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return &note, nil
}

// UpdatePaperNote 更新笔记
func (d *PaperNoteDAO) UpdateById(ctx context.Context, note *model.PaperNote) error {
	return d.Modify(ctx, note)
}

// DeletePaperNote 删除笔记
func (d *PaperNoteDAO) DeleteById(ctx context.Context, id string) error {
	// 由于 BaseDAO 中的 DeleteById 接收的是 int64 类型，这里需要转换
	return d.RemoveById(ctx, id)
}

// ListPaperNotes 列出笔记
func (d *PaperNoteDAO) List(ctx context.Context, limit, offset int) ([]model.PaperNote, error) {
	var notes []model.PaperNote
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&notes)
	if result.Error != nil {
		d.logger.Error("msg", "列出论文笔记失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return notes, nil
}

// CountPaperNotes 获取笔记总数
func (d *PaperNoteDAO) Count(ctx context.Context) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PaperNote{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文笔记总数失败", "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetPaperNoteByID 根据ID获取笔记
func (d *PaperNoteDAO) FindByPdfIdAndUserId(ctx context.Context, pdfId string, userId string) (*model.PaperNote, error) {
	var note model.PaperNote
	result := d.GetDB(ctx).Where("pdf_id = ? and creator_id = ? and is_deleted = false", pdfId, userId).First(&note)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文笔记失败", "pdfId", pdfId, "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return &note, nil
}

// SelectByUserIdLimit 根据用户ID获取笔记列表，并限制返回数量
func (d *PaperNoteDAO) SelectByUserIdLimit(ctx context.Context, userId string, limit int) ([]model.PaperNote, error) {
	var notes []model.PaperNote
	result := d.GetDB(ctx).
		Where("creator_id = ? AND is_deleted = false", userId).
		Order("updated_at DESC"). // 按更新时间倒序排列，获取最近的笔记
		Limit(limit).
		Find(&notes)

	if result.Error != nil {
		d.logger.Error("msg", "根据用户ID获取笔记列表失败", "userId", userId, "limit", limit, "error", result.Error.Error())
		return nil, result.Error
	}
	return notes, nil
}

// CountByUserId 根据用户ID获取笔记数量
func (d *PaperNoteDAO) CountByUserId(ctx context.Context, userId string) (int64, error) {
	var count int64
	result := d.GetDB(ctx).
		Model(&model.PaperNote{}).
		Where("creator_id = ? AND is_deleted = false", userId).
		Count(&count)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户笔记数量失败", "userId", userId, "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetByPaperIdAndUserIdOrderByNoteCount 根据论文ID和用户ID获取笔记，按笔记数量和修改时间排序
func (d *PaperNoteDAO) GetByPaperIdAndUserIdOrderByNoteCount(ctx context.Context, paperId string, userId string) (*model.PaperNote, error) {
	var note model.PaperNote
	result := d.GetDB(ctx).
		Where("paper_id = ? AND is_deleted = false AND creator_id = ? AND note_count > 0", paperId, userId).
		Order("note_count DESC, updated_at DESC").
		Limit(1).
		First(&note)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "根据论文ID和用户ID获取笔记失败", "paperId", paperId, "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}

	return &note, nil
}

func (d *PaperNoteDAO) GetAllNoteByUserId(ctx context.Context, userId string) ([]model.PaperNote, error) {
	var notes []model.PaperNote
	result := d.GetDB(ctx).
		Where("creator_id = ? AND is_deleted = false", userId).Find(&notes)

	if result.Error != nil {
		d.logger.Error("msg", "获取用户笔记列表失败", "userId", userId, "error", result.Error.Error())
		return nil, result.Error
	}
	return notes, nil
}
