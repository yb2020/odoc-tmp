package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfThumbService PDF缩略图服务实现
type PdfThumbService struct {
	pdfThumbDAO *dao.PdfThumbDAO
	logger      logging.Logger
	tracer      opentracing.Tracer
}

// NewPdfThumbService 创建新的PDF缩略图服务
func NewPdfThumbService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfThumbDAO *dao.PdfThumbDAO,
) *PdfThumbService {
	return &PdfThumbService{
		pdfThumbDAO: pdfThumbDAO,
		logger:      logger,
		tracer:      tracer,
	}
}

// CreatePdfThumb 创建PDF缩略图
func (s *PdfThumbService) CreatePdfThumb(ctx context.Context, thumb *model.PdfThumb) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfThumbService.CreatePdfThumb")
	defer span.Finish()

	if err := s.pdfThumbDAO.Create(ctx, thumb); err != nil {
		s.logger.Error("创建PDF缩略图失败", "error", err)
		return "0", errors.Biz("pdf.pdf_thumb.errors.create_failed")
	}

	return thumb.Id, nil
}

// UpdatePdfThumb 更新PDF缩略图
func (s *PdfThumbService) UpdatePdfThumb(ctx context.Context, thumb *model.PdfThumb) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfThumbService.UpdatePdfThumb")
	defer span.Finish()

	// 获取PDF缩略图
	existingThumb, err := s.pdfThumbDAO.FindById(ctx, thumb.Id)
	if err != nil {
		s.logger.Error("获取PDF缩略图失败", "error", err)
		return false, errors.Biz("pdf.pdf_thumb.errors.get_failed")
	}

	if existingThumb == nil {
		return false, errors.Biz("pdf.pdf_thumb.errors.not_found")
	}

	if err := s.pdfThumbDAO.ModifyExcludeNull(ctx, thumb); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfThumbById 删除PDF缩略图
func (s *PdfThumbService) DeletePdfThumbById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfThumbService.DeletePdfThumbById")
	defer span.Finish()

	// 删除PDF缩略图
	if err := s.pdfThumbDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF缩略图失败", "error", err)
		return false, errors.Biz("pdf.pdf_thumb.errors.delete_failed")
	}

	return true, nil
}

// GetPdfThumbById 根据ID获取PDF缩略图
func (s *PdfThumbService) GetPdfThumbById(ctx context.Context, id string) (*model.PdfThumb, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfThumbService.GetPdfThumbById")
	defer span.Finish()

	// 获取PDF缩略图
	thumb, err := s.pdfThumbDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF缩略图失败", "error", err)
		return nil, errors.Biz("pdf.pdf_thumb.errors.get_failed")
	}

	if thumb == nil {
		return nil, errors.Biz("pdf.pdf_thumb.errors.not_found")
	}

	return thumb, nil
}
