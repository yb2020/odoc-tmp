package dao

import (
	"context"
	"errors"

	baseDao "github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/model"
	"gorm.io/gorm"
)

// PaperCommentApprovalDAO 提供论文评论赞同状态数据访问功能
type PaperCommentApprovalDAO struct {
	*baseDao.GormBaseDAO[model.PaperCommentApproval]
	logger logging.Logger
}

// NewPaperCommentApprovalDAO 创建一个新的论文评论赞同状态DAO
func NewPaperCommentApprovalDAO(db *gorm.DB, logger logging.Logger) *PaperCommentApprovalDAO {
	return &PaperCommentApprovalDAO{
		GormBaseDAO: baseDao.NewGormBaseDAO[model.PaperCommentApproval](db, logger),
		logger:      logger,
	}
}

// Create 创建论文评论赞同状态
func (d *PaperCommentApprovalDAO) Create(ctx context.Context, approval *model.PaperCommentApproval) error {
	return d.GetDB(ctx).Create(approval).Error
}

// FindById 根据ID获取论文评论赞同状态
func (d *PaperCommentApprovalDAO) FindById(ctx context.Context, id string) (*model.PaperCommentApproval, error) {
	var approval model.PaperCommentApproval
	result := d.GetDB(ctx).Where("id = ?", id).First(&approval)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		d.logger.Error("msg", "获取论文评论赞同状态失败", "id", id, "error", result.Error.Error())
		return nil, result.Error
	}
	return &approval, nil
}

// FindByPaperCommentId 根据论文评论ID获取论文评论赞同状态列表
func (d *PaperCommentApprovalDAO) FindByPaperCommentId(ctx context.Context, paperCommentId string) ([]model.PaperCommentApproval, error) {
	var approvals []model.PaperCommentApproval
	result := d.GetDB(ctx).Where("paper_comment_id = ?", paperCommentId).Find(&approvals)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文评论赞同状态列表失败", "paper_comment_id", paperCommentId, "error", result.Error.Error())
		return nil, result.Error
	}
	return approvals, nil
}

// FindByApprovalStatus 根据赞同状态获取论文评论赞同状态列表
func (d *PaperCommentApprovalDAO) FindByApprovalStatus(ctx context.Context, approvalStatus string) ([]model.PaperCommentApproval, error) {
	var approvals []model.PaperCommentApproval
	result := d.GetDB(ctx).Where("approval_status = ?", approvalStatus).Find(&approvals)
	if result.Error != nil {
		d.logger.Error("msg", "获取论文评论赞同状态列表失败", "approval_status", approvalStatus, "error", result.Error.Error())
		return nil, result.Error
	}
	return approvals, nil
}

// UpdateById 更新论文评论赞同状态
func (d *PaperCommentApprovalDAO) UpdateById(ctx context.Context, approval *model.PaperCommentApproval) error {
	return d.Modify(ctx, approval)
}

// DeleteById 删除论文评论赞同状态
func (d *PaperCommentApprovalDAO) DeleteById(ctx context.Context, id string) error {
	return d.RemoveById(ctx, id)
}

// DeleteByPaperCommentId 根据论文评论ID删除论文评论赞同状态
func (d *PaperCommentApprovalDAO) DeleteByPaperCommentId(ctx context.Context, paperCommentId string) error {
	result := d.GetDB(ctx).Where("paper_comment_id = ?", paperCommentId).Delete(&model.PaperCommentApproval{})
	if result.Error != nil {
		d.logger.Error("msg", "删除论文评论赞同状态失败", "paper_comment_id", paperCommentId, "error", result.Error.Error())
		return result.Error
	}
	return nil
}
