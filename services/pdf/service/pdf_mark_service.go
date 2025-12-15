package service

import (
	"context"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc-proto/gen/go/common"
	notePb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	proto "github.com/yb2020/odoc/services/note/proto"
	"github.com/yb2020/odoc/services/pdf/dao"
	"github.com/yb2020/odoc/services/pdf/dto"
	pdfInterfaces "github.com/yb2020/odoc/services/pdf/interfaces"
	"github.com/yb2020/odoc/services/pdf/model"
)

// PdfMarkService PDF标记服务实现
// 实现 interfaces.IPdfMarkService 接口
type PdfMarkService struct {
	pdfMarkDAO                *dao.PdfMarkDAO
	logger                    logging.Logger
	tracer                    opentracing.Tracer
	pdfMarkTagRelationService *PdfMarkTagRelationService
	pdfMarkTagService         *PdfMarkTagService
	paperNoteService          noteInterfaces.IPaperNoteService
	paperPdfService           pdfInterfaces.IPaperPdfService
}

// NewPdfMarkService 创建新的PDF标记服务
func NewPdfMarkService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	pdfMarkDAO *dao.PdfMarkDAO,
	pdfMarkTagRelationService *PdfMarkTagRelationService,
	pdfMarkTagService *PdfMarkTagService,
	paperNoteService noteInterfaces.IPaperNoteService,
	paperPdfService pdfInterfaces.IPaperPdfService,

) *PdfMarkService {
	return &PdfMarkService{
		pdfMarkDAO:                pdfMarkDAO,
		logger:                    logger,
		tracer:                    tracer,
		pdfMarkTagRelationService: pdfMarkTagRelationService,
		pdfMarkTagService:         pdfMarkTagService,
		paperNoteService:          paperNoteService,
		paperPdfService:           paperPdfService,
	}
}

// SavePdfMark 保存PDF标记
func (s *PdfMarkService) SavePdfMark(ctx context.Context, mark *model.PdfMark) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.SavePdfMark")
	defer span.Finish()

	mark.Id = idgen.GenerateUUID()
	if err := s.pdfMarkDAO.Create(ctx, mark); err != nil {
		s.logger.Error("保存PDF标记失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark.errors.save_failed")
	}

	return mark.Id, nil
}

// BatchSavePdfMarks 批量保存PDF标记
func (s *PdfMarkService) BatchSavePdfMarks(ctx context.Context, marks []*model.PdfMark) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.BatchSavePdfMarks")
	defer span.Finish()

	if len(marks) == 0 {
		return "0", nil
	}
	// 批量保存
	if err := s.pdfMarkDAO.BatchCreate(ctx, marks); err != nil {
		s.logger.Error("批量保存PDF标记失败", "error", err, "count", len(marks))
		return "0", errors.Biz("pdf.pdf_mark.errors.batch_save_failed")
	}

	return "0", nil
}

// UpdatePdfMark 更新PDF标记
func (s *PdfMarkService) UpdatePdfMark(ctx context.Context, mark *model.PdfMark) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.UpdatePdfMark")
	defer span.Finish()

	// 获取PDF标记
	existingMark, err := s.pdfMarkDAO.FindById(ctx, mark.Id)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	if existingMark == nil {
		return false, errors.Biz("pdf.pdf_mark.errors.not_found")
	}

	if err := s.pdfMarkDAO.ModifyExcludeNull(ctx, mark); err != nil {
		return false, errors.BizWrap(err.Error(), err)
	}

	return true, nil
}

// DeletePdfMarkById 删除PDF标记
func (s *PdfMarkService) DeletePdfMarkById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.DeletePdfMarkById")
	defer span.Finish()

	// 删除PDF标记
	if err := s.pdfMarkDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除PDF标记失败", "error", err)
		return false, errors.Biz("pdf.pdf_mark.errors.delete_failed")
	}

	return true, nil
}

// GetPdfMarkById 根据ID获取PDF标记
func (s *PdfMarkService) GetPdfMarkById(ctx context.Context, id string) (*model.PdfMark, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetPdfMarkById")
	defer span.Finish()

	// 获取PDF标记
	mark, err := s.pdfMarkDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	return mark, nil
}

