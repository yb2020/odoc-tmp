package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
)

// DocClassifyRelationService 文档分类关系服务实现
type DocClassifyRelationService struct {
	docClassifyRelationDAO *dao.DocClassifyRelationDAO
	logger                 logging.Logger
	tracer                 opentracing.Tracer
}

// NewDocClassifyRelationService 创建新的文档分类关系服务
func NewDocClassifyRelationService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	docClassifyRelationDAO *dao.DocClassifyRelationDAO,
) *DocClassifyRelationService {
	return &DocClassifyRelationService{
		logger:                 logger,
		tracer:                 tracer,
		docClassifyRelationDAO: docClassifyRelationDAO,
	}
}

func (s *DocClassifyRelationService) Save(ctx context.Context, docClassifyRelation *model.DocClassifyRelation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DocClassifyRelationService.Save")
	defer span.Finish()

	// 调用DAO层保存文档分类关系
	return s.docClassifyRelationDAO.Save(ctx, docClassifyRelation)
}

// GetDocClassifyRelationsByUserID 根据用户ID获取文档分类关系列表
func (s *DocClassifyRelationService) GetDocClassifyRelationsByUserID(ctx context.Context, userId string) ([]model.DocClassifyRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DocClassifyRelationService.GetDocClassifyRelationsByUserID")
	defer span.Finish()

	// 使用 GORM 直接查询文档分类关系列表
	relations, err := s.docClassifyRelationDAO.GetDocClassifyRelationsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取文档分类关系列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Wrap(err, "获取文档分类关系列表失败")
	}
	// 直接返回模型对象
	return relations, nil
}

func (s *DocClassifyRelationService) DeleteByClassifyId(ctx context.Context, classifyId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DocClassifyRelationService.Delete")
	defer span.Finish()

	// 调用DAO层删除文档分类关系
	return s.docClassifyRelationDAO.DeleteByClassifyId(ctx, classifyId)
}

// 根据用户id，文档id，分类id删除文档分类关系
func (s *DocClassifyRelationService) DeleteByUserIdDocIdClassifyId(ctx context.Context, userId string, docId string, classifyId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DocClassifyRelationService.DeleteByUserIdDocIdClassifyId")
	defer span.Finish()

	// 调用DAO层删除文档分类关系
	return s.docClassifyRelationDAO.DeleteByUserIdDocIdClassifyId(ctx, userId, docId, classifyId)
}

// 根据分类id和文档id查询对应关系数据列表
func (s *DocClassifyRelationService) GetByClassifyIdAndDocId(ctx context.Context, classifyId string, docId string) ([]model.DocClassifyRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DocClassifyRelationService.GetByClassifyIdAndDocId")
	defer span.Finish()

	// 调用DAO层查询文档分类关系
	return s.docClassifyRelationDAO.GetByClassifyIdAndDocId(ctx, classifyId, docId)
}
