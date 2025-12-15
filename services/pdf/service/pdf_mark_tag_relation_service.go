package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfMarkTagRelationService PDF标记标签关系服务实现
type PdfMarkTagRelationService struct {
	pdfMarkTagRelationDAO *dao.PdfMarkTagRelationDAO
	logger                logging.Logger
	tracer                opentracing.Tracer
}

// NewPdfMarkTagRelationService 创建新的PDF标记标签关系服务
func NewPdfMarkTagRelationService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfMarkTagRelationDAO *dao.PdfMarkTagRelationDAO,
) *PdfMarkTagRelationService {
	return &PdfMarkTagRelationService{
		pdfMarkTagRelationDAO: pdfMarkTagRelationDAO,
		logger:                logger,
		tracer:                tracer,
	}
}

// Save 创建PDF标记标签关系
func (s *PdfMarkTagRelationService) Save(ctx context.Context, relation *model.PdfMarkTagRelation) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.CreatePdfMarkTagRelation")
	defer span.Finish()

	if err := s.pdfMarkTagRelationDAO.Create(ctx, relation); err != nil {
		s.logger.Error("创建PDF标记标签关系失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark_tag_relation.errors.create_failed")
	}

	return relation.Id, nil
}

// Update 更新PDF标记标签关系
func (s *PdfMarkTagRelationService) Update(ctx context.Context, relation *model.PdfMarkTagRelation) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.UpdatePdfMarkTagRelation")
	defer span.Finish()

	// 获取PDF标记标签关系
	existingRelation, err := s.pdfMarkTagRelationDAO.FindById(ctx, relation.Id)
	if err != nil {
		s.logger.Error("获取PDF标记标签关系失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
	}

	if existingRelation == nil {
		return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.not_found")
	}

	if err := s.pdfMarkTagRelationDAO.ModifyExcludeNull(ctx, relation); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeleteById 删除PDF标记标签关系
func (s *PdfMarkTagRelationService) DeleteById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.DeletePdfMarkTagRelationById")
	defer span.Finish()

	// 删除PDF标记标签关系
	if err := s.pdfMarkTagRelationDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF标记标签关系失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.delete_failed")
	}

	return true, nil
}

// GetById 根据ID获取PDF标记标签关系
func (s *PdfMarkTagRelationService) GetById(ctx context.Context, id string) (*model.PdfMarkTagRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.GetPdfMarkTagRelationById")
	defer span.Finish()

	// 获取PDF标记标签关系
	relation, err := s.pdfMarkTagRelationDAO.FindExistById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF标记标签关系失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
	}

	return relation, nil
}

// DeletePdfMarkTagRelationById 删除PDF标记标签关系
func (s *PdfMarkTagRelationService) DeleteByMarkId(ctx context.Context, markId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.DeletePdfMarkTagRelationById")
	defer span.Finish()

	markTagRelList, err := s.pdfMarkTagRelationDAO.GetByMarkId(ctx, markId)
	if err != nil {
		s.logger.Error("获取PDF标记标签关系失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
	}

	if len(markTagRelList) == 0 {
		return true, nil
	}

	s.logger.Debug("markTagRelList", markTagRelList)
	// 提取标记标签关系的ID列表
	var ids []string
	for _, rel := range markTagRelList {
		ids = append(ids, rel.Id)
	}

	// 删除PDF标记标签关系
	if err := s.pdfMarkTagRelationDAO.DeleteByIds(ctx, ids); err != nil {
		s.logger.Error("删除PDF标记标签关系失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.delete_failed")
	}

	return true, nil
}

// GetByMarkIdAndTagId 根据MarkId和TagId获取PDF标记标签关系
func (s *PdfMarkTagRelationService) GetByMarkIdAndTagId(ctx context.Context, markId string, tagId string) (*model.PdfMarkTagRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.GetByMarkIdAndTagId")
	defer span.Finish()

	// 获取PDF标记标签关系
	relation, err := s.pdfMarkTagRelationDAO.GetByMarkIdAndTagId(ctx, markId, tagId)
	if err != nil {
		s.logger.Error("获取PDF标记标签关系失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
	}

	return relation, nil
}

// AddTagToPdfMark 添加PDF标记标签关系
func (s *PdfMarkTagRelationService) AddTagToPdfMark(ctx context.Context, markId string, tagId string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.AddTagToPdfMark")
	defer span.Finish()
	relation := &model.PdfMarkTagRelation{
		MarkId: markId,
		TagId:  tagId,
	}
	relation.Id = idgen.GenerateUUID()

	return s.Save(ctx, relation)
}

// DeleteByTagId 删除PDF标记标签关系
func (s *PdfMarkTagRelationService) DeleteByTagId(ctx context.Context, tagId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.DeletePdfMarkTagRelationById")
	defer span.Finish()

	// 删除PDF标记标签关系
	if err := s.pdfMarkTagRelationDAO.DeleteByTagId(ctx, tagId); err != nil {
		s.logger.Error("删除PDF标记标签关系失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.delete_failed")
	}

	return true, nil
}

// GetByMarkId 根据MarkId获取PDF标记标签关系
func (s *PdfMarkTagRelationService) GetByMarkId(ctx context.Context, markId string) ([]model.PdfMarkTagRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.GetByMarkId")
	defer span.Finish()

	// 获取PDF标记标签关系
	relations, err := s.pdfMarkTagRelationDAO.GetByMarkId(ctx, markId)
	if err != nil {
		s.logger.Error("获取PDF标记标签关系失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
	}

	return relations, nil
}

// GetByMarkIds 根据MarkIds获取PDF标记标签关系
func (s *PdfMarkTagRelationService) GetByMarkIds(ctx context.Context, markIds []string) ([]model.PdfMarkTagRelation, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagRelationService.GetByMarkIds")
	defer span.Finish()

	// 获取PDF标记标签关系
	relations, err := s.pdfMarkTagRelationDAO.GetByMarkIds(ctx, markIds)
	if err != nil {
		s.logger.Error("获取PDF标记标签关系失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
	}

	return relations, nil
}
