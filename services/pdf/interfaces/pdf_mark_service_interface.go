package interfaces

import (
	"context"

	"github.com/yb2020/odoc-proto/gen/go/common"
	notePb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/services/pdf/dto"
	"github.com/yb2020/odoc/services/pdf/model"
)

// IPdfMarkService PdfMark服务接口
type IPdfMarkService interface {
	// SavePdfMark 保存PDF标记
	SavePdfMark(ctx context.Context, mark *model.PdfMark) (string, error)

	// UpdatePdfMark 更新PDF标记
	UpdatePdfMark(ctx context.Context, mark *model.PdfMark) (bool, error)

	// DeletePdfMarkById 删除PDF标记
	DeletePdfMarkById(ctx context.Context, id string) (bool, error)

	// GetPdfMarkById 根据ID获取PDF标记
	GetPdfMarkById(ctx context.Context, id string) (*model.PdfMark, error)

	// GetPdfMarksByNoteId 根据笔记ID获取PDF标记列表
	GetPdfMarksByNoteId(ctx context.Context, noteId string) ([]model.PdfMark, error)

	// GetCountPdfMarksByNoteId 根据笔记ID获取PDF标记总数
	GetCountPdfMarksByNoteId(ctx context.Context, noteId string) (int64, error)

	// GetAnnotationRawModelsByNoteId 根据笔记ID获取注释原始模型列表
	GetAnnotationRawModelsByNoteId(ctx context.Context, noteId string) ([]*notePb.AnnotationRawModel, error)

	// GetWebNoteAnnotationModelsByNoteId 根据笔记ID获取WebNoteAnnotationModel列表
	GetWebNoteAnnotationModelsByNoteId(ctx context.Context, noteId string) ([]*notePb.WebNoteAnnotationModel, error)

	// SavePdfMarkByAnnotation 通过注释保存PDF标记
	SavePdfMarkByAnnotation(ctx context.Context, annotation *notePb.WebNoteAnnotationModel) (string, error)

	// SavePdfMarkByBean 通过Bean保存PDF标记
	SavePdfMarkByBean(ctx context.Context, pdfMark *model.PdfMark) (string, error)

	// DeleteByAnnotationPointer 通过注释指针删除PDF标记
	DeleteByAnnotationPointer(ctx context.Context, annotationPointer *common.AnnotationPointer) (bool, error)

	// DeleteMarkTagRelationsById 通过标记ID删除标记标签关系
	DeleteMarkTagRelationsById(ctx context.Context, id string) (bool, error)

	// UpdatePdfMarkByAnnotation 通过注释更新PDF标记
	UpdatePdfMarkByAnnotation(ctx context.Context, annotation *notePb.WebNoteAnnotationModel) (string, error)

	// UpdatePdfMarkByBean 通过Bean更新PDF标记
	UpdatePdfMarkByBean(ctx context.Context, updateMark *model.PdfMark) (string, error)

	// GetAnnotateTagsByMarkId 通过markId获取标签列表
	GetAnnotateTagsByMarkId(ctx context.Context, markId string) ([]*common.AnnotateTag, error)

	// BatchSavePdfMarks 批量保存PDF标记
	BatchSavePdfMarks(ctx context.Context, marks []*model.PdfMark) (string, error)

	// GetMarkTagInfosByNoteIds 根据noteIds获取用户markTagInfos
	GetMarkTagInfosByNoteIds(ctx context.Context, userId string, noteIds []string) ([]*common.PdfMarkTagInfo, error)

	// GetMarkTagInfosByFolderId 根据folderId获取用户markTagInfos
	GetMarkTagInfosByFolderId(ctx context.Context, userId string, folderId string) ([]*common.PdfMarkTagInfo, error)

	// GetUserPdfMarkPage获取用户笔记标记分页
	GetUserPdfMarkPage(ctx context.Context, userId string, searchDto *dto.PdfMarkSearchPageDto) ([]*notePb.WebNoteAnnotationModel, int32, error)
}
