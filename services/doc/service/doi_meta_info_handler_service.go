package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/utils"
	pb "github.com/yb2020/odoc/proto/gen/go/doc"
	"github.com/yb2020/odoc/services/doc/constant"
	"github.com/yb2020/odoc/services/doc/model"
)

// DoiMetaInfoHandlerService 负责处理 DOI 类型的元信息
type DoiMetaInfoHandlerService struct {
	logger                            logging.Logger
	tracer                            opentracing.Tracer
	AbstractDocMetaInfoHandlerService // 嵌入基础结构体，自动获得 Logger 等字段和方法
	doiMetaInfoService                *DoiMetaInfoService
	cfg                               *config.Config
	cache                             cache.Cache
}

// NewDoiMetaInfoHandlerService 创建一个新的 DOI 处理器
func NewDoiMetaInfoHandlerService(logger logging.Logger, tracer opentracing.Tracer, cache cache.Cache, doiMetaInfoService *DoiMetaInfoService, cfg *config.Config) *DoiMetaInfoHandlerService {
	return &DoiMetaInfoHandlerService{
		logger:                            logger,
		tracer:                            tracer,
		cache:                             cache,
		AbstractDocMetaInfoHandlerService: *NewAbstractDocMetaInfoHandlerService(logger, tracer, doiMetaInfoService, cfg),
		doiMetaInfoService:                doiMetaInfoService,
		cfg:                               cfg,
	}
}

// GetName 返回处理器的唯一名称
func (s *DoiMetaInfoHandlerService) GetName() string {
	return constant.DOI // 这是此处理器的唯一标识
}

// GetDocMetaInfo 获取文档元信息
func (s *DoiMetaInfoHandlerService) GetDocMetaInfo(ctx context.Context, req *pb.DocMetaInfoHandlerReq) (*pb.DocMetaInfoSimpleVo, error) {

	//TODO
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoHandlerService.GetDocMetaInfo")
	defer span.Finish()
	s.logger.Info("msg", "get doc meta info", "req", req)

	doiMetaInfoExt, err := s.doiMetaInfoService.GetByDoi(ctx, *req.Doi)
	if err != nil {
		return nil, err
	}

	if doiMetaInfoExt == nil || !doiMetaInfoExt.IsCrawlMeta {
		// TODO
		doiMetaInfoExt, err = s.PersistentMetaInfo(ctx, req)
		if err != nil {
			return nil, err
		}
		if doiMetaInfoExt == nil {
			return nil, nil
		}
	}

	doiMetaInfo, err := s.doiMetaInfoService.GetById(ctx, doiMetaInfoExt.Id)
	if err != nil {
		return nil, err
	}

	// doiMetaInfo 不为空
	vo := s.GetEmptyDocMetaInfo()
	vo.Doi = &doiMetaInfo.Doi
	s.ParseToMetaInfoVo(ctx, vo, doiMetaInfo.Content)
	s.logger.Info("msg", "=== 准备调用 UpdateDbMetaInfo ===", "vo", vo)
	err = s.UpdateDbMetaInfo(ctx, vo, doiMetaInfo)
	if err != nil {
		s.logger.Error("msg", "UpdateDbMetaInfo 调用失败", "error", err)
		return nil, err
	}
	s.logger.Info("msg", "=== GetDocMetaInfo 执行完成 ===")

	return vo, nil
}

