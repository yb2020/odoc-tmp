package interfaces

import (
	"context"

	"github.com/yb2020/odoc/proto/gen/go/common"
	pb "github.com/yb2020/odoc/proto/gen/go/pdf"
	docModel "github.com/yb2020/odoc/services/doc/model"
	"github.com/yb2020/odoc/services/pdf/model"
)

// IPaperPdfService 论文PDF服务接口
type IPaperPdfService interface {
	// SetPdfMarkService a setter method to inject the PdfMarkService dependency.
	SetPdfMarkService(pdfMarkService IPdfMarkService)
	// SavePaperPDF 保存论文PDF
	SavePaperPDF(ctx context.Context, pdf *model.PaperPdf) error
	// ModifyPaperPDF 修改论文PDF
	ModifyPaperPDF(ctx context.Context, pdf *model.PaperPdf) error
	// ListPaperPDFs 列出论文PDF
	ListPaperPDFs(ctx context.Context, req *pb.ListPaperPDFsRequest) (*pb.ListPaperPDFsResponse, error)
	// CountPaperPDFs 获取论文PDF总数
	CountPaperPDFs(ctx context.Context, req *pb.CountPaperPDFsRequest) (*pb.CountPaperPDFsResponse, error)
	// GetById 根据ID获取论文PDF
	GetById(ctx context.Context, id string) (*model.PaperPdf, error)
	// GetByIds 根据IDS获取论文PDF列表
	GetByIds(ctx context.Context, ids []string) ([]model.PaperPdf, error)
	// GetByPaperId 根据论文ID获取PDF
	GetByPaperId(ctx context.Context, paperId string) (*model.PaperPdf, error)
	// GetByFileSHA256 根据SHA256获取PDF
	GetByFileSHA256(ctx context.Context, fileSHA256 string) (*model.PaperPdf, error)
	// GetByUserIdAndSHA256 根据用户ID和SHA256获取PDF
	GetByUserIdAndFileSHA256(ctx context.Context, userId string, fileSHA256 string) (*model.PaperPdf, error)
	// GetDefaultPdfByPaperId 获取默认PDF
	GetDefaultPdfByPaperId(ctx context.Context, paperId string) (*model.PaperPdf, error)
	// GetPdfUrlById 获取PDF文件地址
	GetPdfUrlById(ctx context.Context, pdfId string, duration int) (*string, error)
	// AuthPermissionDenied Pdf权限认证是否许可
	AuthPermissionDenied(ctx context.Context, pdfId string, userId string) (bool, error)
	// GetPdfStatusInfo 获取PDF状态信息
	GetPdfStatusInfo(ctx context.Context, paperId string, userId string) (*pb.GetPdfStatusInfoResponse, error)
	// GetDocById 根据PdfId获取文献信息
	GetDocById(ctx context.Context, id string) (*docModel.UserDoc, error)

	// GetMarkTagInfosByFolderId 根据folderId获取用户markTagInfos
	GetMarkTagInfosByFolderId(ctx context.Context, userId string, folderId string) ([]*common.PdfMarkTagInfo, error)

	//查询用户的所有Pdf文献总量 单位：bit
	GetPdfFileTotalSize(ctx context.Context, userId string) (int64, error)

	// GetCountPdfMarksByNoteId 根据笔记ID获取PDF标记总数
	GetCountPdfMarksByNoteId(ctx context.Context, noteId string) (int64, error)

	// GetPdfMarksByNoteId 根据笔记ID获取PDF标记
	GetPdfMarksByNoteId(ctx context.Context, noteId string) ([]model.PdfMark, error)

	// GetPdfMarkTagsByMarkId 根据markId获取PDF标记标签信息
	GetPdfMarkTagsByMarkId(ctx context.Context, markId string) ([]*common.AnnotateTag, error)
}
