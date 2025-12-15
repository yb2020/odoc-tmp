package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/constants"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperCommentApprovalService 论文评论赞同状态服务实现
type PaperCommentApprovalService struct {
	paperCommentApprovalDAO *dao.PaperCommentApprovalDAO
	logger                  logging.Logger
	tracer                  opentracing.Tracer
}

// NewPaperCommentApprovalService 创建新的论文评论赞同状态服务
func NewPaperCommentApprovalService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperCommentApprovalDAO *dao.PaperCommentApprovalDAO,
) *PaperCommentApprovalService {
	return &PaperCommentApprovalService{
		logger:                  logger,
		tracer:                  tracer,
		paperCommentApprovalDAO: paperCommentApprovalDAO,
	}
}

// CreatePaperCommentApproval 创建论文评论赞同状态
func (s *PaperCommentApprovalService) CreatePaperCommentApproval(ctx context.Context, approval *model.PaperCommentApproval) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.CreatePaperCommentApproval")
	defer span.Finish()

	if err := s.paperCommentApprovalDAO.Create(ctx, approval); err != nil {
		s.logger.Error("创建论文评论赞同状态失败", "error", err)
		return "0", errors.Biz("paper.comment_approval.errors.create_failed")
	}

	return approval.Id, nil
}

// GetPaperCommentApprovalById 根据ID获取论文评论赞同状态
func (s *PaperCommentApprovalService) GetPaperCommentApprovalById(ctx context.Context, id string) (*model.PaperCommentApproval, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.GetPaperCommentApprovalById")
	defer span.Finish()

	approval, err := s.paperCommentApprovalDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文评论赞同状态失败", "id", id, "error", err)
		return nil, errors.Biz("paper.comment_approval.errors.get_failed")
	}

	if approval == nil {
		return nil, errors.Biz("paper.comment_approval.errors.not_found")
	}

	return approval, nil
}

// GetPaperCommentApprovalsByPaperCommentId 根据论文评论ID获取论文评论赞同状态列表
func (s *PaperCommentApprovalService) GetPaperCommentApprovalsByPaperCommentId(ctx context.Context, paperCommentId string) ([]*model.PaperCommentApproval, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.GetPaperCommentApprovalsByPaperCommentId")
	defer span.Finish()

	approvals, err := s.paperCommentApprovalDAO.FindByPaperCommentId(ctx, paperCommentId)
	if err != nil {
		s.logger.Error("获取论文评论赞同状态列表失败", "paper_comment_id", paperCommentId, "error", err)
		return nil, errors.Biz("paper.comment_approval.errors.list_failed")
	}

	result := make([]*model.PaperCommentApproval, 0, len(approvals))
	for i := range approvals {
		result = append(result, &approvals[i])
	}

	return result, nil
}

// GetPaperCommentApprovalsByApprovalStatus 根据赞同状态获取论文评论赞同状态列表
func (s *PaperCommentApprovalService) GetPaperCommentApprovalsByApprovalStatus(ctx context.Context, approvalStatus constants.ApprovalStatus) ([]*model.PaperCommentApproval, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.GetPaperCommentApprovalsByApprovalStatus")
	defer span.Finish()

	approvals, err := s.paperCommentApprovalDAO.FindByApprovalStatus(ctx, approvalStatus.String())
	if err != nil {
		s.logger.Error("获取论文评论赞同状态列表失败", "approval_status", approvalStatus.String(), "error", err)
		return nil, errors.Biz("paper.comment_approval.errors.list_failed")
	}

	result := make([]*model.PaperCommentApproval, 0, len(approvals))
	for i := range approvals {
		result = append(result, &approvals[i])
	}

	return result, nil
}

// UpdatePaperCommentApproval 更新论文评论赞同状态
func (s *PaperCommentApprovalService) UpdatePaperCommentApproval(ctx context.Context, approval *model.PaperCommentApproval) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.UpdatePaperCommentApproval")
	defer span.Finish()

	// 获取论文评论赞同状态
	existingApproval, err := s.paperCommentApprovalDAO.FindById(ctx, approval.Id)
	if err != nil {
		s.logger.Error("获取论文评论赞同状态失败", "id", approval.Id, "error", err)
		return false, errors.Biz("paper.comment_approval.errors.get_failed")
	}

	if existingApproval == nil {
		return false, errors.Biz("paper.comment_approval.errors.not_found")
	}

	// 更新论文评论赞同状态
	if err := s.paperCommentApprovalDAO.UpdateById(ctx, approval); err != nil {
		s.logger.Error("更新论文评论赞同状态失败", "id", approval.Id, "error", err)
		return false, errors.Biz("paper.comment_approval.errors.update_failed")
	}

	return true, nil
}

// DeletePaperCommentApproval 删除论文评论赞同状态
func (s *PaperCommentApprovalService) DeletePaperCommentApproval(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.DeletePaperCommentApproval")
	defer span.Finish()

	// 删除论文评论赞同状态
	if err := s.paperCommentApprovalDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文评论赞同状态失败", "id", id, "error", err)
		return false, errors.Biz("paper.comment_approval.errors.delete_failed")
	}

	return true, nil
}

// DeletePaperCommentApprovalsByPaperCommentId 根据论文评论ID删除论文评论赞同状态
func (s *PaperCommentApprovalService) DeletePaperCommentApprovalsByPaperCommentId(ctx context.Context, paperCommentId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentApprovalService.DeletePaperCommentApprovalsByPaperCommentId")
	defer span.Finish()

	// 删除论文评论赞同状态
	if err := s.paperCommentApprovalDAO.DeleteByPaperCommentId(ctx, paperCommentId); err != nil {
		s.logger.Error("删除论文评论赞同状态失败", "paper_comment_id", paperCommentId, "error", err)
		return false, errors.Biz("paper.comment_approval.errors.delete_failed")
	}

	return true, nil
}
