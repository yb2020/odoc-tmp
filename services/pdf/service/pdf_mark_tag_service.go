package service

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc-proto/gen/go/common"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/interfaces"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfMarkTagService PDF标记标签服务实现
type PdfMarkTagService struct {
	pdfMarkTagDAO             *dao.PdfMarkTagDAO
	logger                    logging.Logger
	tracer                    opentracing.Tracer
	pdfMarkService            interfaces.IPdfMarkService
	pdfMarkTagRelationService *PdfMarkTagRelationService
	paperNoteService          noteInterfaces.IPaperNoteService
}

func (s *PdfMarkTagService) SetPdfMarkService(service interfaces.IPdfMarkService) {
	s.pdfMarkService = service
}

// NewPdfMarkTagService 创建新的PDF标记标签服务
func NewPdfMarkTagService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfMarkTagDAO *dao.PdfMarkTagDAO,
	pdfMarkService interfaces.IPdfMarkService,
	pdfMarkTagRelationService *PdfMarkTagRelationService,
	paperNoteService noteInterfaces.IPaperNoteService,
) *PdfMarkTagService {
	return &PdfMarkTagService{
		pdfMarkTagDAO:             pdfMarkTagDAO,
		logger:                    logger,
		tracer:                    tracer,
		pdfMarkService:            pdfMarkService,
		pdfMarkTagRelationService: pdfMarkTagRelationService,
		paperNoteService:          paperNoteService,
	}
}

// Update 更新PDF标记标签
func (s *PdfMarkTagService) Update(ctx context.Context, tag *model.PdfMarkTag) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.UpdatePdfMarkTag")
	defer span.Finish()

	// 获取PDF标记标签
	existingTag, err := s.pdfMarkTagDAO.FindById(ctx, tag.Id)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	if existingTag == nil {
		return false, errors.Biz("pdf.pdf_mark_tag.errors.not_found")
	}

	if err := s.pdfMarkTagDAO.ModifyExcludeNull(ctx, tag); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeleteById 删除PDF标记标签
func (s *PdfMarkTagService) DeleteById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.DeletePdfMarkTagById")
	defer span.Finish()

	// 删除PDF标记标签
	if err := s.pdfMarkTagDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("delete pdf mark tag failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.delete_failed")
	}

	return true, nil
}

// GetById 根据ID获取PDF标记标签
func (s *PdfMarkTagService) GetById(ctx context.Context, id string) (*model.PdfMarkTag, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.GetPdfMarkTagById")
	defer span.Finish()

	// 获取PDF标记标签
	tag, err := s.pdfMarkTagDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	return tag, nil
}

// GetAllAnnotateTags 获取所有PDF标记标签
func (s *PdfMarkTagService) GetAllAnnotateTags(ctx context.Context) ([]*common.AnnotateTag, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.GetAllAnnotateTags")
	defer span.Finish()

	// 获取PDF标记标签
	tags, err := s.pdfMarkTagDAO.FindExistAll(ctx)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	var annotateTags []*common.AnnotateTag
	for _, tag := range tags {
		annotateTags = append(annotateTags, &common.AnnotateTag{
			TagId:         tag.Id,
			TagName:       tag.TagName,
			LatestUseTime: uint64(tag.UpdatedAt.Unix()),
		})
	}

	return annotateTags, nil
}

func (s *PdfMarkTagService) GetTagsByUserId(ctx context.Context, userId string, onlyUsed bool) ([]*common.AnnotateTag, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.GetTagsByUserId")
	defer span.Finish()

	// 获取PDF标记标签
	tags, err := s.pdfMarkTagDAO.GetTagsByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	// 过滤处理onlyUsed
	if onlyUsed {
		notes, err := s.paperNoteService.GetAllNoteByUserId(ctx, userId)
		if err != nil {
			s.logger.Error("get paper note failed", "error", err)
			return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
		}

		noteIds := make([]string, 0)
		for _, note := range notes {
			noteIds = append(noteIds, note.Id)
		}

		markTagInfos, err := s.pdfMarkService.GetMarkTagInfosByNoteIds(ctx, userId, noteIds)
		if err != nil {
			s.logger.Error("get pdf mark tag failed", "error", err)
			return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
		}

		// 直接从markTagInfos构建返回结果，这些就是已使用的标签
		var annotateTags []*common.AnnotateTag
		for _, markTagInfo := range markTagInfos {
			annotateTags = append(annotateTags, &common.AnnotateTag{
				TagId:         markTagInfo.TagId,
				TagName:       markTagInfo.TagName,
				LatestUseTime: markTagInfo.ModifyDate,
			})
		}
		return annotateTags, nil
	}

	var annotateTags []*common.AnnotateTag
	for _, tag := range tags {
		annotateTags = append(annotateTags, &common.AnnotateTag{
			TagId:         tag.Id,
			TagName:       tag.TagName,
			LatestUseTime: uint64(tag.UpdatedAt.Unix()),
		})
	}

	return annotateTags, nil
}

