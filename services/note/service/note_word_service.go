package service

import (
	"context"
	"encoding/json"
	"strings"

	userContext "github.com/yb2020/odoc/pkg/context"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/note"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	docInterfaces "github.com/yb2020/odoc/services/doc/interfaces"
	"github.com/yb2020/odoc/services/note/bean"
	"github.com/yb2020/odoc/services/note/dao"
	noteInterfaces "github.com/yb2020/odoc/services/note/interfaces"
	"github.com/yb2020/odoc/services/note/model"
)

// NoteWordService 笔记单词服务实现
type NoteWordService struct {
	noteWordDAO           *dao.NoteWordDAO
	logger                logging.Logger
	tracer                opentracing.Tracer
	noteWordConfigService *NoteWordConfigService
	userDocService        docInterfaces.IUserDocService
	paperNoteService      noteInterfaces.IPaperNoteService
}

// NewNoteWordService 创建新的笔记单词服务
func NewNoteWordService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	noteWordDAO *dao.NoteWordDAO,
	noteWordConfigService *NoteWordConfigService,
	userDocService docInterfaces.IUserDocService,
	paperNoteService noteInterfaces.IPaperNoteService,
) *NoteWordService {
	return &NoteWordService{
		noteWordDAO:           noteWordDAO,
		logger:                logger,
		tracer:                tracer,
		noteWordConfigService: noteWordConfigService,
		userDocService:        userDocService,
		paperNoteService:      paperNoteService,
	}
}

// NoteWord2WordInfoPb 将NoteWord转换为WordInfo
func (s *NoteWordService) NoteWord2WordInfoPb(noteWord *model.NoteWord) (*pb.WordInfo, error) {
	// 转换为 WordInfo
	wordInfo := &pb.WordInfo{
		Id:   noteWord.Id,
		Word: noteWord.Word,
		Rectangle: func() []*pb.RectOption {
			if noteWord.Rectangle != "" {
				var rects []*pb.RectOption
				json.Unmarshal([]byte(noteWord.Rectangle), &rects)
				return rects
			}
			return []*pb.RectOption{}
		}(),
		TranslateInfo: func() *pb.TranslateInfo {
			translateInfo := &pb.TranslateInfo{
				BritishSymbol:        noteWord.BritishSymbol,
				AmericaSymbol:        noteWord.AmericaSymbol,
				BritishFormat:        noteWord.BritishFormat,
				BritishPronunciation: noteWord.BritishPronunciation,
				AmericaFormat:        noteWord.AmericaFormat,
				AmericaPronunciation: noteWord.AmericaPronunciation,
			}

			targetContent := []string{}
			if noteWord.TargetContent != "" {
				json.Unmarshal([]byte(noteWord.TargetContent), &targetContent)
			}
			targetResp := []*pb.TargetResp{}
			if noteWord.TargetResp != "" {
				json.Unmarshal([]byte(noteWord.TargetResp), &targetResp)
			}
			translateInfo.TargetContent = targetContent
			translateInfo.TargetResp = targetResp
			return translateInfo
		}(),
	}
	return wordInfo, nil
}

