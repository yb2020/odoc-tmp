package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfAnnotationService PDF注释服务实现
type PdfAnnotationService struct {
	pdfAnnotationDAO *dao.PdfAnnotationDAO
	logger           logging.Logger
	tracer           opentracing.Tracer
}

// NewPdfAnnotationService 创建新的PDF注释服务
func NewPdfAnnotationService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfAnnotationDAO *dao.PdfAnnotationDAO,
) *PdfAnnotationService {
	return &PdfAnnotationService{
		pdfAnnotationDAO: pdfAnnotationDAO,
		logger:           logger,
		tracer:           tracer,
	}
}

// CreatePdfAnnotation 创建PDF注释
func (s *PdfAnnotationService) CreatePdfAnnotation(ctx context.Context, annotation *model.PdfAnnotation) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfAnnotationService.CreatePdfAnnotation")
	defer span.Finish()

	if err := s.pdfAnnotationDAO.Create(ctx, annotation); err != nil {
		s.logger.Error("创建PDF注释失败", "error", err)
		return "0", errors.Biz("pdf.pdf_annotation.errors.create_failed")
	}

	return annotation.Id, nil
}

// UpdatePdfAnnotation 更新PDF注释
func (s *PdfAnnotationService) UpdatePdfAnnotation(ctx context.Context, annotation *model.PdfAnnotation) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfAnnotationService.UpdatePdfAnnotation")
	defer span.Finish()

	// 获取PDF注释
	existingAnnotation, err := s.pdfAnnotationDAO.FindById(ctx, annotation.Id)
	if err != nil {
		s.logger.Error("获取PDF注释失败", "error", err)
		return false, errors.Biz("pdf.pdf_annotation.errors.get_failed")
	}

	if existingAnnotation == nil {
		return false, errors.Biz("pdf.pdf_annotation.errors.not_found")
	}

	if err := s.pdfAnnotationDAO.ModifyExcludeNull(ctx, annotation); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfAnnotationById 删除PDF注释
func (s *PdfAnnotationService) DeletePdfAnnotationById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfAnnotationService.DeletePdfAnnotationById")
	defer span.Finish()

	// 删除PDF注释
	if err := s.pdfAnnotationDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF注释失败", "error", err)
		return false, errors.Biz("pdf.pdf_annotation.errors.delete_failed")
	}

	return true, nil
}

// GetPdfAnnotationById 根据ID获取PDF注释
func (s *PdfAnnotationService) GetPdfAnnotationById(ctx context.Context, id string) (*model.PdfAnnotation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfAnnotationService.GetPdfAnnotationById")
	defer span.Finish()

	// 获取PDF注释
	annotation, err := s.pdfAnnotationDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF注释失败", "error", err)
		return nil, errors.Biz("pdf.pdf_annotation.errors.get_failed")
	}

	if annotation == nil {
		return nil, errors.Biz("pdf.pdf_annotation.errors.not_found")
	}

	return annotation, nil
}
