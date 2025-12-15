package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"

	commonPb "github.com/yb2020/odoc-proto/gen/go/common"
)

// NoteShapeService 笔记形状服务实现
type NoteShapeService struct {
	noteShapeDAO *dao.NoteShapeDAO
	logger       logging.Logger
	tracer       opentracing.Tracer
}

// NewNoteShapeService 创建新的笔记形状服务
func NewNoteShapeService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteShapeDAO *dao.NoteShapeDAO,
) *NoteShapeService {
	return &NoteShapeService{
		noteShapeDAO: noteShapeDAO,
		logger:       logger,
		tracer:       tracer,
	}
}

// CreateNoteShape 创建笔记形状
func (s *NoteShapeService) CreateNoteShape(ctx context.Context, shape *model.NoteShape) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.CreateNoteShape")
	defer span.Finish()

	if err := s.noteShapeDAO.Create(ctx, shape); err != nil {
		s.logger.Error("创建笔记形状失败", "error", err)
		return "0", errors.Biz("note.note_shape.errors.create_failed")
	}

	return shape.Id, nil
}

// UpdateNoteShape 更新笔记形状
func (s *NoteShapeService) UpdateNoteShape(ctx context.Context, shape *model.NoteShape) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.UpdateNoteShape")
	defer span.Finish()

	// 获取笔记形状
	existingShape, err := s.noteShapeDAO.FindById(ctx, shape.Id)
	if err != nil {
		s.logger.Error("获取笔记形状失败", "error", err)
		return false, errors.Biz("note.note_shape.errors.get_failed")
	}

	if existingShape == nil {
		return false, errors.Biz("note.note_shape.errors.not_found")
	}

	// 更新笔记形状
	if err := s.noteShapeDAO.UpdateById(ctx, shape); err != nil {
		s.logger.Error("更新笔记形状失败", "error", err)
		return false, errors.Biz("note.note_shape.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteShapeById 删除笔记形状
func (s *NoteShapeService) DeleteNoteShapeById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.DeleteNoteShapeById")
	defer span.Finish()

	// 删除笔记形状
	if err := s.noteShapeDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记形状失败", "error", err)
		return false, errors.Biz("note.note_shape.errors.delete_failed")
	}

	return true, nil
}

// DeleteNoteShapeByIds 删除笔记形状
func (s *NoteShapeService) DeleteNoteShapeByIds(ctx context.Context, ids []string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.DeleteNoteShapeByIds")
	defer span.Finish()

	// 删除笔记形状
	if err := s.noteShapeDAO.DeleteByIds(ctx, ids); err != nil {
		s.logger.Error("删除笔记形状失败", "error", err)
		return false, errors.Biz("note.note_shape.errors.delete_failed")
	}

	return true, nil
}

// GetNoteShapeById 根据ID获取笔记形状
func (s *NoteShapeService) GetNoteShapeById(ctx context.Context, id string) (*model.NoteShape, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.GetNoteShapeById")
	defer span.Finish()

	// 获取笔记形状
	shape, err := s.noteShapeDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记形状失败", "error", err)
		return nil, errors.Biz("note.note_shape.errors.get_failed")
	}

	if shape == nil {
		return nil, errors.Biz("note.note_shape.errors.not_found")
	}

	return shape, nil
}

// GetNoteShapeByUUID 根据UUID获取笔记形状
func (s *NoteShapeService) GetNoteShapeByUUID(ctx context.Context, uuid string) (*model.NoteShape, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.GetNoteShapeByUUID")
	defer span.Finish()

	// 获取笔记形状
	shape, err := s.noteShapeDAO.GetByUUID(ctx, uuid)
	if err != nil {
		s.logger.Error("获取笔记形状失败", "error", err)
		return nil, errors.Biz("note.note_shape.errors.get_failed")
	}

	if shape == nil {
		return nil, errors.Biz("note.note_shape.errors.not_found")
	}

	return shape, nil
}

// GetNoteShapesByNoteId 根据笔记ID获取笔记形状列表
func (s *NoteShapeService) GetNoteShapesByNoteId(ctx context.Context, noteId string) ([]model.NoteShape, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.GetNoteShapesByNoteId")
	defer span.Finish()

	// 获取笔记形状列表
	shapes, err := s.noteShapeDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记形状列表失败", "error", err)
		return nil, errors.Biz("note.note_shape.errors.list_failed")
	}

	return shapes, nil
}

