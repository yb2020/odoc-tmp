package service

import (
	"context"

	"github.com/opentracing/opentracing-go"

	pb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	docInterfaces "github.com/yb2020/odoc/services/doc/interfaces"
	"github.com/yb2020/odoc/services/note/dao"
	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/model"
	pdfInterfaces "github.com/yb2020/odoc/services/pdf/interfaces"
)

// NoteSummaryService 笔记摘要服务实现
type NoteSummaryService struct {
	noteSummaryDAO   *dao.NoteSummaryDAO
	logger           logging.Logger
	tracer           opentracing.Tracer
	userDocService   docInterfaces.IUserDocService
	paperNoteService noteInterfaces.IPaperNoteService
	pdfService       pdfInterfaces.IPaperPdfService
}

// NewNoteSummaryService 创建新的笔记摘要服务
func NewNoteSummaryService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteSummaryDAO *dao.NoteSummaryDAO,
	userDocService docInterfaces.IUserDocService,
	paperNoteService noteInterfaces.IPaperNoteService,
	pdfService pdfInterfaces.IPaperPdfService,
) *NoteSummaryService {
	return &NoteSummaryService{
		noteSummaryDAO:   noteSummaryDAO,
		logger:           logger,
		tracer:           tracer,
		userDocService:   userDocService,
		paperNoteService: paperNoteService,
		pdfService:       pdfService,
	}
}

// SetPaperPdfService 设置论文PDF服务
func (s *NoteSummaryService) SetPaperPdfService(pdfService pdfInterfaces.IPaperPdfService) error {
	if pdfService == nil {
		return errors.Biz("pdfService cannot be nil")
	}
	s.pdfService = pdfService
	return nil
}

// CreateNoteSummary 创建笔记摘要
func (s *NoteSummaryService) CreateNoteSummary(ctx context.Context, summary *model.NoteSummary) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.CreateNoteSummary")
	defer span.Finish()

	if err := s.noteSummaryDAO.Create(ctx, summary); err != nil {
		s.logger.Error("创建笔记摘要失败", "error", err)
		return "0", errors.Biz("note.note_summary.errors.create_failed")
	}

	return summary.Id, nil
}

// UpdateNoteSummary 更新笔记摘要
func (s *NoteSummaryService) UpdateNoteSummary(ctx context.Context, summary *model.NoteSummary) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.UpdateNoteSummary")
	defer span.Finish()

	// 获取笔记摘要
	existingSummary, err := s.noteSummaryDAO.FindById(ctx, summary.Id)
	if err != nil {
		s.logger.Error("获取笔记摘要失败", "error", err)
		return false, errors.Biz("note.note_summary.errors.get_failed")
	}

	if existingSummary == nil {
		return false, errors.Biz("note.note_summary.errors.not_found")
	}

	// 更新笔记摘要
	if err := s.noteSummaryDAO.UpdateById(ctx, summary); err != nil {
		s.logger.Error("更新笔记摘要失败", "error", err)
		return false, errors.Biz("note.note_summary.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteSummaryById 删除笔记摘要
func (s *NoteSummaryService) DeleteNoteSummaryById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.DeleteNoteSummaryById")
	defer span.Finish()

	// 删除笔记摘要
	if err := s.noteSummaryDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记摘要失败", "error", err)
		return false, errors.Biz("note.note_summary.errors.delete_failed")
	}

	return true, nil
}

// GetNoteSummaryById 根据ID获取笔记摘要
func (s *NoteSummaryService) GetNoteSummaryById(ctx context.Context, id string) (*model.NoteSummary, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.GetNoteSummaryById")
	defer span.Finish()

	// 获取笔记摘要
	summary, err := s.noteSummaryDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记摘要失败", "error", err)
		return nil, errors.Biz("note.note_summary.errors.get_failed")
	}

	if summary == nil {
		return nil, errors.Biz("note.note_summary.errors.not_found")
	}

	return summary, nil
}

// GetNoteSummaryByNoteId 根据笔记ID获取笔记摘要
func (s *NoteSummaryService) GetNoteSummaryByNoteId(ctx context.Context, noteId string) (*model.NoteSummary, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.GetNoteSummaryByNoteId")
	defer span.Finish()

	// 获取笔记摘要
	summary, err := s.noteSummaryDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记摘要失败", "error", err)
		return nil, errors.Biz("note.note_summary.errors.get_failed")
	}

	return summary, nil
}

// GetNoteSummariesByUserId 根据用户ID获取笔记摘要列表
func (s *NoteSummaryService) GetNoteSummariesByUserId(ctx context.Context, userId string) ([]model.NoteSummary, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.GetNoteSummariesByUserId")
	defer span.Finish()

	// 获取笔记摘要列表
	summaries, err := s.noteSummaryDAO.GetByUserId(ctx, userId)
	if err != nil {
		s.logger.Error("获取笔记摘要列表失败", "error", err)
		return nil, errors.Biz("note.note_summary.errors.list_failed")
	}

	return summaries, nil
}