// CreateNoteWord 创建笔记单词
func (s *NoteWordService) CreateNoteWord(ctx context.Context, noteWordPb *pb.SaveNoteWordRequest) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.CreateNoteWord")
	defer span.Finish()

	// 获取用户ID
	userId, _ := userContext.GetUserID(ctx)

	// 获取参数
	noteId := noteWordPb.NoteId
	wordInfo := noteWordPb.WordInfo
	s.logger.Debug("创建笔记单词", "noteId", noteId, "wordInfo", wordInfo)

	if wordInfo == nil {
		return "0", errors.Biz("参数有误")
	}

	translateInfo := wordInfo.TranslateInfo
	if translateInfo == nil {
		return "0", errors.Biz("参数有误")
	}

	rectangles := wordInfo.Rectangle
	s.logger.Debug("rectangles", "rectangles", rectangles)
	word := strings.ToLower(wordInfo.Word)
	s.logger.Debug("word", "word", word)
	wordLength := len(word)
	s.logger.Debug("wordLength", "wordLength", wordLength)
	// 验证单词长度, nacos配置 :TODO
	if wordLength < 2 {
		return "0", errors.Biz("至少输入2个字母才能添加")
	}
	if wordLength > 100 {
		return "0", errors.Biz("最多输入100个字母才能添加")
	}

	dbNoteWord, err := s.noteWordDAO.GetByNoteIdAndWord(ctx, noteId, word)
	if err != nil {
		s.logger.Error("创建笔记单词失败", "error", err)
		return "0", errors.Biz("note.note_word.errors.create_failed")
	}
	if dbNoteWord != nil {
		return "0", errors.Biz("请勿重复添加")
	}

	id := idgen.GenerateUUID()
	// 设置参数
	noteWord := &model.NoteWord{
		NoteId:               noteId,
		Word:                 word,
		UserId:               userId,
		AmericaFormat:        translateInfo.AmericaFormat,
		AmericaPronunciation: translateInfo.AmericaPronunciation,
		AmericaSymbol:        translateInfo.AmericaSymbol,
		BritishFormat:        translateInfo.BritishFormat,
		BritishPronunciation: translateInfo.BritishPronunciation,
		BritishSymbol:        translateInfo.BritishSymbol,
	}
	noteWord.Id = id
	if rectangles != nil {
		rectanglesJson, _ := json.Marshal(rectangles)
		noteWord.Rectangle = string(rectanglesJson)
	}
	if translateInfo.TargetContent != nil {
		targetContentJson, _ := json.Marshal(translateInfo.TargetContent)
		noteWord.TargetContent = string(targetContentJson)
	}
	if translateInfo.TargetResp != nil {
		targetRespJson, _ := json.Marshal(translateInfo.TargetResp)
		noteWord.TargetResp = string(targetRespJson)
	}

	if err := s.noteWordDAO.Create(ctx, noteWord); err != nil {
		s.logger.Error("创建笔记单词失败", "error", err)
		return "0", errors.Biz("note.note_word.errors.create_failed")
	}

	return noteWord.Id, nil
}

// UpdateNoteWord 更新笔记单词
func (s *NoteWordService) UpdateNoteWord(ctx context.Context, noteWord *model.NoteWord) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.UpdateNoteWord")
	defer span.Finish()

	// 获取笔记单词
	word, err := s.noteWordDAO.FindById(ctx, noteWord.Id)
	if err != nil {
		s.logger.Error("获取笔记单词失败", "error", err)
		return false, errors.Biz("note.note_word.errors.get_failed")
	}

	if word == nil {
		return false, errors.Biz("note.note_word.errors.not_found")
	}

	// 更新笔记单词
	if err := s.noteWordDAO.UpdateById(ctx, noteWord); err != nil {
		s.logger.Error("更新笔记单词失败", "error", err)
		return false, errors.Biz("note.note_word.errors.update_failed")
	}

	return true, nil
}

// DeleteNoteWordById 删除笔记单词
func (s *NoteWordService) DeleteNoteWordById(ctx context.Context, id string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.DeleteNoteWordById")
	defer span.Finish()

	// 删除笔记单词
	if err := s.noteWordDAO.DeleteById(ctx, id); err != nil {
		s.logger.Error("删除笔记单词失败", "error", err)
		return false, errors.Biz("note.note_word.errors.delete_failed")
	}

	return true, nil
}

// GetNoteWordById 根据ID获取笔记单词
func (s *NoteWordService) GetNoteWordById(ctx context.Context, id string) (*model.NoteWord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetNoteWordById")
	defer span.Finish()

	// 获取笔记单词
	word, err := s.noteWordDAO.FindById(ctx, id)
	if err != nil {
		s.logger.Error("获取笔记单词失败", "error", err)
		return nil, errors.Biz("note.note_word.errors.get_failed")
	}

	return word, nil
}

// GetNoteWordsByNoteId 根据笔记ID获取笔记单词列表
func (s *NoteWordService) GetNoteWordsByNoteId(ctx context.Context, noteId string) ([]model.NoteWord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetNoteWordsByNoteId")
	defer span.Finish()

	// 获取笔记单词列表
	words, err := s.noteWordDAO.GetByNoteID(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记单词列表失败", "error", err)
		return nil, errors.Biz("note.note_word.errors.list_failed")
	}

	return words, nil
}

// DeleteNoteWordsByNoteId 根据笔记ID删除所有笔记单词
func (s *NoteWordService) DeleteNoteWordsByNoteId(ctx context.Context, noteId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.DeleteNoteWordsByNoteId")
	defer span.Finish()

	// 删除笔记单词
	if err := s.noteWordDAO.DeleteByNoteID(ctx, noteId); err != nil {
		s.logger.Error("删除笔记单词失败", "error", err)
		return false, errors.Biz("note.note_word.errors.delete_failed")
	}

	return true, nil
}