// SavePdfMarkTag 创建PDF标记标签
func (s *PdfMarkTagService) SavePdfMarkTag(ctx context.Context, markId string, tagName string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.SavePdfMarkTag")
	defer span.Finish()

	tag := &model.PdfMarkTag{
		TagName: tagName,
	}
	tag.Id = idgen.GenerateUUID()

	if err := s.pdfMarkTagDAO.Save(ctx, tag); err != nil {
		s.logger.Error("create pdf mark tag failed", "error", err)
		return "0", errors.Biz("pdf.pdf_mark_tag.errors.create_failed")
	}

	return tag.Id, nil
}

// AddTagIdsToAnnotation 标注批量添加标签
func (s *PdfMarkTagService) AddTagIdsToAnnotation(ctx context.Context, markId string, tagIds []string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.AddTagIdsToAnnotation")
	defer span.Finish()

	// pdfMark的校验
	pdfMark, err := s.pdfMarkService.GetPdfMarkById(ctx, markId)
	if err != nil {
		s.logger.Error("get pdf mark annotate failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}
	if pdfMark == nil {
		return false, errors.Biz("note mark annotate is not exits can not add tags")
	}

	for _, tagId := range tagIds {
		s.logger.Debug("tagId", tagId)
		//TODO 保存关联关系逻辑，如果不存在则保存，存在则忽略
		relation, err := s.pdfMarkTagRelationService.GetByMarkIdAndTagId(ctx, markId, tagId)
		if err != nil {
			s.logger.Error("get pdf mark tag relation failed", "error", err)
			return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
		}
		if relation == nil {
			_, err := s.pdfMarkTagRelationService.AddTagToPdfMark(ctx, markId, tagId)
			if err != nil {
				s.logger.Error("add pdf mark tag relation failed", "error", err)
				return false, errors.Biz("add pdf mark tag relation failed")
			}
		} else {
			//ignore
			s.logger.Info("msg", "mark tag relation is exist, ignore", "markId", markId, "tagId", tagId)
		}
	}

	return true, nil
}

// DeleteTagsToAnnotate 标注批量添加标签
func (s *PdfMarkTagService) DeleteTagsToAnnotate(ctx context.Context, markId string, tagIds []string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.DeleteTagsToAnnotate")
	defer span.Finish()

	// pdfMark的校验
	pdfMark, err := s.pdfMarkService.GetPdfMarkById(ctx, markId)
	if err != nil {
		s.logger.Error("get pdf mark annotate failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}
	if pdfMark == nil {
		return false, errors.Biz("note mark annotate is not exits can not delete tags")
	}

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)
	if pdfMark.CreatorId != userId {
		s.logger.Error("msg", "pdf mark annotate is not belong you")
		return false, errors.Biz("pdf mark annotate is not belong you")
	}

	for _, tagId := range tagIds {
		s.logger.Debug("tagId", tagId)
		relation, err := s.pdfMarkTagRelationService.GetByMarkIdAndTagId(ctx, markId, tagId)
		if err != nil {
			s.logger.Error("get pdf mark tag relation failed", "error", err)
			return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.get_failed")
		}
		if relation != nil {
			_, err := s.pdfMarkTagRelationService.DeleteById(ctx, relation.Id)
			if err != nil {
				s.logger.Error("delete pdf mark tag relation failed", "error", err)
				return false, errors.Biz("delete pdf mark tag relation failed")
			}
		} else {
			//ignore
			s.logger.Info("msg", "mark tag relation is exist, ignore", "markId", markId, "tagId", tagId)
		}

	}

	return true, nil
}

// RenameTagById 更新TagName
func (s *PdfMarkTagService) RenameTagById(ctx context.Context, id string, tagName string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.UpdateById")
	defer span.Finish()

	// 获取PDF标记标签
	tag, err := s.pdfMarkTagDAO.FindExistById(ctx, id)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	if tag == nil {
		s.logger.Info("msg", "pdf mark tag is not exist")
		return false, errors.Biz("pdf mark tag is not exist")
	}

	// 判断修改权限
	userId, _ := userContext.GetUserID(ctx)
	if tag.CreatorId != userId {
		s.logger.Error("msg", "pdf mark annotate is not belong you")
		return false, errors.Biz("no permission edit")
	}

	tag.TagName = tagName

	if err := s.pdfMarkTagDAO.ModifyExcludeNull(ctx, tag); err != nil {
		return false, errors.Biz("update failed")
	}

	return true, nil
}

// DeleteTagById 删除Tag需要创建者权限
func (s *PdfMarkTagService) DeleteTagById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.DeleteTagById")
	defer span.Finish()

	// 获取PDF标记标签
	tag, err := s.pdfMarkTagDAO.FindExistById(ctx, id)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	if tag == nil {
		s.logger.Info("msg", "pdf mark tag is not exist")
		return false, errors.Biz("pdf mark tag is not exist")
	}

	// 判断修改权限
	userId, _ := userContext.GetUserID(ctx)
	if tag.CreatorId != userId {
		s.logger.Error("msg", "no permission delete")
		return false, errors.Biz("no permission delete")
	}

	isDelete, err := s.DeleteById(ctx, id)
	if err != nil {
		s.logger.Error("delete pdf mark tag failed", "error", err)
		return false, errors.Biz("pdf.pdf_mark_tag.errors.delete_failed")
	}
	s.logger.Info("msg", "ok", isDelete)

	if isDelete {
		// 删除所有标注标签关联
		_, err := s.pdfMarkTagRelationService.DeleteByTagId(ctx, id)
		if err != nil {
			s.logger.Error("delete pdf mark tag relation failed", "error", err)
			return false, errors.Biz("pdf.pdf_mark_tag_relation.errors.delete_failed")
		}
	}

	return true, nil

}

// GetByMarkId 根据markId获取PDF标记标签
func (s *PdfMarkTagService) GetByMarkId(ctx context.Context, markId string) ([]model.PdfMarkTag, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.GetByMarkId")
	defer span.Finish()

	relations, err := s.pdfMarkTagRelationService.GetByMarkId(ctx, markId)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	if len(relations) == 0 {
		return []model.PdfMarkTag{}, nil
	}

	s.logger.Debug("relations", relations)
	// 提取标记标签关系的ID列表
	var ids []string
	for _, rel := range relations {
		ids = append(ids, rel.TagId)
	}

	// 获取PDF标记标签列表
	tags, err := s.pdfMarkTagDAO.GetByIds(ctx, ids)
	if err != nil {
		s.logger.Error("get pdf mark tags failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_list_failed")
	}

	return tags, nil
}

// GetByMarkIds 根据markIds获取PDF标记标签
func (s *PdfMarkTagService) GetByMarkIds(ctx context.Context, markIds []string) ([]model.PdfMarkTag, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkTagService.GetByMarkIds")
	defer span.Finish()

	relations, err := s.pdfMarkTagRelationService.GetByMarkIds(ctx, markIds)
	if err != nil {
		s.logger.Error("get pdf mark tag failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_failed")
	}

	if len(relations) == 0 {
		return []model.PdfMarkTag{}, nil
	}

	s.logger.Debug("relations", relations)
	// 提取标记标签关系的ID列表
	var ids []string
	for _, rel := range relations {
		ids = append(ids, rel.TagId)
	}

	// 获取PDF标记标签列表
	tags, err := s.pdfMarkTagDAO.GetByIds(ctx, ids)
	if err != nil {
		s.logger.Error("get pdf mark tags failed", "error", err)
		return nil, errors.Biz("pdf.pdf_mark_tag.errors.get_list_failed")
	}

	return tags, nil
}