// GetDocMetaInfosByTitle 获取文档元信息
func (s *DoiMetaInfoHandlerService) GetDocMetaInfosByTitle(ctx context.Context, title string) ([]*pb.DocMetaInfoSimpleVo, error) {
	//TODO
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoHandlerService.GetDocMetaInfosByTitle")
	defer span.Finish()
	s.logger.Info("msg", "get doc meta infos by title", "title", title)

	if title == "" {
		return nil, errors.New("title is empty")
	}

	titleMd5 := utils.MD5(title)
	cacheKey := fmt.Sprintf("%s%s", constant.DocSearchResultByTitleMd5KeyPrefix, titleMd5)
	s.logger.Info("msg", "get doc meta infos by title", "title", title, "titleMd5", titleMd5)

	var results []*pb.DocMetaInfoSimpleVo
	//从redis中获取
	exists, err := s.cache.GetNotBizPrefix(ctx, cacheKey, &results)
	if err != nil {
		s.logger.Error("msg", "get doc meta infos by title", "title", title, "titleMd5", titleMd5, "error", err)
		return nil, err
	}
	if exists {
		s.logger.Info("msg", "get doc meta infos by title", "title", title, "titleMd5", titleMd5, "results", results)
		return results, nil
	}

	//1. 通过Squid代理Http请求从doiapi获取元信息列表，
	// apiURL := s.cfg.DocMetaInfoSearch.Doi.QueryPaperDoiInfoUrl + url.QueryEscape(title) + "&rows=10"

	// // 2. 通过代理获取元数据
	// body, err := s.squidProxyService.ProxyHttp(ctx, apiURL)
	// if err != nil {
	// 	s.logger.Error("msg", "通过代理获取DOI元信息失败", "doi", title, "error", err)
	// 	return nil, err
	// }
	// s.logger.Info("msg", "通过代理获取DOI元信息成功", "doi", string(body))

	// // 3. 判断body是否为空并且是否为json格式
	// if len(body) == 0 {
	// 	s.logger.Warn("msg", "通过代理获取的DOI元信息为空", "doi", title)
	// 	return nil, errors.New("empty response from CrossRef API")
	// }

	// if !json.Valid(body) {
	// 	s.logger.Error("msg", "通过代理获取的DOI元信息不是有效的JSON格式", "doi", title, "body", string(body))
	// 	return nil, errors.New("invalid JSON response from CrossRef API")
	// }

	// // 4. 解析json成数组List
	// // 定义用于解析列表响应的结构体
	// type CrossRefListResponse struct {
	// 	Message struct {
	// 		Items []json.RawMessage `json:"items"`
	// 	} `json:"message"`
	// }

	// var listResp CrossRefListResponse
	// if err := json.Unmarshal(body, &listResp); err != nil {
	// 	s.logger.Error("msg", "解析DOI列表元信息失败", "error", err, "body", string(body))
	// 	return nil, err
	// }

	// itemsInfo := listResp.Message.Items
	// s.logger.Info("msg", "成功提取items列表", "count", len(itemsInfo))

	// for _, item := range itemsInfo {
	// 	// ParseToMetaInfoVo 需要一个完整的 CrossRefResponse 格式的 JSON，
	// 	// 而列表中的 item 只是 message 的内容，所以我们手动包装一下。
	// 	wrappedJSON := fmt.Sprintf(`{"status":"ok","message":%s}`, string(item))

	// 	vo := s.GetEmptyDocMetaInfo()
	// 	s.ParseToMetaInfoVo(ctx, vo, wrappedJSON)
	// 	results = append(results, vo)
	// }
	// //判断results不是空的数组，保存到redis中
	// if len(results) > 0 {
	// 	s.cache.SetNotBizPrefix(ctx, cacheKey, results, 3*time.Hour) // 3小时过期
	// }

	return results, nil
}

// InitMetaFieldFromContent 初始化元信息字段
func (s *DoiMetaInfoHandlerService) InitMetaFieldFromContent(fieldName string) {
	//TODO
	s.logger.Info("msg", "init meta field from content", "fieldName", fieldName)
}

// GetEmptyDocMetaInfo 获取空的文档元信息
func (s *DoiMetaInfoHandlerService) GetEmptyDocMetaInfo() *pb.DocMetaInfoSimpleVo {
	lang := constant.ENGLISH
	return &pb.DocMetaInfoSimpleVo{
		Language: &lang,
	}
}

