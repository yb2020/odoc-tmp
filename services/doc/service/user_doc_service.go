package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/dao"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	osspb "github.com/yb2020/odoc/proto/gen/go/oss"
	docBean "github.com/yb2020/odoc/services/doc/bean"
	"github.com/yb2020/odoc/services/doc/constant"
	docDao "github.com/yb2020/odoc/services/doc/dao"
	helper "github.com/yb2020/odoc/services/doc/helper"
	"github.com/yb2020/odoc/services/doc/interfaces"
	"github.com/yb2020/odoc/services/doc/model"
	membershipInterface "github.com/yb2020/odoc/services/membership/interfaces"
	noteService "github.com/yb2020/odoc/services/note/interfaces"
	noteModel "github.com/yb2020/odoc/services/note/model"
	ossModel "github.com/yb2020/odoc/services/oss/model"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperModel "github.com/yb2020/odoc/services/paper/model"
	paperService "github.com/yb2020/odoc/services/paper/service"
	parseConstant "github.com/yb2020/odoc/services/parse/constant"
	pdfService "github.com/yb2020/odoc/services/pdf/interfaces"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
	"gorm.io/gorm"
)

// 确保 UserDocService 实现了 IUserDocService 接口
var _ interfaces.IUserDocService = (*UserDocService)(nil)

// UserDocService 用户文档服务实现
type UserDocService struct {
	userDocDAO                   *docDao.UserDocDAO
	logger                       logging.Logger
	tracer                       opentracing.Tracer
	cache                        cache.Cache
	userDocAttachmentService     *UserDocAttachmentService
	userDocClassifyService       *UserDocClassifyService
	userDocFolderService         *UserDocFolderService
	docClassifyRelationService   *DocClassifyRelationService
	userDocFolderRelationService *UserDocFolderRelationService
	paperJcrService              *paperService.PaperJcrService
	paperService                 *paperService.PaperService
	noteService                  noteService.IPaperNoteService       // 笔记服务，使用interface{}避免循环依赖
	paperPdfService              pdfService.IPaperPdfService         // PDF服务，使用interface{}避免循环依赖
	paperPdfParsedService        *paperService.PaperPdfParsedService // PDF解析服务
	ossService                   ossService.OssServiceInterface
	config                       *config.Config // 添加配置字段
	membershipService            membershipInterface.IMembershipService
	squidProxyService            *SquidProxyService
}

// NewUserDocService 创建新的用户文档服务
func NewUserDocService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	cache cache.Cache,
	userDocDAO *docDao.UserDocDAO,
	userDocAttachmentService *UserDocAttachmentService,
	userDocClassifyService *UserDocClassifyService,
	userDocFolderService *UserDocFolderService,
	docClassifyRelationService *DocClassifyRelationService,
	userDocFolderRelationService *UserDocFolderRelationService,
	paperJcrService *paperService.PaperJcrService,
	paperService *paperService.PaperService,
	ossService ossService.OssServiceInterface,
	config *config.Config,
	membershipService membershipInterface.IMembershipService,
	paperPdfParsedService *paperService.PaperPdfParsedService,
) *UserDocService {
	return &UserDocService{
		logger:                       logger,
		tracer:                       tracer,
		cache:                        cache,
		userDocDAO:                   userDocDAO,
		userDocAttachmentService:     userDocAttachmentService,
		userDocClassifyService:       userDocClassifyService,
		userDocFolderService:         userDocFolderService,
		docClassifyRelationService:   docClassifyRelationService,
		userDocFolderRelationService: userDocFolderRelationService,
		paperJcrService:              paperJcrService,
		paperService:                 paperService,
		ossService:                   ossService,
		config:                       config,
		membershipService:            membershipService,
		paperPdfParsedService:        paperPdfParsedService,
	}
}

// SetNoteService 设置笔记服务，用于解决循环依赖问题
func (s *UserDocService) SetNoteService(noteService noteService.IPaperNoteService) error {
	if noteService == nil {
		return errors.Biz("noteService cannot be nil")
	}
	s.noteService = noteService
	return nil
}

// SetPaperPdfService 设置pdf服务，用于解决循环依赖问题
func (s *UserDocService) SetPaperPdfService(paperPdfService pdfService.IPaperPdfService) error {
	if paperPdfService == nil {
		return errors.Biz("noteService cannot be nil")
	}
	s.paperPdfService = paperPdfService
	return nil
}

// SetSquidProxyService 设置squid代理服务，用于解决循环依赖问题
func (s *UserDocService) SetSquidProxyService(squidProxyService *SquidProxyService) error {
	if squidProxyService == nil {
		return errors.Biz("squidProxyService cannot be nil")
	}
	s.squidProxyService = squidProxyService
	return nil
}

// GetPaperJcrEntity 根据venue查找PaperJcrEntity
func (s *UserDocService) GetPaperJcrEntity(ctx context.Context, venue string) (*paperModel.PaperJcrEntity, error) {
	if venue == "" {
		return nil, nil
	}
	return s.paperJcrService.GetPaperJcrEntityByVenue(ctx, venue)
}

// GetUserDocsByFolderId 获取指定文件夹及其所有子文件夹下的文档列表
func (s *UserDocService) GetUserDocsByFolderId(ctx context.Context, userId string, folderId string) ([]model.UserDoc, error) {
	// 1. 获取该用户的所有文档
	allUserDocs, err := s.userDocDAO.GetUserDocsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户文档失败", "userId", userId, "error", err.Error())
		return nil, err
	}

	// 2. 如果folderId为空或为0，直接返回所有文档
	if folderId == "0" {
		return allUserDocs, nil
	}

	// 3. 获取该用户的所有文件夹
	userFolders, err := s.userDocFolderService.GetUserDocFolders(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户文件夹失败", "userId", userId, "error", err.Error())
		return nil, err
	}

	// 4. 递归获取所有子文件夹ID
	descendantFolderIds := make([]string, 0)
	helper.GetDescendantFolderIds(folderId, userFolders, &descendantFolderIds)
	descendantFolderIds = append(descendantFolderIds, folderId)

	// 5. 查询这些文件夹下的所有文档关系
	docFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByFolderIDs(ctx, descendantFolderIds)
	if err != nil {
		s.logger.Error("msg", "获取文档文件夹关系失败", "folderIds", descendantFolderIds, "error", err.Error())
		return nil, err
	}

	// 6. 提取所有文档ID，去重
	docIdMap := make(map[string]struct{})
	for _, rel := range docFolderRelations {
		docIdMap[rel.DocId] = struct{}{}
	}
	// 7. 返回属于这些文档ID的文档列表
	result := make([]model.UserDoc, 0)
	for _, doc := range allUserDocs {
		if _, ok := docIdMap[doc.Id]; ok {
			result = append(result, doc)
		}
	}
	return result, nil
}

// GetAllJcrPartionsByFolderId 获取指定文件夹及其所有子文件夹下所有文档的JCR分区（去重+优先用户编辑）
func (s *UserDocService) GetAllJcrPartionsByFolderId(ctx context.Context, folderId string, userID string) (*docpb.JcrPartionsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetAllJcrPartionsByFolderId")
	defer span.Finish()
	// 获取指定文件夹下的所有文档
	userDocs, err := s.GetUserDocsByFolderId(ctx, userID, folderId)
	if err != nil {
		s.logger.Error("msg", "获取文件夹下文档失败", "folderId", folderId, "userId", userID, "error", err.Error())
		return &docpb.JcrPartionsResponse{JcrPartions: []string{}}, nil
	}
	if len(userDocs) == 0 {
		return &docpb.JcrPartionsResponse{JcrPartions: []string{}}, nil
	}
	// 提取JCR分区并去重
	partionSet := make(map[string]struct{})
	for _, doc := range userDocs {
		var partion string
		// 优先使用用户编辑的分区
		if doc.UserEditedJcrPartion != "" {
			partion = doc.UserEditedJcrPartion
		} else {
			// 否则查询JCR实体获取分区
			entity, err := s.GetPaperJcrEntity(ctx, doc.Venue)
			if err == nil && entity != nil && entity.JcrPartion != "" {
				partion = entity.JcrPartion
			}
		}

		// 只添加非空分区
		if partion != "" {
			partionSet[partion] = struct{}{}
		}
	}
	// 将去重后的分区转换为切片
	result := make([]string, 0, len(partionSet))
	for k := range partionSet {
		result = append(result, k)
	}
	return &docpb.JcrPartionsResponse{JcrPartions: result}, nil
}

// GetById 根据ID获取文献
func (s *UserDocService) GetById(ctx context.Context, id string) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetById")
	defer span.Finish()

	// 获取文档
	doc, err := s.userDocDAO.FindExistById(ctx, id)
	if err != nil {
		s.logger.Error("msg", "get user doc by user id and paper id failed", "id", id, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}

	return doc, nil
}

// GetByPdId 根据PdfId获取文献
func (s *UserDocService) GetByPdId(ctx context.Context, pdfId string) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetById")
	defer span.Finish()

	// 获取文档
	doc, err := s.userDocDAO.GetUserDocByPdfId(ctx, pdfId)
	if err != nil {
		s.logger.Error("msg", "get user doc by user id and paper id failed", "pdfId", pdfId, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}

	return doc, nil
}

// GetByUserIDAndPaperID 根据用户ID和论文ID获取用户文档
func (s *UserDocService) GetByUserIDAndPaperID(ctx context.Context, userId, paperId string) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetByUserIDAndPaperID")
	defer span.Finish()

	// 获取文档
	doc, err := s.userDocDAO.GetUserDocByUserIDAndPaperID(ctx, userId, paperId)
	if err != nil {
		s.logger.Error("msg", "get user doc by user id and paper id failed", "userId", userId, "paperId", paperId, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}

	return doc, nil
}

// GetByUserIdAndPdfId 根据用户ID和PDF ID获取用户文档
func (s *UserDocService) GetByUserIdAndPdfId(ctx context.Context, userId, pdfId string) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetByUserIdAndPdfId")
	defer span.Finish()

	// 获取文档
	doc, err := s.userDocDAO.GetUserDocByUserIdAndPdfId(ctx, userId, pdfId)
	if err != nil {
		s.logger.Error("msg", "get user doc by user id and pdf id failed", "userId", userId, "pdfId", pdfId, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}

	return doc, nil
}

// SaveUserDoc 直接保存用户文档对象
func (s *UserDocService) SaveUserDoc(ctx context.Context, userDoc *model.UserDoc) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.SaveUserDoc")
	defer span.Finish()
	// 检查必填字段
	if userDoc.UserId == "0" {
		return errors.Biz("doc.user_doc.errors.user_id_required")
	}
	// 保存文档到数据库
	if err := s.userDocDAO.Save(ctx, userDoc); err != nil {
		s.logger.Error("msg", "save user doc failed", "error", err, "userId", userDoc.UserId)
		return errors.Biz("doc.user_doc.errors.save_failed")
	}
	s.logger.Info("msg", "save user doc success", "docId", userDoc.Id, "userId", userDoc.UserId)
	return nil
}

// 根据笔记id获取用户文档
func (s *UserDocService) GetUserDocByNoteId(ctx context.Context, noteId string) (*model.UserDoc, error) {
	// 获取文档
	doc, err := s.userDocDAO.GetUserDocByNoteId(ctx, noteId)
	if err != nil {
		s.logger.Error("msg", "get user doc by note id failed", "noteId", noteId, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}

	return doc, nil
}

// GetDocRelatedVenueList 获取文档相关的发表场所列表
func (s *UserDocService) GetDocRelatedVenueList(ctx context.Context, req *docpb.GetDocRelatedVenueListReq) (*docpb.GetDocRelatedVenueListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocRelatedVenueList")
	defer span.Finish()

	// 创建响应对象
	response := &docpb.GetDocRelatedVenueListResponse{
		VenueInfos: make([]string, 0),
	}

	// 查询用户所有的文献数据
	userDocs, err := s.userDocDAO.GetAllUserDocsByUserID(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get user docs by user id failed")
	}

	// 如果指定了文件夹ID，按文件夹过滤文档
	if req.FolderId != nil && req.GetFolderId() != "" {
		// 获取用户文件夹关系
		userDocFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, req.GetUserId())
		if err != nil {
			s.logger.Error("msg", "get user doc folder relations by user id failed", "userId", req.GetUserId(), "error", err)
			return nil, errors.Wrap(err, "get user doc folder relations by user id failed")
		}

		// 创建一个临时的GetDocListReq对象用于过滤
		tempReq := &docpb.GetDocListReq{}
		tempReq.FolderId = req.FolderId

		// 使用现有的过滤函数过滤文档
		userDocs = helper.FilterUserDocsByFolder(tempReq, userDocs, userDocFolderRelations)
	}

	// 从所有文档中提取发表场所信息
	venuesMap := make(map[string]bool)
	for i := range userDocs {
		// 获取文档的发表场所列表
		venues := helper.GetUserDocDisplayVenueInfos(&userDocs[i])
		for _, venue := range venues {
			venuesMap[venue] = true
		}
	}

	// 将发表场所信息转换为列表
	venueInfos := make([]string, 0, len(venuesMap))
	for venue := range venuesMap {
		venueInfos = append(venueInfos, venue)
	}

	// 设置响应
	response.Total = uint32(len(venueInfos))
	response.VenueInfos = venueInfos
	return response, nil
}