// GetNoteWordsByNoteId 根据笔记ID获取笔记单词列表
func (s *NoteWordService) GetNoteWordsByNoteIdWithMinLoadedId(ctx context.Context, noteId string, minLoadedId string, limit int) ([]model.NoteWord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetNoteWordsByNoteIdWithMinLoadedId")
	defer span.Finish()

	// 获取笔记单词列表
	words, err := s.noteWordDAO.GetByNoteIDWithMinLoadedId(ctx, noteId, minLoadedId, limit)
	if err != nil {
		s.logger.Error("获取笔记单词列表失败", "error", err)
		return nil, errors.Biz("note.note_word.errors.list_failed")
	}

	return words, nil
}

// GetListByNoteId 根据笔记ID获取笔记单词列表
func (s *NoteWordService) GetListByNoteId(ctx context.Context, noteId string, limit, offset int) ([]model.NoteWord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetListByNoteId")
	defer span.Finish()

	// 获取笔记单词列表
	words, err := s.noteWordDAO.GetListByNoteID(ctx, noteId, limit, offset)
	if err != nil {
		s.logger.Error("获取笔记单词列表失败", "error", err)
		return nil, errors.Biz("note.note_word.errors.list_failed")
	}

	return words, nil
}

// GetListByNoteIds 根据笔记Ids列表获取笔记单词列表
func (s *NoteWordService) GetListByNoteIds(ctx context.Context, noteIds []string, limit, offset int) ([]model.NoteWord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetListByNoteIds")
	defer span.Finish()

	// 获取笔记单词列表
	words, err := s.noteWordDAO.GetListByNoteIds(ctx, noteIds, limit, offset)
	if err != nil {
		s.logger.Error("获取笔记单词列表失败", "error", err)
		return nil, errors.Biz("note.note_word.errors.list_failed")
	}

	return words, nil
}

func (s *NoteWordService) GetCountByNoteId(ctx context.Context, noteId string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetCountByNoteId")
	defer span.Finish()

	// 获取笔记单词总数
	count, err := s.noteWordDAO.GetCountByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("获取笔记单词总数失败", "error", err)
		return 0, errors.Biz("note.note_word.errors.count_failed")
	}

	return count, nil
}

func (s *NoteWordService) GetCountByNoteIds(ctx context.Context, noteIds []string) (int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetCountByNoteIds")
	defer span.Finish()

	// 获取笔记单词总数
	count, err := s.noteWordDAO.GetCountByNoteIds(ctx, noteIds)
	if err != nil {
		s.logger.Error("获取笔记单词总数失败", "error", err)
		return 0, errors.Biz("note.note_word.errors.count_failed")
	}

	return count, nil
}

// GetNoteWordInfoPageByNoteIds 根据NoteIds获取笔记单词分页
func (s *NoteWordService) GetNoteWordInfoPageByNoteIds(ctx context.Context, noteIds []string, currentPage int32, pageSize int32) ([]*pb.WordInfo, int64, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetNoteWordInfoPageByNoteIds")
	defer span.Finish()

	// 获取笔记单词总数
	total, err := s.GetCountByNoteIds(ctx, noteIds)
	if err != nil {
		s.logger.Error("获取笔记单词总数失败", "error", err)
		return nil, 0, errors.Biz("note.note_word.errors.count_failed")
	}

	if total == 0 {
		return []*pb.WordInfo{}, 0, nil
	}

	// 获取笔记单词列表
	noteWords, err := s.GetListByNoteIds(ctx, noteIds, int(pageSize), int((currentPage-1)*pageSize))
	if err != nil {
		s.logger.Error("获取笔记单词列表失败", "error", err)
		return nil, 0, errors.Biz("note.note_word.errors.list_failed")
	}

	// 转换为 WordInfo
	var wordInfos []*pb.WordInfo
	for _, noteWord := range noteWords {
		wordInfo, err := s.NoteWord2WordInfoPb(&noteWord)
		if err != nil {
			s.logger.Error("转换笔记单词失败", "error", err)
			return nil, 0, errors.Biz("note.note_word.errors.convert_failed")
		}
		wordInfos = append(wordInfos, wordInfo)
	}
	return wordInfos, total, nil
}