// PersistentMetaInfo 持久化元信息
func (s *DoiMetaInfoHandlerService) PersistentMetaInfo(ctx context.Context, req *pb.DocMetaInfoHandlerReq) (*model.DoiMetaInfo, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoHandlerService.PersistentMetaInfo")
	defer span.Finish()
	s.logger.Info("msg", "persistent meta info", "req", req)
	doiMetaInfo, err := s.doiMetaInfoService.GetByDoi(ctx, *req.Doi)
	if err != nil {
		return nil, err
	}

	if doiMetaInfo != nil && doiMetaInfo.IsCrawlMeta {
		return doiMetaInfo, nil
	}

	// // 1. 从配置构建 API URL
	// apiURL := s.cfg.DocMetaInfoSearch.Doi.QueryDoiInfoUrl + *req.Doi

	// // 2. 通过代理获取元数据
	// body, err := s.squidProxyService.ProxyHttp(ctx, apiURL)
	// if err != nil {
	// 	s.logger.Error("msg", "通过代理获取DOI元信息失败", "doi", *req.Doi, "error", err)
	// 	return nil, err
	// }

	// s.logger.Info("msg", "通过代理获取DOI元信息成功", "doi", string(body))

	// // 3. 判断body是否为空并且是否为json格式
	// if len(body) == 0 {
	// 	s.logger.Warn("msg", "通过代理获取的DOI元信息为空", "doi", *req.Doi)
	// 	return nil, errors.New("empty response from CrossRef API")
	// }

	// if !json.Valid(body) {
	// 	s.logger.Error("msg", "通过代理获取的DOI元信息不是有效的JSON格式", "doi", *req.Doi, "body", string(body))
	// 	return nil, errors.New("invalid JSON response from CrossRef API")
	// }
	// // 4.更新或者保存到数据库
	// if doiMetaInfo != nil {
	// 	// 保持原有记录的ID和其他字段，只更新Content和IsCrawlMeta
	// 	doiMetaInfo.Content = string(body)
	// 	doiMetaInfo.IsCrawlMeta = true
	// 	if err := s.doiMetaInfoService.Update(ctx, doiMetaInfo); err != nil {
	// 		return nil, err
	// 	}
	// } else {
	// 	doiMetaInfo = &model.DoiMetaInfo{
	// 		Doi:         *req.Doi,
	// 		IsCrawlMeta: true,
	// 	}
	// 	doiMetaInfo.Content = string(body)
	// 	createdId, err := s.doiMetaInfoService.CreateDoiMetaInfo(ctx, doiMetaInfo)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	doiMetaInfo.Id = createdId // 设置正确的ID
	// }

	return doiMetaInfo, nil
}