// GetDocRelatedAuthorList 获取文档相关的作者列表
func (s *UserDocService) GetDocRelatedAuthorList(ctx context.Context, req *docpb.GetDocRelatedAuthorListReq) (*docpb.GetDocRelatedAuthorListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocRelatedAuthorList")
	defer span.Finish()

	// 创建响应对象
	response := &docpb.GetDocRelatedAuthorListResponse{
		AuthorInfos: make([]string, 0),
	}

	// 查询用户所有的文献数据
	userDocs, err := s.userDocDAO.GetAllUserDocsByUserID(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get user docs by user id failed")
	}

	// 如果指定了文件夹ID，按文件夹过滤文档
	if req.FolderId != nil && req.GetFolderId() != "" {
		// 获取用户文件夹关系
		userDocFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, req.GetUserId())
		if err != nil {
			s.logger.Error("msg", "get user doc folder relations by user id failed", "userId", req.GetUserId(), "error", err)
			return nil, errors.Wrap(err, "get user doc folder relations by user id failed")
		}

		// 创建一个临时的GetDocListReq对象用于过滤
		tempReq := &docpb.GetDocListReq{}
		tempReq.FolderId = req.FolderId

		// 使用现有的过滤函数过滤文档
		userDocs = helper.FilterUserDocsByFolder(tempReq, userDocs, userDocFolderRelations)
	}
	var authorInfos []string
	for _, userDoc := range userDocs {
		authors := helper.GetUserDocDisplayAuthors(&userDoc)
		for _, author := range authors {
			authorInfos = append(authorInfos, author)
		}
	}
	// 设置响应
	response.Total = uint32(len(authorInfos))
	response.AuthorInfos = authorInfos
	return response, nil
}

// GetDocRelatedClassifyList 获取文档相关的分类列表
func (s *UserDocService) GetDocRelatedClassifyList(ctx context.Context, req *docpb.GetDocRelatedClassifyListReq) (*docpb.GetDocRelatedClassifyListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocRelatedClassifyList")
	defer span.Finish()

	// 创建响应对象
	response := &docpb.GetDocRelatedClassifyListResponse{
		ClassifyInfos: make([]*docpb.UserDocClassifyInfo, 0),
	}

	// 获取用户所有分类
	allClassifies, err := s.userDocClassifyService.GetUserDocClassifies(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get user doc classify list failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get user doc classify list failed")
	}

	// 获取该用户下所有文献对应的标签关系
	allDocClassifyRelations, err := s.docClassifyRelationService.GetDocClassifyRelationsByUserID(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get doc classify relations by user id failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get doc classify relations by user id failed")
	}

	// 获取用户所有文档ID
	userDocIds, err := s.userDocDAO.GetDocIdsByUserID(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get doc ids by user id failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get doc ids by user id failed")
	}

	// 过滤出用户文档相关的分类关系
	validDocClassifyRelations := make([]model.DocClassifyRelation, 0)
	for _, relation := range allDocClassifyRelations {
		for _, docId := range userDocIds {
			if relation.DocId == docId {
				validDocClassifyRelations = append(validDocClassifyRelations, relation)
				break
			}
		}
	}

	// 如果没有指定文件夹或文件夹ID为0，返回所有有效分类
	if req.FolderId == nil || req.GetFolderId() == "0" {
		// 获取所有有效的分类ID
		validClassifyIds := make(map[string]bool)
		for _, relation := range validDocClassifyRelations {
			validClassifyIds[relation.ClassifyId] = true
		}

		// 按排序过滤并构建响应
		validClassifies := make([]model.UserDocClassify, 0)
		for _, classify := range allClassifies {
			if validClassifyIds[classify.Id] {
				validClassifies = append(validClassifies, classify)
			}
		}

		// 按Sort字段排序
		sort.Slice(validClassifies, func(i, j int) bool {
			return validClassifies[i].Sort < validClassifies[j].Sort
		})

		// 构建响应
		for _, classify := range validClassifies {
			classifyInfo := &docpb.UserDocClassifyInfo{
				ClassifyId:   classify.Id,
				ClassifyName: classify.Name,
			}
			response.ClassifyInfos = append(response.ClassifyInfos, classifyInfo)
		}

		response.Total = uint32(len(validClassifies))
		return response, nil
	}

	// 获取用户所有文件夹
	userFolders, err := s.userDocFolderService.GetUserDocFolders(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get user doc folders failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get user doc folders failed")
	}

	// 获取指定文件夹及其所有子文件夹
	descendantFolderIds := make([]string, 0)
	helper.GetDescendantFolderIds(req.GetFolderId(), userFolders, &descendantFolderIds)
	descendantFolderIds = append(descendantFolderIds, req.GetFolderId())

	// 获取这些文件夹下的文档关系
	docFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByFolderIDs(ctx, descendantFolderIds)
	if err != nil {
		s.logger.Error("msg", "get doc folder relations by folder ids failed", "folderIds", descendantFolderIds, "error", err)
		return nil, errors.Wrap(err, "get doc folder relations by folder ids failed")
	}

	// 获取这些文件夹下的文档ID
	docIds := make(map[string]bool)
	for _, relation := range docFolderRelations {
		// 检查是否是用户的文档
		for _, userDocId := range userDocIds {
			if relation.DocId == userDocId {
				docIds[relation.DocId] = true
				break
			}
		}
	}

	// 获取这些文档相关的分类ID
	classifyIds := make(map[string]bool)
	for _, relation := range validDocClassifyRelations {
		if docIds[relation.DocId] {
			classifyIds[relation.ClassifyId] = true
		}
	}

	// 构建响应
	for _, classify := range allClassifies {
		if classifyIds[classify.Id] {
			classifyInfo := &docpb.UserDocClassifyInfo{
				ClassifyId:   classify.Id,
				ClassifyName: classify.Name,
			}
			response.ClassifyInfos = append(response.ClassifyInfos, classifyInfo)
		}
	}

	response.Total = uint32(len(response.ClassifyInfos))
	return response, nil
}

// 查询用户上传文件的解析状态
func (s *UserDocService) GetUserUploadParseStatus(ctx context.Context, userId string, statusToken string) (*docpb.GetUserDocCreateStatusResponse, error) {
	resp := &docpb.GetUserDocCreateStatusResponse{
		Status: docpb.UserDocParsedStatusEnum_READY, // 默认设置为 READY 准备中
		Token:  statusToken,
	}
	// 创建缓存键
	tokenCacheKey := fmt.Sprintf("%s%s", constant.ParsePDFStatusTokenKeyPrefix, statusToken)

	// 从缓存中获取状态和是否存在sha256的key
	tokenCacheObj := docpb.UserDocParsedStatusObj{}
	exists, err := s.cache.GetNotBizPrefix(ctx, tokenCacheKey, &tokenCacheObj)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.Biz("file parse status not found")
	}
	var fileSHA256 string
	if tokenCacheObj.HasSha256File {
		fileSHA256 = tokenCacheObj.GetSha256Value()
		var cacheValue docpb.UserDocParsedStatusEnum
		sha256CacheKey := fmt.Sprintf("%s%s", constant.ParsePDFStatusKeyPrefix, fileSHA256)
		exists, err := s.cache.GetNotBizPrefix(ctx, sha256CacheKey, &cacheValue)
		if err != nil {
			return nil, err
		}
		if !exists {
			//如果不存这个缓存key 直接使用tokenCacheObj中的状态
			resp.Status = tokenCacheObj.Status
		} else {
			resp.Status = cacheValue
		}
	} else {
		resp.Status = tokenCacheObj.Status
	}
	// 根据状态判断是否需要设置额外的响应内容
	hasDocInfo := false
	if resp.Status >= docpb.UserDocParsedStatusEnum_HEADER_DATA_PARSED && resp.Status <= docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSE_FAILED {
		hasDocInfo = true
	}

	if hasDocInfo {
		// 根据用户id和file_sha256查询paperPdf数据
		paperPdf, err := s.paperPdfService.GetByUserIdAndFileSHA256(ctx, userId, fileSHA256)
		if err != nil {
			s.logger.Error("msg", "查询paperPdf失败", "error", err, "userId", userId, "fileSHA256", fileSHA256)
			return resp, nil
		}
		if paperPdf == nil {
			return resp, nil
		}
		resp.PdfId = paperPdf.Id
		//根据paperPdfId查询userDoc数据
		userDoc, err := s.GetByUserIDAndPaperID(ctx, userId, paperPdf.PaperId)
		if err != nil {
			s.logger.Error("msg", "查询userDoc失败", "error", err, "userId", userId, "paperId", paperPdf.PaperId)
			return resp, nil
		}
		if userDoc == nil {
			return resp, nil
		}
		resp.DocInfo = &docpb.MyCollectedDocInfo{
			Id:      userDoc.Id,
			PaperId: userDoc.PaperId,
			PdfId:   paperPdf.Id,
			UserId:  userId,
		}
	}
	return resp, nil
}

// 查询左侧的文献列表
func (s *UserDocService) GetDocIndex(ctx context.Context, userId string) (*docpb.GetDocIndexResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocIndex")
	defer span.Finish()
	// 2. 通过用户ID调用DAO的根据用户ID查询未删除文献列表
	allUserDocs, err := s.userDocDAO.GetAllUserDocsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", userId, "error", err)
		return nil, err
	}

	// 3. 创建文献ID到文献对象的映射
	userDocMap := make(map[string]*model.UserDoc)
	for i := range allUserDocs {
		userDocMap[allUserDocs[i].Id] = &allUserDocs[i]
	}

	// 4. 创建响应对象
	resp := &docpb.GetDocIndexResponse{}

	// 5. 设置总文献数量
	resp.TotalDocCount = uint32(len(allUserDocs))

	// 6. 构建文件夹目录树和未分类文献列表

	// 构建 docId-userDoc Map
	userDocMap2 := make(map[string]*model.UserDoc)
	for i := range allUserDocs {
		userDocMap2[allUserDocs[i].Id] = &allUserDocs[i]
	}

	// 获取用户所有文件夹
	allUserFolders, err := s.userDocFolderService.GetUserDocFolders(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get user doc folders failed", "userId", userId, "error", err)
		return nil, err
	}

	// 获取用户文献和文件夹关联关系
	allDocFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get doc folder relations by user id failed", "userId", userId, "error", err)
		return nil, err
	}

	// 构建文件夹-关联关系Map
	folderAndDocRelationsMap := make(map[string][]*model.UserDocFolderRelation)
	for i := range allDocFolderRelations {
		folderId := allDocFolderRelations[i].FolderId
		if _, ok := folderAndDocRelationsMap[folderId]; !ok {
			folderAndDocRelationsMap[folderId] = make([]*model.UserDocFolderRelation, 0)
		}
		folderAndDocRelationsMap[folderId] = append(folderAndDocRelationsMap[folderId], &allDocFolderRelations[i])
	}

	// 构造文件夹目录树
	folderInfos := s.buildSimpleFolderTree(ctx, userDocMap2, allUserFolders, folderAndDocRelationsMap)

	// 排序（按照 Sort 字段降序排列）
	sort.Slice(folderInfos, func(i, j int) bool {
		return folderInfos[i].Sort > folderInfos[j].Sort
	})

	// 设置文件夹信息
	resp.FolderInfos = folderInfos

	// 设置未分类文献（不在任何文件夹中的文献）
	// 首先收集所有已分类的文献 ID（仅统计 folderId != 0 的关联）
	classifiedDocIds := make(map[string]bool)
	for folderId, relations := range folderAndDocRelationsMap {
		if folderId != "0" { // 非0文件夹的文献才视为已分类
			for _, relation := range relations {
				classifiedDocIds[relation.DocId] = true
			}
		}
	}

	// 构建未分类文献列表：包含 folderId = 0 的文献 + 完全没有关联关系的文献
	unclassifiedDocs := make([]*docpb.SimpleUserDocInfo, 0)
	addedUnclassified := make(map[string]struct{})

	// 先添加 folderId = 0 的文献
	if unclassifiedRelations, exists := folderAndDocRelationsMap["0"]; exists {
		for _, relation := range unclassifiedRelations {
			if doc, ok := userDocMap2[relation.DocId]; ok {
				docInfo := &docpb.SimpleUserDocInfo{
					Sort:    uint32(doc.Sort),
					DocName: doc.DocName,
					DocId:   doc.Id,
					PdfId:   doc.PdfId,
					PaperId: doc.PaperId,
					NoteId:  doc.NoteId,
				}
				unclassifiedDocs = append(unclassifiedDocs, docInfo)
				addedUnclassified[doc.Id] = struct{}{}
			}
		}
	}

	// 再添加完全没有关联关系（未被计入已分类）的文献
	for _, doc := range allUserDocs {
		if !classifiedDocIds[doc.Id] { // 未被任何非0文件夹归类
			if _, exists := addedUnclassified[doc.Id]; exists { // 已经作为 folderId=0 添加过
				continue
			}
			docInfo := &docpb.SimpleUserDocInfo{
				Sort:    uint32(doc.Sort),
				DocName: doc.DocName,
				DocId:   doc.Id,
				PdfId:   doc.PdfId,
				PaperId: doc.PaperId,
				NoteId:  doc.NoteId,
			}
			unclassifiedDocs = append(unclassifiedDocs, docInfo)
		}
	}

	// 按照 Sort 字段降序排序
	sort.Slice(unclassifiedDocs, func(i, j int) bool {
		return unclassifiedDocs[i].Sort > unclassifiedDocs[j].Sort
	})
	// 设置未分类文献列表
	resp.UnclassifiedDocInfos = unclassifiedDocs
	return resp, nil
}

