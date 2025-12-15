package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PaperPdfSelectRecordService 论文PDF选择记录服务实现
type PaperPdfSelectRecordService struct {
	paperPdfSelectRecordDAO *dao.PaperPdfSelectRecordDAO
	logger                  logging.Logger
	tracer                  opentracing.Tracer
}

// NewPaperPdfSelectRecordService 创建新的论文PDF选择记录服务
func NewPaperPdfSelectRecordService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperPdfSelectRecordDAO *dao.PaperPdfSelectRecordDAO,
) *PaperPdfSelectRecordService {
	return &PaperPdfSelectRecordService{
		paperPdfSelectRecordDAO: paperPdfSelectRecordDAO,
		logger:                  logger,
		tracer:                  tracer,
	}
}

// CreatePaperPdfSelectRecord 创建论文PDF选择记录
func (s *PaperPdfSelectRecordService) CreatePaperPdfSelectRecord(ctx context.Context, record *model.PaperPdfSelectRecord) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfSelectRecordService.CreatePaperPdfSelectRecord")
	defer span.Finish()

	if err := s.paperPdfSelectRecordDAO.Create(ctx, record); err != nil {
		s.logger.Error("创建论文PDF选择记录失败", "error", err)
		return "0", errors.Biz("pdf.paper_pdf_select_record.errors.create_failed")
	}

	return record.Id, nil
}

// UpdatePaperPdfSelectRecord 更新论文PDF选择记录
func (s *PaperPdfSelectRecordService) UpdatePaperPdfSelectRecord(ctx context.Context, record *model.PaperPdfSelectRecord) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfSelectRecordService.UpdatePaperPdfSelectRecord")
	defer span.Finish()

	// 获取论文PDF选择记录
	existingRecord, err := s.paperPdfSelectRecordDAO.FindById(ctx, record.Id)
	if err != nil {
		s.logger.Error("获取论文PDF选择记录失败", "error", err)
		return false, errors.Biz("pdf.paper_pdf_select_record.errors.get_failed")
	}

	if existingRecord == nil {
		return false, errors.Biz("pdf.paper_pdf_select_record.errors.not_found")
	}

	if err := s.paperPdfSelectRecordDAO.ModifyExcludeNull(ctx, record); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePaperPdfSelectRecordById 删除论文PDF选择记录
func (s *PaperPdfSelectRecordService) DeletePaperPdfSelectRecordById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfSelectRecordService.DeletePaperPdfSelectRecordById")
	defer span.Finish()

	// 删除论文PDF选择记录
	if err := s.paperPdfSelectRecordDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文PDF选择记录失败", "error", err)
		return false, errors.Biz("pdf.paper_pdf_select_record.errors.delete_failed")
	}

	return true, nil
}

// GetPaperPdfSelectRecordById 根据ID获取论文PDF选择记录
func (s *PaperPdfSelectRecordService) GetPaperPdfSelectRecordById(ctx context.Context, id string) (*model.PaperPdfSelectRecord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfSelectRecordService.GetPaperPdfSelectRecordById")
	defer span.Finish()

	// 获取论文PDF选择记录
	record, err := s.paperPdfSelectRecordDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文PDF选择记录失败", "error", err)
		return nil, errors.Biz("pdf.paper_pdf_select_record.errors.get_failed")
	}

	if record == nil {
		return nil, errors.Biz("pdf.paper_pdf_select_record.errors.not_found")
	}

	return record, nil
}
