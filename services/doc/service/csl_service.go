package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"net/http"

	"github.com/opentracing/opentracing-go"
	userContext "github.com/yb2020/odoc/pkg/context"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	pb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/constant"
	"github.com/yb2020/odoc/services/doc/dao"
	"github.com/yb2020/odoc/services/doc/model"
	"github.com/yb2020/odoc/services/doc/util"
)

// CslService 引用样式服务实现
type CslService struct {
	cslDAO                 *dao.CslDAO
	logger                 logging.Logger
	tracer                 opentracing.Tracer
	userCslRelationService *UserCslRelationService
	userDocService         *UserDocService
}

// NewCslService 创建新的引用样式服务
func NewCslService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	cslDAO *dao.CslDAO,
	userCslRelationService *UserCslRelationService,
	userDocService *UserDocService,
) *CslService {
	return &CslService{
		logger:                 logger,
		tracer:                 tracer,
		cslDAO:                 cslDAO,
		userCslRelationService: userCslRelationService,
		userDocService:         userDocService,
	}
}

// GetDefaultCslList 获取默认的引用样式列表
func (s *CslService) GetDefaultCslList(ctx context.Context, isI18n bool) (*pb.CetDefaultCslListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.GetDefaultCslList")
	defer span.Finish()

	// 从数据库获取默认引用样式列表
	defaultCslList, err := s.cslDAO.GetDefaultCsl(ctx)
	if err != nil {
		s.logger.Error("get default csl list failed", "error", err.Error())
		return nil, err
	}

	response := &pb.CetDefaultCslListResponse{}
	response.Total = uint32(len(defaultCslList))

	// 获取默认顺序
	defaultCslTitleOrderStr := constant.DefaultCslTitleOrder
	defaultTitles := strings.Split(defaultCslTitleOrderStr, constant.SPLIT_SEMICOLON)

	// 按照默认顺序构建响应列表
	var list []*pb.CslItem
	for _, defaultTitle := range defaultTitles {
		// 查找匹配的引用样式
		var matchedCsl *model.Csl
		for _, csl := range defaultCslList {
			if strings.EqualFold(csl.Title, defaultTitle) {
				matchedCsl = &csl
				break
			}
		}
		if matchedCsl == nil {
			s.logger.Error("default csl config error", "title", defaultTitle)
			return nil, errors.Biz("default csl config error")
		}
		// 构建响应项
		cslItem := &pb.CslItem{
			Id:                matchedCsl.Id,
			Title:             matchedCsl.Title,
			ShortTitle:        &matchedCsl.ShortTitle,
			FileUrl:           matchedCsl.FileUrl,
			CustomDefineTitle: &matchedCsl.CustomDefineTitle,
		}
		// 设置语言
		s.setCslItemLang(cslItem)
		list = append(list, cslItem)
	}

	// 针对国际版过滤掉特定的引用样式
	if isI18n {
		var filteredList []*pb.CslItem
		for _, item := range list {
			if !strings.EqualFold(item.Title, constant.I18nFilterCslTitle) {
				filteredList = append(filteredList, item)
			}
		}
		list = filteredList
		response.Total = uint32(len(list))
	}
	response.List = list
	return response, nil
}

// setCslItemLang 设置引用样式的语言
func (s *CslService) setCslItemLang(item *pb.CslItem) {
	defaultLang := constant.CslDefaultLang
	cslDefaultLangMapRel := constant.CslDefaultLangMapRel

	// 解析JSON映射
	var jsonObject map[string]interface{}
	err := json.Unmarshal([]byte(cslDefaultLangMapRel), &jsonObject)
	if err != nil {
		s.logger.Error("parse csl default lang map failed", "error", err.Error())
		item.Lang = &defaultLang
		return
	}
	// 获取语言
	idStr := fmt.Sprintf("%d", item.Id)
	if lang, ok := jsonObject[idStr]; ok {
		langStr := fmt.Sprintf("%v", lang)
		item.Lang = &langStr
	} else {
		item.Lang = &defaultLang
	}
}

