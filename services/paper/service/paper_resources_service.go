package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperResourcesService 论文资源服务实现
type PaperResourcesService struct {
	paperResourcesDAO *dao.PaperResourcesDAO
	logger            logging.Logger
	tracer            opentracing.Tracer
}

// NewPaperResourcesService 创建新的论文资源服务
func NewPaperResourcesService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperResourcesDAO *dao.PaperResourcesDAO,
) *PaperResourcesService {
	return &PaperResourcesService{
		logger:            logger,
		tracer:            tracer,
		paperResourcesDAO: paperResourcesDAO,
	}
}

// CreatePaperResources 创建论文资源
func (s *PaperResourcesService) CreatePaperResources(ctx context.Context, resources *model.PaperResources) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.CreatePaperResources")
	defer span.Finish()

	if err := s.paperResourcesDAO.Create(ctx, resources); err != nil {
		s.logger.Error("创建论文资源失败", "error", err)
		return "0", errors.Biz("paper.resources.errors.create_failed")
	}

	return resources.Id, nil
}

// GetPaperResourcesById 根据ID获取论文资源
func (s *PaperResourcesService) GetPaperResourcesById(ctx context.Context, id string) (*model.PaperResources, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.GetPaperResourcesById")
	defer span.Finish()

	resources, err := s.paperResourcesDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文资源失败", "id", id, "error", err)
		return nil, errors.Biz("paper.resources.errors.get_failed")
	}

	if resources == nil {
		return nil, errors.Biz("paper.resources.errors.not_found")
	}

	return resources, nil
}

// GetPaperResourcesByPaperId 根据论文ID获取论文资源列表
func (s *PaperResourcesService) GetPaperResourcesByPaperId(ctx context.Context, paperId string) ([]*model.PaperResources, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.GetPaperResourcesByPaperId")
	defer span.Finish()

	resourcesList, err := s.paperResourcesDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文资源列表失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.resources.errors.list_failed")
	}

	result := make([]*model.PaperResources, 0, len(resourcesList))
	for i := range resourcesList {
		result = append(result, &resourcesList[i])
	}

	return result, nil
}

// GetPaperResourcesByPaperTitle 根据论文标题获取论文资源列表
func (s *PaperResourcesService) GetPaperResourcesByPaperTitle(ctx context.Context, paperTitle string) ([]*model.PaperResources, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.GetPaperResourcesByPaperTitle")
	defer span.Finish()

	resourcesList, err := s.paperResourcesDAO.FindByPaperTitle(ctx, paperTitle)
	if err != nil {
		s.logger.Error("获取论文资源列表失败", "paper_title", paperTitle, "error", err)
		return nil, errors.Biz("paper.resources.errors.list_failed")
	}

	result := make([]*model.PaperResources, 0, len(resourcesList))
	for i := range resourcesList {
		result = append(result, &resourcesList[i])
	}

	return result, nil
}

// GetPaperResourcesByResourceTitle 根据资源标题获取论文资源列表
func (s *PaperResourcesService) GetPaperResourcesByResourceTitle(ctx context.Context, resourceTitle string) ([]*model.PaperResources, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.GetPaperResourcesByResourceTitle")
	defer span.Finish()

	resourcesList, err := s.paperResourcesDAO.FindByResourceTitle(ctx, resourceTitle)
	if err != nil {
		s.logger.Error("获取论文资源列表失败", "resource_title", resourceTitle, "error", err)
		return nil, errors.Biz("paper.resources.errors.list_failed")
	}

	result := make([]*model.PaperResources, 0, len(resourcesList))
	for i := range resourcesList {
		result = append(result, &resourcesList[i])
	}

	return result, nil
}

// UpdatePaperResources 更新论文资源
func (s *PaperResourcesService) UpdatePaperResources(ctx context.Context, resources *model.PaperResources) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.UpdatePaperResources")
	defer span.Finish()

	// 获取论文资源
	existingResources, err := s.paperResourcesDAO.FindById(ctx, resources.Id)
	if err != nil {
		s.logger.Error("获取论文资源失败", "id", resources.Id, "error", err)
		return false, errors.Biz("paper.resources.errors.get_failed")
	}

	if existingResources == nil {
		return false, errors.Biz("paper.resources.errors.not_found")
	}

	// 更新论文资源
	if err := s.paperResourcesDAO.UpdateById(ctx, resources); err != nil {
		s.logger.Error("更新论文资源失败", "id", resources.Id, "error", err)
		return false, errors.Biz("paper.resources.errors.update_failed")
	}

	return true, nil
}

// DeletePaperResources 删除论文资源
func (s *PaperResourcesService) DeletePaperResources(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.DeletePaperResources")
	defer span.Finish()

	// 删除论文资源
	if err := s.paperResourcesDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文资源失败", "id", id, "error", err)
		return false, errors.Biz("paper.resources.errors.delete_failed")
	}

	return true, nil
}

// DeletePaperResourcesByPaperId 根据论文ID删除论文资源
func (s *PaperResourcesService) DeletePaperResourcesByPaperId(ctx context.Context, paperId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperResourcesService.DeletePaperResourcesByPaperId")
	defer span.Finish()

	// 删除论文资源
	if err := s.paperResourcesDAO.DeleteByPaperId(ctx, paperId); err != nil {
		s.logger.Error("删除论文资源失败", "paper_id", paperId, "error", err)
		return false, errors.Biz("paper.resources.errors.delete_failed")
	}

	return true, nil
}