// GetPdfMarksById 根据ID获取PDF标记
func (s *PdfMarkService) GetPdfMarksByNoteId(ctx context.Context, noteId string) ([]model.PdfMark, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetPdfMarksById")
	defer span.Finish()

	// 获取PDF标记
	marks, err := s.pdfMarkDAO.GetPdfMarksByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	return marks, nil
}

// GetCountPdfMarksByNoteId 根据笔记ID获取PDF标记总数
func (s *PdfMarkService) GetCountPdfMarksByNoteId(ctx context.Context, noteId string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetCountPdfMarksByNoteId")
	defer span.Finish()

	// 获取PDF标记
	count, err := s.pdfMarkDAO.GetCountPdfMarksByNoteIdWithoutIsHighlight(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return 0, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	return count, nil
}

// GetAnnotationRawModelsByNoteId 根据ID获取PDF标记
func (s *PdfMarkService) GetAnnotationRawModelsByNoteId(ctx context.Context, noteId string) ([]*notePb.AnnotationRawModel, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetAnnotationRawModelsByNoteId")
	defer span.Finish()

	// 获取PDF标记
	marks, err := s.GetPdfMarksByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	annotations := make([]*notePb.AnnotationRawModel, 0, len(marks))
	for _, mark := range marks {
		annotationModelTool := &proto.AnnotationModelTool{}
		annotation, _ := annotationModelTool.ToAnnotationRawModel(&mark)
		annotations = append(annotations, annotation)
	}

	// 修改返回值
	return annotations, nil
}

// GetWebNoteAnnotationModelsByNoteId 根据笔记ID获取WebNoteAnnotationModel列表
func (s *PdfMarkService) GetWebNoteAnnotationModelsByNoteId(ctx context.Context, noteId string) ([]*notePb.WebNoteAnnotationModel, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetWebNoteAnnotationModelsByNoteId")
	defer span.Finish()

	// 查询
	annotateionRawModels, err := s.GetAnnotationRawModelsByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取注释原始模型列表失败", "error", err)
		return nil, errors.Biz("pdf.pdf_mark.errors.get_failed")
	}
	s.logger.Debug("list", annotateionRawModels)

	webAnnotationModels := make([]*notePb.WebNoteAnnotationModel, 0, len(annotateionRawModels))
	for _, annotation := range annotateionRawModels {
		webNoteProtoTransformer := &proto.WebNoteProtoTransformer{}
		webAnnotation, err := webNoteProtoTransformer.WebAnnotationModel(annotation)
		if err != nil || webAnnotation == nil {
			s.logger.Warn("转换注释模型失败", "error", err, "annotation", annotation)
			continue
		}

		//TODO mock
		// 1.设置DeleteAuthority属性判断和设置
		webAnnotation.DeleteAuthority = true
		// 2.设置关联的tags
		tags, _ := s.GetAnnotateTagsByMarkId(ctx, webAnnotation.Id)
		webAnnotation.Tags = tags
		// 3.设置CommentatorInfoView信息

		webAnnotationModels = append(webAnnotationModels, webAnnotation)
	}

	// 修改返回值
	return webAnnotationModels, nil
}

func (s *PdfMarkService) SavePdfMarkByAnnotation(ctx context.Context, annotation *notePb.WebNoteAnnotationModel) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.SavePdfMarkByAnnotation")
	defer span.Finish()

	//1. 判断参数

	//2. 转换proto
	webNoteProtoTransformer := proto.WebNoteProtoTransformer{}
	annotationRawModel, err := webNoteProtoTransformer.AnnotationRawModel(annotation)
	if err != nil {
		return "0", err
	}
	s.logger.Info("annotationRawModel", "annotationRawModel", annotationRawModel)
	//3. BaseNoteProtoChecker.checkAnnotationRawModel(annotation); //检查TODO

	//3. AnnotationModelTool处理，原逻辑需要根据annotation.GroupId分类不同处理，新版本去除了Group功能，只需要处理SINGLE_PERSON逻辑
	annotationModelTool := &proto.AnnotationModelTool{}
	noteId, isPresent := annotationModelTool.GetNoteIdFromAnnotationRawModel(annotationRawModel)
	if !isPresent {
		return "0", errors.Biz("note.id.absent")
	}
	s.logger.Debug("noteId", "noteId", noteId)

	pageNumber, isPresent := annotationModelTool.GetPageNumberFromAnnotationRawModel(annotationRawModel)
	if !isPresent {
		return "0", errors.Biz("note.annotation.page.absent")
	}
	s.logger.Debug("pageNumber", "pageNumber", pageNumber)

	//4. 必须先重排序，然后保存。否则会因为新增了数据，校验不通过

	//5. 构造保存bean TODO AnnotationRawModel --》SavePdfMark BEAN
	//savePdfMarkBean := &bean.SavePdfMarkBean{}
	savePdfMark, err := annotationModelTool.ToPdfMark(annotationRawModel)
	if err != nil {
		s.logger.Error("保存PDF标记失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark.errors.save_failed")
	}
	s.logger.Debug("savePdfMMark", savePdfMark)

	id, err := s.SavePdfMarkByBean(ctx, savePdfMark)
	if err != nil {
		s.logger.Error("保存PDF标记失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark.errors.save_failed")
	}

	return id, nil
}