// GetDocMetaInfo 获取文档元数据信息
func (s *CslService) GetDocMetaInfo(ctx context.Context, req *pb.DocMetaInfoHandlerReq) (*pb.DocMetaInfoSimpleVo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.GetDocMetaInfo")
	defer span.Finish()

	// 创建空的响应对象
	metaInfoVo := &pb.DocMetaInfoSimpleVo{}
	// 获取用户文档
	userId, _ := userContext.GetUserID(ctx)
	userDoc, err := s.userDocService.GetUserDocId(ctx, req.GetPaperId(), req.GetPdfId(), userId)
	if err != nil {
		s.logger.Error("get user doc failed", "error", err.Error())
		return nil, errors.Biz("user doc not found")
	}
	if userDoc == nil {
		return nil, errors.Biz("user doc not found")
	}
	metaInfoVo.Title = userDoc.DocName
	metaInfoVo.PublishTimestamp = uint64(util.GetTimestampByDate(userDoc.PublishDate))
	metaInfoVo.PublishDateStr = &userDoc.PublishDate
	metaInfoVo.Page = &userDoc.UserEditedPage
	if userDoc.UserEditedDocType != "" {
		metaInfoVo.DocType = userDoc.UserEditedDocType
		// 从常量配置中获取文档类型列表的JSON字符串
		docTypeInfoListJsonDesc := constant.DocTypeInfoListJsonDesc
		// 解析JSON字符串为DocTypeInfo列表
		var docTypeInfos []*pb.DocTypeInfo
		err = json.Unmarshal([]byte(docTypeInfoListJsonDesc), &docTypeInfos)
		if err != nil {
			s.logger.Error("parse doc type info list failed", "error", err.Error())
			return nil, err
		}
		for _, docTypeInfo := range docTypeInfos {
			if docTypeInfo.Code == metaInfoVo.DocType {
				metaInfoVo.DocTypeName = &docTypeInfo.Name
				break
			}
		}
	}
	metaInfoVo.Doi = &userDoc.UserEditedDoi
	metaInfoVo.Volume = &userDoc.UserEditedVolume
	metaInfoVo.Issue = &userDoc.UserEditedIssue
	metaInfoVo.Partition = &userDoc.UserEditedPartition
	if userDoc.VenueEdited && userDoc.Venue != "" {
		containerTitle := []string{userDoc.Venue}
		metaInfoVo.ContainerTitle = containerTitle
	}
	if userDoc.AuthorDesc != "" {
		authorList := []*pb.AuthorInfo{}
		// 将AuthorDesc转换成AuthorList
		err = json.Unmarshal([]byte(userDoc.AuthorDesc), &authorList)
		if err != nil {
			s.logger.Error("parse author list failed", "error", err.Error())
			return nil, err
		}
		metaInfoVo.AuthorList = authorList
	}
	if metaInfoVo.Doi != nil {
		url := "http://dx.doi.org" + *metaInfoVo.Doi
		metaInfoVo.Url = &url
	}

	return metaInfoVo, nil
}

// SelectByUserId 根据用户ID获取用户的引用样式列表
func (s *CslService) SelectByUserId(ctx context.Context, userId string) ([]model.Csl, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.SelectByUserId")
	defer span.Finish()

	return s.cslDAO.SelectByUserId(ctx, userId)
}