// 根据文件的file_sha256值查询整个pdf上传的数据
func (s *UserDocService) GetUserUploadBaseDataBySHA256AndUserId(ctx context.Context, userId string, fileSHA256 string) (*model.UserDoc, *pdfModel.PaperPdf, *paperModel.Paper, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUserUploadBaseDataByFileSHA256AndUserId")
	defer span.Finish()

	paperPdf, err := s.paperPdfService.GetByUserIdAndFileSHA256(ctx, userId, fileSHA256)
	if err != nil {
		return nil, nil, nil, err
	}
	if paperPdf == nil {
		return nil, nil, nil, errors.Biz("paper pdf not found")
	}

	userDoc, err := s.userDocDAO.GetUserDocByPdfId(ctx, paperPdf.Id)
	if err != nil {
		return nil, nil, nil, err
	}
	if userDoc == nil {
		return nil, nil, nil, errors.Biz("user doc not found")
	}

	paper, err := s.paperService.GetPaperById(ctx, paperPdf.PaperId)
	if err != nil {
		return nil, nil, nil, err
	}
	if paper == nil {
		return nil, nil, nil, errors.Biz("paper not found")
	}

	return userDoc, paperPdf, paper, nil
}

// 根据文件的file_sha256值查询未解析完成的用户列表
func (s *UserDocService) GetUsersNotParsedByFileSHA256(ctx context.Context, fileSHA256 string) ([]string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUsersNotParsedByFileSHA256")
	defer span.Finish()

	// 调用DAO层方法查询未解析完成的用户ID列表
	userIds, err := s.userDocDAO.GetUsersNotParsedByFileSHA256(ctx, fileSHA256, int32(docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSED))
	if err != nil {
		s.logger.Error("查询未解析完成的用户列表失败", "error", err, "fileSHA256", fileSHA256)
		return nil, errors.BizWrap("GetUsersNotParsedByFileSHA256 failed", err)
	}
	return userIds, nil
}

// GetDocList 获取用户文档列表
func (s *UserDocService) GetDocList(ctx context.Context, req *docpb.GetDocListReq) (*docpb.GetDocListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocList")
	defer span.Finish()

	// 创建响应对象
	response := &docpb.GetDocListResponse{
		DocList: make([]*docpb.UserDocInfo, 0),
		Total:   0,
	}

	// 查询用户所有的文献数据
	userDocs, err := s.userDocDAO.GetAllUserDocsByUserID(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get user docs by user id failed")
	}

	if len(userDocs) == 0 {
		return response, nil
	}

	// 1. 使用 sort.Slice 对 userDocs 进行降序排序
	//    将 LastReadTime 最晚的记录排在最前面
	sort.Slice(userDocs, func(i, j int) bool {
		// 使用 Before 来实现降序，因为我们希望大的（晚的）在前
		return userDocs[j].LastReadTime.Before(userDocs[i].LastReadTime)
	})

	// 2. 查找第一个有效的记录
	var lastReadDoc *model.UserDoc
	for i := range userDocs {
		// 找到第一个 LastReadTime 不是零值的记录，它就是最晚的
		if !userDocs[i].LastReadTime.IsZero() {
			lastReadDoc = &userDocs[i]
			break // 找到就立即跳出循环
		}
	}

	// 获取用户所有的文件夹关系信息
	userDocFolderRelations := []model.UserDocFolderRelation{}
	if req.FolderId != nil && *req.FolderId != "0" {
		// 只有当需要按文件夹过滤时，才查询文件夹关系
		var err error
		userDocFolderRelations, err = s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, req.GetUserId())
		if err != nil {
			s.logger.Error("msg", "get user doc folder relations by user id failed", "userId", req.GetUserId(), "error", err)
			return nil, errors.Wrap(err, "get user doc folder relations by user id failed")
		}
	}

	// 获取用户文献和标签关联关系
	docClassifyRelations, err := s.docClassifyRelationService.GetDocClassifyRelationsByUserID(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get doc classify relations by user id failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get doc classify relations by user id failed")
	}

	// 获取用户创建的所有标签
	userDocClassifies, err := s.userDocClassifyService.GetUserDocClassifies(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("msg", "get user doc classify list failed", "userId", req.GetUserId(), "error", err)
		return nil, errors.Wrap(err, "get user doc classify list failed")
	}

	// 构建文献-标签信息Map
	docClassifiesMap := helper.GetDocClassifyInfoMap(docClassifyRelations, userDocClassifies)

	// 条件过滤
	// 如果是仅看有pdf的
	if req.GetOnlyPdf() {
		filteredDocs := make([]model.UserDoc, 0)
		for _, doc := range userDocs {
			if doc.PdfId != "0" {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		userDocs = filteredDocs
	}
	// 过滤标签
	userDocs = helper.FilterUserDocsByClassify(req, userDocs, docClassifyRelations)
	// 过滤作者
	userDocs = helper.FilterUserDocsByAuthor(req, userDocs)
	// 过滤收录情况
	userDocs = helper.FilterUserDocsByVenue(req, userDocs)
	if len(userDocs) == 0 {
		return response, nil
	}
	// 过滤文件夹
	userDocs = helper.FilterUserDocsByFolder(req, userDocs, userDocFolderRelations)
	// 忽略影响因子和JCR分区相关的过滤逻辑
	if len(userDocs) == 0 {
		return response, nil
	}

	// 文献id-搜索内容Map
	docSearchResultMap := make(map[string]*docpb.DocSearchResult)

	// 搜索过滤相关
	userDocs = helper.FilterUserDocsBySearchContent(req, userDocs, docSearchResultMap)

	// 文献列表排序
	helper.SortDocList(userDocs, req)

	// 设置总数
	response.Total = uint32(len(userDocs))

	// 分页
	if req.GetCurrentPage() > 0 && req.GetPageSize() > 0 {
		start := (req.GetCurrentPage() - 1) * req.GetPageSize()
		end := start + req.GetPageSize()
		if start >= uint32(len(userDocs)) {
			return response, nil
		}
		if end > uint32(len(userDocs)) {
			end = uint32(len(userDocs))
		}
		userDocs = userDocs[start:end]
	}

	// 组装返回给前端的数据
	userDocInfoList := make([]*docpb.UserDocInfo, 0, len(userDocs))

	userDocIds := make([]string, 0, len(userDocs))
	for _, doc := range userDocs {
		userDocIds = append(userDocIds, doc.Id)
	}

	for _, doc := range userDocs {
		// 使用辅助函数获取完整的UserDocInfo对象
		userDocInfo := helper.GetAllInfoPbUserDocInfo(ctx, &doc, docClassifiesMap, docSearchResultMap, lastReadDoc, s.paperJcrService)
		// 设置Embedding状态
		embeddingStatus := docpb.UserDocParsedStatusEnum_EMBEDDING_FAILED
		userDocInfo.EmbeddingStatus = embeddingStatus
		parsedProgress := helper.GetUserDocParsedStatusEnum(userDocInfo)
		userDocInfo.ParsedProgress = &parsedProgress
		userDocInfoList = append(userDocInfoList, userDocInfo)

	}

	response.DocList = userDocInfoList
	return response, nil
}

// 根据笔记id获取文献详情
func (s *UserDocService) GetUserDocDetailInfoById(ctx context.Context, noteId string) (*docpb.DocDetailInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUserDocById")
	defer span.Finish()
	// 获取笔记信息
	paperNote, err := s.noteService.GetPaperNoteById(ctx, noteId)
	if err != nil {
		return nil, errors.Biz("note.not.exist")
	}
	if paperNote == nil {
		return nil, errors.Biz("note.not.exist")
	}
	// 根据笔记获取用户文档
	userDoc, err := s.getUserDocByPaperNote(ctx, paperNote)
	if err != nil {
		return nil, err
	}
	if userDoc == nil {
		return nil, errors.Biz("user.doc.not.exist")
	}
	// 获取用户文档信息
	userDocInfo, err := s.getSingleDocInfo(ctx, userDoc)
	if err != nil {
		return nil, err
	}
	userDocDetailInfo, err := helper.GetUserDocDetailInfo(ctx, userDocInfo, userDoc)
	if err != nil {
		return nil, err
	}
	// 查询Embedding状态
	userDocDetailInfo.EmbeddingStatus = docpb.UserDocParsedStatusEnum_EMBEDDING_FAILED
	// 返回结果
	return userDocDetailInfo, nil
}

func (s *UserDocService) GetUserDocStatusByIds(ctx context.Context, docIds []string) (*docpb.UserDocStatusByIdsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUserDocStatusByIds")
	defer span.Finish()

	if len(docIds) == 0 {
		return &docpb.UserDocStatusByIdsResponse{
			Items: []*docpb.UserDocStatusByIdsItem{},
		}, nil
	}

	// 批量查询用户文档
	userDocs, err := s.userDocDAO.GetUserDocsByIds(ctx, docIds)
	if err != nil {
		return nil, err
	}
	if len(userDocs) == 0 {
		return &docpb.UserDocStatusByIdsResponse{
			Items: []*docpb.UserDocStatusByIdsItem{},
		}, nil
	}
	// 获取paperIds和docIds用于批量查询
	paperIds := make([]string, 0, len(userDocs))
	userDocIds := make([]string, 0, len(userDocs))
	for _, userDoc := range userDocs {
		paperIds = append(paperIds, userDoc.PdfId)
		userDocIds = append(userDocIds, userDoc.Id)
	}

	// 批量查询PDF信息
	paperPdfList, err := s.paperPdfService.GetByIds(ctx, paperIds)
	if err != nil {
		return nil, err
	}
	// 构建pdfId到PaperPdf的FileSHA256映射
	paperPdfMap := make(map[string]string, len(paperPdfList))
	for i := range paperPdfList {
		paperPdfMap[paperPdfList[i].Id] = paperPdfList[i].FileSHA256
	}

	// 构建结果列表
	docStatusList := make([]*docpb.UserDocStatusByIdsItem, 0, len(userDocs))
	for _, userDoc := range userDocs {
		// 查询Embedding状态
		embeddingStatus := docpb.UserDocParsedStatusEnum_EMBEDDING_FAILED
		// 构建状态信息
		docStatusInfo := &docpb.UserDocStatusByIdsItem{
			DocId:           userDoc.Id,
			Status:          docpb.UserDocParsedStatusEnum(userDoc.ParseStatus),
			EmbeddingStatus: embeddingStatus,
		}
		// 判断redis中是否存在状态 ，如果存在，则优先使用redis中的状态
		//获取redis中的状态
		cacheKey := fmt.Sprintf("%s%s", constant.ParsePDFStatusKeyPrefix, paperPdfMap[userDoc.PdfId])
		var cacheValue docpb.UserDocParsedStatusEnum
		exists, err := s.cache.GetNotBizPrefix(ctx, cacheKey, &cacheValue)
		if err != nil {
			continue
		}
		if exists {
			docStatusInfo.Status = cacheValue
		}

		docStatusList = append(docStatusList, docStatusInfo)
	}

	// 返回结果
	return &docpb.UserDocStatusByIdsResponse{
		Items: docStatusList,
	}, nil
}

// 根据笔记对象获取用户文档
func (s *UserDocService) getUserDocByPaperNote(ctx context.Context, paperNote *noteModel.PaperNote) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.getUserDocByPaperNote")
	defer span.Finish()
	userId := paperNote.CreatorId
	var userDoc *model.UserDoc
	var err error
	if paperNote.PaperId != "" {
		userDoc, err = s.GetByUserIDAndPaperID(ctx, userId, paperNote.PaperId)
		if err != nil || userDoc == nil {
			return nil, errors.Biz("personal.doc.absent.by.paper")
		}
	} else if paperNote.PdfId != "" {
		userDoc, err = s.GetByUserIdAndPdfId(ctx, userId, paperNote.PdfId)
		if err != nil || userDoc == nil {
			return nil, errors.Biz("personal.doc.absent.by.pdf")
		}
	} else {
		return nil, errors.Biz("personal.doc.not.found.by.paperAndPdf")
	}
	return userDoc, nil
}

