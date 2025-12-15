package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperCommentService 论文评论服务实现
type PaperCommentService struct {
	paperCommentDAO *dao.PaperCommentDAO
	logger          logging.Logger
	tracer          opentracing.Tracer
}

// NewPaperCommentService 创建新的论文评论服务
func NewPaperCommentService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperCommentDAO *dao.PaperCommentDAO,
) *PaperCommentService {
	return &PaperCommentService{
		logger:          logger,
		tracer:          tracer,
		paperCommentDAO: paperCommentDAO,
	}
}

// CreatePaperComment 创建论文评论
func (s *PaperCommentService) CreatePaperComment(ctx context.Context, comment *model.PaperComment) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.CreatePaperComment")
	defer span.Finish()

	if err := s.paperCommentDAO.Create(ctx, comment); err != nil {
		s.logger.Error("创建论文评论失败", "error", err)
		return "0", errors.Biz("paper.comment.errors.create_failed")
	}

	return comment.Id, nil
}

// GetPaperCommentById 根据ID获取论文评论
func (s *PaperCommentService) GetPaperCommentById(ctx context.Context, id string) (*model.PaperComment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.GetPaperCommentById")
	defer span.Finish()

	comment, err := s.paperCommentDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文评论失败", "id", id, "error", err)
		return nil, errors.Biz("paper.comment.errors.get_failed")
	}

	if comment == nil {
		return nil, errors.Biz("paper.comment.errors.not_found")
	}

	return comment, nil
}

// GetPaperCommentsByPaperId 根据论文ID获取论文评论列表
func (s *PaperCommentService) GetPaperCommentsByPaperId(ctx context.Context, paperId string) ([]*model.PaperComment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.GetPaperCommentsByPaperId")
	defer span.Finish()

	comments, err := s.paperCommentDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文评论列表失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.comment.errors.list_failed")
	}

	result := make([]*model.PaperComment, 0, len(comments))
	for i := range comments {
		result = append(result, &comments[i])
	}

	return result, nil
}

// GetPaperCommentsByCommentLevel 根据评论等级获取论文评论列表
func (s *PaperCommentService) GetPaperCommentsByCommentLevel(ctx context.Context, commentLevel string) ([]*model.PaperComment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.GetPaperCommentsByCommentLevel")
	defer span.Finish()

	comments, err := s.paperCommentDAO.FindByCommentLevel(ctx, commentLevel)
	if err != nil {
		s.logger.Error("获取论文评论列表失败", "comment_level", commentLevel, "error", err)
		return nil, errors.Biz("paper.comment.errors.list_failed")
	}

	result := make([]*model.PaperComment, 0, len(comments))
	for i := range comments {
		result = append(result, &comments[i])
	}

	return result, nil
}

// UpdatePaperComment 更新论文评论
func (s *PaperCommentService) UpdatePaperComment(ctx context.Context, comment *model.PaperComment) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.UpdatePaperComment")
	defer span.Finish()

	// 获取论文评论
	existingComment, err := s.paperCommentDAO.FindById(ctx, comment.Id)
	if err != nil {
		s.logger.Error("获取论文评论失败", "id", comment.Id, "error", err)
		return false, errors.Biz("paper.comment.errors.get_failed")
	}

	if existingComment == nil {
		return false, errors.Biz("paper.comment.errors.not_found")
	}

	// 更新论文评论
	if err := s.paperCommentDAO.UpdateById(ctx, comment); err != nil {
		s.logger.Error("更新论文评论失败", "id", comment.Id, "error", err)
		return false, errors.Biz("paper.comment.errors.update_failed")
	}

	return true, nil
}

// DeletePaperComment 删除论文评论
func (s *PaperCommentService) DeletePaperComment(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.DeletePaperComment")
	defer span.Finish()

	// 删除论文评论
	if err := s.paperCommentDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文评论失败", "id", id, "error", err)
		return false, errors.Biz("paper.comment.errors.delete_failed")
	}

	return true, nil
}

// DeletePaperCommentsByPaperId 根据论文ID删除论文评论
func (s *PaperCommentService) DeletePaperCommentsByPaperId(ctx context.Context, paperId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperCommentService.DeletePaperCommentsByPaperId")
	defer span.Finish()

	// 删除论文评论
	if err := s.paperCommentDAO.DeleteByPaperId(ctx, paperId); err != nil {
		s.logger.Error("删除论文评论失败", "paper_id", paperId, "error", err)
		return false, errors.Biz("paper.comment.errors.delete_failed")
	}

	return true, nil
}
