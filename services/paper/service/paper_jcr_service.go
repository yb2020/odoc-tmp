package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperJcrService 论文JCR服务
type PaperJcrService struct {
	paperJcrDAO *dao.PaperJcrDAO
	logger      logging.Logger
	tracer      opentracing.Tracer
}

// NewPaperJcrService 创建一个新的论文JCR服务
func NewPaperJcrService(logger logging.Logger, tracer opentracing.Tracer, paperJcrDAO *dao.PaperJcrDAO) *PaperJcrService {
	return &PaperJcrService{
		paperJcrDAO: paperJcrDAO,
		logger:      logger,
		tracer:      tracer,
	}
}

// GetPaperJcrEntityByVenue 根据venue查找PaperJcrEntity
func (s *PaperJcrService) GetPaperJcrEntityByVenue(ctx context.Context, venue string) (*model.PaperJcrEntity, error) {
	if venue == "" {
		return nil, nil
	}
	return s.paperJcrDAO.GetPaperJcrEntityByVenue(ctx, venue)
}