// buildSimpleFolderTree 构建简化的文件夹目录树
func (s *UserDocService) buildSimpleFolderTree(ctx context.Context, userDocMap map[string]*model.UserDoc, allUserFolders []model.UserDocFolder, folderAndDocRelationsMap map[string][]*model.UserDocFolderRelation) []*docpb.UserDocFolderInfo {
	// 创建根文件夹列表
	rootFolders := make([]*docpb.UserDocFolderInfo, 0)

	// 构建文件夹ID到文件夹对象的映射
	folderMap := make(map[string]*model.UserDocFolder)
	for i := range allUserFolders {
		folderMap[allUserFolders[i].Id] = &allUserFolders[i]
	}

	// 首先找出所有根文件夹（父ID为0的文件夹）
	for i := range allUserFolders {
		if allUserFolders[i].ParentId == "0" {
			// 创建文件夹信息对象
			folderInfo := &docpb.UserDocFolderInfo{
				FolderId:        allUserFolders[i].Id,
				Name:            allUserFolders[i].Name,
				ParentId:        allUserFolders[i].ParentId,
				Level:           0, // 根据文件夹层级设置，这里是根文件夹，所以是0
				Sort:            uint32(allUserFolders[i].Sort),
				ChildrenFolders: make([]*docpb.UserDocFolderInfo, 0),
				DocInfos:        make([]*docpb.SimpleUserDocInfo, 0),
			}

			// 设置该文件夹下面的文献，并去重
			if relations, ok := folderAndDocRelationsMap[folderInfo.FolderId]; ok && len(relations) > 0 {
				// 去重处理，每个文献ID只保留一个关联关系
				uniqueDocs := make(map[string]*model.UserDocFolderRelation)
				for _, relation := range relations {
					docId := relation.DocId
					if _, exists := uniqueDocs[docId]; !exists {
						uniqueDocs[docId] = relation
					}
				}

				// 将去重后的关联关系转换为切片
				uniqueRelations := make([]*model.UserDocFolderRelation, 0, len(uniqueDocs))
				for _, relation := range uniqueDocs {
					uniqueRelations = append(uniqueRelations, relation)
				}

				// 包装文献信息
				folderInfo.DocInfos = s.wrapAndSortDocInfo(ctx, userDocMap, uniqueRelations)
			}

			// 递归构建子文件夹
			s.buildChildFolders(ctx, folderInfo, allUserFolders, folderMap, userDocMap, folderAndDocRelationsMap)

			// 设置文件夹及其所有子文件夹中的文档总数
			s.setDescendantFolderIds(ctx, &allUserFolders[i], folderInfo, allUserFolders, userDocMap, folderAndDocRelationsMap)

			// 添加到根文件夹列表
			rootFolders = append(rootFolders, folderInfo)
		}
	}

	// 按Sort降序排序
	sort.Slice(rootFolders, func(i, j int) bool {
		return rootFolders[i].Sort > rootFolders[j].Sort
	})

	return rootFolders
}

// buildChildFolders 递归构建子文件夹
func (s *UserDocService) buildChildFolders(ctx context.Context, parentFolder *docpb.UserDocFolderInfo, allFolders []model.UserDocFolder, folderMap map[string]*model.UserDocFolder, userDocMap map[string]*model.UserDoc, folderAndDocRelationsMap map[string][]*model.UserDocFolderRelation) {
	// 添加当前文件夹下的文档，并去重
	if relations, ok := folderAndDocRelationsMap[parentFolder.FolderId]; ok && len(relations) > 0 {
		// 去重处理，每个文献ID只保留一个关联关系
		uniqueDocs := make(map[string]*model.UserDocFolderRelation)
		for _, relation := range relations {
			docId := relation.DocId
			if _, exists := uniqueDocs[docId]; !exists {
				uniqueDocs[docId] = relation
			}
		}

		// 将去重后的关联关系转换为切片
		uniqueRelations := make([]*model.UserDocFolderRelation, 0, len(uniqueDocs))
		for _, relation := range uniqueDocs {
			uniqueRelations = append(uniqueRelations, relation)
		}

		// 包装文献信息
		parentFolder.DocInfos = s.wrapAndSortDocInfo(ctx, userDocMap, uniqueRelations)
	}

	// 查找当前文件夹的所有子文件夹
	for i := range allFolders {
		if allFolders[i].ParentId == parentFolder.FolderId {
			// 创建子文件夹信息对象
			childFolder := &docpb.UserDocFolderInfo{
				FolderId:        allFolders[i].Id,
				Name:            allFolders[i].Name,
				ParentId:        allFolders[i].ParentId,
				Level:           uint32(int64(parentFolder.Level) + 1), // 子文件夹层级是父文件夹层级+1
				Sort:            uint32(allFolders[i].Sort),
				ChildrenFolders: make([]*docpb.UserDocFolderInfo, 0),
				DocInfos:        make([]*docpb.SimpleUserDocInfo, 0),
			}

			// 递归构建子文件夹的子文件夹
			s.buildChildFolders(ctx, childFolder, allFolders, folderMap, userDocMap, folderAndDocRelationsMap)

			// 设置文件夹及其所有子文件夹中的文档总数
			s.setDescendantFolderIds(ctx, &allFolders[i], childFolder, allFolders, userDocMap, folderAndDocRelationsMap)

			// 添加到父文件夹的子文件夹列表
			parentFolder.ChildrenFolders = append(parentFolder.ChildrenFolders, childFolder)
		}
	}

	// 子文件夹按Sort降序排序
	if len(parentFolder.ChildrenFolders) > 0 {
		sort.Slice(parentFolder.ChildrenFolders, func(i, j int) bool {
			return parentFolder.ChildrenFolders[i].Sort > parentFolder.ChildrenFolders[j].Sort
		})
	}
}

// wrapAndSortDocInfo 包装并排序文档信息
func (s *UserDocService) wrapAndSortDocInfo(ctx context.Context, userDocMap map[string]*model.UserDoc, relations []*model.UserDocFolderRelation) []*docpb.SimpleUserDocInfo {
	// 创建简化的文献信息列表
	docInfos := make([]*docpb.SimpleUserDocInfo, 0, len(relations))

	// 处理每个关联关系
	for _, relation := range relations {
		// 获取文档对象
		docId := relation.DocId
		doc, ok := userDocMap[docId]
		if !ok {
			continue // 如果文档不存在，跳过
		}

		// 创建简化的文献信息
		docInfo := &docpb.SimpleUserDocInfo{
			Sort:    uint32(relation.Sort), // 使用关联关系中的排序值
			DocName: doc.DocName,
			DocId:   doc.Id,
			PdfId:   doc.PdfId,
			PaperId: doc.PaperId,
			NoteId:  doc.NoteId,
		}
		docInfos = append(docInfos, docInfo)
	}

	// 按Sort降序排序
	sort.Slice(docInfos, func(i, j int) bool {
		return docInfos[i].Sort > docInfos[j].Sort
	})

	return docInfos
}

// setDescendantFolderIds 设置文件夹及其所有子文件夹中的文档总数
func (s *UserDocService) setDescendantFolderIds(ctx context.Context, folder *model.UserDocFolder, folderInfo *docpb.UserDocFolderInfo, allFolders []model.UserDocFolder, userDocMap map[string]*model.UserDoc, folderAndDocRelationsMap map[string][]*model.UserDocFolderRelation) {
	// 获取所有后代文件夹ID
	descendantFolderIds := make([]string, 0)
	helper.GetDescendantFolderIds(folder.Id, allFolders, &descendantFolderIds)

	// 收集当前文件夹及其所有后代文件夹中的文档关联关系
	descendantRelations := make([]*model.UserDocFolderRelation, 0)

	// 添加当前文件夹的文档关联关系
	if relations, ok := folderAndDocRelationsMap[folder.Id]; ok && len(relations) > 0 {
		descendantRelations = append(descendantRelations, relations...)
	}

	// 添加所有后代文件夹的文档关联关系
	for _, folderId := range descendantFolderIds {
		if relations, ok := folderAndDocRelationsMap[folderId]; ok && len(relations) > 0 {
			descendantRelations = append(descendantRelations, relations...)
		}
	}

	// 统计不重复的文档数量
	docIdMap := make(map[string]bool)
	for _, relation := range descendantRelations {
		docId := relation.DocId
		if _, exists := userDocMap[docId]; exists {
			docIdMap[docId] = true
		}
	}

	// 设置文档总数
	folderInfo.DocCount = uint32(len(docIdMap))
}

func (s *UserDocService) getSingleDocInfo(ctx context.Context, userDoc *model.UserDoc) (*docpb.UserDocInfo, error) {
	// 获取用户文献和标签关联关系
	docClassifyRelations, err := s.docClassifyRelationService.GetDocClassifyRelationsByUserID(ctx, userDoc.UserId)
	if err != nil {
		return nil, errors.Biz("get doc classify relations by user id failed")
	}
	// 获取用户创建的所有标签
	userDocClassifies, err := s.userDocClassifyService.GetUserDocClassifies(ctx, userDoc.UserId)
	if err != nil {
		return nil, errors.Biz("get user doc classify list failed")
	}
	// 构建文献-标签信息Map
	docClassifiesMap := helper.GetDocClassifyInfoMap(docClassifyRelations, userDocClassifies)
	// 使用辅助函数获取UserDocInfo对象
	return helper.GetPbUserDocInfo(ctx, userDoc, docClassifiesMap, s.paperJcrService)
}

// DeleteUserDocs 删除用户文献
func (s *UserDocService) DeleteUserDocs(ctx context.Context, docIds []string, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.DeleteUserDocs")
	defer span.Finish()
	if len(docIds) == 0 {
		return nil
	}

	// 使用事务执行删除操作
	return s.userDocDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		for _, docId := range docIds {
			// 1. 验证文献所有权
			userDoc, err := s.userDocDAO.FindById(ctx, docId)
			if err != nil {
				return errors.Biz("get user doc by id failed")
			}
			if userDoc == nil {
				continue
			}
			if userDoc.UserId != userId {
				return errors.Biz("doc.user_doc.errors.doc_not_belong_to_current_user")
			}
			// 2. 逻辑删除文献记录
			userDoc.IsDeleted = true
			s.userDocDAO.ModifyExcludeNull(ctx, userDoc)
			// 3. 物理删除文献文件夹关系
			err = s.userDocFolderRelationService.DeleteRelationsByUserIdAndDocIds(ctx, userId, docId)
			if err != nil {
				return errors.Biz("delete user doc folder relation failed")
			}

			s.logger.Info("msg", "成功删除文献", "docId", docId, "userId", userId)
		}

		return nil
	})
}

// RenameUserDoc 重命名用户文献
func (s *UserDocService) RenameUserDoc(ctx context.Context, docId string, docName string, userId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.RenameUserDoc")
	defer span.Finish()
	if docId == "" || docName == "" {
		return errors.Biz("doc.user_doc.errors.invalid_params")
	}
	// 1. 验证文献所有权
	userDoc, err := s.userDocDAO.FindById(ctx, docId)
	if err != nil {
		return errors.Biz("get user doc by id failed")
	}
	if userDoc == nil {
		return errors.Biz("doc.user_doc.errors.doc_not_found")
	}
	if userDoc.UserId != userId {
		return errors.Biz("doc.user_doc.errors.doc_not_belong_to_current_user")
	}
	//2. 验证是否存在同名文献
	// exists, err := s.CheckDocExistsByFileName(ctx, docName, userId)
	// if err != nil {
	// 	return errors.Biz("check doc exists by file name failed")
	// }
	// if exists {
	// 	return errors.Biz("doc.user_doc.errors.doc_name_exists")
	// }
	// 3. 更新文献名称
	userDoc.DocName = docName
	userDoc.UserEditedDocName = docName
	userDoc.DocNameEdited = true
	// 4. 保存更新
	err = s.userDocDAO.Modify(ctx, userDoc)
	if err != nil {
		return errors.Biz("update user doc name failed")
	}
	return nil
}

// CheckDocExistsByFileName 检查指定用户是否已有同名文档
func (s *UserDocService) CheckDocExistsByFileName(ctx context.Context, fileName string, userId string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.CheckDocExistsByFileName")
	defer span.Finish()

	if fileName == "" || userId == "" {
		return false, errors.Biz("doc.user_doc.errors.invalid_params")
	}

	// 获取用户所有未删除的文档
	userDoc, err := s.userDocDAO.GetUserDocByUserIdAndFileName(ctx, userId, fileName)
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", userId, "error", err.Error())
		return false, errors.Biz("get user docs by user id failed")
	}
	if userDoc != nil {
		return true, nil
	}
	return false, nil
}

