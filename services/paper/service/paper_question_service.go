package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperQuestionService 论文问题服务实现
type PaperQuestionService struct {
	paperQuestionDAO *dao.PaperQuestionDAO
	logger           logging.Logger
	tracer           opentracing.Tracer
}

// NewPaperQuestionService 创建新的论文问题服务
func NewPaperQuestionService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperQuestionDAO *dao.PaperQuestionDAO,
) *PaperQuestionService {
	return &PaperQuestionService{
		logger:           logger,
		tracer:           tracer,
		paperQuestionDAO: paperQuestionDAO,
	}
}

// CreatePaperQuestion 创建论文问题
func (s *PaperQuestionService) CreatePaperQuestion(ctx context.Context, question *model.PaperQuestion) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.CreatePaperQuestion")
	defer span.Finish()

	if err := s.paperQuestionDAO.Create(ctx, question); err != nil {
		s.logger.Error("创建论文问题失败", "error", err)
		return "0", errors.Biz("paper.question.errors.create_failed")
	}

	return question.Id, nil
}

// GetPaperQuestionById 根据ID获取论文问题
func (s *PaperQuestionService) GetPaperQuestionById(ctx context.Context, id string) (*model.PaperQuestion, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.GetPaperQuestionById")
	defer span.Finish()

	question, err := s.paperQuestionDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文问题失败", "id", id, "error", err)
		return nil, errors.Biz("paper.question.errors.get_failed")
	}

	if question == nil {
		return nil, errors.Biz("paper.question.errors.not_found")
	}

	return question, nil
}

// GetPaperQuestionsByPaperId 根据论文ID获取论文问题列表
func (s *PaperQuestionService) GetPaperQuestionsByPaperId(ctx context.Context, paperId string) ([]*model.PaperQuestion, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.GetPaperQuestionsByPaperId")
	defer span.Finish()

	questions, err := s.paperQuestionDAO.FindByPaperId(ctx, paperId)
	if err != nil {
		s.logger.Error("获取论文问题列表失败", "paper_id", paperId, "error", err)
		return nil, errors.Biz("paper.question.errors.list_failed")
	}

	result := make([]*model.PaperQuestion, 0, len(questions))
	for i := range questions {
		result = append(result, &questions[i])
	}

	return result, nil
}

// GetPaperQuestionsByPdfId 根据PDF ID获取论文问题列表
func (s *PaperQuestionService) GetPaperQuestionsByPdfId(ctx context.Context, pdfId string) ([]*model.PaperQuestion, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.GetPaperQuestionsByPdfId")
	defer span.Finish()

	questions, err := s.paperQuestionDAO.FindByPdfId(ctx, pdfId)
	if err != nil {
		s.logger.Error("获取论文问题列表失败", "pdf_id", pdfId, "error", err)
		return nil, errors.Biz("paper.question.errors.list_failed")
	}

	result := make([]*model.PaperQuestion, 0, len(questions))
	for i := range questions {
		result = append(result, &questions[i])
	}

	return result, nil
}

// UpdatePaperQuestion 更新论文问题
func (s *PaperQuestionService) UpdatePaperQuestion(ctx context.Context, question *model.PaperQuestion) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.UpdatePaperQuestion")
	defer span.Finish()

	// 获取论文问题
	existingQuestion, err := s.paperQuestionDAO.FindById(ctx, question.Id)
	if err != nil {
		s.logger.Error("获取论文问题失败", "id", question.Id, "error", err)
		return false, errors.Biz("paper.question.errors.get_failed")
	}

	if existingQuestion == nil {
		return false, errors.Biz("paper.question.errors.not_found")
	}

	// 更新论文问题
	if err := s.paperQuestionDAO.UpdateById(ctx, question); err != nil {
		s.logger.Error("更新论文问题失败", "id", question.Id, "error", err)
		return false, errors.Biz("paper.question.errors.update_failed")
	}

	return true, nil
}

// IncrementViewCount 增加查看次数
func (s *PaperQuestionService) IncrementViewCount(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.IncrementViewCount")
	defer span.Finish()

	// 增加查看次数
	if err := s.paperQuestionDAO.IncrementViewCount(ctx, id); err != nil {
		s.logger.Error("增加论文问题查看次数失败", "id", id, "error", err)
		return false, errors.Biz("paper.question.errors.increment_view_count_failed")
	}

	return true, nil
}

// DeletePaperQuestion 删除论文问题
func (s *PaperQuestionService) DeletePaperQuestion(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.DeletePaperQuestion")
	defer span.Finish()

	// 删除论文问题
	if err := s.paperQuestionDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文问题失败", "id", id, "error", err)
		return false, errors.Biz("paper.question.errors.delete_failed")
	}

	return true, nil
}

// DeletePaperQuestionsByPaperId 根据论文ID删除论文问题
func (s *PaperQuestionService) DeletePaperQuestionsByPaperId(ctx context.Context, paperId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperQuestionService.DeletePaperQuestionsByPaperId")
	defer span.Finish()

	// 删除论文问题
	if err := s.paperQuestionDAO.DeleteByPaperId(ctx, paperId); err != nil {
		s.logger.Error("删除论文问题失败", "paper_id", paperId, "error", err)
		return false, errors.Biz("paper.question.errors.delete_failed")
	}

	return true, nil
}