// DeleteNoteSummaryByNoteId 根据笔记ID删除笔记摘要
func (s *NoteSummaryService) DeleteNoteSummaryByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.DeleteNoteSummaryByNoteId")
	defer span.Finish()

	// 删除笔记摘要
	if err := s.noteSummaryDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记摘要失败", "error", err)
		return false, errors.Biz("note.note_summary.errors.delete_failed")
	}

	return true, nil
}

// 笔记管理TAB——总结
func (s *NoteSummaryService) GetSummaryListByUserId(ctx context.Context, userId string) (*pb.GetSummaryListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.GetListByUserId")
	defer span.Finish()

	folderInfo, err := s.paperNoteService.GetNoteDocTreeNodeByFolderId(ctx, userId, "0")
	if err != nil {
		s.logger.Error("获取笔记列表失败", "error", err, "userId", userId)
		return nil, errors.Biz("note.paper_note.errors.list_failed")
	}

	if folderInfo == nil {
		return &pb.GetSummaryListResponse{
			Total:                0,
			FolderInfos:          nil,
			UnclassifiedDocInfos: nil,
		}, nil
	}

	//对folderInfo过滤掉Count为0文件夹或子文件夹
	s.filterEmptySubFold(folderInfo)

	return &pb.GetSummaryListResponse{
		Total:                uint32(folderInfo.Count),
		FolderInfos:          folderInfo.ChildrenFolders,
		UnclassifiedDocInfos: folderInfo.DocInfos,
	}, nil
}

// filterEmptySubFold 递归过滤掉 count 为 0 的文件夹
func (s *NoteSummaryService) filterEmptySubFold(folder *pb.NoteManageFolderInfo) {
	if folder == nil {
		return
	}

	// 递归过滤子文件夹
	filteredChildren := make([]*pb.NoteManageFolderInfo, 0, len(folder.ChildrenFolders))
	for _, child := range folder.ChildrenFolders {
		s.filterEmptySubFold(child)
		// 如果子文件夹在过滤后仍有内容（其Count > 0），则保留
		if child.Count > 0 {
			filteredChildren = append(filteredChildren, child)
		}
	}
	folder.ChildrenFolders = filteredChildren
}

// 笔记管理TAB——摘要
func (s *NoteSummaryService) GetExtractListByUserId(ctx context.Context, userId string) (*pb.GetExtractListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.GetExtractListByUserId")
	defer span.Finish()

	folderInfo, err := s.paperNoteService.GetNoteDocTreeNodeByFolderId(ctx, userId, "0")
	if err != nil {
		s.logger.Error("获取笔记列表失败", "error", err, "userId", userId)
		return nil, errors.Biz("note.paper_note.errors.list_failed")
	}

	if folderInfo == nil {
		return &pb.GetExtractListResponse{
			Total:                0,
			FolderInfos:          nil,
			UnclassifiedDocInfos: nil,
		}, nil
	}

	//对folderInfo过滤掉NoteAnnotateCount为0的文件和文件夹或子文件夹
	s.filterEmptyItems(folderInfo)

	return &pb.GetExtractListResponse{
		Total:                uint32(folderInfo.NoteAnnotateCount),
		FolderInfos:          folderInfo.ChildrenFolders,
		UnclassifiedDocInfos: folderInfo.DocInfos,
	}, nil
}

// filterEmptyItems 递归过滤掉 noteAnnotateCount 为 0 的文档和文件夹
func (s *NoteSummaryService) filterEmptyItems(folder *pb.NoteManageFolderInfo) {
	if folder == nil {
		return
	}

	// 过滤当前文件夹下的文档
	filteredDocs := make([]*pb.NoteManageDocInfo, 0, len(folder.DocInfos))
	for _, doc := range folder.DocInfos {
		if doc.NoteAnnotateCount > 0 {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	folder.DocInfos = filteredDocs

	// 递归过滤子文件夹
	filteredChildren := make([]*pb.NoteManageFolderInfo, 0, len(folder.ChildrenFolders))
	for _, child := range folder.ChildrenFolders {
		s.filterEmptyItems(child)
		// 如果子文件夹在过滤后仍有内容（其NoteAnnotateCount > 0），则保留
		if child.NoteAnnotateCount > 0 {
			filteredChildren = append(filteredChildren, child)
		}
	}
	folder.ChildrenFolders = filteredChildren
}

// 笔记管理TAB——根据folderId获取总结列表
func (s *NoteSummaryService) GetUserSummaryListByFolderId(ctx context.Context, userId string, req *pb.GetSummaryListByFolderIdReq) (*pb.GetSummaryListByFolderIdResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteSummaryService.GetUserSummaryListByFolderId")
	defer span.Finish()

	folderInfo, total, err := s.paperNoteService.GetAllNoteDocInfosPageByFolderId(ctx, userId, *req.FolderId, int32(req.CurrentPage), int32(req.PageSize))
	if err != nil {
		s.logger.Error("获取笔记列表失败", "error", err, "userId", userId)
		return nil, errors.Biz("note.paper_note.errors.list_failed")
	}
	if folderInfo == nil {
		return &pb.GetSummaryListByFolderIdResponse{}, nil
	}
	var resp = &pb.GetSummaryListByFolderIdResponse{
		DocInfos: folderInfo,
		Total:    uint32(total),
	}
	return resp, nil
}