// 保存文件上传的相关记录
func (s *UserDocService) SavePdfUploadRecords(ctx context.Context, paperPdf *pdfModel.PaperPdf, paper *paperModel.Paper, userDoc *model.UserDoc, relation *model.UserDocFolderRelation) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.SaveUploadRecords")
	defer span.Finish()
	// 使用事务保存所有记录
	err := s.userDocDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，确保内部操作使用同一个事务连接
		// 这对于 SQLite 尤其重要，因为 SQLite 使用数据库级锁
		txCtx := context.WithValue(ctx, dao.TransactionContextKey, tx)

		// 1. 保存 PaperPdf 记录
		if err := s.paperPdfService.SavePaperPDF(txCtx, paperPdf); err != nil {
			s.logger.Error("msg", "save paper pdf failed", "error", err.Error(), "userId", userDoc.UserId)
			return err
		}
		// 2. 保存 Paper 记录
		if err := s.paperService.SavePaper(txCtx, paper); err != nil {
			s.logger.Error("msg", "save paper failed", "error", err.Error(), "userId", userDoc.UserId)
			return err
		}
		// 3. 保存用户文档记录
		if err := s.userDocDAO.Save(txCtx, userDoc); err != nil {
			s.logger.Error("msg", "save user doc failed", "error", err.Error(), "userId", userDoc.UserId)
			return err
		}
		// 4. 如果有文件夹关系对象，创建文档与文件夹的关系
		if relation != nil {
			if err := s.userDocFolderRelationService.CreateUserDocFolderRelation(txCtx, relation); err != nil {
				s.logger.Error("msg", "save doc folder relation failed", "error", err.Error(), "userId", userDoc.UserId, "docId", userDoc.Id, "folderId", relation.FolderId)
				return err
			}
		}
		return nil
	})
	if err != nil {
		return errors.Biz("save user doc failed")
	}
	return nil
}

// 根据用户id和文件名查询文件上传的相关记录
func (s *UserDocService) GetUserDocUploadRecords(ctx context.Context, userId string, fileName string) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUserDocUploadRecords")
	defer span.Finish()

	//根据用户和文件名称查找文件记录
	userDoc, err := s.userDocDAO.GetUserDocByUserIdAndFileName(ctx, userId, fileName)
	if err != nil {
		return nil, errors.Biz("get user doc by user id and file name failed")
	}
	return userDoc, nil
}

// SaveUploadRecords 保存上传记录（简化版，用于本地上传）
func (s *UserDocService) SaveUploadRecords(ctx context.Context, paperPdf *pdfModel.PaperPdf, paper *paperModel.Paper, userDoc *model.UserDoc, folderId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.SaveUploadRecords")
	defer span.Finish()

	// 创建文件夹关系
	relation := &model.UserDocFolderRelation{
		UserId:   userDoc.UserId,
		DocId:    userDoc.Id,
		FolderId: folderId,
	}
	if folderId == "" {
		relation.FolderId = "0"
	}
	relation.CreatorId = userDoc.UserId
	relation.ModifierId = userDoc.UserId
	relation.CreatedAt = time.Now()

	return s.SavePdfUploadRecords(ctx, paperPdf, paper, userDoc, relation)
}

// 修改文件上传的相关记录
func (s *UserDocService) ModifyUploadRecords(ctx context.Context, paperPdf *pdfModel.PaperPdf, paper *paperModel.Paper, userDoc *model.UserDoc) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.ModifyUploadRecords")
	defer span.Finish()
	// 使用事务保存所有记录
	err := s.userDocDAO.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 将事务对象放入 context，确保内部操作使用同一个事务连接
		// 这对于 SQLite 尤其重要，因为 SQLite 使用数据库级锁
		txCtx := context.WithValue(ctx, dao.TransactionContextKey, tx)

		// 1. 修改 PaperPdf 记录
		if err := s.paperPdfService.ModifyPaperPDF(txCtx, paperPdf); err != nil {
			s.logger.Error("msg", "modify paper pdf failed", "error", err.Error(), "userId", userDoc.UserId)
			return err
		}
		// 2. 修改 Paper 记录
		if err := s.paperService.ModifyPaper(txCtx, paper); err != nil {
			s.logger.Error("msg", "modify paper failed", "error", err.Error(), "userId", userDoc.UserId)
			return err
		}
		// 3. 修改用户文档记录
		if err := s.userDocDAO.Modify(txCtx, userDoc); err != nil {
			s.logger.Error("msg", "modify user doc failed", "error", err.Error(), "userId", userDoc.UserId)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 处理文件秒传的逻辑
func (s *UserDocService) HandleFileFastUpload(ctx context.Context, userId string, fileName string, userDoc *model.UserDoc, paperPdf *pdfModel.PaperPdf, paper *paperModel.Paper) (*docpb.HandleFileFastUploadResp, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.HandleFileFastUpload")
	defer span.Finish()
	//定义返回值
	resp := &docpb.HandleFileFastUploadResp{}

	queryUserDoc, queryPaperPdf, _, shaErr := s.GetUserUploadBaseDataBySHA256AndUserId(ctx, userId, paperPdf.FileSHA256)
	if shaErr == nil {
		resp.PdfId = queryUserDoc.PdfId
		resp.DocInfo = &docpb.MyCollectedDocInfo{
			Id:      queryUserDoc.Id,
			PaperId: queryPaperPdf.PaperId,
			PdfId:   queryPaperPdf.Id,
			UserId:  userId,
		}
		resp.FileSHA256 = paperPdf.FileSHA256
		return resp, nil
	}

	paperId := idgen.GenerateUUID()
	paperPdfId := idgen.GenerateUUID()
	// 创建PDF记录对象
	paperPdfNew := paperPdf
	paperPdfNew.PaperId = paperId
	paperPdfNew.BaseModel.Id = paperPdfId
	paperPdfNew.BaseModel.CreatorId = userId
	paperPdfNew.BaseModel.ModifierId = userId
	paperPdfNew.BaseModel.CreatedAt = time.Now()
	paperPdfNew.BaseModel.UpdatedAt = time.Time{}
	// 创建用户文档对象
	userDocNew := userDoc
	userDocNew.PaperId = paperId
	userDocNew.PdfId = paperPdfId
	userDocNew.NewPaper = false
	userDocNew.UserId = userId
	userDocId := idgen.GenerateUUID()
	userDocNew.BaseModel.Id = userDocId
	userDocNew.BaseModel.CreatorId = userId
	userDocNew.BaseModel.ModifierId = userId
	userDocNew.BaseModel.CreatedAt = time.Now()
	userDocNew.BaseModel.UpdatedAt = time.Time{}
	//创建paper对象并保存到数据库    todo: 这里存储之后缺少一步DOC_INFO_CHANGE的逻辑
	paperNew := paper
	paperNew.PaperId = paperId
	paperNew.OwnerId = userId
	paperNew.BaseModel.Id = paperId
	paperNew.BaseModel.CreatorId = userId
	paperNew.BaseModel.ModifierId = userId
	paperNew.BaseModel.CreatedAt = time.Now()
	paperNew.BaseModel.UpdatedAt = time.Time{}

	// 对应关系
	userDocFolderRelation := &model.UserDocFolderRelation{
		UserId:   userId,
		DocId:    userDocId,
		FolderId: "0",
	}
	userDocFolderRelation.CreatorId = userId
	userDocFolderRelation.ModifierId = userId
	userDocFolderRelation.CreatedAt = time.Now()
	userDocFolderRelation.UpdatedAt = time.Time{}
	// 使用事务处理所有数据库操作
	err := s.SavePdfUploadRecords(ctx, paperPdfNew, paperNew, userDocNew, userDocFolderRelation)
	if err != nil {
		s.logger.Error("msg", "file fast upload failed", "userId", userId, "fileName", fileName, "error", err)
		return nil, err
	}
	resp.PdfId = paperPdfId
	resp.DocInfo = &docpb.MyCollectedDocInfo{
		Id:      userDocId,
		PaperId: paperId,
		PdfId:   paperPdfId,
		UserId:  userId,
	}
	resp.FileSHA256 = paperPdf.FileSHA256
	return resp, nil
}

// GetUserDocId 根据论文ID或PDF ID获取用户文档ID
func (s *UserDocService) GetUserDocId(ctx context.Context, paperId, pdfId, userId string) (*model.UserDoc, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUserDocId")
	defer span.Finish()

	var userDocResult *model.UserDoc
	// 如果有论文ID，根据用户ID和论文ID获取文档
	if paperId != "0" {
		userDoc, err := s.GetByUserIDAndPaperID(ctx, userId, paperId)
		if err != nil {
			s.logger.Error("msg", "get user doc by user id and paper id failed", "error", err.Error(), "userId", userId, "paperId", paperId)
			return nil, err
		}
		userDocResult = userDoc
	} else if pdfId != "0" {
		// 如果有PDF ID，根据用户ID和PDF ID获取文档
		userDoc, err := s.GetByUserIdAndPdfId(ctx, userId, pdfId)
		if err != nil {
			s.logger.Error("msg", "get user doc by user id and pdf id failed", "error", err.Error(), "userId", userId, "pdfId", pdfId)
			return nil, err
		}
		userDocResult = userDoc
	}
	return userDocResult, nil
}

// ManualUpdateDocCiteInfo 手动更新文档引用信息
func (s *UserDocService) ManualUpdateDocCiteInfo(ctx context.Context, userId string, req *docpb.ManualUpdateDocCiteInfoReq) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.ManualUpdateDocCiteInfo")
	defer span.Finish()

	return s.ManualUpdateDocCiteInfoByBoolean(ctx, userId, req, false)
}

// ManualUpdateDocCiteInfoByBoolean 手动更新文档引用信息，根据来源不同有差异
func (s *UserDocService) ManualUpdateDocCiteInfoByBoolean(ctx context.Context, userId string, req *docpb.ManualUpdateDocCiteInfoReq, isFromSearchUpdate bool) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.ManualUpdateDocCiteInfoByBoolean")
	defer span.Finish()

	userDoc, err := s.GetUserDocId(ctx, *req.PaperId, *req.PdfId, userId)
	if err != nil {
		return errors.Wrap(err, "get user doc by user id and paper id failed")
	}
	if userDoc == nil {
		return errors.Biz("user doc not found")
	}
	// 处理发布日期 - 根据来源不同有差异
	if isFromSearchUpdate {
		userDoc.PublishDate = req.GetPublishDate()
		userDoc.PublishDateEdited = req.GetPublishDate() != ""
	} else {
		if req.GetPublishDate() == "" {
			userDoc.PublishDate = ""
		} else {
			userDoc.PublishDate = req.GetPublishDate()
		}
		userDoc.PublishDateEdited = true
	}

	// 更新发表场所 - 两种情况相同
	userDoc.Venue = req.GetVenue()
	userDoc.VenueEdited = true

	// 处理作者列表 - 根据来源不同有差异
	if len(req.GetAuthorList()) > 0 {
		// 将作者列表转换为JSON字符串 json字符串不要包含id属性
		// authorDescBytes, err := json.Marshal(req.GetAuthorList())
		// if err != nil {
		// 	return errors.Wrap(err, "marshal author list failed")
		// }
		// userDoc.AuthorDesc = string(authorDescBytes)
		// 创建一个新的切片来存储不包含Id的作者信息  这里其实应该使用上面注释的内容，但是因为前端的问题，这里存入了id会导致错误
		authorListWithoutId := make([]map[string]interface{}, len(req.GetAuthorList()))
		// 遍历authorList，复制除Id外的所有字段
		for i, author := range req.GetAuthorList() {
			// 创建新的map
			authorMap := make(map[string]interface{})
			// 复制需要的字段
			if author.Literal != "" {
				authorMap["literal"] = author.Literal
			}
			authorMap["isAuthentication"] = author.GetIsAuthentication()
			// 添加其他需要的字段...
			authorListWithoutId[i] = authorMap
		}
		// 将不包含Id的作者列表转换为JSON
		authorDescBytes, err := json.Marshal(authorListWithoutId)
		if err != nil {
			return errors.Wrap(err, "marshal author list failed")
		}
		userDoc.AuthorDesc = string(authorDescBytes)

	} else if isFromSearchUpdate {
		// 只有在搜索更新且作者列表为空时才清空
		userDoc.AuthorDesc = ""
	}
	// 两种情况都需要标记作者描述已编辑
	userDoc.AuthorDescEdited = true

	// 以下字段在两种情况下的更新逻辑相同
	// 更新文档名称
	userDoc.DocName = req.GetDocName()
	userDoc.UserEditedDocName = req.GetDocName()
	userDoc.DocNameEdited = true

	// 更新文档类型
	userDoc.UserEditedDocType = req.GetDocType()
	userDoc.DocTypeEdited = true

	// 更新DOI
	userDoc.UserEditedDoi = req.GetDoi()
	userDoc.DoiEdited = true

	// 更新期号
	userDoc.UserEditedIssue = req.GetIssue()
	userDoc.IssueEdited = true

	// 更新页码
	userDoc.UserEditedPage = req.GetPage()
	userDoc.PageEdited = true

	// 更新分区
	userDoc.UserEditedPartition = req.GetPartition()
	userDoc.PartitionEdited = true

	// 更新卷号
	userDoc.UserEditedVolume = req.GetVolume()
	userDoc.VolumeEdited = true

	// 调用修改方法 - 两种情况使用相同的修改方法
	return s.userDocDAO.Modify(ctx, userDoc)
}

