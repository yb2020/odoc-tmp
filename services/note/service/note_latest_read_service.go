package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"
)

// NoteLatestReadService 笔记最近阅读服务实现
type NoteLatestReadService struct {
	noteLatestReadDAO *dao.NoteLatestReadDAO
	logger            logging.Logger
	tracer            opentracing.Tracer
}

// NewNoteLatestReadService 创建新的笔记最近阅读服务
func NewNoteLatestReadService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteLatestReadDAO *dao.NoteLatestReadDAO,
) *NoteLatestReadService {
	return &NoteLatestReadService{
		noteLatestReadDAO: noteLatestReadDAO,
		logger:            logger,
		tracer:            tracer,
	}
}

// CreateNoteLatestRead 创建笔记最近阅读记录
func (s *NoteLatestReadService) CreateNoteLatestRead(ctx context.Context, latestRead *model.NoteLatestRead) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteLatestReadService.CreateNoteLatestRead")
	defer span.Finish()

	if err := s.noteLatestReadDAO.Create(ctx, latestRead); err != nil {
		s.logger.Error("创建笔记最近阅读记录失败", "error", err)
		return "0", errors.Biz("note.note_latest_read.errors.create_failed")
	}

	return latestRead.Id, nil
}

// UpdateNoteLatestRead 更新笔记最近阅读记录
func (s *NoteLatestReadService) UpdateNoteLatestRead(ctx context.Context, latestRead *model.NoteLatestRead) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteLatestReadService.UpdateNoteLatestRead")
	defer span.Finish()

	// 获取笔记最近阅读记录
	existingLatestRead, err := s.noteLatestReadDAO.FindById(ctx, latestRead.Id)
	if err != nil {
		s.logger.Error("获取笔记最近阅读记录失败", "error", err)
		return false, errors.Biz("note.note_latest_read.errors.get_failed")
	}

	if existingLatestRead == nil {
		return false, errors.Biz("note.note_latest_read.errors.not_found")
	}

	// 更新笔记最近阅读记录
	if err := s.noteLatestReadDAO.UpdateById(ctx, latestRead); err != nil {
		s.logger.Error("更新笔记最近阅读记录失败", "error", err)
		return false, errors.Biz("note.note_latest_read.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteLatestReadById 删除笔记最近阅读记录
func (s *NoteLatestReadService) DeleteNoteLatestReadById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteLatestReadService.DeleteNoteLatestReadById")
	defer span.Finish()

	// 删除笔记最近阅读记录
	if err := s.noteLatestReadDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记最近阅读记录失败", "error", err)
		return false, errors.Biz("note.note_latest_read.errors.delete_failed")
	}

	return true, nil
}

// GetNoteLatestReadById 根据ID获取笔记最近阅读记录
func (s *NoteLatestReadService) GetNoteLatestReadById(ctx context.Context, id string) (*model.NoteLatestRead, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteLatestReadService.GetNoteLatestReadById")
	defer span.Finish()

	// 获取笔记最近阅读记录
	latestRead, err := s.noteLatestReadDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记最近阅读记录失败", "error", err)
		return nil, errors.Biz("note.note_latest_read.errors.get_failed")
	}

	if latestRead == nil {
		return nil, errors.Biz("note.note_latest_read.errors.not_found")
	}

	return latestRead, nil
}

// GetNoteLatestReadByNoteId 根据笔记ID获取笔记最近阅读记录
func (s *NoteLatestReadService) GetNoteLatestReadByNoteId(ctx context.Context, noteId string) (*model.NoteLatestRead, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteLatestReadService.GetNoteLatestReadByNoteId")
	defer span.Finish()

	// 获取笔记最近阅读记录
	latestRead, err := s.noteLatestReadDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记最近阅读记录失败", "error", err)
		return nil, errors.Biz("note.note_latest_read.errors.get_failed")
	}

	if latestRead == nil {
		return nil, errors.Biz("note.note_latest_read.errors.not_found")
	}

	return latestRead, nil
}

// DeleteNoteLatestReadByNoteId 根据笔记ID删除笔记最近阅读记录
func (s *NoteLatestReadService) DeleteNoteLatestReadByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteLatestReadService.DeleteNoteLatestReadByNoteId")
	defer span.Finish()

	// 删除笔记最近阅读记录
	if err := s.noteLatestReadDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记最近阅读记录失败", "error", err)
		return false, errors.Biz("note.note_latest_read.errors.delete_failed")
	}

	return true, nil
}
