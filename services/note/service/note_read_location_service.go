package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"
)

// NoteReadLocationService 笔记阅读位置服务实现
type NoteReadLocationService struct {
	noteReadLocationDAO *dao.NoteReadLocationDAO
	logger              logging.Logger
	tracer              opentracing.Tracer
}

// NewNoteReadLocationService 创建新的笔记阅读位置服务
func NewNoteReadLocationService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteReadLocationDAO *dao.NoteReadLocationDAO,
) *NoteReadLocationService {
	return &NoteReadLocationService{
		noteReadLocationDAO: noteReadLocationDAO,
		logger:              logger,
		tracer:              tracer,
	}
}

// CreateNoteReadLocation 创建笔记阅读位置
func (s *NoteReadLocationService) CreateNoteReadLocation(ctx context.Context, location *model.NoteReadLocation) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.CreateNoteReadLocation")
	defer span.Finish()

	if err := s.noteReadLocationDAO.Create(ctx, location); err != nil {
		s.logger.Error("创建笔记阅读位置失败", "error", err)
		return "0", errors.Biz("note.note_read_location.errors.create_failed")
	}

	return location.Id, nil
}

// UpdateNoteReadLocation 更新笔记阅读位置
func (s *NoteReadLocationService) UpdateNoteReadLocation(ctx context.Context, location *model.NoteReadLocation) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.UpdateNoteReadLocation")
	defer span.Finish()

	// 获取笔记阅读位置
	existingLocation, err := s.noteReadLocationDAO.FindById(ctx, location.Id)
	if err != nil {
		s.logger.Error("获取笔记阅读位置失败", "error", err)
		return false, errors.Biz("note.note_read_location.errors.get_failed")
	}

	if existingLocation == nil {
		return false, errors.Biz("note.note_read_location.errors.not_found")
	}

	// 更新笔记阅读位置
	if err := s.noteReadLocationDAO.UpdateById(ctx, location); err != nil {
		s.logger.Error("更新笔记阅读位置失败", "error", err)
		return false, errors.Biz("note.note_read_location.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteReadLocationById 删除笔记阅读位置
func (s *NoteReadLocationService) DeleteNoteReadLocationById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.DeleteNoteReadLocationById")
	defer span.Finish()

	// 删除笔记阅读位置
	if err := s.noteReadLocationDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记阅读位置失败", "error", err)
		return false, errors.Biz("note.note_read_location.errors.delete_failed")
	}

	return true, nil
}

// GetNoteReadLocationById 根据ID获取笔记阅读位置
func (s *NoteReadLocationService) GetNoteReadLocationById(ctx context.Context, id string) (*model.NoteReadLocation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.GetNoteReadLocationById")
	defer span.Finish()

	// 获取笔记阅读位置
	location, err := s.noteReadLocationDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记阅读位置失败", "error", err)
		return nil, errors.Biz("note.note_read_location.errors.get_failed")
	}

	if location == nil {
		return nil, errors.Biz("note.note_read_location.errors.not_found")
	}

	return location, nil
}

// GetNoteReadLocationByNoteId 根据笔记ID获取笔记阅读位置
func (s *NoteReadLocationService) GetNoteReadLocationByNoteId(ctx context.Context, noteId string) (*model.NoteReadLocation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.GetNoteReadLocationByNoteId")
	defer span.Finish()

	// 获取笔记阅读位置
	location, err := s.noteReadLocationDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记阅读位置失败", "error", err)
		return nil, errors.Biz("note.note_read_location.errors.get_failed")
	}

	if location == nil {
		return nil, errors.Biz("note.note_read_location.errors.not_found")
	}

	return location, nil
}

// DeleteNoteReadLocationByNoteId 根据笔记ID删除笔记阅读位置
func (s *NoteReadLocationService) DeleteNoteReadLocationByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.DeleteNoteReadLocationByNoteId")
	defer span.Finish()

	// 删除笔记阅读位置
	if err := s.noteReadLocationDAO.DeleteByNoteId(ctx, noteId); err != nil {
		s.logger.Error("删除笔记阅读位置失败", "error", err)
		return false, errors.Biz("note.note_read_location.errors.delete_failed")
	}

	return true, nil
}

// GetNoteReadLocationByNoteId 根据笔记ID获取笔记阅读位置
func (s *NoteReadLocationService) GetNoteReadLocationByUserIdAndNoteId(ctx context.Context, userId string, noteId string) (*model.NoteReadLocation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteReadLocationService.GetNoteReadLocationByUserIdAndNoteId")
	defer span.Finish()

	// 获取笔记阅读位置
	location, err := s.noteReadLocationDAO.GetByUserIdAndNoteId(ctx, userId, noteId)
	if err != nil {
		s.logger.Error("获取笔记阅读位置失败", "error", err)
		return nil, errors.Biz("note.note_read_location.errors.get_failed")
	}

	return location, nil
}