// GetDocTreeRootNodeOfDocsByUserId 根据userId查询用户未分类的所有文献; isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
func (s *UserDocService) GetDocTreeRootNodeOfDocsByUserId(ctx context.Context, userId string, isOnlyShowNote bool) ([]*docpb.SimpleUserDocInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUnclassifiedDocsByUserId")
	defer span.Finish()

	docFolderRelas, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "获取用户文档文件夹关系列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Biz("")
	}

	s.logger.Info("msg", "docFolderRelas", docFolderRelas)
	// 从关系数组中提取文档ID
	var docIds []string
	for _, relation := range docFolderRelas {
		if relation.FolderId == "0" {
			docIds = append(docIds, relation.DocId)
		}
	}

	docs, err := s.userDocDAO.GetByUserIdAndWithIds(ctx, userId, docIds)
	if err != nil {
		s.logger.Error("msg", "获取用户文档列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Biz("")
	}

	// 如果只显示关联NoteId的文献, 过滤掉NoteId为0的文献
	if isOnlyShowNote {
		var filteredDocs []model.UserDoc
		for _, doc := range docs {
			if doc.NoteId != "0" {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		docs = filteredDocs
	}

	// 将docs转换为docsInfos数组
	docsInfos := make([]*docpb.SimpleUserDocInfo, 0, len(docs))
	for _, doc := range docs {
		docInfo := &docpb.SimpleUserDocInfo{
			DocName:    doc.DocName,
			DocId:      doc.Id,
			PdfId:      doc.PdfId,
			PaperId:    doc.PaperId,
			NoteId:     doc.NoteId,
			ModifyDate: uint64(doc.UpdatedAt.UnixMilli()),
		}
		docsInfos = append(docsInfos, docInfo)
	}

	return docsInfos, nil
}

// GetFolderInfosByUserId 根据userId查询用户已分类的所有文献夹-文献列表
func (s *UserDocService) GetFolderInfosByUserId(ctx context.Context, userId string) ([]*docpb.UserDocFolderInfo, int32, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUnclassifiedDocsByUserId")
	defer span.Finish()

	// 2. 通过用户ID调用DAO的根据用户ID查询未删除文献列表
	userDocs, err := s.userDocDAO.GetAllUserDocsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", userId, "error", err)
		return nil, 0, errors.Wrap(err, "get user docs by user id failed")
	}

	// 3. 创建文献ID到文献对象的映射
	// userDocMap := make(map[string]*model.UserDoc)
	// for i := range userDocs {
	// 	userDocMap[userDocs[i].Id] = &userDocs[i]
	// }

	// 5. 设置总文献数量
	var totalDocCount = int32(len(userDocs))

	// 6. 构建文件夹目录树和未分类文献列表

	// 构建 docId-userDoc Map
	docId2UserDocMap := make(map[string]*model.UserDoc)
	for i := range userDocs {
		docId2UserDocMap[userDocs[i].Id] = &userDocs[i]
	}

	// 获取用户所有文件夹
	userFolders, err := s.userDocFolderService.GetUserDocFolders(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get user doc folders failed", "userId", userId, "error", err)
		return nil, 0, errors.Wrap(err, "get user doc folders failed")
	}

	// 获取用户文献和文件夹关联关系
	docFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get doc folder relations by user id failed", "userId", userId, "error", err)
		return nil, 0, errors.Wrap(err, "get doc folder relations by user id failed")
	}

	// 构建文件夹-关联关系Map
	folderAndDocRelationsMap := make(map[string][]*model.UserDocFolderRelation)
	for i := range docFolderRelations {
		folderId := docFolderRelations[i].FolderId
		if _, ok := folderAndDocRelationsMap[folderId]; !ok {
			folderAndDocRelationsMap[folderId] = make([]*model.UserDocFolderRelation, 0)
		}
		folderAndDocRelationsMap[folderId] = append(folderAndDocRelationsMap[folderId], &docFolderRelations[i])
	}

	// 构造文件夹目录树
	folderInfos := s.buildSimpleFolderTree(ctx, docId2UserDocMap, userFolders, folderAndDocRelationsMap)

	// 排序（按照 Sort 字段降序排列）
	sort.Slice(folderInfos, func(i, j int) bool {
		return folderInfos[i].Sort > folderInfos[j].Sort
	})

	return folderInfos, totalDocCount, nil
}

// GetFolderInfosByUserIdAndIds 根据userId查询用户已分类的所有文献夹-文献列表
func (s *UserDocService) GetFolderInfosByUserIdAndIds(ctx context.Context, userId string, docIds []string) ([]*docpb.UserDocFolderInfo, int32, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetUnclassifiedDocsByUserId")
	defer span.Finish()

	// 2. 通过用户ID调用DAO的根据用户ID查询未删除文献列表
	userDocs, err := s.userDocDAO.GetByUserIdAndWithIds(ctx, userId, docIds)
	if err != nil {
		s.logger.Error("msg", "get user docs by user id failed", "userId", userId, "error", err)
		return nil, 0, errors.Wrap(err, "get user docs by user id failed")
	}

	// 3. 创建文献ID到文献对象的映射
	// userDocMap := make(map[string]*model.UserDoc)
	// for i := range userDocs {
	// 	userDocMap[userDocs[i].Id] = &userDocs[i]
	// }

	// 5. 设置总文献数量
	var totalDocCount = int32(len(userDocs))

	// 6. 构建文件夹目录树和未分类文献列表

	// 构建 docId-userDoc Map
	docId2UserDocMap := make(map[string]*model.UserDoc)
	for i := range userDocs {
		docId2UserDocMap[userDocs[i].Id] = &userDocs[i]
	}

	// 获取用户所有文件夹
	userFolders, err := s.userDocFolderService.GetUserDocFolders(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get user doc folders failed", "userId", userId, "error", err)
		return nil, 0, errors.Wrap(err, "get user doc folders failed")
	}

	// 获取用户文献和文件夹关联关系
	docFolderRelations, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByUserID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "get doc folder relations by user id failed", "userId", userId, "error", err)
		return nil, 0, errors.Wrap(err, "get doc folder relations by user id failed")
	}

	// 构建文件夹-关联关系Map
	folderAndDocRelationsMap := make(map[string][]*model.UserDocFolderRelation)
	for i := range docFolderRelations {
		folderId := docFolderRelations[i].FolderId
		if _, ok := folderAndDocRelationsMap[folderId]; !ok {
			folderAndDocRelationsMap[folderId] = make([]*model.UserDocFolderRelation, 0)
		}
		folderAndDocRelationsMap[folderId] = append(folderAndDocRelationsMap[folderId], &docFolderRelations[i])
	}

	// 构造文件夹目录树
	folderInfos := s.buildSimpleFolderTree(ctx, docId2UserDocMap, userFolders, folderAndDocRelationsMap)

	// 排序（按照 Sort 字段降序排列）
	sort.Slice(folderInfos, func(i, j int) bool {
		return folderInfos[i].Sort > folderInfos[j].Sort
	})

	return folderInfos, totalDocCount, nil
}

// GetDocTreeNodeByFolderId 根据userId和文件Id查询用户的文献夹信息;参数isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
func (s *UserDocService) GetDocTreeNodeByFolderId(ctx context.Context, userId string, folderId string, isOnlyShowNote bool) (*docpb.UserDocFolderInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocTreeNodeByFolderId")
	defer span.Finish()

	// 获取文件夹信息
	docFolder, err := s.userDocFolderService.GetByIdAndUserId(ctx, folderId, userId)
	if err != nil {
		s.logger.Error("msg", "get user doc folders failed", "userId", userId, "error", err)
		return nil, errors.Wrap(err, "get user doc folders failed")
	}

	if docFolder == nil {
		return nil, nil
	}

	docFolderInfo := &docpb.UserDocFolderInfo{
		Name:     docFolder.Name,
		DocCount: 0,
		ParentId: docFolder.ParentId,
		FolderId: docFolder.Id,
		Level:    0,
		Sort:     uint32(docFolder.Sort),
	}

	// 设置文件列表
	docFolderInfo.DocInfos, err = s.GetDocsByFolderId(ctx, userId, docFolderInfo.FolderId, isOnlyShowNote)
	if err != nil {
		s.logger.Error("msg", "get docs by folder id failed", "userId", userId, "folderId", folderId, "error", err)
		return nil, errors.Wrap(err, "get docs by folder id failed")
	}

	// 设置子文件列表
	childrenFolders, err := s.userDocFolderService.GetChildrenFoldersByFolderId(ctx, docFolderInfo.FolderId)
	if err != nil {
		s.logger.Error("msg", "get children folders by folder id failed", "userId", userId, "folderId", folderId, "error", err)
		return nil, errors.Wrap(err, "get children folders by folder id failed")
	}

	// 递归处理子文件夹
	if len(childrenFolders) > 0 {
		docFolderInfo.ChildrenFolders = make([]*docpb.UserDocFolderInfo, 0, len(childrenFolders))
		for _, childFolder := range childrenFolders {
			// 递归调用获取子文件夹的完整信息
			childFolderInfo, err := s.GetDocTreeNodeByFolderId(ctx, userId, childFolder.Id, isOnlyShowNote)
			if err != nil {
				s.logger.Error("msg", "get child folder info failed", "userId", userId, "childFolderId", childFolder.Id, "error", err)
				continue
			}
			if childFolderInfo != nil {
				docFolderInfo.ChildrenFolders = append(docFolderInfo.ChildrenFolders, childFolderInfo)
			}
		}
	}

	// 设置总数 - 当前文件夹中的文档数量 + 所有子文件夹中的文档数量
	docCount := uint32(len(docFolderInfo.DocInfos))

	// 递归累加所有子文件夹中的文档数量
	if len(docFolderInfo.ChildrenFolders) > 0 {
		for _, childFolder := range docFolderInfo.ChildrenFolders {
			docCount += childFolder.DocCount
		}
	}

	// 设置最终的文档总数
	docFolderInfo.DocCount = docCount

	return docFolderInfo, nil
}

