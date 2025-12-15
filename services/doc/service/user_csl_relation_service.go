package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
)

// UserCslRelationService 用户引用样式关联服务
type UserCslRelationService struct {
	dao    *dao.UserCslRelationDAO
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewUserCslRelationService 创建用户引用样式关联服务
func NewUserCslRelationService(dao *dao.UserCslRelationDAO, logger logging.Logger, tracer opentracing.Tracer) *UserCslRelationService {
	return &UserCslRelationService{
		dao:    dao,
		logger: logger,
		tracer: tracer,
	}
}

// IsAddedCsl 检查用户是否已添加引用样式
func (s *UserCslRelationService) IsAddedCsl(ctx context.Context, userId string) bool {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserCslRelationService.IsAddedCsl")
	defer span.Finish()

	relation, err := s.dao.IsAddedCsl(ctx, userId)
	if err != nil {
		s.logger.Error("check user csl relation failed", "error", err.Error())
		return false
	}
	return relation != nil
}

// Save 保存用户引用样式关联
func (s *UserCslRelationService) Save(ctx context.Context, relation *model.UserCslRelation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserCslRelationService.Save")
	defer span.Finish()

	return s.dao.Save(ctx, relation)
}
