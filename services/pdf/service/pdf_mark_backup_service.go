package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfMarkBackupService PDF标记备份服务实现
type PdfMarkBackupService struct {
	pdfMarkBackupDAO *dao.PdfMarkBackupDAO
	logger           logging.Logger
	tracer           opentracing.Tracer
}

// NewPdfMarkBackupService 创建新的PDF标记备份服务
func NewPdfMarkBackupService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfMarkBackupDAO *dao.PdfMarkBackupDAO,
) *PdfMarkBackupService {
	return &PdfMarkBackupService{
		pdfMarkBackupDAO: pdfMarkBackupDAO,
		logger:           logger,
		tracer:           tracer,
	}
}

// CreatePdfMarkBackup 创建PDF标记备份
func (s *PdfMarkBackupService) CreatePdfMarkBackup(ctx context.Context, backup *model.PdfMarkBackup) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkBackupService.CreatePdfMarkBackup")
	defer span.Finish()

	if err := s.pdfMarkBackupDAO.Create(ctx, backup); err != nil {
		s.logger.Error("创建PDF标记备份失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark_backup.errors.create_failed")
	}

	return backup.Id, nil
}

// UpdatePdfMarkBackup 更新PDF标记备份
func (s *PdfMarkBackupService) UpdatePdfMarkBackup(ctx context.Context, backup *model.PdfMarkBackup) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkBackupService.UpdatePdfMarkBackup")
	defer span.Finish()

	// 获取PDF标记备份
	existingBackup, err := s.pdfMarkBackupDAO.FindById(ctx, backup.Id)
	if err != nil {
		s.logger.Error("获取PDF标记备份失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_backup.errors.get_failed")
	}

	if existingBackup == nil {
		return false, errors.Biz("pdf.pdf_mark_backup.errors.not_found")
	}

	if err := s.pdfMarkBackupDAO.ModifyExcludeNull(ctx, backup); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfMarkBackupById 删除PDF标记备份
func (s *PdfMarkBackupService) DeletePdfMarkBackupById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkBackupService.DeletePdfMarkBackupById")
	defer span.Finish()

	// 删除PDF标记备份
	if err := s.pdfMarkBackupDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF标记备份失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_backup.errors.delete_failed")
	}

	return true, nil
}

// GetPdfMarkBackupById 根据ID获取PDF标记备份
func (s *PdfMarkBackupService) GetPdfMarkBackupById(ctx context.Context, id string) (*model.PdfMarkBackup, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkBackupService.GetPdfMarkBackupById")
	defer span.Finish()

	// 获取PDF标记备份
	backup, err := s.pdfMarkBackupDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF标记备份失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_backup.errors.get_failed")
	}

	if backup == nil {
		return nil, errors.Biz("pdf.pdf_mark_backup.errors.not_found")
	}

	return backup, nil
}