// ParseToMetaInfoVo 解析元信息, metaInfoContent是json格式的元信息
func (s *DoiMetaInfoHandlerService) ParseToMetaInfoVo(ctx context.Context, vo *pb.DocMetaInfoSimpleVo, metaInfoContent string) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DoiMetaInfoHandlerService.ParseToMetaInfoVo")
	defer span.Finish()
	if metaInfoContent == "" {
		return
	}

	// 1. 定义用于解析 JSON 的结构体
	type CrossRefAuthor struct {
		Given  string `json:"given"`
		Family string `json:"family"`
	}
	type CrossRefDate struct {
		DateParts [][]int `json:"date-parts"`
	}
	type CrossRefEvent struct {
		Name     string       `json:"name"`
		Location string       `json:"location"`
		Start    CrossRefDate `json:"start"`
		End      CrossRefDate `json:"end"`
	}
	type CrossRefCreated struct {
		Timestamp int64 `json:"timestamp"`
	}
	type CrossRefMessage struct {
		Title          []string         `json:"title"`
		Author         []CrossRefAuthor `json:"author"`
		PublishedPrint CrossRefDate     `json:"published-print"`
		Created        CrossRefCreated  `json:"created"`
		Publisher      string           `json:"publisher"`
		Volume         string           `json:"volume"`
		Issued         CrossRefDate     `json:"issued"`
		Page           string           `json:"page"`
		ContainerTitle []string         `json:"container-title"`
		Type           string           `json:"type"`
		DOI            string           `json:"DOI"`
		Event          CrossRefEvent    `json:"event"`
		Language       string           `json:"language"`
	}
	type CrossRefResponse struct {
		Status  string          `json:"status"`
		Message CrossRefMessage `json:"message"`
	}

	// 2. 解析 JSON 响应
	var crossRefResp CrossRefResponse
	if err := json.Unmarshal([]byte(metaInfoContent), &crossRefResp); err != nil {
		s.logger.Error("msg", "解析元信息内容失败", "error", err, "metaInfoContent", metaInfoContent)
		return
	}

	if crossRefResp.Status != "ok" {
		s.logger.Warn("msg", "元信息内容状态不是ok", "status", crossRefResp.Status)
		return
	}

	msg := crossRefResp.Message

	// 3. 将解析的数据映射到 DocMetaInfoSimpleVo
	vo.Doi = &msg.DOI
	if transformedType, ok := s.cfg.DocMetaInfoSearch.TransformDoiDocTypeMapRel[msg.Type]; ok {
		vo.DocType = transformedType
	} else {
		vo.DocType = msg.Type
	}
	if len(msg.Title) > 0 {
		vo.Title = msg.Title[0]
	}
	vo.ContainerTitle = msg.ContainerTitle
	vo.Volume = &msg.Volume
	if len(msg.Issued.DateParts) > 0 && len(msg.Issued.DateParts[0]) > 0 {
		parts := msg.Issued.DateParts[0]
		var dateStr string
		if len(parts) == 1 {
			dateStr = fmt.Sprintf("%d", parts[0])
		} else if len(parts) == 2 {
			dateStr = fmt.Sprintf("%d-%02d", parts[0], parts[1])
		} else if len(parts) >= 3 {
			dateStr = fmt.Sprintf("%d-%02d-%02d", parts[0], parts[1], parts[2])
		}
		if dateStr != "" {
			vo.Issue = &dateStr
		}
	}
	vo.Page = &msg.Page

	// 解析发布日期
	if len(msg.PublishedPrint.DateParts) > 0 && len(msg.PublishedPrint.DateParts[0]) > 0 {
		parts := msg.PublishedPrint.DateParts[0]
		var dateStr string
		var t time.Time
		var err error

		if len(parts) == 1 {
			dateStr = fmt.Sprintf("%d", parts[0])
			t, err = time.Parse("2006", dateStr)
		} else if len(parts) == 2 {
			dateStr = fmt.Sprintf("%d-%02d", parts[0], parts[1])
			t, err = time.Parse("2006-01", dateStr)
		} else if len(parts) >= 3 {
			dateStr = fmt.Sprintf("%d-%02d-%02d", parts[0], parts[1], parts[2])
			t, err = time.Parse("2006-01-02", dateStr)
		}

		if err == nil {
			vo.PublishTimestamp = uint64(t.UnixMilli())
		}
		if dateStr != "" {
			vo.PublishDateStr = &dateStr
		}
	} else if msg.Created.Timestamp > 0 {
		vo.PublishTimestamp = uint64(msg.Created.Timestamp)
	}

	// 解析language
	if msg.Language != "" {
		if transformedLang, ok := s.cfg.DocMetaInfoSearch.TransformDoiLanguageMapRel[msg.Language]; ok {
			vo.Language = &transformedLang
		} else {
			lang := msg.Language
			vo.Language = &lang
		}
	}

	var authorList []*pb.AuthorInfo
	// 匹配中文字符的正则表达式
	chinesePattern := regexp.MustCompile(`[\p{Han}]`)

	for _, a := range msg.Author {
		given := a.Given
		family := a.Family

		isChinese := chinesePattern.MatchString(given) || chinesePattern.MatchString(family)

		var literal string
		if isChinese {
			literal = family + given
		} else {
			literal = given
			if family != "" {
				if literal != "" {
					literal += " "
				}
				literal += family
			}
		}

		authorInfo := &pb.AuthorInfo{
			Given:   &given,
			Family:  &family,
			Literal: literal,
		}
		authorList = append(authorList, authorInfo)
	}
	vo.AuthorList = authorList

	// 4. 解析 Event Info
	if msg.Event.Name != "" {
		vo.EventTitle = &msg.Event.Name
	}
	if msg.Event.Location != "" {
		vo.EventPlace = &msg.Event.Location
	}

	var eventDates []string
	// 格式化日期的辅助函数
	formatDate := func(dateParts [][]int) string {
		if len(dateParts) > 0 && len(dateParts[0]) > 0 {
			parts := dateParts[0]
			if len(parts) == 1 {
				return fmt.Sprintf("%d", parts[0])
			} else if len(parts) == 2 {
				return fmt.Sprintf("%d-%02d", parts[0], parts[1])
			} else if len(parts) >= 3 {
				return fmt.Sprintf("%d-%02d-%02d", parts[0], parts[1], parts[2])
			}
		}
		return ""
	}

	startDate := formatDate(msg.Event.Start.DateParts)
	if startDate != "" {
		eventDates = append(eventDates, startDate)
	}

	endDate := formatDate(msg.Event.End.DateParts)
	if endDate != "" {
		eventDates = append(eventDates, endDate)
	}

	if len(eventDates) > 0 {
		vo.EventDate = eventDates
	}

	// 设置通用字段
	s.CommonDataWrap(vo)
}

