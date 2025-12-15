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

// PdfCommentDAO GORM实现的PDF评论DAO
type PdfCommentDAO struct {
	*baseDao.GormBaseDAO[model.PdfComment]
	logger logging.Logger
}

// NewPdfCommentDAO 创建一个新的PDF评论DAO
func NewPdfCommentDAO(db *gorm.DB, logger logging.Logger) *PdfCommentDAO {
	return &PdfCommentDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PdfComment](db, logger),
		logger:      logger,
	}
}

// Create 创建PDF评论
func (d *PdfCommentDAO) Create(ctx context.Context, comment *model.PdfComment) error {
	return d.GetDB(ctx).Create(comment).Error
}

// GetPdfCommentByID 根据ID获取PDF评论
func (d *PdfCommentDAO) GetPdfCommentByID(ctx context.Context, id string) (*model.PdfComment, error) {
	var comment model.PdfComment
	result := d.GetDB(ctx).Where("id = ?", id).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取PDF评论失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &comment, nil
}

// GetPdfCommentsByPaperID 根据论文ID获取PDF评论列表
func (d *PdfCommentDAO) GetPdfCommentsByPaperID(ctx context.Context, paperID string) ([]model.PdfComment, error) {
	var comments []model.PdfComment
	result := d.GetDB(ctx).Where("paper_id = ?", paperID).Find(&comments)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF评论列表失败", "paper_id", paperID, "error", result.Error.Error())
		return nil, result.Error
	}
	return comments, nil
}

// GetPdfCommentsByParentID 根据父评论ID获取PDF评论列表
func (d *PdfCommentDAO) GetPdfCommentsByParentID(ctx context.Context, parentID string) ([]model.PdfComment, error) {
	var comments []model.PdfComment
	result := d.GetDB(ctx).Where("parent_id = ?", parentID).Find(&comments)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF评论列表失败", "parent_id", parentID, "error", result.Error.Error())
		return nil, result.Error
	}
	return comments, nil
}

// UpdatePdfComment 更新PDF评论
func (d *PdfCommentDAO) UpdatePdfComment(ctx context.Context, comment *model.PdfComment) error {
	return d.Modify(ctx, comment)
}

// DeletePdfComment 删除PDF评论
func (d *PdfCommentDAO) DeletePdfComment(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// ListPdfComments 列出PDF评论
func (d *PdfCommentDAO) ListPdfComments(ctx context.Context, limit, offset int) ([]model.PdfComment, error) {
	var comments []model.PdfComment
	result := d.GetDB(ctx).Offset(offset).Limit(limit).Order("id DESC").Find(&comments)
	if result.Error != nil {
		d.logger.Error("msg", "列出PDF评论失败", "error", result.Error.Error())
		return nil, result.Error
	}
	return comments, nil
}

// CountPdfComments 获取PDF评论总数
func (d *PdfCommentDAO) CountPdfComments(ctx context.Context) (string, error) {
	var count int64
	result := d.GetDB(ctx).Model(&model.PdfComment{}).Count(&count)
	if result.Error != nil {
		d.logger.Error("msg", "获取PDF评论总数失败", "error", result.Error.Error())
		return "0", result.Error
	}
	return fmt.Sprintf("%d", count), nil
}
