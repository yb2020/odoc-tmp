package service

import (
	"context"
	"net/url"
	"regexp"
	"strings"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/doc"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/services/doc/constant"
	"github.com/yb2020/odoc/services/doc/factory"
)

// DocCiteSearchService 文档引用搜索服务
type DocCiteSearchService struct {
	logger                           logging.Logger
	tracer                           opentracing.Tracer
	docMetaInfoHandlerServiceFactory *factory.DocMetaInfoHandlerServiceFactory
}

// NewDocCiteSearchService 创建新的文档引用搜索服务
func NewDocCiteSearchService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	docMetaInfoHandlerServiceFactory *factory.DocMetaInfoHandlerServiceFactory,
) *DocCiteSearchService {
	return &DocCiteSearchService{
		logger:                           logger,
		tracer:                           tracer,
		docMetaInfoHandlerServiceFactory: docMetaInfoHandlerServiceFactory,
	}
}

// EnDocCiteSearch 英文文档引用搜索
func (s *DocCiteSearchService) EnDocCiteSearch(ctx context.Context, req *pb.EnDocCiteSearchReq) (*pb.EnDocCiteSearchResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "DocCiteSearchService.EnDocCiteSearch")
	defer span.Finish()

	searchContent := req.GetSearchContent()
	// 如果启用了URL路径提取，则提取路径
	if constant.EnDocCiteSearchExtractPathSwitch {
		extractedPath := s.extractPath(searchContent)
		if extractedPath != "" {
			searchContent = extractedPath
		}
	}
	// 判断是DOI还是标题
	isDoi := s.validateDoi(searchContent)
	var result []*pb.DocMetaInfoSimpleVo
	// 这部分逻辑暂时留空，如用户要求
	if isDoi {
		// 使用DOI搜索
		s.logger.Info("searching by DOI", "doi", searchContent)
		result = s.EnDocCiteSearchByDoi(ctx, req)
	} else {
		// 使用标题搜索
		s.logger.Info("searching by title", "title", searchContent)
		result = s.EnDocCiteSearchByTitle(ctx, req)
	}

	return &pb.EnDocCiteSearchResponse{
		Result: result,
	}, nil
}

// extractPath 从URL中提取路径
func (s *DocCiteSearchService) extractPath(input string) string {
	// 判断输入是否是URL
	if s.isURL(input) {
		parsedURL, err := url.Parse(input)
		if err != nil {
			s.logger.Error("failed to parse URL", "error", err.Error(), "input", input)
			return input
		}

		path := parsedURL.Path
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}
		return path
	}
	// 如果不是URL，直接返回输入
	return input
}

// isURL 判断字符串是否是URL
func (s *DocCiteSearchService) isURL(str string) bool {
	// 简单的URL验证
	return strings.HasPrefix(str, "http://") || strings.HasPrefix(str, "https://")
}

// validateDoi 验证DOI
func (s *DocCiteSearchService) validateDoi(doi string) bool {
	// 正则表达式，匹配以 "10." 开头，后面至少有一个字符，接着是 "/"，然后至少有一个字符
	regex := `^10\..+/.+$`
	match, err := regexp.MatchString(regex, doi)
	if err != nil {
		s.logger.Error("DOI validation regex error", "error", err.Error())
		return false
	}
	return match
}

// zhDocCiteSearch 通过DOI搜索文档引用
func (s *DocCiteSearchService) ZhDocCiteSearch(ctx context.Context, req *pb.ZhDocCiteSearchReq) (*pb.ZhDocCiteSearchResponse, error) {
	// todo: 这部分逻辑暂时留空，将在后续实现  因为这里有一个工厂模式  需要先实现
	s.logger.Info("DOI search implementation pending", "doi", req)
	return &pb.ZhDocCiteSearchResponse{
		Result: []*pb.DocMetaInfoSimpleVo{},
	}, nil
}

// enDocCiteSearchByDoi 通过DOI搜索文档引用
func (s *DocCiteSearchService) EnDocCiteSearchByDoi(ctx context.Context, req *pb.EnDocCiteSearchReq) []*pb.DocMetaInfoSimpleVo {
	// todo: 这部分逻辑暂时留空，将在后续实现  因为这里有一个工厂模式  需要先实现
	s.logger.Info("DOI search implementation pending", "doi", req)
	doiMetaInfoHandlerService, err := s.docMetaInfoHandlerServiceFactory.GetProvider(constant.DOI)
	if err != nil {
		s.logger.Error("failed to get provider", "error", err.Error())
		return nil
	}
	s.logger.Info("get provider", "provider", doiMetaInfoHandlerService)

	docMetaInfoHandlerReq := &pb.DocMetaInfoHandlerReq{
		Doi: &req.SearchContent,
	}
	result, err := doiMetaInfoHandlerService.GetDocMetaInfo(ctx, docMetaInfoHandlerReq)
	if err != nil {
		s.logger.Error("failed to get doc meta info", "error", err.Error())
		return nil
	}
	return []*pb.DocMetaInfoSimpleVo{result}
}

// enDocCiteSearchByTitle 通过标题搜索文档引用
func (s *DocCiteSearchService) EnDocCiteSearchByTitle(ctx context.Context, req *pb.EnDocCiteSearchReq) []*pb.DocMetaInfoSimpleVo {
	// todo: 这部分逻辑暂时留空，将在后续实现 因为这里有一个工厂模式  需要先实现
	s.logger.Info("Title search implementation pending", "title", req)
	doiMetaInfoHandlerService, err := s.docMetaInfoHandlerServiceFactory.GetProvider(constant.DOI)
	if err != nil {
		s.logger.Error("failed to get provider", "error", err.Error())
		return nil
	}
	s.logger.Info("get provider", "provider", doiMetaInfoHandlerService)

	result, err := doiMetaInfoHandlerService.GetDocMetaInfosByTitle(ctx, req.SearchContent)
	if err != nil {
		s.logger.Error("failed to get doc meta infos by title", "error", err.Error())
		return nil
	}
	return result
}