// GetNoteShapesByNoteIdAndPage 根据笔记ID和页码获取笔记形状列表
func (s *NoteShapeService) GetNoteShapesByNoteIdAndPage(ctx context.Context, noteId string, pageNumber int) ([]model.NoteShape, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.GetNoteShapesByNoteIdAndPage")
	defer span.Finish()

	// 获取笔记形状列表
	shapes, err := s.noteShapeDAO.GetByNoteIDAndPage(ctx, noteId, pageNumber)
	if err != nil {
		s.logger.Error("获取笔记形状列表失败", "error", err)
		return nil, errors.Biz("note.note_shape.errors.list_failed")
	}

	return shapes, nil
}

// DeleteNoteShapeByUUID 根据UUID删除笔记形状
func (s *NoteShapeService) DeleteNoteShapeByUUID(ctx context.Context, uuid string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.DeleteNoteShapeByUUID")
	defer span.Finish()

	// 删除笔记形状
	if err := s.noteShapeDAO.DeleteByUUID(ctx, uuid); err != nil {
		s.logger.Error("删除笔记形状失败", "error", err)
		return false, errors.Biz("note.note_shape.errors.delete_failed")
	}

	return true, nil
}

// DeleteNoteShapesByNoteId 根据笔记ID删除所有笔记形状
func (s *NoteShapeService) DeleteNoteShapesByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.DeleteNoteShapesByNoteId")
	defer span.Finish()

	// 删除笔记形状
	if err := s.noteShapeDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记形状失败", "error", err)
		return false, errors.Biz("note.note_shape.errors.delete_failed")
	}

	return true, nil
}

// SaveNoteShape 创建笔记形状
func (s *NoteShapeService) SaveNoteShape(ctx context.Context, noteId string, shapeAnnotation *commonPb.ShapeAnnotation) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.CreateNoteShape")
	defer span.Finish()

	// 实现保存笔记形状的逻辑
	noteShape := &model.NoteShape{
		NoteId:      noteId,
		UUID:        shapeAnnotation.Uuid,
		Type:        shapeAnnotation.Type.String(),
		X:           shapeAnnotation.X,
		Y:           shapeAnnotation.Y,
		StrokeColor: shapeAnnotation.StrokeColor.String(),
		Width:       shapeAnnotation.Width,
		Height:      shapeAnnotation.Height,
		RadiusX:     shapeAnnotation.RadiusX,
		RadiusY:     shapeAnnotation.RadiusY,
		EndX:        shapeAnnotation.EndX,
		EndY:        shapeAnnotation.EndY,
		PageNumber:  int(shapeAnnotation.PageNumber),
	}
	noteShape.Id = idgen.GenerateUUID()

	noteShapeId, err := s.CreateNoteShape(ctx, noteShape)
	if err != nil {
		s.logger.Error("msg", "保存笔记形状失败", "error", err.Error())
		return "0", errors.Biz("note.note_shape.errors.create_failed")
	}

	return noteShapeId, nil
}

// UpdateNoteShape 更新笔记形状
func (s *NoteShapeService) UpdateNoteShapeAnnotations(ctx context.Context, annotations []*commonPb.ShapeAnnotation) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteShapeService.UpdateNoteShapeAnnotations")
	defer span.Finish()

	// TODO: 实现更新笔记形状的逻辑
	for _, annotation := range annotations {
		// 获取笔记形状
		noteShape, err := s.noteShapeDAO.FindById(ctx, *annotation.ShapeId)
		if err != nil {
			s.logger.Error("get note shape failed", "error", err)
			return false, errors.Biz("note.note_shape.errors.get_failed")
		}

		if noteShape == nil {
			return false, errors.Biz("note.note_shape.errors.not_found")
		}

		// 更新笔记形状
		noteShape.Type = annotation.Type.String()
		noteShape.X = annotation.X
		noteShape.Y = annotation.Y
		noteShape.StrokeColor = annotation.StrokeColor.String()
		noteShape.Width = annotation.Width
		noteShape.Height = annotation.Height
		noteShape.RadiusX = annotation.RadiusX
		noteShape.RadiusY = annotation.RadiusY
		noteShape.EndX = annotation.EndX
		noteShape.EndY = annotation.EndY
		noteShape.PageNumber = int(annotation.PageNumber)

		if err := s.noteShapeDAO.UpdateById(ctx, noteShape); err != nil {
			s.logger.Error("update note shape failed", "error", err)
			return false, errors.Biz("note.note_shape.errors.update_failed")
		}
	}
	return true, nil
}
