package interfaces

import (
	"context"

	pb "github.com/yb2020/odoc/proto/gen/go/note"
	docModel "github.com/yb2020/odoc/services/doc/model"
	"github.com/yb2020/odoc/services/note/model"
	pdfInterfaces "github.com/yb2020/odoc/services/pdf/interfaces"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

// PaperNoteService 论文笔记服务接口
type IPaperNoteService interface {
	SetNoteWordService(noteWordService INoteWordService) error
	CreatePaperNote(ctx context.Context, note *model.PaperNote) (string, error)
	UpdatePaperNote(ctx context.Context, paperNote *model.PaperNote) (bool, error)
	DeletePaperNoteById(ctx context.Context, id string) (bool, error)
	GetPaperNoteById(ctx context.Context, id string) (*model.PaperNote, error)
	// GetByIds 根据IDS获取论文笔记列表
	GetByIds(ctx context.Context, ids []string) ([]model.PaperNote, error)
	GetPaperNoteIdByPdfIdAndUserId(ctx context.Context, pdfId string, userId string) (string, error)
	GetPaperNoteByPdfIdAndUserId(ctx context.Context, pdfId string, userId string) (*model.PaperNote, error)
	GetPaperNoteBaseInfoById(ctx context.Context, noteId string) (*pb.PaperNoteBaseInfoResponse, error)
	GetByPaperIdAndUserIdOrderByNoteCount(ctx context.Context, paperId string, userId string) (*model.PaperNote, error)
	// GetPdfByNoteId 根据NotId获取Pdf信息
	GetPdfByNoteId(ctx context.Context, noteId string) (*pdfModel.PaperPdf, error)
	// GetDocByNoteId 根据NotId获取文献信息
	GetDocByNoteId(ctx context.Context, noteId string) (*docModel.UserDoc, error)
	// GetNoteManageFolderInfosByDocIds 根据文献Ids获取文献信息
	GetNoteManageFolderInfosByDocIds(ctx context.Context, userId string, docIds []string) ([]*pb.NoteManageFolderInfo, int32, error)
	// GetUnclassifiedDocInfosByDocIds 根据文献Ids获取未分类的文献列表
	GetUnclassifiedDocInfosByDocIds(ctx context.Context, userId string, docIds []string) ([]*pb.NoteManageDocInfo, error)
	// GetNoteDocTreeNodeByFolderId 根据文件夹ID获取文件夹，folderId=0时获取用户根文件夹
	GetNoteDocTreeNodeByFolderId(ctx context.Context, userId string, folderId string) (*pb.NoteManageFolderInfo, error)

	// GetAllNoteDocInfosByFolderId 根据文件夹ID获取文件夹及子文件夹下的所有文献，folderId=0时获取用户的所有文献；返回文献列表
	GetAllNoteDocInfosByFolderId(ctx context.Context, userId string, folderId string) ([]*pb.NoteManageDocInfo, error)
	// GetAllNoteDocInfosPageByFolderId 根据文件夹ID获取文件夹及子文件夹下的所有文献分页，folderId=0时获取用户的所有文献分页；返回文献列表和总记录数
	GetAllNoteDocInfosPageByFolderId(ctx context.Context, userId string, folderId string, page int32, pageSize int32) ([]*pb.NoteManageDocInfo, int64, error)

	// SetPaperPdfService 设置论文PDF服务
	SetPaperPdfService(pdfService pdfInterfaces.IPaperPdfService) error
	// GetOwnerPaperNoteBaseInfo 获取用户自己的笔记基础信息，如果不存在则创建
	GetOwnerPaperNoteBaseInfo(ctx context.Context, userId string, pdfId string) (*pb.PaperNoteBaseInfoResponse, error)
	// SelectByUserIdLimit 根据用户ID获取笔记列表
	SelectByUserIdLimit(ctx context.Context, userId string, limit int) ([]model.PaperNote, error)
	// GetUserPaperNoteCount 根据用户ID获取笔记数量
	GetUserPaperNoteCount(ctx context.Context, userId string) (int64, error)

	// 生成下载笔记标签PDF
	GetDownloadNoteMarkPdf(ctx context.Context, noteId string) ([]byte, error)
	// 生成下载笔记标签Markdown
	GetDownloadNoteMarkMarkdown(ctx context.Context, noteId string) ([]byte, error)

	GetAllNoteByUserId(ctx context.Context, userId string) ([]model.PaperNote, error)
}
