package interfaces

import (
	"context"

	pb "github.com/yb2020/odoc/proto/gen/go/note"
	"github.com/yb2020/odoc/services/note/bean"
	"github.com/yb2020/odoc/services/note/model"
)

// INoteWordService 笔记单词服务接口
type INoteWordService interface {
	// NoteWord2WordInfoPb 将NoteWord转换为WordInfo
	NoteWord2WordInfoPb(noteWord *model.NoteWord) (*pb.WordInfo, error)
	// CreateNoteWord 创建笔记单词
	CreateNoteWord(ctx context.Context, noteWordPb *pb.SaveNoteWordRequest) (string, error)
	// UpdateNoteWord 更新笔记单词
	UpdateNoteWord(ctx context.Context, noteWord *model.NoteWord) (bool, error)
	// DeleteNoteWordById 删除笔记单词
	DeleteNoteWordById(ctx context.Context, id string) (bool, error)
	// GetNoteWordById 根据ID获取笔记单词
	GetNoteWordById(ctx context.Context, id string) (*model.NoteWord, error)
	// GetNoteWordsByNoteId 根据笔记ID获取笔记单词列表
	GetNoteWordsByNoteId(ctx context.Context, noteId string) ([]model.NoteWord, error)
	// DeleteNoteWordsByNoteId 根据笔记ID删除所有笔记单词
	DeleteNoteWordsByNoteId(ctx context.Context, noteId string) (bool, error)
	// GetNoteWordsByNoteIdWithMinLoadedId 根据笔记ID和最小加载ID获取笔记单词列表
	GetNoteWordsByNoteIdWithMinLoadedId(ctx context.Context, noteId string, minLoadedId string, limit int) ([]model.NoteWord, error)
	// GetListByNoteId 根据笔记ID获取笔记单词列表
	GetListByNoteId(ctx context.Context, noteId string, limit, offset int) ([]model.NoteWord, error)
	// GetListByNoteIds 根据笔记Ids列表获取笔记单词列表
	GetListByNoteIds(ctx context.Context, noteIds []string, limit, offset int) ([]model.NoteWord, error)
	// GetCountByNoteId 根据笔记ID获取笔记单词总数
	GetCountByNoteId(ctx context.Context, noteId string) (int64, error)
	// GetCountByNoteIds 根据笔记Ids列表获取笔记单词总数
	GetCountByNoteIds(ctx context.Context, noteIds []string) (int64, error)
	// GetNoteWordInfoPageByNoteIds 根据NoteIds获取笔记单词分页
	GetNoteWordInfoPageByNoteIds(ctx context.Context, noteIds []string, currentPage int32, pageSize int32) ([]*pb.WordInfo, int64, error)
	// GetNoteWordsByNoteQuery 根据笔记ID和查询条件获取笔记单词列表
	GetNoteWordsByNoteQuery(ctx context.Context, query *bean.NoteWordQuery) (*pb.GetNoteWordsByNoteIdResponse, error)
	// SaveOrUpdateNoteWordConfig 保存或更新生词配置
	SaveOrUpdateNoteWordConfig(ctx context.Context, noteId string, color pb.WordColor, displayMode pb.DISPLAY_MODE) (bool, error)
	// UpdateNoteWordTargetContent 更新笔记单词目标内容
	UpdateNoteWordTargetContent(ctx context.Context, pbTargetContent *pb.UpdateNoteWordRequest) (bool, error)
	// GetListByUserId 获取用户单词列表
	GetListByUserId(ctx context.Context, userId string) (*pb.GetWordListResponse, error)
	// GetUserWordListByFolderId 根据文件夹ID获取用户单词列表
	GetUserWordListByFolderId(ctx context.Context, userId string, req *pb.GetWordListByFolderIdReq) (*pb.GetWordListByFolderIdResponse, error)
}