// SavePdfMark 保存PDF标记
func (s *PdfMarkService) SavePdfMarkByBean(ctx context.Context, pdfMark *model.PdfMark) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.SavePdfMark")
	defer span.Finish()

	pdfMark.Id = idgen.GenerateUUID()
	if err := s.pdfMarkDAO.Create(ctx, pdfMark); err != nil {
		s.logger.Error("保存PDF标记失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark.errors.save_failed")
	}

	// TODO 查找文献信息，同步到userDoc ES

	return pdfMark.Id, nil
}

// DeleteByAnnotationPointer 删除PDF标记Annotation
func (s *PdfMarkService) DeleteByAnnotationPointer(ctx context.Context, annotationPointer *common.AnnotationPointer) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.DeleteByAnnotationPointer")
	defer span.Finish()

	// 新版本删除不用处理group相关的区分，已经去除group分组概念
	// 1.删除PDF标记
	if err := s.pdfMarkDAO.DeleteById(ctx, annotationPointer.Id); err != nil {
		s.logger.Error("删除PDF标记Annotation", "error", err)
		return false, errors.Biz("pdf.pdf_mark.errors.delete_failed")
	}

	// 2.markUnusedCdnUrl(id); TODO

	// 3.删除mark tag relation关系
	_, err := s.DeleteMarkTagRelationsById(ctx, annotationPointer.Id)
	if err != nil {
		s.logger.Error("删除PDF标记Annotation", "error", err)
		return false, errors.Biz("pdf.pdf_mark.errors.delete_failed")
	}

	return true, nil
}

// 通过markId删除marktag关联
func (s *PdfMarkService) DeleteMarkTagRelationsById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.DeleteMarkTagRelationsById")
	defer span.Finish()

	return s.pdfMarkTagRelationService.DeleteByMarkId(ctx, id)

}

func (s *PdfMarkService) UpdatePdfMarkByAnnotation(ctx context.Context, annotation *notePb.WebNoteAnnotationModel) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.UpdatePdfMarkByAnnotation")
	defer span.Finish()

	//1. 判断参数

	//2. 转换proto
	webNoteProtoTransformer := proto.WebNoteProtoTransformer{}
	annotationRawModel, err := webNoteProtoTransformer.AnnotationRawModel(annotation)
	if err != nil {
		return "0", err
	}

	//3. AnnotationRawModel --》SavePdfMark
	annotationModelTool := &proto.AnnotationModelTool{}
	updatePdfMark, err := annotationModelTool.ToPdfMark(annotationRawModel)
	if err != nil {
		return "0", err
	}
	s.logger.Info("annotationRawModel", "annotationRawModel", annotationRawModel)

	return s.UpdatePdfMarkByBean(ctx, updatePdfMark)
}