// GetMyCslList 获取我的引用样式列表
func (s *CslService) GetMyCslList(ctx context.Context, req *pb.GetMyCslListReq, userId string) (*pb.GetMyCslListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.GetMyCslList")
	defer span.Finish()

	if userId == "" {
		s.logger.Error("user not login", "userId", userId)
		return nil, errors.Biz("user not login")
	}

	response := &pb.GetMyCslListResponse{}

	// 检查用户是否已添加引用样式
	isAddedCsl := s.userCslRelationService.IsAddedCsl(ctx, userId)
	var list []*pb.MyCslItem

	// 第一次进来添加默认格式
	if !isAddedCsl {
		// TODO: 这里应该添加分布式锁，暂时简化处理
		lockKey := fmt.Sprintf("lock:addDefaultCsl:%d", userId)
		s.logger.Info("addDefaultCsl", "userId", userId, "lockKey", lockKey)

		// 再次检查，防止并发问题
		if s.userCslRelationService.IsAddedCsl(ctx, userId) {
			return response, nil
		}

		// 获取默认引用样式列表
		defaultCslList, err := s.cslDAO.GetDefaultCsl(ctx)
		if err != nil {
			s.logger.Error("get default csl list failed", "error", err.Error())
			return nil, err
		}

		// 获取默认顺序
		defaultCslTitleOrderStr := constant.DefaultCslTitleOrder
		defaultTitles := strings.Split(defaultCslTitleOrderStr, constant.SPLIT_SEMICOLON)

		// 针对国际版过滤掉特定的引用样式
		// todo:  暂时使用固定值，默认英文
		isI18n := true // 固定为 true，默认使用英文
		if isI18n {
			var filteredTitles []string
			for _, title := range defaultTitles {
				if !strings.EqualFold(title, constant.I18nFilterCslTitle) {
					filteredTitles = append(filteredTitles, title)
				}
			}
			defaultTitles = filteredTitles
		}

		// 为用户添加默认引用样式
		for i, defaultTitle := range defaultTitles {
			// 查找匹配的引用样式
			var defaultCsl *model.Csl
			for _, csl := range defaultCslList {
				if strings.EqualFold(csl.Title, defaultTitle) {
					defaultCsl = &csl
					break
				}
			}
			if defaultCsl == nil {
				return nil, errors.Biz("default csl config error")
			}

			// 创建用户引用样式关联
			relation := &model.UserCslRelation{
				UserId: userId,
				CslId:  defaultCsl.Id,
				Sort:   int32(i),
			}

			// 保存关联
			err := s.userCslRelationService.Save(ctx, relation)
			if err != nil {
				s.logger.Error("save user csl relation failed", "error", err.Error())
				return nil, err
			}

			// 构建响应项
			shortTitle := defaultCsl.ShortTitle
			customDefineTitle := defaultCsl.CustomDefineTitle
			cslItem := &pb.MyCslItem{
				Id:                defaultCsl.Id,
				Title:             defaultCsl.Title,
				ShortTitle:        &shortTitle,
				FileUrl:           defaultCsl.FileUrl,
				CustomDefineTitle: &customDefineTitle,
			}
			list = append(list, cslItem)
		}
		response.Total = uint32(len(defaultTitles))
	} else {
		// 获取用户引用样式列表
		userCslList, err := s.cslDAO.SelectByUserId(ctx, userId)
		if err != nil {
			s.logger.Error("select by user id failed", "error", err.Error())
			return nil, err
		}

		response.Total = uint32(len(userCslList))

		// 构建响应列表
		for _, csl := range userCslList {
			shortTitle := csl.ShortTitle
			customDefineTitle := csl.CustomDefineTitle
			cslItem := &pb.MyCslItem{
				Id:                csl.Id,
				Title:             csl.Title,
				ShortTitle:        &shortTitle,
				FileUrl:           csl.FileUrl,
				CustomDefineTitle: &customDefineTitle,
			}
			list = append(list, cslItem)
		}
	}

	// 设置语言
	s.setMyCslItemLang(list)
	response.List = list
	return response, nil
}

// setMyCslItemLang 设置我的引用样式的语言
func (s *CslService) setMyCslItemLang(items []*pb.MyCslItem) {
	defaultLang := constant.CslDefaultLang
	cslDefaultLangMapRel := constant.CslDefaultLangMapRel

	// 解析JSON映射
	var jsonObject map[string]interface{}
	err := json.Unmarshal([]byte(cslDefaultLangMapRel), &jsonObject)
	if err != nil {
		s.logger.Error("parse csl default lang map failed", "error", err.Error())
		for _, item := range items {
			item.Lang = &defaultLang
		}
		return
	}

	// 设置语言
	for _, item := range items {
		idStr := fmt.Sprintf("%d", item.Id)
		if lang, ok := jsonObject[idStr]; ok {
			langStr := fmt.Sprintf("%v", lang)
			item.Lang = &langStr
		} else {
			item.Lang = &defaultLang
		}
	}
}

