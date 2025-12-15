package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
)

type DoiMetaInfoService struct {
	logger         logging.Logger
	tracer         opentracing.Tracer
	doiMetaInfoDAO *dao.DoiMetaInfoDAO
}

func NewDoiMetaInfoService(logger logging.Logger, tracer opentracing.Tracer, doiMetaInfoDAO *dao.DoiMetaInfoDAO) *DoiMetaInfoService {
	return &DoiMetaInfoService{
		logger:         logger,
		tracer:         tracer,
		doiMetaInfoDAO: doiMetaInfoDAO,
	}
}

func (s *DoiMetaInfoService) GetById(ctx context.Context, id string) (*model.DoiMetaInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoService.GetById")
	defer span.Finish()
	return s.doiMetaInfoDAO.FindExistById(ctx, id)
}

func (s *DoiMetaInfoService) CreateDoiMetaInfo(ctx context.Context, doiMetaInfo *model.DoiMetaInfo) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoService.CreateDoiMetaInfo")
	defer span.Finish()
	doiMetaInfo.Id = idgen.GenerateUUID()
	return doiMetaInfo.Id, s.doiMetaInfoDAO.SaveExcludeNull(ctx, doiMetaInfo)
}

func (s *DoiMetaInfoService) Update(ctx context.Context, doiMetaInfo *model.DoiMetaInfo) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoService.Update")
	defer span.Finish()
	s.logger.Info("msg", "update doi meta info", "doiMetaInfo", doiMetaInfo)
	return s.doiMetaInfoDAO.ModifyExcludeNull(ctx, doiMetaInfo)
}

// 获取DOI元信息By Doi
func (s *DoiMetaInfoService) GetByDoi(ctx context.Context, doi string) (*model.DoiMetaInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoService.GetByDoi")
	defer span.Finish()
	return s.doiMetaInfoDAO.GetByDoi(ctx, doi)
}

// 获取DOI元信息By PaperId
func (s *DoiMetaInfoService) GetByPaperId(ctx context.Context, paperId string) (*model.DoiMetaInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoService.GetByPaperId")
	defer span.Finish()
	return s.doiMetaInfoDAO.GetByPaperId(ctx, paperId)
}

// 获取DOI元信息某个字段
func (s *DoiMetaInfoService) GetForInitMetaField(ctx context.Context, batchSize int, latestId *string) ([]*model.DoiMetaInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoService.GetForInitMetaField")
	defer span.Finish()
	return s.doiMetaInfoDAO.SelectForInitMetaField(ctx, batchSize, latestId)
}
