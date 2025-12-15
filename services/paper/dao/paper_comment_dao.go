package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperCommentDAO 提供论文评论数据访问功能
type PaperCommentDAO struct {
	*baseDao.GormBaseDAO[model.PaperComment]
	logger logging.Logger
}

// NewPaperCommentDAO 创建一个新的论文评论DAO
func NewPaperCommentDAO(db *gorm.DB, logger logging.Logger) *PaperCommentDAO {
	return &PaperCommentDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperComment](db, logger),
		logger:      logger,
	}
}

// Create 创建论文评论
func (d *PaperCommentDAO) Create(ctx context.Context, comment *model.PaperComment) error {
	return d.GetDB(ctx).Create(comment).Error
}

// FindById 根据ID获取论文评论
func (d *PaperCommentDAO) FindById(ctx context.Context, id string) (*model.PaperComment, error) {
	var comment model.PaperComment
	result := d.GetDB(ctx).Where("id = ?", id).First(&comment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文评论失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &comment, nil
}

// FindByPaperId 根据论文ID获取论文评论列表
func (d *PaperCommentDAO) FindByPaperId(ctx context.Context, paperId string) ([]model.PaperComment, error) {
	var comments []model.PaperComment
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Find(&comments)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文评论列表失败", "paper_id", paperId, "error", result.Error.Error())
		return nil, result.Error
	}
	return comments, nil
}

// FindByCommentLevel 根据评论等级获取论文评论列表
func (d *PaperCommentDAO) FindByCommentLevel(ctx context.Context, commentLevel string) ([]model.PaperComment, error) {
	var comments []model.PaperComment
	result := d.GetDB(ctx).Where("comment_level = ?", commentLevel).Find(&comments)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文评论列表失败", "comment_level", commentLevel, "error", result.Error.Error())
		return nil, result.Error
	}
	return comments, nil
}

// UpdateById 更新论文评论
func (d *PaperCommentDAO) UpdateById(ctx context.Context, comment *model.PaperComment) error {
	return d.Modify(ctx, comment)
}

// DeleteById 删除论文评论
func (d *PaperCommentDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByPaperId 根据论文ID删除论文评论
func (d *PaperCommentDAO) DeleteByPaperId(ctx context.Context, paperId string) error {
	result := d.GetDB(ctx).Where("paper_id = ?", paperId).Delete(&model.PaperComment{})
	if result.Error != nil {
		d.logger.Error("msg", "删除论文评论失败", "paper_id", paperId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
