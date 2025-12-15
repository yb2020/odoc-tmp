package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"
)

// NoteExportHistoryService 笔记导出历史服务实现
type NoteExportHistoryService struct {
	noteExportHistoryDAO *dao.NoteExportHistoryDAO
	logger               logging.Logger
	tracer               opentracing.Tracer
}

// NewNoteExportHistoryService 创建新的笔记导出历史服务
func NewNoteExportHistoryService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteExportHistoryDAO *dao.NoteExportHistoryDAO,
) *NoteExportHistoryService {
	return &NoteExportHistoryService{
		noteExportHistoryDAO: noteExportHistoryDAO,
		logger:               logger,
		tracer:               tracer,
	}
}

// CreateNoteExportHistory 创建笔记导出历史
func (s *NoteExportHistoryService) CreateNoteExportHistory(ctx context.Context, history *model.NoteExportHistory) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.CreateNoteExportHistory")
	defer span.Finish()

	if err := s.noteExportHistoryDAO.Create(ctx, history); err != nil {
		s.logger.Error("创建笔记导出历史失败", "error", err)
		return "0", errors.Biz("note.note_export_history.errors.create_failed")
	}

	return history.Id, nil
}

// UpdateNoteExportHistory 更新笔记导出历史
func (s *NoteExportHistoryService) UpdateNoteExportHistory(ctx context.Context, history *model.NoteExportHistory) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.UpdateNoteExportHistory")
	defer span.Finish()

	// 获取笔记导出历史
	existingHistory, err := s.noteExportHistoryDAO.FindById(ctx, history.Id)
	if err != nil {
		s.logger.Error("获取笔记导出历史失败", "error", err)
		return false, errors.Biz("note.note_export_history.errors.get_failed")
	}

	if existingHistory == nil {
		return false, errors.Biz("note.note_export_history.errors.not_found")
	}

	// 更新笔记导出历史
	if err := s.noteExportHistoryDAO.UpdateById(ctx, history); err != nil {
		s.logger.Error("更新笔记导出历史失败", "error", err)
		return false, errors.Biz("note.note_export_history.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteExportHistoryById 删除笔记导出历史
func (s *NoteExportHistoryService) DeleteNoteExportHistoryById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.DeleteNoteExportHistoryById")
	defer span.Finish()

	// 删除笔记导出历史
	if err := s.noteExportHistoryDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记导出历史失败", "error", err)
		return false, errors.Biz("note.note_export_history.errors.delete_failed")
	}

	return true, nil
}

// GetNoteExportHistoryById 根据ID获取笔记导出历史
func (s *NoteExportHistoryService) GetNoteExportHistoryById(ctx context.Context, id string) (*model.NoteExportHistory, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.GetNoteExportHistoryById")
	defer span.Finish()

	// 获取笔记导出历史
	history, err := s.noteExportHistoryDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记导出历史失败", "error", err)
		return nil, errors.Biz("note.note_export_history.errors.get_failed")
	}

	if history == nil {
		return nil, errors.Biz("note.note_export_history.errors.not_found")
	}

	return history, nil
}

// GetNoteExportHistoriesByNoteId 根据笔记ID获取笔记导出历史列表
func (s *NoteExportHistoryService) GetNoteExportHistoriesByNoteId(ctx context.Context, noteId string) ([]model.NoteExportHistory, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.GetNoteExportHistoriesByNoteId")
	defer span.Finish()

	// 获取笔记导出历史列表
	histories, err := s.noteExportHistoryDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记导出历史列表失败", "error", err)
		return nil, errors.Biz("note.note_export_history.errors.list_failed")
	}

	return histories, nil
}

// GetLatestNoteExportHistoryByNoteId 根据笔记ID获取最新的笔记导出历史
func (s *NoteExportHistoryService) GetLatestNoteExportHistoryByNoteId(ctx context.Context, noteId string) (*model.NoteExportHistory, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.GetLatestNoteExportHistoryByNoteId")
	defer span.Finish()

	// 获取最新的笔记导出历史
	history, err := s.noteExportHistoryDAO.GetLatestByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取最新的笔记导出历史失败", "error", err)
		return nil, errors.Biz("note.note_export_history.errors.get_failed")
	}

	if history == nil {
		return nil, errors.Biz("note.note_export_history.errors.not_found")
	}

	return history, nil
}

// GetNoteExportHistoryByNoteIdAndVersion 根据笔记ID和版本获取笔记导出历史
func (s *NoteExportHistoryService) GetNoteExportHistoryByNoteIdAndVersion(ctx context.Context, noteId string, version int64) (*model.NoteExportHistory, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.GetNoteExportHistoryByNoteIdAndVersion")
	defer span.Finish()

	// 获取笔记导出历史
	history, err := s.noteExportHistoryDAO.GetByNoteIDAndVersion(ctx, noteId, version)
	if err != nil {
		s.logger.Error("获取笔记导出历史失败", "error", err)
		return nil, errors.Biz("note.note_export_history.errors.get_failed")
	}

	if history == nil {
		return nil, errors.Biz("note.note_export_history.errors.not_found")
	}

	return history, nil
}

// DeleteNoteExportHistoriesByNoteId 根据笔记ID删除所有笔记导出历史
func (s *NoteExportHistoryService) DeleteNoteExportHistoriesByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteExportHistoryService.DeleteNoteExportHistoriesByNoteId")
	defer span.Finish()

	// 删除笔记导出历史
	if err := s.noteExportHistoryDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记导出历史失败", "error", err)
		return false, errors.Biz("note.note_export_history.errors.delete_failed")
	}

	return true, nil
}
