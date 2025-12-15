package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/note/dao"
	"github.com/yb2020/odoc/services/note/model"
)

// PaperNoteAccessService 论文笔记访问记录服务实现
type PaperNoteAccessService struct {
	paperNoteAccessDAO *dao.PaperNoteAccessDAO
	logger             logging.Logger
	tracer             opentracing.Tracer
}

// NewPaperNoteAccessService 创建新的论文笔记访问记录服务
func NewPaperNoteAccessService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	paperNoteAccessDAO *dao.PaperNoteAccessDAO,
) *PaperNoteAccessService {
	return &PaperNoteAccessService{
		paperNoteAccessDAO: paperNoteAccessDAO,
		logger:             logger,
		tracer:             tracer,
	}
}

// CreatePaperNoteAccess 创建论文笔记访问记录
func (s *PaperNoteAccessService) CreatePaperNoteAccess(ctx context.Context, access *model.PaperNoteAccess) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteAccessService.CreatePaperNoteAccess")
	defer span.Finish()

	if err := s.paperNoteAccessDAO.Create(ctx, access); err != nil {
		s.logger.Error("创建论文笔记访问记录失败", "error", err)
		return "0", errors.Biz("note.paper_note_access.errors.create_failed")
	}

	return access.Id, nil
}

// GetPaperNoteAccessById 根据ID获取论文笔记访问记录
func (s *PaperNoteAccessService) GetPaperNoteAccessById(ctx context.Context, id string) (*model.PaperNoteAccess, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteAccessService.GetPaperNoteAccessById")
	defer span.Finish()

	// 获取论文笔记访问记录
	access, err := s.paperNoteAccessDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取论文笔记访问记录失败", "error", err)
		return nil, errors.Biz("note.paper_note_access.errors.get_failed")
	}

	if access == nil {
		return nil, errors.Biz("note.paper_note_access.errors.not_found")
	}

	return access, nil
}

// GetPaperNoteAccessByNoteId 根据笔记ID获取论文笔记访问记录
func (s *PaperNoteAccessService) GetPaperNoteAccessByNoteId(ctx context.Context, noteId string) (*model.PaperNoteAccess, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteAccessService.GetPaperNoteAccessByNoteId")
	defer span.Finish()

	// 获取论文笔记访问记录
	access, err := s.paperNoteAccessDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("根据笔记ID获取论文笔记访问记录失败", "noteId", noteId, "error", err)
		return nil, errors.Biz("note.paper_note_access.errors.get_by_note_id_failed")
	}

	return access, nil
}

// UpdatePaperNoteAccess 更新论文笔记访问记录
func (s *PaperNoteAccessService) UpdatePaperNoteAccess(ctx context.Context, access *model.PaperNoteAccess) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteAccessService.UpdatePaperNoteAccess")
	defer span.Finish()

	// 获取论文笔记访问记录
	existingAccess, err := s.paperNoteAccessDAO.FindById(ctx, access.Id)
	if err != nil {
		s.logger.Error("获取论文笔记访问记录失败", "error", err)
		return false, errors.Biz("note.paper_note_access.errors.get_failed")
	}

	if existingAccess == nil {
		return false, errors.Biz("note.paper_note_access.errors.not_found")
	}

	// 更新论文笔记访问记录
	if err := s.paperNoteAccessDAO.UpdateById(ctx, access); err != nil {
		s.logger.Error("更新论文笔记访问记录失败", "error", err)
		return false, errors.Biz("note.paper_note_access.errors.update_failed")
	}

	return true, nil
}

// DeletePaperNoteAccessById 删除论文笔记访问记录
func (s *PaperNoteAccessService) DeletePaperNoteAccessById(ctx context.Context, id string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteAccessService.DeletePaperNoteAccessById")
	defer span.Finish()

	// 删除论文笔记访问记录
	if err := s.paperNoteAccessDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除论文笔记访问记录失败", "error", err)
		return errors.Biz("note.paper_note_access.errors.delete_failed")
	}

	return nil
}

// DeletePaperNoteAccessByNoteId 根据笔记ID删除论文笔记访问记录
func (s *PaperNoteAccessService) DeletePaperNoteAccessByNoteId(ctx context.Context, noteId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PaperNoteAccessService.DeletePaperNoteAccessByNoteId")
	defer span.Finish()

	// 删除论文笔记访问记录
	if err := s.paperNoteAccessDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("根据笔记ID删除论文笔记访问记录失败", "noteId", noteId, "error", err)
		return errors.Biz("note.paper_note_access.errors.delete_by_note_id_failed")
	}

	return nil
}