// GetDocsByFolderId 根据userId查询folderId文件夹下的文献, folderId=0时返回未分类文献; 参数isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
func (s *UserDocService) GetDocsByFolderId(ctx context.Context, userId string, folderId string, isOnlyShowNote bool) ([]*docpb.SimpleUserDocInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetDocsByFolderId")
	defer span.Finish()

	if folderId == "0" {
		// 根目录,TODO 返回未分类文献
		docInfos, err := s.GetDocTreeRootNodeOfDocsByUserId(ctx, userId, isOnlyShowNote)
		if err != nil {
			s.logger.Error("msg", "获取用户未分类文献列表失败", "userId", userId, "error", err.Error())
			return nil, errors.Biz("")
		}
		return docInfos, nil
	}

	docFolderRelas, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByFolderId(ctx, userId, folderId)
	if err != nil {
		s.logger.Error("msg", "获取用户文档文件夹关系列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Biz("获取用户文档文件夹关系列表失败")
	}

	s.logger.Info("msg", "docFolderRelas", docFolderRelas)
	// 从关系数组中提取文档ID
	var docIds []string
	for _, relation := range docFolderRelas {
		docIds = append(docIds, relation.DocId)
	}

	docs, err := s.userDocDAO.GetUserDocsByIds(ctx, docIds)
	if err != nil {
		s.logger.Error("msg", "获取用户文档列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Wrap(err, "获取用户文档列表失败")
	}

	// 如果只显示关联NoteId的文献, 过滤掉NoteId为0的文献
	if isOnlyShowNote {
		var filteredDocs []model.UserDoc
		for _, doc := range docs {
			if doc.NoteId != "0" {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		docs = filteredDocs
	}

	// 将docs转换为docsInfos数组
	docsInfos := make([]*docpb.SimpleUserDocInfo, 0, len(docs))
	for _, doc := range docs {
		docInfo := &docpb.SimpleUserDocInfo{
			DocName:    doc.DocName,
			DocId:      doc.Id,
			PdfId:      doc.PdfId,
			PaperId:    doc.PaperId,
			NoteId:     doc.NoteId,
			ModifyDate: uint64(doc.UpdatedAt.UnixMilli()),
		}
		docsInfos = append(docsInfos, docInfo)
	}

	return docsInfos, nil
}

// GetAllDocsByFolderId 根据userId查询folderId文件夹及其子文件夹下的文献, folderId=0时返回所有文献; isOnlyShowNote 是否只显示已经关联noteId的文献 true:只显示关联NoteId的文献 false: 显示全部文献
func (s *UserDocService) GetAllDocsByFolderId(ctx context.Context, userId string, folderId string, isOnlyShowNote bool) ([]*docpb.SimpleUserDocInfo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetAllDocsByFolderId")
	defer span.Finish()

	if folderId == "0" {
		// 根目录,TODO 返回所有文献
		docs, err := s.userDocDAO.GetUserDocsByUserID(ctx, userId)
		if err != nil {
			s.logger.Error("msg", "获取用户文档列表失败", "userId", userId, "error", err.Error())
			return nil, errors.Wrap(err, "获取用户文档列表失败")
		}

		// 如果只显示关联NoteId的文献, 过滤掉NoteId为0的文献
		if isOnlyShowNote {
			var filteredDocs []model.UserDoc
			for _, doc := range docs {
				if doc.NoteId != "0" {
					filteredDocs = append(filteredDocs, doc)
				}
			}
			docs = filteredDocs
		}

		// 将docs转换为docsInfos数组
		docsInfos := make([]*docpb.SimpleUserDocInfo, 0, len(docs))
		for _, doc := range docs {
			docInfo := &docpb.SimpleUserDocInfo{
				DocName:    doc.DocName,
				DocId:      doc.Id,
				PdfId:      doc.PdfId,
				PaperId:    doc.PaperId,
				NoteId:     doc.NoteId,
				ModifyDate: uint64(doc.UpdatedAt.UnixMilli()),
			}
			docsInfos = append(docsInfos, docInfo)
		}
		return docsInfos, nil
	}

	// 获取文件夹信息
	folders, err := s.userDocFolderService.GetAllByIdAndUserId(ctx, folderId, userId)
	if err != nil {
		s.logger.Error("msg", "获取文件夹信息失败", "folderId", folderId, "userId", userId, "error", err.Error())
		return nil, errors.Wrap(err, "获取文件夹信息失败")
	}

	var folderIds []string
	for _, folder := range folders {
		folderIds = append(folderIds, folder.Id)
	}

	docFolderRelas, err := s.userDocFolderRelationService.GetUserDocFolderRelationsByFolderIDs(ctx, folderIds)
	if err != nil {
		s.logger.Error("msg", "获取用户文档文件夹关系列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Biz("获取用户文档文件夹关系列表失败")
	}

	var docIds []string
	for _, docFolderRela := range docFolderRelas {
		docIds = append(docIds, docFolderRela.DocId)
	}

	docs, err := s.userDocDAO.GetUserDocsByIds(ctx, docIds)
	if err != nil {
		s.logger.Error("msg", "获取用户文档列表失败", "userId", userId, "error", err.Error())
		return nil, errors.Wrap(err, "获取用户文档列表失败")
	}

	// 如果只显示关联NoteId的文献, 过滤掉NoteId为0的文献
	if isOnlyShowNote {
		var filteredDocs []model.UserDoc
		for _, doc := range docs {
			if doc.NoteId != "0" {
				filteredDocs = append(filteredDocs, doc)
			}
		}
		docs = filteredDocs
	}

	// 将docs转换为docsInfos数组
	docsInfos := make([]*docpb.SimpleUserDocInfo, 0, len(docs))
	for _, doc := range docs {
		docInfo := &docpb.SimpleUserDocInfo{
			DocName:    doc.DocName,
			DocId:      doc.Id,
			PdfId:      doc.PdfId,
			PaperId:    doc.PaperId,
			NoteId:     doc.NoteId,
			ModifyDate: uint64(doc.UpdatedAt.UnixMilli()),
		}
		docsInfos = append(docsInfos, docInfo)
	}

	return docsInfos, nil
}

// UpdateUserDocNoteIdByPdfId 根据pdfId更新UserDoc表中的noteId
func (s *UserDocService) UpdateUserDocNoteIdByPdfId(ctx context.Context, pdfId string, noteId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.UpdateUserDocNoteIdByPdfId")
	defer span.Finish()

	doc, err := s.userDocDAO.GetUserDocByPdfId(ctx, pdfId)
	if err != nil {
		s.logger.Error("msg", "get user doc by pdf id failed", "pdfId", pdfId, "error", err)
		return errors.Biz("doc.user_doc.errors.get_failed")
	}

	if doc == nil {
		return nil
	}

	doc.NoteId = noteId

	return s.userDocDAO.ModifyExcludeNull(ctx, doc)
}

func (s *UserDocService) GetAuthors(ctx context.Context, docId string) (*docpb.GetAuthorsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetAuthors")
	defer span.Finish()

	resp := &docpb.GetAuthorsResponse{}
	//
	userDoc, err := s.userDocDAO.FindById(ctx, docId)
	if err != nil {
		s.logger.Error("msg", "get user doc by id failed", "docId", docId, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}
	if userDoc == nil {
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}
	displaySimpleAuthor := helper.GetDisplaySimpleAuthorByUserDoc(userDoc)
	resp.DisplayAuthor = displaySimpleAuthor
	return resp, nil
}

// 统一的更新方法
func (s *UserDocService) UpdateUserDocByCustomType(ctx context.Context, userId string, docId string, req *docBean.UserDocUpdateRequest) (*docBean.UserDocUpdateResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.UpdateUserDocByCustomType")
	defer span.Finish()

	// 1. 获取并验证用户文档
	userDoc, err := s.getUserDocAndValidate(ctx, docId, userId)
	if err != nil {
		return nil, err
	}

	// 2. 根据更新类型执行不同的更新逻辑
	if err := s.updateUserDocField(ctx, userDoc, req); err != nil {
		return nil, err
	}

	// 3. 执行更新
	if err := s.userDocDAO.Modify(ctx, userDoc); err != nil {
		s.logger.Error("msg", "update user doc failed", "docId", docId, "updateType", req.UpdateType, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.update_failed")
	}

	// 4. 准备响应
	return s.prepareUpdateResponse(ctx, userDoc, req)
}

// getUserDocAndValidate 获取并验证用户文档
func (s *UserDocService) getUserDocAndValidate(ctx context.Context, docId, userId string) (*model.UserDoc, error) {
	userDoc, err := s.userDocDAO.FindById(ctx, docId)
	if err != nil {
		s.logger.Error("msg", "get user doc by id failed", "docId", docId, "error", err)
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}
	if userDoc == nil {
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}
	if userDoc.CreatorId != userId {
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}
	return userDoc, nil
}

// updateUserDocField 根据更新类型更新用户文档字段
func (s *UserDocService) updateUserDocField(ctx context.Context, userDoc *model.UserDoc, req *docBean.UserDocUpdateRequest) error {
	switch req.UpdateType {
	case docBean.UpdateTypeAuthors:
		return s.updateAuthors(ctx, userDoc, req.Authors)
	case docBean.UpdateTypePublishDate:
		s.updatePublishDate(userDoc, req.PublishDate)
	case docBean.UpdateTypeVenue:
		s.updateVenue(userDoc, req.Venue)
	case docBean.UpdateTypeJcrPartion:
		userDoc.UserEditedJcrPartion = req.JcrPartion
	case docBean.UpdateTypeRemark:
		userDoc.Remark = req.Remark
	case docBean.UpdateTypeImpactFactor:
		if req.ImpactFactor != nil {
			userDoc.UserEditedImpactOfFactor = *req.ImpactFactor
		} else {
			userDoc.UserEditedImpactOfFactor = 0
		}
	case docBean.UpdateTypeImportanceScore:
		userDoc.ImportanceScore = int(req.ImportanceScore)
	default:
		return errors.Biz("doc.user_doc.errors.invalid_update_type")
	}
	return nil
}

// updateAuthors 更新作者信息
func (s *UserDocService) updateAuthors(ctx context.Context, userDoc *model.UserDoc, authors []*docpb.BaseAuthorInfo) error {
	if len(authors) == 0 {
		userDoc.AuthorDesc = userDoc.MetaAuthors
		userDoc.AuthorDescEdited = false
		return nil
	}

	authorInfos := helper.BaseAuthorInfoToUserDocAuthorInfo(authors)
	authorInfoStr, err := json.Marshal(authorInfos)
	if err != nil {
		s.logger.Error("msg", "marshal author info failed", "error", err)
		return errors.Biz("doc.user_doc.errors.update_failed")
	}

	userDoc.AuthorDesc = string(authorInfoStr)
	userDoc.AuthorDescEdited = true
	return nil
}

// updatePublishDate 更新发布日期
func (s *UserDocService) updatePublishDate(userDoc *model.UserDoc, publishDate string) {
	if publishDate == "" {
		userDoc.PublishDate = userDoc.MetaPublishDate
		userDoc.PublishDateEdited = false
	} else {
		userDoc.PublishDate = publishDate
		userDoc.PublishDateEdited = true
	}
}

// updateVenue 更新收录情况
func (s *UserDocService) updateVenue(userDoc *model.UserDoc, venue string) {
	if venue == "" {
		userDoc.Venue = userDoc.MetaVenues
		userDoc.VenueEdited = false
	} else {
		userDoc.Venue = venue
		userDoc.VenueEdited = true
	}
}

// prepareUpdateResponse 准备更新响应
func (s *UserDocService) prepareUpdateResponse(ctx context.Context, userDoc *model.UserDoc, req *docBean.UserDocUpdateRequest) (*docBean.UserDocUpdateResponse, error) {
	response := &docBean.UserDocUpdateResponse{}

	switch req.UpdateType {
	case docBean.UpdateTypeAuthors:
		displaySimpleAuthor := helper.GetDisplaySimpleAuthorByUserDoc(userDoc)
		response.AuthorsResponse = &docpb.UpdateAuthorsResponse{
			NewAuthor: displaySimpleAuthor,
		}
	case docBean.UpdateTypePublishDate:
		displayPublishDate := helper.GetPublishDateStrByUserDoc(userDoc)
		response.PublishDateResponse = &docpb.UpdatePublishDateResponse{
			NewPublishDate: displayPublishDate,
		}
	case docBean.UpdateTypeVenue:
		displayVenue := helper.GetVenueStrByUserDoc(userDoc)
		response.VenueResponse = &docpb.UpdateVenueResponse{
			NewVenue: displayVenue,
		}
	case docBean.UpdateTypeJcrPartion, docBean.UpdateTypeImpactFactor:
		if err := s.prepareJcrOrImpactResponse(ctx, userDoc, req, response); err != nil {
			return nil, err
		}
	}

	return response, nil
}

// prepareJcrOrImpactResponse 准备JCR分区或影响因子响应
func (s *UserDocService) prepareJcrOrImpactResponse(ctx context.Context, userDoc *model.UserDoc, req *docBean.UserDocUpdateRequest, response *docBean.UserDocUpdateResponse) error {
	venue := userDoc.Venue
	if venue == "" {
		venue = userDoc.MetaVenues
	}

	jcrPartionEntity, err := s.paperJcrService.GetPaperJcrEntityByVenue(ctx, venue)
	if err != nil {
		s.logger.Error("msg", "get paper jcr entity by venue failed", "venue", venue, "error", err)
		return errors.Biz("doc.user_doc.errors.get_failed")
	}

	if req.UpdateType == docBean.UpdateTypeJcrPartion {
		response.JcrPartionUpdateResponse = s.createJcrPartionResponse(jcrPartionEntity, req.JcrPartion)
	} else if req.ImpactFactor != nil { // UpdateTypeImpactFactor
		response.ImpactFactorResponse = s.createImpactFactorResponse(jcrPartionEntity, *req.ImpactFactor)
	}

	return nil
}

// createJcrPartionResponse 创建JCR分区响应
func (s *UserDocService) createJcrPartionResponse(jcrEntity *paperModel.PaperJcrEntity, currentJcrPartion string) *docpb.JcrPartionUpdateResponse {
	response := &docpb.JcrPartionUpdateResponse{
		CurrentJcrPartion: currentJcrPartion,
		RollbackEnable:    false,
	}

	if jcrEntity == nil || jcrEntity.JcrPartion == "" {
		response.RollbackEnable = currentJcrPartion != ""
	} else {
		response.OriginalJcrPartion = jcrEntity.JcrPartion
		response.RollbackEnable = jcrEntity.JcrPartion != currentJcrPartion
	}

	return response
}

// createImpactFactorResponse 创建影响因子响应
func (s *UserDocService) createImpactFactorResponse(jcrEntity *paperModel.PaperJcrEntity, impactFactor float32) *docpb.ImpactFactorResponse {
	response := &docpb.ImpactFactorResponse{
		CurrentImpactOfFactor: impactFactor,
		RollbackEnable:        false,
	}

	if jcrEntity == nil || jcrEntity.JcrPartion == "" {
		response.RollbackEnable = impactFactor != 0
	} else {
		response.OriginalImpactOfFactor = jcrEntity.ImpactOfFactor
		response.RollbackEnable = jcrEntity.ImpactOfFactor != impactFactor
	}

	return response
}

func (s *UserDocService) AttachDocToClassify(ctx context.Context, userId string, docId string, classifyId string, classifyName string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.AttachDocToClassify")
	defer span.Finish()

	//如果分类没有找到，则默认是新增分类的id
	if classifyId == "" {
		classifyId = idgen.GenerateUUID()
		userDocClassify := &model.UserDocClassify{
			UserId: userId,
			Name:   classifyName,
		}
		userDocClassify.Id = classifyId
		err := s.userDocClassifyService.Save(ctx, userDocClassify)
		if err != nil {
			return errors.Biz("doc.user_doc.errors.save_failed")
		}
	} else {
		//验证分类id
		userDocClassify, err := s.userDocClassifyService.FindById(ctx, classifyId)
		if err != nil {
			return errors.Biz("doc.user_doc.errors.get_failed")
		}
		if userDocClassify == nil {
			return errors.Biz("doc.user_doc.errors.get_failed")
		}
	}
	//查询该文档是否存在此分类
	classifyRelationList, err := s.docClassifyRelationService.GetByClassifyIdAndDocId(ctx, classifyId, docId)
	if err != nil {
		return errors.Biz("doc.user_doc.errors.get_failed")
	}
	if len(classifyRelationList) > 0 {
		return nil
	}
	docClassifyRelation := &model.DocClassifyRelation{
		UserId:     userId,
		DocId:      docId,
		ClassifyId: classifyId,
	}
	err = s.docClassifyRelationService.Save(ctx, docClassifyRelation)
	if err != nil {
		return errors.Biz("doc.user_doc.errors.save_failed")
	}
	return nil
}

func (s *UserDocService) RemoveDocFromClassify(ctx context.Context, userId string, docId string, classifyId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.RemoveDocFromClassify")
	defer span.Finish()

	err := s.docClassifyRelationService.DeleteByUserIdDocIdClassifyId(ctx, userId, docId, classifyId)
	if err != nil {
		return errors.Biz("doc.user_doc.errors.update_failed")
	}
	return nil
}

func (s *UserDocService) UpdateReadStatus(ctx context.Context, userId string, req *docpb.UpdateReadStatusRequest) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.UpdateReadStatus")
	defer span.Finish()

	//查找用户文献
	userDoc, err := s.GetUserDocId(ctx, req.GetPaperId(), req.GetPdfId(), userId)
	if err != nil {
		return errors.Biz("doc.user_doc.errors.get_failed")
	}
	if userDoc == nil {
		return errors.Biz("doc is not exist")
	}
	if req.Status != nil {
		// 根据状态设置阅读进度
		var progress int32
		status := req.GetStatus()
		// 根据状态值设置进度
		switch status {
		case docpb.DocReadingStatus_READING:
			userDoc.ReadingStatus = docpb.DocReadingStatus_READING.String()
			progress = 1
		case docpb.DocReadingStatus_READ:
			userDoc.ReadingStatus = docpb.DocReadingStatus_READ.String()
			progress = 100
		}
		userDoc.Progress = int(progress)
		err = s.userDocDAO.ModifyExcludeNull(ctx, userDoc)
		if err != nil {
			return errors.Biz("doc.user_doc.errors.update_failed")
		}
		return nil
	}
	// 根据请求中的进度值设置阅读状态
	if req.Progress != nil {
		progress := req.GetProgress()
		var status docpb.DocReadingStatus
		if progress >= 100 {
			status = docpb.DocReadingStatus_READ
		} else {
			status = docpb.DocReadingStatus_READING
		}
		userDoc.ReadingStatus = status.String()
		userDoc.Progress = int(progress)
		userDoc.LastReadTime = time.Now()
		err = s.userDocDAO.ModifyExcludeNull(ctx, userDoc)
		if err != nil {
			return errors.Biz("doc.user_doc.errors.update_failed")
		}
	}
	return nil
}

func (s *UserDocService) GetLatestReadDocList(ctx context.Context, userId string, latestReadSize int) (*docpb.GetLatestReadDocListResp, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.GetLatestReadDocList")
	defer span.Finish()

	resp := &docpb.GetLatestReadDocListResp{}
	resp.CurrentTime = uint64(time.Now().Unix())
	// 获取用户阅读过的文献，latest_read_time不为空  按时间顺序排序
	userDocs, err := s.userDocDAO.GetUserDocsByUserIDAndLatestReadTimeIsNotNull(ctx, userId, latestReadSize)
	if err != nil {
		return nil, errors.Biz("doc.user_doc.errors.get_failed")
	}
	if len(userDocs) == 0 {
		return resp, nil
	}
	for _, userDoc := range userDocs {
		// 过滤掉无效的时间数据（零值或异常早期时间）
		if userDoc.LastReadTime.IsZero() || userDoc.LastReadTime.Year() < 1900 {
			continue
		}

		latestReadDocInfo := &docpb.LatestReadDocInfo{
			PaperId:      userDoc.PaperId,
			PdfId:        userDoc.PdfId,
			DocName:      userDoc.DocName,
			PublishDate:  userDoc.PublishDate,
			Venue:        userDoc.Venue,
			Remark:       userDoc.Remark,
			LastReadTime: uint64(userDoc.LastReadTime.Unix()),
		}
		resp.DocInfos = append(resp.DocInfos, latestReadDocInfo)
	}
	return resp, nil
}

/**
 * 上传pdf文件(需要走解析流程)
 * fileName: 原始文件名
 * fileSHA256: 文件SHA256值
 * reader: 文件内容读取器
 * objectSize: 对象大小（字节）
 * filePage: 文件页数
 * userID: 用户ID
 * return ossRecord: 上传成功后的oss记录
 * return needParsed: 是否需要解析，不需要解析就是走了秒传的流程
 * return error: 错误信息
 */
func (s *UserDocService) UploadPdf(ctx context.Context, reader io.Reader, fileName string, fileSHA256 string, objectSize int64, filePage int64, userID string, uploadToken string) (*ossModel.OssRecord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.UploadPdf")
	defer span.Finish()

	//获取文件名的后缀
	fileExt := filepath.Ext(fileName)
	//判断后缀是否是.pdf的文件
	if fileExt != ".pdf" {
		return nil, errors.Biz("doc.user_doc.errors.upload_failed")
	}
	//生成uuid
	uuid := idgen.GenerateUUID()
	// 重新赋值文件名
	pdfFileName := fmt.Sprintf("%s%s", uuid, fileExt)
	//设置桶
	bucketType := osspb.OSSBucketEnum_PDF

	//判断是否已经存在相同的已经解析成功文件数据
	ossRecord, err := s.ossService.GetSuccessRecordByFileSHA256(ctx, fileSHA256)
	if err != nil {
		return nil, err
	}
	useStorageCapacity, err := s.paperPdfService.GetPdfFileTotalSize(ctx, userID)
	if err != nil {
		return nil, err
	}
	if ossRecord != nil {
		//如果存在记录的情况
		userDoc, paperPdf, paper, hasExitErr := s.GetUserUploadBaseDataBySHA256AndUserId(ctx, userID, fileSHA256)
		if hasExitErr != nil {
			//这里代表数据库中没有当前用户的记录
			err = s.membershipService.CreditFunDocsUpload(ctx, objectSize, int32(filePage), useStorageCapacity, func(xctx context.Context, sessionId string) error {
				s.HandleFileFastUpload(ctx, userID, fileName, userDoc, paperPdf, paper)
				return nil
			}, true)
			if err != nil {
				return nil, err
			}
		}
		//修改状态为解析完成
		// s.ChangeUploadTokenStatus(ctx, uploadToken, docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSED)
		cacheValue := docpb.UserDocParsedStatusObj{}
		cacheValue.Status = docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSED
		cacheValue.HasSha256File = true
		cacheValue.Sha256Value = fileSHA256
		parseCacheKey := fmt.Sprintf("%s%s", constant.ParsePDFStatusTokenKeyPrefix, uploadToken)
		s.cache.SetNotBizPrefix(ctx, parseCacheKey, &cacheValue, 30*time.Minute)
		return ossRecord, nil
	}
	//判断是否存在缓存
	parseCacheKey := fmt.Sprintf("%s%s", constant.ParsePDFStatusTokenKeyPrefix, uploadToken)
	cacheValue := docpb.UserDocParsedStatusObj{}
	_, cacheErr := s.cache.GetNotBizPrefix(ctx, parseCacheKey, &cacheValue)
	if cacheErr != nil {
		return nil, cacheErr
	}
	//写入缓存
	cacheValue.Status = docpb.UserDocParsedStatusEnum_UPLOADING
	cacheValue.HasSha256File = true
	cacheValue.Sha256Value = fileSHA256
	s.cache.SetNotBizPrefix(ctx, parseCacheKey, &cacheValue, 30*time.Minute)

	//不存在记录 直接走上传流程
	var uploadedRecord *ossModel.OssRecord
	err = s.membershipService.CreditFunDocsUpload(ctx, objectSize, int32(filePage), useStorageCapacity, func(xctx context.Context, sessionId string) error {
		var uploadErr error
		uploadedRecord, uploadErr = s.releaseUploadPDF(ctx, fileSHA256, pdfFileName, reader, objectSize, bucketType)
		if uploadErr != nil {
			return uploadErr
		}
		return nil
	}, true)
	if err != nil {
		return nil, err
	}
	return uploadedRecord, nil
}

// 真实的上传文件的方法
func (s *UserDocService) releaseUploadPDF(ctx context.Context, fileSHA256 string, pdfFileName string, uploadReader io.Reader, objectSize int64, bucketType osspb.OSSBucketEnum) (*ossModel.OssRecord, error) {

	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.releaseUploadPDF")
	defer span.Finish()
	//设置topic名称  该topic会对论文进行解析
	topicName := s.config.RocketMQ.Topic.UploadCallback.Name
	//自定义的key生成规则
	customObjectKeyGenerator := func(uniqueId, fileName string) string {
		return fmt.Sprintf("%s/%s/%s", fileSHA256, parseConstant.SourcePdfCatalog, pdfFileName)
	}
	metadata := map[string]string{
		"fileSHA256": fileSHA256,
		"fileName":   pdfFileName,
	}
	// 上传文件
	ossRecord, err := s.ossService.UploadObjectAndSaveRecord(
		ctx, bucketType, pdfFileName, fileSHA256, uploadReader, objectSize,
		customObjectKeyGenerator, topicName, metadata,
	)
	if err != nil {
		return nil, err
	}
	return ossRecord, nil
}

// 更新redis 论文状态的value
func (s *UserDocService) ChangeUploadTokenStatus(ctx context.Context, token string, status docpb.UserDocParsedStatusEnum) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.changeUploadTokenStatus")
	defer span.Finish()

	//更新redis 论文状态的value
	key := fmt.Sprintf("%s%s", constant.ParsePDFStatusTokenKeyPrefix, token)
	cacheValue := docpb.UserDocParsedStatusObj{}
	cacheValue.Status = status
	err := s.cache.SetNotBizPrefix(ctx, key, &cacheValue, 30*time.Minute)
	if err != nil {
		s.logger.Error("mq.ParseUploadPdfService error,redis status update failed", "error", err)
		return errors.Biz("redis status update failed")
	}
	s.logger.Info("消费者调用： 更新redis成功")
	return nil
}

func (s *UserDocService) ValidateUrl(ctx context.Context, url string) (string, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "UserDocService.validateUrl")
	defer span.Finish()

	// 验证是否为 arxiv url
	isArxivUrl, err := s.squidProxyService.ValidateIsArxivUrl(ctx, url)
	if err != nil {
		return "", errors.BizWrap("validate arxiv url failed", err)
	}
	if isArxivUrl {
		s.logger.Info("msg", "url is arxiv url", "url", url)
		// 将 arxiv 详情页面 url 替换为 pdf url
		convertedUrl, err := s.squidProxyService.ReplaceArxivUrlWithPdfUrl(ctx, url)
		if err != nil {
			return "", errors.BizWrap("replace arxiv url failed", err)
		}
		return convertedUrl, nil
	}

	// 验证是否为 github url
	isGithubUrl, err := s.squidProxyService.ValidateIsGithubUrl(ctx, url)
	if err != nil {
		return "", errors.BizWrap("validate github url failed", err)
	}
	if isGithubUrl {
		s.logger.Info("msg", "url is github url", "url", url)
		// 将 github 详情页面 url 替换为 pdf url
		convertedUrl, err := s.squidProxyService.ReplaceGithubUrlWithPdfUrl(ctx, url)
		if err != nil {
			return "", errors.BizWrap("replace github url failed", err)
		}
		return convertedUrl, nil
	}

	// 如果不是特殊 URL，返回原 URL
	return url, nil
}