// UpdatePdfMarkByBean 更新PDF标记
func (s *PdfMarkService) UpdatePdfMarkByBean(ctx context.Context, updateMark *model.PdfMark) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.UpdatePdfMarkByBean")
	defer span.Finish()

	// 获取PDF标记
	mark, err := s.pdfMarkDAO.FindExistById(ctx, updateMark.Id)
	if err != nil {
		s.logger.Error("获取PDF标记失败", "error", err)
		return "0", errors.Biz("pdf.pdf_mark.errors.get_failed")
	}

	if mark == nil {
		return "0", errors.Biz("pdf.pdf_mark.errors.not_found")
	}

	// TODO: annotationRaw to Mark fill UserId
	// 1.PdfMark权限判断
	// if updateMark.CreatorId != mark.CreatorId {
	// 	return false, errors.Biz("the content of the note not belong to you")
	// }

	// 2.PaperNote权限判断
	paperNote, err := s.paperNoteService.GetPaperNoteById(ctx, updateMark.NoteId)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return "0", errors.Biz("pdf.paper_pdf.errors.get_failed")
	}
	if paperNote == nil || paperNote.CreatorId != mark.CreatorId {
		return "0", errors.Biz("this note does not belong to you")
	}

	// 3.PaperPdf权限判断
	paperPdf, err := s.paperPdfService.GetById(ctx, paperNote.PdfId)
	if err != nil {
		s.logger.Error("获取论文PDF失败", "error", err)
		return "0", errors.Biz("pdf.paper_pdf.errors.get_failed")
	}
	if paperPdf == nil {
		return "0", errors.Biz("pdf.paper_pdf.errors.not_found")
	}

	// 4.设置更新属性
	mark.Idea = updateMark.Idea
	mark.HtmlIdea = updateMark.HtmlIdea
	mark.Content = updateMark.Content
	mark.PicUrl = updateMark.PicUrl
	mark.KeyContent = updateMark.KeyContent

	mark.StyleId = updateMark.StyleId
	mark.Page = updateMark.Page
	mark.IsHighlight = updateMark.IsHighlight
	if updateMark.Sort > 0 {
		mark.Sort = updateMark.Sort
	}
	mark.RectContent = updateMark.RectContent
	mark.TextBoxContent = updateMark.TextBoxContent
	mark.CommentContent = updateMark.CommentContent

	if err := s.pdfMarkDAO.ModifyExcludeNull(ctx, mark); err != nil {
		return "0", errors.BizWrap(err.Error(), err)
	}

	// 5. TODO sync markInfo to ES

	return mark.Id, nil
}

// GetAnnotateTagsByMarkId 根据markId获取标签列表
func (s *PdfMarkService) GetAnnotateTagsByMarkId(ctx context.Context, markId string) ([]*common.AnnotateTag, error) {
	tags, err := s.pdfMarkTagService.GetByMarkId(ctx, markId)
	if err != nil {
		return []*common.AnnotateTag{}, errors.Biz("get annotate tags failed")
	}
	if len(tags) == 0 {
		return []*common.AnnotateTag{}, nil
	}

	var annotateTags []*common.AnnotateTag
	for _, tag := range tags {
		annotateTag := &common.AnnotateTag{
			TagId:         tag.Id,
			TagName:       tag.TagName,
			LatestUseTime: uint64(tag.Idempotent),
		}
		annotateTags = append(annotateTags, annotateTag)
	}

	return annotateTags, nil
}

// GetMarkTagInfosByNoteIds 根据noteIds获取用户markTagInfos
func (s *PdfMarkService) GetMarkTagInfosByNoteIds(ctx context.Context, userId string, noteIds []string) ([]*common.PdfMarkTagInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetMarkTagInfosByNoteIds")
	defer span.Finish()

	marks, err := s.pdfMarkDAO.GetPdfMarksByNoteIds(ctx, userId, noteIds)
	if err != nil {
		return nil, errors.Biz("get pdf marks failed")
	}
	if len(marks) == 0 {
		return []*common.PdfMarkTagInfo{}, nil
	}

	var markIds []string
	for _, mark := range marks {
		markIds = append(markIds, mark.Id)
	}

	tags, err := s.pdfMarkTagService.GetByMarkIds(ctx, markIds)
	if err != nil {
		return nil, errors.Biz("get pdf mark tag infos failed")
	}
	if len(tags) == 0 {
		return []*common.PdfMarkTagInfo{}, nil
	}

	var markTagInfos []*common.PdfMarkTagInfo
	for _, tag := range tags {
		markTagInfo := &common.PdfMarkTagInfo{
			TagId:      tag.Id,
			TagName:    tag.TagName,
			MarkCount:  0,
			ModifyDate: uint64(tag.UpdatedAt.UnixMilli()),
		}
		markTagInfos = append(markTagInfos, markTagInfo)
	}

	return markTagInfos, nil
}