// GetDocTypeList 获取文档类型列表
func (s *CslService) GetDocTypeList(ctx context.Context) (*pb.GetDocTypeListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.GetDocTypeList")
	defer span.Finish()

	// 从常量配置中获取文档类型列表的JSON字符串
	docTypeInfoListJsonDesc := constant.DocTypeInfoListJsonDesc

	// 解析JSON字符串为DocTypeInfo列表
	var docTypeInfos []*pb.DocTypeInfo
	err := json.Unmarshal([]byte(docTypeInfoListJsonDesc), &docTypeInfos)
	if err != nil {
		s.logger.Error("parse doc type info list failed", "error", err.Error())
		return nil, err
	}

	// 将 docTypeInfos 包装在 GetDocTypeListResponse 中返回
	response := &pb.GetDocTypeListResponse{
		DocTypeInfos: docTypeInfos,
	}

	return response, nil
}

// DownloadBibTex 下载BibTex格式的引用
func (s *CslService) DownloadBibTex(ctx context.Context, docIds []string, w http.ResponseWriter) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.DownloadBibTex")
	defer span.Finish()

	s.logger.Info("downloadBibTex", "docIds", docIds)

	// 构建请求列表
	reqs := make([]*pb.DocMetaInfoHandlerReq, 0, len(docIds))
	for _, docId := range docIds {
		docIdCopy := docId // 创建副本以避免闭包问题
		req := &pb.DocMetaInfoHandlerReq{
			DocId: &docIdCopy,
			CralwDataImmediately: func() *bool {
				b := false
				return &b
			}(),
		}
		reqs = append(reqs, req)
	}

	// 获取文档元数据映射
	metaInfoMap, err := s.GetDocMetaInfoMapByDocIds(ctx, reqs)
	if err != nil {
		s.logger.Error("get doc meta info map by doc ids failed", "error", err.Error())
		return err
	}
	// 设置HTTP响应头
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment;filename=export.bib")

	// 为每个文档生成BibTeX条目
	var buffer strings.Builder
	separator := ""

	for _, docId := range docIds {
		metaInfoSimpleVo := metaInfoMap[docId]
		if metaInfoSimpleVo == nil {
			continue
		}

		// 使用我们的辅助方法获取必要的信息
		docType := util.GetType(metaInfoSimpleVo)
		authors := util.GetAuthor(metaInfoSimpleVo)
		year := util.GetYear(metaInfoSimpleVo)
		month := util.GetMonth(metaInfoSimpleVo)
		containerTitle := util.GetContainerTitle(metaInfoSimpleVo)

		// 生成BibTeX ID
		bibID := fmt.Sprintf("doc%d", docId)

		// 开始构建BibTeX条目
		buffer.WriteString(separator)
		buffer.WriteString("@")
		buffer.WriteString(string(docType))
		buffer.WriteString("{")
		buffer.WriteString(bibID)
		buffer.WriteString(",\n")

		// 添加标题
		if metaInfoSimpleVo.GetTitle() != "" {
			buffer.WriteString("  title = {")
			buffer.WriteString(escapeSpecialChars(metaInfoSimpleVo.GetTitle()))
			buffer.WriteString("},\n")
		}

		// 添加作者
		if len(authors) > 0 {
			buffer.WriteString("  author = {")
			for i, author := range authors {
				if i > 0 {
					buffer.WriteString(" and ")
				}

				// 根据作者信息构建作者字符串
				if author.Family != "" {
					buffer.WriteString(escapeSpecialChars(author.Family))
					if author.Given != "" {
						buffer.WriteString(", ")
						buffer.WriteString(escapeSpecialChars(author.Given))
					}
				} else if author.Given != "" {
					buffer.WriteString(escapeSpecialChars(author.Given))
				} else if author.Literal != "" {
					buffer.WriteString(escapeSpecialChars(author.Literal))
				}
			}
			buffer.WriteString("},\n")
		}

		// 添加年份
		if year > 0 {
			buffer.WriteString("  year = {")
			buffer.WriteString(fmt.Sprintf("%d", year))
			buffer.WriteString("},\n")
		}

		// 添加月份
		if month > 0 {
			buffer.WriteString("  month = {")
			buffer.WriteString(fmt.Sprintf("%d", month))
			buffer.WriteString("},\n")
		}

		// 添加URL
		if metaInfoSimpleVo.GetUrl() != "" {
			buffer.WriteString("  url = {")
			buffer.WriteString(escapeSpecialChars(metaInfoSimpleVo.GetUrl()))
			buffer.WriteString("},\n")
		}

		// 添加DOI
		if metaInfoSimpleVo.GetDoi() != "" {
			buffer.WriteString("  doi = {")
			buffer.WriteString(escapeSpecialChars(metaInfoSimpleVo.GetDoi()))
			buffer.WriteString("},\n")
		}

		// 添加页码
		if metaInfoSimpleVo.GetPage() != "" {
			buffer.WriteString("  pages = {")
			buffer.WriteString(escapeSpecialChars(metaInfoSimpleVo.GetPage()))
			buffer.WriteString("},\n")
		}

		// 添加卷号
		if metaInfoSimpleVo.GetVolume() != "" {
			buffer.WriteString("  volume = {")
			buffer.WriteString(escapeSpecialChars(metaInfoSimpleVo.GetVolume()))
			buffer.WriteString("},\n")
		}

		// 添加期号
		if metaInfoSimpleVo.GetIssue() != "" {
			buffer.WriteString("  number = {")
			buffer.WriteString(escapeSpecialChars(metaInfoSimpleVo.GetIssue()))
			buffer.WriteString("},\n")
		}

		// 添加容器标题（如期刊名称）
		if containerTitle != "" {
			buffer.WriteString("  journal = {")
			buffer.WriteString(escapeSpecialChars(containerTitle))
			buffer.WriteString("},\n")
		}

		// 结束BibTeX条目
		buffer.WriteString("}\n")

		separator = "\n"
	}

	// 将生成的BibTeX写入响应
	_, err = w.Write([]byte(buffer.String()))
	if err != nil {
		s.logger.Error("write bibtex to response failed", "error", err.Error())
		return err
	}
	return nil
}

