package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/doc"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/model"
)

// AbstractDocMetaInfoHandlerService 是所有文档元信息处理器的基础结构体，类似于 Java 中的抽象基类。
// 它通过“组合”的方式，为具体的处理器提供共享的依赖（如 logger）和通用的辅助方法。
type AbstractDocMetaInfoHandlerService struct {
	Logger             logging.Logger
	tracer             opentracing.Tracer
	doiMetaInfoService *DoiMetaInfoService
	Config             *config.Config
}

// NewAbstractDocMetaInfoHandlerService 创建一个基础处理器
func NewAbstractDocMetaInfoHandlerService(logger logging.Logger, tracer opentracing.Tracer, doiMetaInfoService *DoiMetaInfoService, cfg *config.Config) *AbstractDocMetaInfoHandlerService {
	return &AbstractDocMetaInfoHandlerService{
		Logger:             logger,
		tracer:             tracer,
		doiMetaInfoService: doiMetaInfoService,
		Config:             cfg,
	}
}

// ParseToMetaInfoVo 解析元信息Vo
func (s *AbstractDocMetaInfoHandlerService) ParseToMetaInfoVo(vo *pb.DocMetaInfoSimpleVo, metaInfoContent string) {
	//TODO
	s.Logger.Info("msg", "parse to meta info vo", "vo", vo, "metaInfoContent", metaInfoContent)
}

// SetMetaInfoFromDbData 设置元信息从数据库数据
func (s *AbstractDocMetaInfoHandlerService) SetMetaInfoFromDbData(vo *pb.DocMetaInfoSimpleVo, dbMetaInfo *model.DoiMetaInfo) {
	//TODO
	s.Logger.Info("msg", "set meta info from db data", "vo", vo, "dbMetaInfo", dbMetaInfo)
}

// UpdateDbMetaInfo 更新数据库元信息
// func (s *AbstractDocMetaInfoHandlerService) UpdateDbMetaInfo(ctx context.Context, vo *pb.DocMetaInfoSimpleVo, dbMetaInfo *model.DoiMetaInfo) {
// 	dbMetaInfo.IsParse = true
// 	if vo.Language != nil {
// 		dbMetaInfo.Language = *vo.Language
// 	}
// 	dbMetaInfo.DocType = vo.DocType
// 	dbMetaInfo.Title = vo.Title

// 	if len(vo.AuthorList) > 0 {
// 		authorListJSON, err := json.Marshal(vo.AuthorList)
// 		if err == nil {
// 			dbMetaInfo.AuthorList = string(authorListJSON)
// 		} else {
// 			s.Logger.Error("msg", "failed to marshal author list", "error", err)
// 		}
// 	}

// 	if vo.PublishTimestamp > 0 {
// 		dbMetaInfo.PublishTime = time.Unix(int64(vo.PublishTimestamp/1000), 0)
// 	}

// 	s.doiMetaInfoService.Update(ctx, dbMetaInfo)
// }

// CommonDataWrap 公共数据包装，处理一些通用的数据转换和填充
func (s *AbstractDocMetaInfoHandlerService) CommonDataWrap(vo *pb.DocMetaInfoSimpleVo) {
	// 1. 根据 DocType 设置 DocTypeName
	if vo.DocType != "" {
		docTypeInfoListJsonDesc := s.Config.DocMetaInfoSearch.DocTypeInfoListJsonDesc
		if docTypeInfoListJsonDesc != "" {
			// 定义一个临时的 struct 来解析 JSON
			type DocTypeInfo struct {
				Code string `json:"code"`
				Name string `json:"name"`
			}

			var docTypeInfos []DocTypeInfo
			if err := json.Unmarshal([]byte(docTypeInfoListJsonDesc), &docTypeInfos); err == nil {
				for _, docTypeInfo := range docTypeInfos {
					if strings.EqualFold(docTypeInfo.Code, vo.DocType) {
						vo.DocTypeName = &docTypeInfo.Name
						break
					}
				}
			} else {
				s.Logger.Error("msg", "failed to unmarshal docTypeInfoListJsonDesc", "error", err)
			}
		}
	}

	// 2. 如果有时间戳但没有格式化日期字符串，则进行格式化
	if vo.PublishTimestamp > 0 && (vo.PublishDateStr == nil || *vo.PublishDateStr == "") {
		// Java 的 new Date(timestamp) 接收的是毫秒，所以这里也按毫秒处理
		t := time.UnixMilli(int64(vo.PublishTimestamp))
		dateStr := t.Format("2006-01-02")
		vo.PublishDateStr = &dateStr
	}
}
