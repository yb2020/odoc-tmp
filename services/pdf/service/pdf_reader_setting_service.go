package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfReaderSettingService PDF阅读器设置服务实现
type PdfReaderSettingService struct {
	pdfReaderSettingDAO *dao.PdfReaderSettingDAO
	logger              logging.Logger
	tracer              opentracing.Tracer
}

// NewPdfReaderSettingService 创建新的PDF阅读器设置服务
func NewPdfReaderSettingService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfReaderSettingDAO *dao.PdfReaderSettingDAO,
) *PdfReaderSettingService {
	return &PdfReaderSettingService{
		pdfReaderSettingDAO: pdfReaderSettingDAO,
		logger:              logger,
		tracer:              tracer,
	}
}

// CreatePdfReaderSetting 创建PDF阅读器设置
func (s *PdfReaderSettingService) CreatePdfReaderSetting(ctx context.Context, setting *model.PdfReaderSetting) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfReaderSettingService.CreatePdfReaderSetting")
	defer span.Finish()

	if err := s.pdfReaderSettingDAO.Create(ctx, setting); err != nil {
		s.logger.Error("创建PDF阅读器设置失败", "error", err)
		return "0", errors.Biz("pdf.pdf_reader_setting.errors.create_failed")
	}

	return setting.Id, nil
}

// UpdatePdfReaderSetting 更新PDF阅读器设置
func (s *PdfReaderSettingService) UpdatePdfReaderSetting(ctx context.Context, setting *model.PdfReaderSetting) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfReaderSettingService.UpdatePdfReaderSetting")
	defer span.Finish()

	// 获取PDF阅读器设置
	existingSetting, err := s.pdfReaderSettingDAO.FindById(ctx, setting.Id)
	if err != nil {
		s.logger.Error("获取PDF阅读器设置失败", "error", err)
		return false, errors.Biz("pdf.pdf_reader_setting.errors.get_failed")
	}

	if existingSetting == nil {
		return false, errors.Biz("pdf.pdf_reader_setting.errors.not_found")
	}

	if err := s.pdfReaderSettingDAO.ModifyExcludeNull(ctx, setting); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfReaderSettingById 删除PDF阅读器设置
func (s *PdfReaderSettingService) DeletePdfReaderSettingById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfReaderSettingService.DeletePdfReaderSettingById")
	defer span.Finish()

	// 删除PDF阅读器设置
	if err := s.pdfReaderSettingDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF阅读器设置失败", "error", err)
		return false, errors.Biz("pdf.pdf_reader_setting.errors.delete_failed")
	}

	return true, nil
}

// GetPdfReaderSettingById 根据ID获取PDF阅读器设置
func (s *PdfReaderSettingService) GetPdfReaderSettingById(ctx context.Context, id string) (*model.PdfReaderSetting, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfReaderSettingService.GetPdfReaderSettingById")
	defer span.Finish()

	// 获取PDF阅读器设置
	setting, err := s.pdfReaderSettingDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF阅读器设置失败", "error", err)
		return nil, errors.Biz("pdf.pdf_reader_setting.errors.get_failed")
	}

	if setting == nil {
		return nil, errors.Biz("pdf.pdf_reader_setting.errors.not_found")
	}

	return setting, nil
}

// GetPdfReaderSettingById 根据ID获取PDF阅读器设置
func (s *PdfReaderSettingService) GetByUserId(ctx context.Context, userId string, clientType int) (*model.PdfReaderSetting, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfReaderSettingService.GetPdfReaderSettingById")
	defer span.Finish()

	// 获取PDF阅读器设置
	setting, err := s.pdfReaderSettingDAO.GetPdfReaderSettingUserIdByClientType(ctx, userId, clientType)
	if err != nil {
		s.logger.Error("获取PDF阅读器设置失败", "error", err)
		return nil, errors.Biz("pdf.pdf_reader_setting.errors.get_failed")
	}

	return setting, nil
}
