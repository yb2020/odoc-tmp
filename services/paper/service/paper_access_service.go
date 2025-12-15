package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperAccessService 论文访问权限服务实现
type PaperAccessService struct {
	paperAccessDAO *dao.PaperAccessDAO
	logger         logging.Logger
	tracer         opentracing.Tracer
}

// NewPaperAccessService 创建新的论文访问权限服务
func NewPaperAccessService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperAccessDAO *dao.PaperAccessDAO,
) *PaperAccessService {
	return &PaperAccessService{
		logger:         logger,
		tracer:         tracer,
		paperAccessDAO: paperAccessDAO,
	}
}

// CreatePaperAccess 创建论文访问权限
func (s *PaperAccessService) CreatePaperAccess(ctx context.Context, access *model.PaperAccess) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.CreatePaperAccess")
	defer span.Finish()

	if err := s.paperAccessDAO.Create(ctx, access); err != nil {
		s.logger.Error("创建论文访问权限失败", "error", err)
		return "0", errors.Biz("paper.access.errors.create_failed")
	}

	return access.Id, nil
}

// GetPaperAccessById 根据ID获取论文访问权限
func (s *PaperAccessService) GetPaperAccessById(ctx context.Context, id string) (*model.PaperAccess, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.GetPaperAccessById")
	defer span.Finish()

	access, err := s.paperAccessDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文访问权限失败", "id", id, "error", err)
		return nil, errors.Biz("paper.access.errors.get_failed")
	}

	if access == nil {
		return nil, errors.Biz("paper.access.errors.not_found")
	}

	return access, nil
}

// GetPaperAccessesByPaperId 根据论文ID获取论文访问权限列表
func (s *PaperAccessService) GetPaperAccessesByPaperId(ctx context.Context, paperId string) ([]*model.PaperAccess, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.GetPaperAccessesByPaperId")
	defer span.Finish()

	accesses, err := s.paperAccessDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文访问权限列表失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.access.errors.list_failed")
	}

	result := make([]*model.PaperAccess, 0, len(accesses))
	for i := range accesses {
		result = append(result, &accesses[i])
	}

	return result, nil
}

// GetPaperAccessesByUserId 根据用户ID获取论文访问权限列表
func (s *PaperAccessService) GetPaperAccessesByUserId(ctx context.Context, userId string) ([]*model.PaperAccess, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.GetPaperAccessesByUserId")
	defer span.Finish()

	accesses, err := s.paperAccessDAO.FindByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("获取论文访问权限列表失败", "user_id", userId, "error", err)
		return nil, errors.Biz("paper.access.errors.list_failed")
	}

	result := make([]*model.PaperAccess, 0, len(accesses))
	for i := range accesses {
		result = append(result, &accesses[i])
	}

	return result, nil
}

// GetPaperAccessByPaperIdAndUserId 根据论文ID和用户ID获取论文访问权限
func (s *PaperAccessService) GetPaperAccessByPaperIdAndUserId(ctx context.Context, paperId, userId string) (*model.PaperAccess, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.GetPaperAccessByPaperIdAndUserId")
	defer span.Finish()

	access, err := s.paperAccessDAO.FindByPaperIdAndUserId(ctx, paperId, userId)
	if err != nil {
		s.logger.Error("获取论文访问权限失败", "paper_id", paperId, "user_id", userId, "error", err)
		return nil, errors.Biz("paper.access.errors.get_failed")
	}

	if access == nil {
		return nil, errors.Biz("paper.access.errors.not_found")
	}

	return access, nil
}

// DeletePaperAccess 删除论文访问权限
func (s *PaperAccessService) DeletePaperAccess(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.DeletePaperAccess")
	defer span.Finish()

	if err := s.paperAccessDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文访问权限失败", "id", id, "error", err)
		return false, errors.Biz("paper.access.errors.delete_failed")
	}

	return true, nil
}

// DeletePaperAccessByPaperIdAndUserId 根据论文ID和用户ID删除论文访问权限
func (s *PaperAccessService) DeletePaperAccessByPaperIdAndUserId(ctx context.Context, paperId, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.DeletePaperAccessByPaperIdAndUserId")
	defer span.Finish()

	if err := s.paperAccessDAO.DeleteByPaperIdAndUserId(ctx, paperId, userId); err != nil {
		s.logger.Error("删除论文访问权限失败", "paper_id", paperId, "user_id", userId, "error", err)
		return false, errors.Biz("paper.access.errors.delete_failed")
	}

	return true, nil
}

// HasAccess 检查用户是否有访问论文的权限
func (s *PaperAccessService) HasAccess(ctx context.Context, paperId, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.HasAccess")
	defer span.Finish()

	access, err := s.paperAccessDAO.FindByPaperIdAndUserId(ctx, paperId, userId)
	if err != nil {
		s.logger.Error("检查访问权限失败", "paper_id", paperId, "user_id", userId, "error", err)
		return false, errors.Biz("paper.access.errors.check_failed")
	}

	return access != nil, nil
}

// 判断用户是否有权限打开某论文
func (s *PaperAccessService) IsUserHasPaperAccess(ctx context.Context, userId string, paperId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAccessService.IsUserHasPaperAccess")
	defer span.Finish()

	// TODO:实现权限判断逻辑

	return false, nil
}