// GetDocMetaInfoMapByDocIds 根据文档ID列表获取文档元数据映射
func (s *CslService) GetDocMetaInfoMapByDocIds(ctx context.Context, reqs []*pb.DocMetaInfoHandlerReq) (map[string]*pb.DocMetaInfoSimpleVo, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "CslService.GetDocMetaInfoMapByDocIds")
	defer span.Finish()
	// 创建结果映射
	result := make(map[string]*pb.DocMetaInfoSimpleVo)
	// 遍历所有请求
	for _, req := range reqs {
		metaInfo, err := s.GetDocMetaInfo(ctx, req)
		if err != nil {
			s.logger.Error("get doc meta info failed", "error", err.Error())
			continue
		}
		if metaInfo != nil {
			// 将结果添加到映射中
			result[req.GetDocId()] = metaInfo
		}
	}
	return result, nil
}

// escapeSpecialChars 转义BibTeX中的特殊字符
func escapeSpecialChars(input string) string {
	// BibTeX中需要转义的特殊字符
	replacements := map[string]string{
		"&":  "\\&",
		"%":  "\\%",
		"$":  "\\$",
		"#":  "\\#",
		"_":  "\\_",
		"{":  "\\{",
		"}":  "\\}",
		"~":  "\\~{}",
		"^":  "\\^{}",
		"\\": "\\textbackslash{}",
	}

	result := input
	for char, replacement := range replacements {
		result = strings.ReplaceAll(result, char, replacement)
	}

	return result
}