// GetNoteWordsByNoteQuery 根据笔记ID和加载ID获取笔记单词列表
func (s *NoteWordService) GetNoteWordsByNoteQuery(ctx context.Context, query *bean.NoteWordQuery) (*pb.GetNoteWordsByNoteIdResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetNoteWordsByNoteQuery")
	defer span.Finish()

	// TODO: nacos配置读取默认currentPage参数设置

	// 获取笔记单词配置
	noteWordConfig, err := s.noteWordConfigService.GetNoteWordConfigByNoteId(ctx, query.NoteId)
	if err != nil {
		s.logger.Error("msg", "Failed to get note word config", "error", err.Error())
		return nil, errors.Biz("note.note_word.errors.get_failed")
	}

	var wordColor pb.WordColor
	var displayMode pb.DISPLAY_MODE
	if noteWordConfig != nil {
		// 转换为 WordColor
		if colorValue, ok := pb.WordColor_value[noteWordConfig.Color]; ok {
			wordColor = pb.WordColor(colorValue)
		}
		// 转换为 DisplayMode
		if displayModeValue, ok := pb.DISPLAY_MODE_value[noteWordConfig.DisplayMode]; ok {
			displayMode = pb.DISPLAY_MODE(displayModeValue)
		}
	}

	wordTotal, err := s.GetCountByNoteId(ctx, query.NoteId)
	if err != nil {
		s.logger.Error("msg", "获取笔记单词总数失败", "error", err.Error())
		return nil, errors.Biz("note.note_word.errors.count_failed")
	}

	// 获取笔记单词列表
	// 直接使用请求的上下文，它已经在中间件中被更新
	var noteWords []model.NoteWord
	if query.MinLoadedId > "0" {
		// 如果指定了最小加载ID，则使用分页加载方式获取数据
		noteWords, err = s.GetNoteWordsByNoteIdWithMinLoadedId(ctx, query.NoteId, query.MinLoadedId, query.PageSize)
		if err != nil {
			s.logger.Error("msg", "获取笔记单词列表失败", "error", err.Error())
			return nil, errors.Biz("note.note_word.errors.list_failed")
		}
	} else {
		noteWords, err = s.GetListByNoteId(ctx, query.NoteId, query.PageSize, (query.CurrentPage-1)*query.PageSize)
		if err != nil {
			s.logger.Error("msg", "获取笔记单词列表失败", "error", err.Error())
			return nil, errors.Biz("note.note_word.errors.list_failed")
		}
	}

	// 转换为 WordInfo
	var wordInfos []*pb.WordInfo
	for _, noteWord := range noteWords {
		wordInfo, err := s.NoteWord2WordInfoPb(&noteWord)
		if err != nil {
			s.logger.Error("转换笔记单词失败", "error", err)
			return nil, errors.Biz("note.note_word.errors.convert_failed")
		}
		wordInfos = append(wordInfos, wordInfo)
	}

	noteWordsInfo := &pb.GetNoteWordsByNoteIdResponse{
		Words:       wordInfos,
		Color:       wordColor,
		DisplayMode: displayMode,
		Total:       uint32(wordTotal),
	}

	return noteWordsInfo, nil
}

// SaveOrUpdateNoteWordConfig 保存或更新生词配置
func (s *NoteWordService) SaveOrUpdateNoteWordConfig(ctx context.Context, noteId string, color pb.WordColor, displayMode pb.DISPLAY_MODE) (bool, error) {

	// TODO: 保存或更新生词配置
	config, err := s.noteWordConfigService.GetNoteWordConfigByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("msg", "Failed to get note word config", "error", err.Error())
		return false, errors.Biz("note.note_word.errors.get_failed")
	}

	if config == nil {
		// 创建新配置
		config = &model.NoteWordConfig{
			NoteId:      noteId,
			Color:       color.String(),
			DisplayMode: displayMode.String(),
		}
		if _, err := s.noteWordConfigService.CreateNoteWordConfig(ctx, config); err != nil {
			s.logger.Error("msg", "Failed to save note word config", "error", err.Error())
			return false, errors.Biz("note.note_word.errors.save_failed")
		}
	} else {
		// 更新配置
		config.Color = color.String()
		config.DisplayMode = displayMode.String()
		if _, err := s.noteWordConfigService.UpdateNoteWordConfig(ctx, config); err != nil {
			s.logger.Error("msg", "Failed to update note word config", "error", err.Error())
			return false, errors.Biz("note.note_word.errors.update_failed")
		}
	}

	return true, nil
}

