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

// PdfMarkDAO GORM实现的PDF标记DAO
type PdfMarkDAO struct {
	*baseDao.GormBaseDAO[model.PdfMark]
	logger logging.Logger
}

// NewPdfMarkDAO 创建一个新的PDF标记DAO
func NewPdfMarkDAO(db *gorm.DB, logger logging.Logger) *PdfMarkDAO {
	return &PdfMarkDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfMark](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF标记
func (d *PdfMarkDAO) Create(ctx context.Context, mark *model.PdfMark) error {
	return d.GetDB(ctx).Create(mark).Error
}

// GetPdfMarkByID 根据ID获取PDF标记
func (d *PdfMarkDAO) GetPdfMarkByID(ctx context.Context, id string) (*model.PdfMark, error) {
	var mark model.PdfMark
	result := d.GetDB(ctx).Where("id = ? and is_deleted = false", id).First(&mark)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &mark, nil
}

// GetPdfMarksByPaperID 根据论文ID获取PDF标记列表
func (d *PdfMarkDAO) GetPdfMarksByPaperID(ctx context.Context, paperID string) ([]model.PdfMark, error) {
	var marks []model.PdfMark
	result := d.GetDB(ctx).Where("paper_id = ? and is_deleted = false", paperID).Find(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记列表失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return marks, nil
}

// GetCountPdfMarksByNoteIdWithoutIsHighlight 根据笔记ID获取PDF标记总数
func (d *PdfMarkDAO) GetCountPdfMarksByNoteIdWithoutIsHighlight(ctx context.Context, noteId string) (int64, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfMark{}).Where("note_id = ? and is_highlight = false and is_deleted = false and type != 4", noteId).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记总数失败", "note_id", noteId, "error", result.Error.Error())
		return 0, result.Error
	}
	return count, nil
}

// GetPdfMarksByNoteId 根据笔记ID获取PDF标记列表
func (d *PdfMarkDAO) GetPdfMarksByNoteId(ctx context.Context, noteId string) ([]model.PdfMark, error) {
	var marks []model.PdfMark
	result := d.GetDB(ctx).Where("note_id = ? and is_deleted = false", noteId).Find(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记列表失败", "note_id", noteId, "error", result.Error.Error())
		return nil, result.Error
	}
	return marks, nil
}

// GetPdfMarksByNoteIds 根据笔记Ids列表获取PDF标记列表
func (d *PdfMarkDAO) GetPdfMarksByNoteIds(ctx context.Context, userId string, noteIds []string) ([]model.PdfMark, error) {
	var marks []model.PdfMark
	result := d.GetDB(ctx).Where("note_id in ? and is_deleted = false and creator_id = ?", noteIds, userId).Find(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记列表失败", "note_id", noteIds, "error", result.Error.Error())
		return nil, result.Error
	}
	return marks, nil
}

// GetPdfMarksByNoteIds 根据笔记Ids列表获取PDF标记列表， 按照sortExpression表达式排序
func (d *PdfMarkDAO) GetPdfMarksByNoteIdsWithSortWithoutIsHighlight(ctx context.Context, userId string, noteIds []string, sortExpression string) ([]model.PdfMark, error) {
	var marks []model.PdfMark
	result := d.GetDB(ctx).Where("note_id in ? and is_deleted = false and is_highlight = false and creator_id = ? and type != 4", noteIds, userId).Order(sortExpression).Find(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记列表失败", "note_id", noteIds, "error", result.Error.Error())
		return nil, result.Error
	}
	return marks, nil
}

// GetPdfMarksByPage 根据页码获取PDF标记列表
func (d *PdfMarkDAO) GetPdfMarksByPage(ctx context.Context, paperId string, page int) ([]model.PdfMark, error) {
	var marks []model.PdfMark
	result := d.GetDB(ctx).Where("paper_id = ? AND page = ?", paperId, page).Find(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记列表失败", "paper_id", paperId, "page", page, "error", result.Error.Error())
		return nil, result.Error
	}
	return marks, nil
}

// GetPdfMarkByQuadMd5 根据四边形MD5获取PDF标记
func (d *PdfMarkDAO) GetPdfMarkByQuadMd5(ctx context.Context, quadMd5 string) (*model.PdfMark, error) {
	var mark model.PdfMark
	result := d.GetDB(ctx).Where("quad_md5 = ?", quadMd5).First(&mark)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF标记失败", "quad_md5", quadMd5, "error", result.Error.Error())
		return nil, result.Error
	}
	return &mark, nil
}

// UpdatePdfMark 更新PDF标记
func (d *PdfMarkDAO) UpdatePdfMark(ctx context.Context, mark *model.PdfMark) error {
	return d.Modify(ctx, mark)
}

// DeletePdfMark 删除PDF标记
func (d *PdfMarkDAO) DeletePdfMark(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPdfMarks 列出PDF标记
func (d *PdfMarkDAO) ListPdfMarks(ctx context.Context, limit, offset int) ([]model.PdfMark, error) {
	var marks []model.PdfMark
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF标记失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return marks, nil
}

// CountPdfMarks 获取PDF标记总数
func (d *PdfMarkDAO) CountPdfMarks(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfMark{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF标记总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}

// BatchCreate 批量创建PDF标记
func (d *PdfMarkDAO) BatchCreate(ctx context.Context, marks []*model.PdfMark) error {
	if len(marks) == 0 {
		return nil
	}

	result := d.GetDB(ctx).Create(&marks)
	if result.Error != nil {
		d.logger.Error("msg", "批量创建PDF标记失败", "error", result.Error.Error())
		return result.Error
	}

	return nil
}
