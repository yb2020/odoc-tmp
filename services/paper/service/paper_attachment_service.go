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

// PaperAttachmentService 论文附件服务实现
type PaperAttachmentService struct {
	paperAttachmentDAO *dao.PaperAttachmentDAO
	logger             logging.Logger
	tracer             opentracing.Tracer
}

// NewPaperAttachmentService 创建新的论文附件服务
func NewPaperAttachmentService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperAttachmentDAO *dao.PaperAttachmentDAO,
) *PaperAttachmentService {
	return &PaperAttachmentService{
		logger:             logger,
		tracer:             tracer,
		paperAttachmentDAO: paperAttachmentDAO,
	}
}

// CreatePaperAttachment 创建论文附件
func (s *PaperAttachmentService) CreatePaperAttachment(ctx context.Context, attachment *model.PaperAttachment) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.CreatePaperAttachment")
	defer span.Finish()

	if err := s.paperAttachmentDAO.Create(ctx, attachment); err != nil {
		s.logger.Error("创建论文附件失败", "error", err)
		return "0", errors.Biz("paper.attachment.errors.create_failed")
	}

	return attachment.Id, nil
}

// GetPaperAttachmentById 根据ID获取论文附件
func (s *PaperAttachmentService) GetPaperAttachmentById(ctx context.Context, id string) (*model.PaperAttachment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.GetPaperAttachmentById")
	defer span.Finish()

	attachment, err := s.paperAttachmentDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文附件失败", "id", id, "error", err)
		return nil, errors.Biz("paper.attachment.errors.get_failed")
	}

	if attachment == nil {
		return nil, errors.Biz("paper.attachment.errors.not_found")
	}

	return attachment, nil
}

// GetPaperAttachmentsByPaperId 根据论文ID获取论文附件列表
func (s *PaperAttachmentService) GetPaperAttachmentsByPaperId(ctx context.Context, paperId string) ([]*model.PaperAttachment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.GetPaperAttachmentsByPaperId")
	defer span.Finish()

	attachments, err := s.paperAttachmentDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文附件列表失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.attachment.errors.list_failed")
	}

	result := make([]*model.PaperAttachment, 0, len(attachments))
	for i := range attachments {
		result = append(result, &attachments[i])
	}

	return result, nil
}

// GetPaperAttachmentsByType 根据附件类型获取论文附件列表
func (s *PaperAttachmentService) GetPaperAttachmentsByType(ctx context.Context, attachType constants.PaperAttachType) ([]*model.PaperAttachment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.GetPaperAttachmentsByType")
	defer span.Finish()

	attachments, err := s.paperAttachmentDAO.FindByType(ctx, attachType.Value())
	if err != nil {
		s.logger.Error("获取论文附件列表失败", "type", attachType.String(), "error", err)
		return nil, errors.Biz("paper.attachment.errors.list_failed")
	}

	result := make([]*model.PaperAttachment, 0, len(attachments))
	for i := range attachments {
		result = append(result, &attachments[i])
	}

	return result, nil
}

// UpdatePaperAttachment 更新论文附件
func (s *PaperAttachmentService) UpdatePaperAttachment(ctx context.Context, attachment *model.PaperAttachment) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.UpdatePaperAttachment")
	defer span.Finish()

	// 获取论文附件
	existingAttachment, err := s.paperAttachmentDAO.FindById(ctx, attachment.Id)
	if err != nil {
		s.logger.Error("获取论文附件失败", "id", attachment.Id, "error", err)
		return false, errors.Biz("paper.attachment.errors.get_failed")
	}

	if existingAttachment == nil {
		return false, errors.Biz("paper.attachment.errors.not_found")
	}

	// 更新论文附件
	if err := s.paperAttachmentDAO.UpdateById(ctx, attachment); err != nil {
		s.logger.Error("更新论文附件失败", "id", attachment.Id, "error", err)
		return false, errors.Biz("paper.attachment.errors.update_failed")
	}

	return true, nil
}

// DeletePaperAttachment 删除论文附件
func (s *PaperAttachmentService) DeletePaperAttachment(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.DeletePaperAttachment")
	defer span.Finish()

	// 删除论文附件
	if err := s.paperAttachmentDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文附件失败", "id", id, "error", err)
		return false, errors.Biz("paper.attachment.errors.delete_failed")
	}

	return true, nil
}

// DeletePaperAttachmentsByPaperId 根据论文ID删除论文附件
func (s *PaperAttachmentService) DeletePaperAttachmentsByPaperId(ctx context.Context, paperId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAttachmentService.DeletePaperAttachmentsByPaperId")
	defer span.Finish()

	// 删除论文附件
	if err := s.paperAttachmentDAO.DeleteByPaperId(ctx, paperId); err != nil {
		s.logger.Error("删除论文附件失败", "paper_id", paperId, "error", err)
		return false, errors.Biz("paper.attachment.errors.delete_failed")
	}

	return true, nil
}
