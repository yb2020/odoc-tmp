package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/paper/dao"
	"github.com/yb2020/odoc/services/paper/model"
)

// PaperAnswerService 论文问答答案服务实现
type PaperAnswerService struct {
	paperAnswerDAO *dao.PaperAnswerDAO
	logger         logging.Logger
	tracer         opentracing.Tracer
}

// NewPaperAnswerService 创建新的论文问答答案服务
func NewPaperAnswerService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperAnswerDAO *dao.PaperAnswerDAO,
) *PaperAnswerService {
	return &PaperAnswerService{
		logger:         logger,
		tracer:         tracer,
		paperAnswerDAO: paperAnswerDAO,
	}
}

// CreatePaperAnswer 创建论文问答答案
func (s *PaperAnswerService) CreatePaperAnswer(ctx context.Context, answer *model.PaperAnswer) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.CreatePaperAnswer")
	defer span.Finish()

	if err := s.paperAnswerDAO.Create(ctx, answer); err != nil {
		s.logger.Error("创建论文问答答案失败", "error", err)
		return "0", errors.Biz("paper.answer.errors.create_failed")
	}

	return answer.Id, nil
}

// GetPaperAnswerById 根据ID获取论文问答答案
func (s *PaperAnswerService) GetPaperAnswerById(ctx context.Context, id string) (*model.PaperAnswer, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.GetPaperAnswerById")
	defer span.Finish()

	answer, err := s.paperAnswerDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文问答答案失败", "id", id, "error", err)
		return nil, errors.Biz("paper.answer.errors.get_failed")
	}

	if answer == nil {
		return nil, errors.Biz("paper.answer.errors.not_found")
	}

	return answer, nil
}

// GetPaperAnswersByQuestionId 根据问题ID获取论文问答答案列表
func (s *PaperAnswerService) GetPaperAnswersByQuestionId(ctx context.Context, questionId string) ([]*model.PaperAnswer, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.GetPaperAnswersByQuestionId")
	defer span.Finish()

	answers, err := s.paperAnswerDAO.FindByQuestionId(ctx, questionId)
	if err != nil {
		s.logger.Error("获取论文问答答案列表失败", "question_id", questionId, "error", err)
		return nil, errors.Biz("paper.answer.errors.list_failed")
	}

	result := make([]*model.PaperAnswer, 0, len(answers))
	for i := range answers {
		result = append(result, &answers[i])
	}

	return result, nil
}

// GetPaperAnswersByReplyUserId 根据回复用户ID获取论文问答答案列表
func (s *PaperAnswerService) GetPaperAnswersByReplyUserId(ctx context.Context, replyUserId string) ([]*model.PaperAnswer, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.GetPaperAnswersByReplyUserId")
	defer span.Finish()

	answers, err := s.paperAnswerDAO.FindByReplyUserId(ctx, replyUserId)
	if err != nil {
		s.logger.Error("获取论文问答答案列表失败", "reply_user_id", replyUserId, "error", err)
		return nil, errors.Biz("paper.answer.errors.list_failed")
	}

	result := make([]*model.PaperAnswer, 0, len(answers))
	for i := range answers {
		result = append(result, &answers[i])
	}

	return result, nil
}

// GetPaperAnswersByReplyAnswerId 根据回复答案ID获取论文问答答案列表
func (s *PaperAnswerService) GetPaperAnswersByReplyAnswerId(ctx context.Context, replyAnswerId string) ([]*model.PaperAnswer, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.GetPaperAnswersByReplyAnswerId")
	defer span.Finish()

	answers, err := s.paperAnswerDAO.FindByReplyAnswerId(ctx, replyAnswerId)
	if err != nil {
		s.logger.Error("获取论文问答答案列表失败", "reply_answer_id", replyAnswerId, "error", err)
		return nil, errors.Biz("paper.answer.errors.list_failed")
	}

	result := make([]*model.PaperAnswer, 0, len(answers))
	for i := range answers {
		result = append(result, &answers[i])
	}

	return result, nil
}

// UpdatePaperAnswer 更新论文问答答案
func (s *PaperAnswerService) UpdatePaperAnswer(ctx context.Context, answer *model.PaperAnswer) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.UpdatePaperAnswer")
	defer span.Finish()

	// 获取论文问答答案
	existingAnswer, err := s.paperAnswerDAO.FindById(ctx, answer.Id)
	if err != nil {
		s.logger.Error("获取论文问答答案失败", "id", answer.Id, "error", err)
		return false, errors.Biz("paper.answer.errors.get_failed")
	}

	if existingAnswer == nil {
		return false, errors.Biz("paper.answer.errors.not_found")
	}

	// 更新论文问答答案
	if err := s.paperAnswerDAO.UpdateById(ctx, answer); err != nil {
		s.logger.Error("更新论文问答答案失败", "id", answer.Id, "error", err)
		return false, errors.Biz("paper.answer.errors.update_failed")
	}

	return true, nil
}

// DeletePaperAnswer 删除论文问答答案
func (s *PaperAnswerService) DeletePaperAnswer(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperAnswerService.DeletePaperAnswer")
	defer span.Finish()

	// 删除论文问答答案
	if err := s.paperAnswerDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文问答答案失败", "id", id, "error", err)
		return false, errors.Biz("paper.answer.errors.delete_failed")
	}

	return true, nil
}