// GetMarkTagInfosByFolderId 根据folderId获取用户markTagInfos
func (s *PdfMarkService) GetMarkTagInfosByFolderId(ctx context.Context, userId string, folderId string) ([]*common.PdfMarkTagInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetMarkTagInfosByFolderId")
	defer span.Finish()

	folderInfo, err := s.paperNoteService.GetAllNoteDocInfosByFolderId(ctx, userId, folderId)
	if err != nil {
		return nil, errors.Biz("get pdf marks failed")
	}
	if folderInfo == nil {
		return []*common.PdfMarkTagInfo{}, nil
	}

	var noteIds []string
	for _, noteInfo := range folderInfo {
		noteIds = append(noteIds, noteInfo.NoteId)
	}

	return s.GetMarkTagInfosByNoteIds(ctx, userId, noteIds)
}

// GetUserPdfMarkPage获取用户笔记标记分页
func (s *PdfMarkService) GetUserPdfMarkPage(ctx context.Context, userId string, searchDto *dto.PdfMarkSearchPageDto) ([]*notePb.WebNoteAnnotationModel, int32, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfMarkService.GetUserPdfMarkPage")
	defer span.Finish()

	//1.获取文件夹下的文献
	docInfos, err := s.paperNoteService.GetAllNoteDocInfosByFolderId(ctx, userId, searchDto.FolderId)
	if err != nil {
		return nil, 0, errors.Biz("get pdf marks failed")
	}
	if docInfos == nil {
		return []*notePb.WebNoteAnnotationModel{}, 0, nil
	}

	var noteIds []string
	//2.筛选文献。如果searchDto.DocId=0，表示获取所有文献的笔记标记
	if searchDto.DocId == "0" {
		for _, noteInfo := range docInfos {
			noteIds = append(noteIds, noteInfo.NoteId)
		}
	} else {
		for _, noteInfo := range docInfos {
			if noteInfo.DocId == searchDto.DocId {
				noteIds = append(noteIds, noteInfo.NoteId)
				continue
			}
		}
	}

	//查询排序方式
	sortExpression := ""
	switch searchDto.SortType {
	case 0:
		sortExpression = "created_at DESC"
	case 1:
		sortExpression = "created_at ASC"
	case 2:
		sortExpression = "note_id DESC"
	default:
		sortExpression = "created_at DESC"
	}

	//3.根据noteIds获取mark
	marks, err := s.pdfMarkDAO.GetPdfMarksByNoteIdsWithSortWithoutIsHighlight(ctx, userId, noteIds, sortExpression)
	if err != nil {
		return nil, 0, errors.Biz("get pdf marks failed")
	}
	if len(marks) == 0 {
		return []*notePb.WebNoteAnnotationModel{}, 0, nil
	}

	//4.筛选标记marks数组，SearchContent不为空时，如果SearchContent中间有空格需拆分keywords数组，单个keywords进行模糊查询mark的idea或keyContent包含SearchContent。过滤掉不符合的mark
	searchDto.SearchContent = strings.TrimSpace(searchDto.SearchContent)
	if searchDto.SearchContent != "" {
		filteredMarks := make([]model.PdfMark, 0)
		searchContent := strings.ToLower(searchDto.SearchContent)
		keywords := strings.Split(searchContent, " ")
		for _, mark := range marks {
			// Check if Idea or KeyContent contains the search content (case-insensitive)
			for _, keyword := range keywords {
				if strings.Contains(strings.ToLower(mark.Idea), keyword) || strings.Contains(strings.ToLower(mark.KeyContent), keyword) {
					//mark.KeyContent字段中搜索到keyword，将keyword高亮
					mark.CommentContent = strings.ReplaceAll(mark.CommentContent, keyword, "<em>"+keyword+"</em>")
					filteredMarks = append(filteredMarks, mark)
					break
				}
			}
		}
		marks = filteredMarks // Update marks with the filtered results
		if len(marks) == 0 {
			return []*notePb.WebNoteAnnotationModel{}, 0, nil
		}
	}

	//---marks转换为webAnnotationModels---//
	webAnnotationModels := make([]*notePb.WebNoteAnnotationModel, 0, len(marks))
	for _, mark := range marks {
		annotationModelTool := &proto.AnnotationModelTool{}
		annotation, _ := annotationModelTool.ToAnnotationRawModel(&mark)

		webNoteProtoTransformer := &proto.WebNoteProtoTransformer{}
		webAnnotation, _ := webNoteProtoTransformer.WebAnnotationModel(annotation)

		// 1.设置DeleteAuthority属性判断和设置
		webAnnotation.DeleteAuthority = true
		// 2.设置关联的tags
		//var markId int64 = 2840323191133811456
		tags, _ := s.GetAnnotateTagsByMarkId(ctx, webAnnotation.Id)
		webAnnotation.Tags = tags
		// 3.设置CommentatorInfoView信息
		// 4.设置DocName信息
		doc, _ := s.paperNoteService.GetDocByNoteId(ctx, webAnnotation.NoteId)
		if doc != nil {
			webAnnotation.DocName = doc.DocName
		}
		webAnnotationModels = append(webAnnotationModels, webAnnotation)
	}

	if len(webAnnotationModels) == 0 {
		return []*notePb.WebNoteAnnotationModel{}, 0, nil
	}

	//查询条件: 当searchDto.tagIdList不为空时，筛选webAnnotationModels，只保留tags字段中tagId包含searchDto.tagIdList的项
	if len(searchDto.TagIdList) > 0 {
		filteredModels := make([]*notePb.WebNoteAnnotationModel, 0)
		tagIdSet := make(map[string]struct{}, len(searchDto.TagIdList))
		for _, tagId := range searchDto.TagIdList {
			tagIdSet[tagId] = struct{}{}
		}

		for _, model := range webAnnotationModels {
			for _, tag := range model.Tags {
				if _, ok := tagIdSet[tag.TagId]; ok {
					filteredModels = append(filteredModels, model)
					break // Found a matching tag, no need to check other tags for this model
				}
			}
		}
		webAnnotationModels = filteredModels
	}

	if len(webAnnotationModels) == 0 {
		return []*notePb.WebNoteAnnotationModel{}, 0, nil
	}

	//查询条件: 当searchDto.styleIdList不为空时，筛选webAnnotationModels，只保留rect/select字段中styleId包含searchDto.styleIdList的项
	if len(searchDto.StyleIdList) > 0 {
		filteredModels := make([]*notePb.WebNoteAnnotationModel, 0)
		styleIdSet := make(map[int64]struct{}, len(searchDto.StyleIdList))
		for _, styleId := range searchDto.StyleIdList {
			styleIdSet[styleId] = struct{}{}
		}

		for _, model := range webAnnotationModels {
			if model.Rect != nil {
				//styleIdSet包含model.Rect.StyleId
				if _, ok := styleIdSet[int64(model.Rect.StyleId)]; ok {
					filteredModels = append(filteredModels, model)
				}
			}
			if model.Select != nil {
				//styleIdSet包含model.Select.StyleId
				if _, ok := styleIdSet[int64(model.Select.StyleId)]; ok {
					filteredModels = append(filteredModels, model)
				}
			}
		}
		webAnnotationModels = filteredModels
	}

	if len(webAnnotationModels) == 0 {
		return []*notePb.WebNoteAnnotationModel{}, 0, nil
	}

	//---此前代码用于获取数据和查询过滤数据，之后的过程用于构建响应数据---//
	//4.计算总记录数
	total := int32(len(webAnnotationModels))

	//5.计算分页的起始和结束索引
	startIndex := (searchDto.CurrentPage - 1) * searchDto.PageSize
	endIndex := startIndex + searchDto.PageSize

	//6.确保索引不越界
	if startIndex < 0 {
		startIndex = 0
	}
	if startIndex > total {
		startIndex = total
	}
	if endIndex > total {
		endIndex = total
	}

	//7.应用分页
	webAnnotationModels = webAnnotationModels[startIndex:endIndex]
	return webAnnotationModels, total, nil
}
