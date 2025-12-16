package interfaces

import (
	"context"

	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/model"
)

// IUserDocService 用户文档服务接口
type IUserDocService interface {
	// GetById 根据ID获取文献
	GetById(ctx context.Context, id string) (*model.UserDoc, error)

	// GetByPdId 根据PdfId获取文献
	GetByPdId(ctx context.Context, pdfId string) (*model.UserDoc, error)

	// GetByUserIdAndPdfId 根据用户ID和PDF ID获取用户文档
	GetByUserIdAndPdfId(ctx context.Context, userId, pdfId string) (*model.UserDoc, error)

	// GetByUserIDAndPaperID 根据用户ID和论文ID获取用户文档
	GetByUserIDAndPaperID(ctx context.Context, userId, paperId string) (*model.UserDoc, error)

	// CheckDocExistsByFileName 检查指定用户是否已有同名文档
	CheckDocExistsByFileName(ctx context.Context, fileName string, userId string) (bool, error)

	// SaveUserDoc 直接保存用户文档对象
	SaveUserDoc(ctx context.Context, userDoc *model.UserDoc) error

	// DeleteUserDocs 删除用户文献
	DeleteUserDocs(ctx context.Context, docIds []string, userId string) error

	// RenameUserDoc 重命名用户文献
	RenameUserDoc(ctx context.Context, docId string, docName string, userId string) error

	// 查询左侧的文献列表
	GetDocIndex(ctx context.Context, userId string) (*docpb.GetDocIndexResponse, error)

	// GetDocTreeRootNodeOfDocsByUserId 根据userId查询用户未分类的所有文献; isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
	GetDocTreeRootNodeOfDocsByUserId(ctx context.Context, userId string, isOnlyShowNote bool) ([]*docpb.SimpleUserDocInfo, error)

	// GetFolderInfosByUserId 根据userId查询用户已分类的所有文献夹-文献列表
	GetFolderInfosByUserId(ctx context.Context, userId string) ([]*docpb.UserDocFolderInfo, int32, error)
	// GetFolderInfosByUserIdAndIds 根据userId查询用户已分类的所有文献夹-文献列表
	GetFolderInfosByUserIdAndIds(ctx context.Context, userId string, docIds []string) ([]*docpb.UserDocFolderInfo, int32, error)

	// GetDocTreeNodeByFolderId 根据userId和文件Id查询用户的文献夹信息;参数isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
	GetDocTreeNodeByFolderId(ctx context.Context, userId string, folderId string, isOnlyShowNote bool) (*docpb.UserDocFolderInfo, error)

	// GetDocsByFolderId 根据userId查询folderId文件夹下的文献, folderId=0时返回未分类文献; isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
	GetDocsByFolderId(ctx context.Context, userId string, folderId string, isOnlyShowNote bool) ([]*docpb.SimpleUserDocInfo, error)

	// GetAllDocsByFolderId 根据userId查询folderId文件夹及其子文件夹下的文献, folderId=0时返回所有文献; isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
	GetAllDocsByFolderId(ctx context.Context, userId string, folderId string, isOnlyShowNote bool) ([]*docpb.SimpleUserDocInfo, error)

	// UpdateUserDocNoteIdByPdfId 根据pdfId更新UserDoc表中的noteId
	UpdateUserDocNoteIdByPdfId(ctx context.Context, pdfId string, noteId string) error
}
