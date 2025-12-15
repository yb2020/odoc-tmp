package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfCommentService PDF评论服务实现
type PdfCommentService struct {
	pdfCommentDAO *dao.PdfCommentDAO
	logger        logging.Logger
	tracer        opentracing.Tracer
}

// NewPdfCommentService 创建新的PDF评论服务
func NewPdfCommentService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfCommentDAO *dao.PdfCommentDAO,
) *PdfCommentService {
	return &PdfCommentService{
		pdfCommentDAO: pdfCommentDAO,
		logger:        logger,
		tracer:        tracer,
	}
}

// CreatePdfComment 创建PDF评论
func (s *PdfCommentService) CreatePdfComment(ctx context.Context, comment *model.PdfComment) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfCommentService.CreatePdfComment")
	defer span.Finish()

	if err := s.pdfCommentDAO.Create(ctx, comment); err != nil {
		s.logger.Error("创建PDF评论失败", "error", err)
		return "0", errors.Biz("pdf.pdf_comment.errors.create_failed")
	}

	return comment.Id, nil
}

// UpdatePdfComment 更新PDF评论
func (s *PdfCommentService) UpdatePdfComment(ctx context.Context, comment *model.PdfComment) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfCommentService.UpdatePdfComment")
	defer span.Finish()

	// 获取PDF评论
	existingComment, err := s.pdfCommentDAO.FindById(ctx, comment.Id)
	if err != nil {
		s.logger.Error("获取PDF评论失败", "error", err)
		return false, errors.Biz("pdf.pdf_comment.errors.get_failed")
	}

	if existingComment == nil {
		return false, errors.Biz("pdf.pdf_comment.errors.not_found")
	}

	if err := s.pdfCommentDAO.ModifyExcludeNull(ctx, comment); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfCommentById 删除PDF评论
func (s *PdfCommentService) DeletePdfCommentById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfCommentService.DeletePdfCommentById")
	defer span.Finish()

	// 删除PDF评论
	if err := s.pdfCommentDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF评论失败", "error", err)
		return false, errors.Biz("pdf.pdf_comment.errors.delete_failed")
	}

	return true, nil
}

// GetPdfCommentById 根据ID获取PDF评论
func (s *PdfCommentService) GetPdfCommentById(ctx context.Context, id string) (*model.PdfComment, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfCommentService.GetPdfCommentById")
	defer span.Finish()

	// 获取PDF评论
	comment, err := s.pdfCommentDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF评论失败", "error", err)
		return nil, errors.Biz("pdf.pdf_comment.errors.get_failed")
	}

	if comment == nil {
		return nil, errors.Biz("pdf.pdf_comment.errors.not_found")
	}

	return comment, nil
}
