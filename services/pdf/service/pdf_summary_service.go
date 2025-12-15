package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfSummaryService PDF摘要服务实现
type PdfSummaryService struct {
	pdfSummaryDAO *dao.PdfSummaryDAO
	logger        logging.Logger
	tracer        opentracing.Tracer
}

// NewPdfSummaryService 创建新的PDF摘要服务
func NewPdfSummaryService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfSummaryDAO *dao.PdfSummaryDAO,
) *PdfSummaryService {
	return &PdfSummaryService{
		pdfSummaryDAO: pdfSummaryDAO,
		logger:        logger,
		tracer:        tracer,
	}
}

// CreatePdfSummary 创建PDF摘要
func (s *PdfSummaryService) CreatePdfSummary(ctx context.Context, summary *model.PdfSummary) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfSummaryService.CreatePdfSummary")
	defer span.Finish()

	if err := s.pdfSummaryDAO.Save(ctx, summary); err != nil {
		s.logger.Error("创建PDF摘要失败", "error", err)
		return "0", errors.Biz("pdf.pdf_summary.errors.create_failed")
	}

	return summary.Id, nil
}

// UpdatePdfSummary 更新PDF摘要
func (s *PdfSummaryService) UpdatePdfSummary(ctx context.Context, summary *model.PdfSummary) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfSummaryService.UpdatePdfSummary")
	defer span.Finish()

	// 获取PDF摘要
	existingSummary, err := s.pdfSummaryDAO.FindById(ctx, summary.Id)
	if err != nil {
		s.logger.Error("获取PDF摘要失败", "error", err)
		return false, errors.Biz("pdf.pdf_summary.errors.get_failed")
	}

	if existingSummary == nil {
		return false, errors.Biz("pdf.pdf_summary.errors.not_found")
	}

	if err := s.pdfSummaryDAO.ModifyExcludeNull(ctx, summary); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfSummaryById 删除PDF摘要
func (s *PdfSummaryService) DeletePdfSummaryById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfSummaryService.DeletePdfSummaryById")
	defer span.Finish()

	// 删除PDF摘要
	if err := s.pdfSummaryDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF摘要失败", "error", err)
		return false, errors.Biz("pdf.pdf_summary.errors.delete_failed")
	}

	return true, nil
}

// GetBySourceFromAndFileSHA256AndVersionAndLang 根据source_from,file_sha256,version,lang获取PDF摘要
func (s *PdfSummaryService) GetBySourceFromAndFileSHA256AndVersionAndLang(ctx context.Context, sourceFrom, fileSHA256, version, lang string) (*model.PdfSummary, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfSummaryService.GetBySourceFromAndFileSHA256AndVersionAndLang")
	defer span.Finish()

	// 获取PDF摘要
	return s.pdfSummaryDAO.FindExistBySourceFromAndFileSHA256AndVersionAndLang(ctx, sourceFrom, fileSHA256, version, lang)
}