// UpdateNoteWordTargetContent 更新笔记单词目标内容
func (s *NoteWordService) UpdateNoteWordTargetContent(ctx context.Context, pbTargetContent *pb.UpdateNoteWordRequest) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.UpdateNoteWordTargetContent")
	defer span.Finish()

	// 先获取原始数据
	noteWord, err := s.GetNoteWordById(ctx, pbTargetContent.WordId)
	if err != nil {
		s.logger.Error("msg", "获取笔记单词失败", "error", err.Error())
		return false, errors.Biz("note.note_word.errors.get_failed")
	}
	if noteWord == nil {
		s.logger.Error("msg", "笔记单词不存在")
		return false, errors.Biz("note.note_word.errors.not_found")
	}

	targetContent := pbTargetContent.WordInfo.TargetContent
	// 更新字段
	if pbTargetContent.WordInfo.TargetContent != nil {
		targetContentJson, _ := json.Marshal(targetContent)
		noteWord.TargetContent = string(targetContentJson)
		noteWord.TargetResp = ""
	}

	// 更新笔记单词
	if err := s.noteWordDAO.UpdateById(ctx, noteWord); err != nil {
		s.logger.Error("更新笔记单词失败", "error", err)
		return false, errors.Biz("note.note_word.errors.update_failed")
	}

	return true, nil
}

func (s *NoteWordService) GetListByUserId(ctx context.Context, userId string) (*pb.GetWordListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetListByUserId")
	defer span.Finish()
	folderInfo, err := s.paperNoteService.GetNoteDocTreeNodeByFolderId(ctx, userId, "0")
	if err != nil {
		s.logger.Error("获取笔记列表失败", "error", err, "userId", userId)
		return nil, errors.Biz("note.paper_note.errors.list_failed")
	}

	if folderInfo == nil {
		return &pb.GetWordListResponse{
			Total:                0,
			FolderInfos:          nil,
			UnclassifiedDocInfos: nil,
		}, nil
	}

	//对folderInfo过滤掉NoteWordCount为0的文件和文件夹或子文件夹
	s.filterEmptyItems(folderInfo)

	return &pb.GetWordListResponse{
		Total:                uint32(folderInfo.NoteWordCount),
		FolderInfos:          folderInfo.ChildrenFolders,
		UnclassifiedDocInfos: folderInfo.DocInfos,
	}, nil
}

// filterEmptyItems 递归过滤掉 word count 为 0 的文档和文件夹
func (s *NoteWordService) filterEmptyItems(folder *pb.NoteManageFolderInfo) {
	if folder == nil {
		return
	}

	// 过滤当前文件夹下的文档
	filteredDocs := make([]*pb.NoteManageDocInfo, 0, len(folder.DocInfos))
	for _, doc := range folder.DocInfos {
		if doc.NoteWordCount > 0 {
			filteredDocs = append(filteredDocs, doc)
		}
	}
	folder.DocInfos = filteredDocs

	// 递归过滤子文件夹
	filteredChildren := make([]*pb.NoteManageFolderInfo, 0, len(folder.ChildrenFolders))
	for _, child := range folder.ChildrenFolders {
		s.filterEmptyItems(child)
		// 如果子文件夹在过滤后仍有内容（其NoteWordCount > 0），则保留
		if child.NoteWordCount > 0 {
			filteredChildren = append(filteredChildren, child)
		}
	}
	folder.ChildrenFolders = filteredChildren
}

func (s *NoteWordService) GetUserWordListByFolderId(ctx context.Context, userId string, req *pb.GetWordListByFolderIdReq) (*pb.GetWordListByFolderIdResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "NoteWordService.GetUserWordListByFolderId")
	defer span.Finish()

	folderInfo, err := s.paperNoteService.GetAllNoteDocInfosByFolderId(ctx, userId, *req.FolderId)
	if err != nil {
		s.logger.Error("获取笔记列表失败", "error", err, "userId", userId)
		return nil, errors.Biz("note.paper_note.errors.list_failed")
	}
	if folderInfo == nil {
		return &pb.GetWordListByFolderIdResponse{}, nil
	}

	var noteIds []string
	for _, docInfo := range folderInfo {
		noteIds = append(noteIds, docInfo.NoteId)
	}

	if len(noteIds) == 0 {
		return &pb.GetWordListByFolderIdResponse{}, nil
	}

	wordList, total, err := s.GetNoteWordInfoPageByNoteIds(ctx, noteIds, int32(req.CurrentPage), int32(req.PageSize))
	if err != nil {
		s.logger.Error("获取笔记单词列表失败", "error", err)
		return nil, errors.Biz("note.note_word.errors.list_failed")
	}

	var resp = &pb.GetWordListByFolderIdResponse{
		Total: uint32(total),
		Words: wordList,
	}
	return resp, nil
}