// UpdateDbMetaInfo 更新数据库中的元信息
func (s *DoiMetaInfoHandlerService) UpdateDbMetaInfo(ctx context.Context, vo *pb.DocMetaInfoSimpleVo, dbMetaInfo *model.DoiMetaInfo) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(context.Background(), s.tracer, "DoiMetaInfoHandlerService.UpdateDbMetaInfo")
	defer span.Finish()

	// 设置通用字段（原 AbstractDocMetaInfoHandlerService.UpdateDbMetaInfo 的逻辑）
	dbMetaInfo.IsParse = true
	if vo.Language != nil {
		dbMetaInfo.Language = *vo.Language
	}
	dbMetaInfo.DocType = vo.DocType
	dbMetaInfo.Title = vo.Title

	if len(vo.AuthorList) > 0 {
		authorListJSON, err := json.Marshal(vo.AuthorList)
		if err == nil {
			dbMetaInfo.AuthorList = string(authorListJSON)
		} else {
			s.logger.Error("msg", "failed to marshal author list", "error", err)
		}
	}

	if vo.PublishTimestamp > 0 {
		dbMetaInfo.PublishTime = time.Unix(int64(vo.PublishTimestamp/1000), 0)
	}

	// 设置 DoiMetaInfo 特有字段
	if vo.Page != nil {
		dbMetaInfo.Page = *vo.Page
	}
	if vo.Volume != nil {
		dbMetaInfo.Volume = *vo.Volume
	}
	if vo.Issue != nil {
		dbMetaInfo.Issue = *vo.Issue
	}
	if vo.Url != nil {
		dbMetaInfo.Url = *vo.Url
	}
	if vo.EventPlace != nil {
		dbMetaInfo.EventPlace = *vo.EventPlace
	}
	if vo.EventTitle != nil {
		dbMetaInfo.EventTitle = *vo.EventTitle
	}

	if len(vo.EventDate) > 0 {
		eventDateJSON, err := json.Marshal(vo.EventDate)
		if err == nil {
			dbMetaInfo.EventDate = string(eventDateJSON)
		}
	}

	if len(vo.ContainerTitle) > 0 {
		containerTitleJSON, err := json.Marshal(vo.ContainerTitle)
		if err == nil {
			dbMetaInfo.ContainerTitle = string(containerTitleJSON)
		}
	}

	//s.logger.Info("msg", "准备更新数据库", "dbMetaInfo", dbMetaInfo, "id", dbMetaInfo.Id, "title", dbMetaInfo.Title, "page", dbMetaInfo.Page)
	s.logger.Info("msg", "更新数据库", "id", dbMetaInfo.Id)
	err := s.doiMetaInfoService.Update(ctx, dbMetaInfo)
	if err != nil {
		s.logger.Error("msg", "数据库更新失败", "error", err)
	} else {
		s.logger.Info("msg", "数据库更新成功", "id", dbMetaInfo.Id)
	}
	return err
}
