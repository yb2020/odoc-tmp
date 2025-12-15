package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperPdfParsedService 提供PaperPdfParsed服务
type PaperPdfParsedService struct {
	logger            logging.Logger
	tracer            opentracing.Tracer
	paperPdfParsedDAO *dao.PaperPdfParsedDAO
}

// NewPaperPdfParsedService 创建一个新的PaperPdfParsed服务
func NewPaperPdfParsedService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperPdfParsedDAO *dao.PaperPdfParsedDAO,
) *PaperPdfParsedService {
	return &PaperPdfParsedService{
		logger:            logger,
		tracer:            tracer,
		paperPdfParsedDAO: paperPdfParsedDAO,
	}
}

// Save 保存PaperPdfParsed
func (s *PaperPdfParsedService) Save(ctx context.Context, record *model.PaperPdfParsed) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfParsedService.Save")
	defer span.Finish()

	return s.paperPdfParsedDAO.Save(ctx, record)
}

// BatchSave 批量保存PaperPdfParsed
func (s *PaperPdfParsedService) BatchSave(ctx context.Context, records []*model.PaperPdfParsed) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfParsedService.BatchSave")
	defer span.Finish()

	return s.paperPdfParsedDAO.BatchSave(ctx, records)
}

// 根据原pdf文件的fileSHA256查询paper_pdf_parsed记录列表
func (s *PaperPdfParsedService) GetBySourcePdfFileSHA256AndVersion(ctx context.Context, sha256, version string) ([]*model.PaperPdfParsed, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfParsedService.GetBySourcePdfFileSHA256AndVersion")
	defer span.Finish()

	return s.paperPdfParsedDAO.GetBySourcePdfFileSHA256AndVersion(ctx, sha256, version)
}

func (s *PaperPdfParsedService) GetBySourcePdfFileSHA256AndTypeAndVersion(ctx context.Context, sha256, fileType, version string) (*model.PaperPdfParsed, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfParsedService.GetBySourcePdfFileSHA256AndTypeAndVersion")
	defer span.Finish()

	return s.paperPdfParsedDAO.GetBySourcePdfFileSHA256AndTypeAndVersion(ctx, sha256, fileType, version)
}

// TODO: 增加缓存，减少查库
func (s *PaperPdfParsedService) HasExistBySourcePdfFileSHA256AndVersion(ctx context.Context, sha256, version string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperPdfParsedService.HasExistBySourcePdfFileSHA256AndVersion")
	defer span.Finish()

	return s.paperPdfParsedDAO.HasExistBySourcePdfFileSHA256AndVersion(ctx, sha256, version)
}
