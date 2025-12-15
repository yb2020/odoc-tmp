package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"
)

// NoteWordConfigService 笔记单词配置服务实现
type NoteWordConfigService struct {
	noteWordConfigDAO *dao.NoteWordConfigDAO
	logger            logging.Logger
	tracer            opentracing.Tracer
}

// NewNoteWordConfigService 创建新的笔记单词配置服务
func NewNoteWordConfigService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteWordConfigDAO *dao.NoteWordConfigDAO,
) *NoteWordConfigService {
	return &NoteWordConfigService{
		noteWordConfigDAO: noteWordConfigDAO,
		logger:            logger,
		tracer:            tracer,
	}
}

// CreateNoteWordConfig 创建笔记单词配置
func (s *NoteWordConfigService) CreateNoteWordConfig(ctx context.Context, config *model.NoteWordConfig) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordConfigService.CreateNoteWordConfig")
	defer span.Finish()

	if err := s.noteWordConfigDAO.Create(ctx, config); err != nil {
		s.logger.Error("创建笔记单词配置失败", "error", err)
		return "0", errors.Biz("note.note_word_config.errors.create_failed")
	}

	return config.Id, nil
}

// UpdateNoteWordConfig 更新笔记单词配置
func (s *NoteWordConfigService) UpdateNoteWordConfig(ctx context.Context, config *model.NoteWordConfig) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordConfigService.UpdateNoteWordConfig")
	defer span.Finish()

	// 获取笔记单词配置
	existingConfig, err := s.noteWordConfigDAO.FindById(ctx, config.Id)
	if err != nil {
		s.logger.Error("获取笔记单词配置失败", "error", err)
		return false, errors.Biz("note.note_word_config.errors.get_failed")
	}

	if existingConfig == nil {
		return false, errors.Biz("note.note_word_config.errors.not_found")
	}

	// 更新笔记单词配置
	if err := s.noteWordConfigDAO.UpdateById(ctx, config); err != nil {
		s.logger.Error("更新笔记单词配置失败", "error", err)
		return false, errors.Biz("note.note_word_config.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteWordConfigById 删除笔记单词配置
func (s *NoteWordConfigService) DeleteNoteWordConfigById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordConfigService.DeleteNoteWordConfigById")
	defer span.Finish()

	// 删除笔记单词配置
	if err := s.noteWordConfigDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记单词配置失败", "error", err)
		return false, errors.Biz("note.note_word_config.errors.delete_failed")
	}

	return true, nil
}

// GetNoteWordConfigById 根据ID获取笔记单词配置
func (s *NoteWordConfigService) GetNoteWordConfigById(ctx context.Context, id string) (*model.NoteWordConfig, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordConfigService.GetNoteWordConfigById")
	defer span.Finish()

	// 获取笔记单词配置
	config, err := s.noteWordConfigDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记单词配置失败", "error", err)
		return nil, errors.Biz("note.note_word_config.errors.get_failed")
	}

	return config, nil
}

// GetNoteWordConfigByNoteId 根据笔记ID获取笔记单词配置
func (s *NoteWordConfigService) GetNoteWordConfigByNoteId(ctx context.Context, noteId string) (*model.NoteWordConfig, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordConfigService.GetNoteWordConfigByNoteId")
	defer span.Finish()

	// 获取笔记单词配置
	config, err := s.noteWordConfigDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记单词配置失败", "error", err)
		return nil, errors.Biz("note.note_word_config.errors.get_failed")
	}

	return config, nil
}

// DeleteNoteWordConfigByNoteId 根据笔记ID删除笔记单词配置
func (s *NoteWordConfigService) DeleteNoteWordConfigByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordConfigService.DeleteNoteWordConfigByNoteId")
	defer span.Finish()

	// 删除笔记单词配置
	if err := s.noteWordConfigDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记单词配置失败", "error", err)
		return false, errors.Biz("note.note_word_config.errors.delete_failed")
	}

	return true, nil
}
