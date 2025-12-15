package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"
)

// NoteDrawEntityService 笔记绘制实体服务实现
type NoteDrawEntityService struct {
	noteDrawEntityDAO *dao.NoteDrawEntityDAO
	logger            logging.Logger
	tracer            opentracing.Tracer
}

// NewNoteDrawEntityService 创建新的笔记绘制实体服务
func NewNoteDrawEntityService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteDrawEntityDAO *dao.NoteDrawEntityDAO,
) *NoteDrawEntityService {
	return &NoteDrawEntityService{
		noteDrawEntityDAO: noteDrawEntityDAO,
		logger:            logger,
		tracer:            tracer,
	}
}

// CreateNoteDrawEntity 创建笔记绘制实体
func (s *NoteDrawEntityService) CreateNoteDrawEntity(ctx context.Context, drawEntity *model.NoteDrawEntity) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.CreateNoteDrawEntity")
	defer span.Finish()

	if err := s.noteDrawEntityDAO.Create(ctx, drawEntity); err != nil {
		s.logger.Error("创建笔记绘制实体失败", "error", err)
		return "0", errors.Biz("note.note_draw_entity.errors.create_failed")
	}

	return drawEntity.Id, nil
}

// UpdateNoteDrawEntity 更新笔记绘制实体
func (s *NoteDrawEntityService) UpdateNoteDrawEntity(ctx context.Context, drawEntity *model.NoteDrawEntity) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.UpdateNoteDrawEntity")
	defer span.Finish()

	// 获取笔记绘制实体
	existingDrawEntity, err := s.noteDrawEntityDAO.FindById(ctx, drawEntity.Id)
	if err != nil {
		s.logger.Error("获取笔记绘制实体失败", "error", err)
		return false, errors.Biz("note.note_draw_entity.errors.get_failed")
	}

	if existingDrawEntity == nil {
		return false, errors.Biz("note.note_draw_entity.errors.not_found")
	}

	// 更新笔记绘制实体
	if err := s.noteDrawEntityDAO.UpdateById(ctx, drawEntity); err != nil {
		s.logger.Error("更新笔记绘制实体失败", "error", err)
		return false, errors.Biz("note.note_draw_entity.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteDrawEntityById 删除笔记绘制实体
func (s *NoteDrawEntityService) DeleteNoteDrawEntityById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.DeleteNoteDrawEntityById")
	defer span.Finish()

	// 删除笔记绘制实体
	if err := s.noteDrawEntityDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记绘制实体失败", "error", err)
		return false, errors.Biz("note.note_draw_entity.errors.delete_failed")
	}

	return true, nil
}

// GetNoteDrawEntityById 根据ID获取笔记绘制实体
func (s *NoteDrawEntityService) GetNoteDrawEntityById(ctx context.Context, id string) (*model.NoteDrawEntity, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.GetNoteDrawEntityById")
	defer span.Finish()

	// 获取笔记绘制实体
	drawEntity, err := s.noteDrawEntityDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记绘制实体失败", "error", err)
		return nil, errors.Biz("note.note_draw_entity.errors.get_failed")
	}

	if drawEntity == nil {
		return nil, errors.Biz("note.note_draw_entity.errors.not_found")
	}

	return drawEntity, nil
}

// GetNoteDrawEntitiesByNoteId 根据笔记ID获取笔记绘制实体列表
func (s *NoteDrawEntityService) GetNoteDrawEntitiesByNoteId(ctx context.Context, noteId string) ([]model.NoteDrawEntity, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.GetNoteDrawEntitiesByNoteId")
	defer span.Finish()

	// 获取笔记绘制实体列表
	drawEntities, err := s.noteDrawEntityDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记绘制实体列表失败", "error", err)
		return nil, errors.Biz("note.note_draw_entity.errors.list_failed")
	}

	return drawEntities, nil
}

// GetNoteDrawEntitiesByNoteIdAndPage 根据笔记ID和页码获取笔记绘制实体列表
func (s *NoteDrawEntityService) GetNoteDrawEntitiesByNoteIdAndPage(ctx context.Context, noteId string, pageNumber int) ([]model.NoteDrawEntity, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.GetNoteDrawEntitiesByNoteIdAndPage")
	defer span.Finish()

	// 获取笔记绘制实体列表
	drawEntities, err := s.noteDrawEntityDAO.GetByNoteIDAndPage(ctx, noteId, pageNumber)
	if err != nil {
		s.logger.Error("获取笔记绘制实体列表失败", "error", err)
		return nil, errors.Biz("note.note_draw_entity.errors.list_failed")
	}

	return drawEntities, nil
}

// DeleteNoteDrawEntitiesByNoteId 根据笔记ID删除所有笔记绘制实体
func (s *NoteDrawEntityService) DeleteNoteDrawEntitiesByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteDrawEntityService.DeleteNoteDrawEntitiesByNoteId")
	defer span.Finish()

	// 删除笔记绘制实体
	if err := s.noteDrawEntityDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记绘制实体失败", "error", err)
		return false, errors.Biz("note.note_draw_entity.errors.delete_failed")
	}

	return true, nil
}
