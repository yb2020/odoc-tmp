package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	commonPb "github.com/yb2020/odoc-proto/gen/go/common"
	docpb "github.com/yb2020/odoc-proto/gen/go/doc"
	parsedPb "github.com/yb2020/odoc-proto/gen/go/parsed"
	pb "github.com/yb2020/odoc-proto/gen/go/pdf"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/cache"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/mq/rocketmq"
	"github.com/yb2020/odoc/pkg/mq/rocketmq/producer"
	docConstant "github.com/yb2020/odoc/services/doc/constant"
	userDocService "github.com/yb2020/odoc/services/doc/service"
	ossConstant "github.com/yb2020/odoc/services/oss/constant"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperModel "github.com/yb2020/odoc/services/paper/model"
	paperService "github.com/yb2020/odoc/services/paper/service"
	parseConstant "github.com/yb2020/odoc/services/parse/constant"
	util "github.com/yb2020/odoc/services/parse/util"
)

// PdfParseService PDF解析服务
type PdfParseService struct {
	tracer                opentracing.Tracer
	logger                logging.Logger
	paperPdfService       *PaperPdfService
	paperPdfParsedService *paperService.PaperPdfParsedService
	ossService            ossService.OssServiceInterface
	userDocService        *userDocService.UserDocService
	producer              *producer.RocketMQProducer
	config                *config.Config
	cache                 cache.Cache
}

// NewPdfParseService 创建新的PDF解析服务
func NewPdfParseService(
	paperPdfService *PaperPdfService,
	paperPdfParsedService *paperService.PaperPdfParsedService,
	ossService ossService.OssServiceInterface,
	userDocService *userDocService.UserDocService,
	producer *producer.RocketMQProducer,
	cache cache.Cache,
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
) *PdfParseService {
	return &PdfParseService{
		paperPdfService:       paperPdfService,
		paperPdfParsedService: paperPdfParsedService,
		ossService:            ossService,
		userDocService:        userDocService,
		producer:              producer,
		cache:                 cache,
		config:                config,
		logger:                logger,
		tracer:                tracer,
	}
}

// GetReferenceMarkers 获取参考文献标记
func (s *PdfParseService) GetReferenceMarkers(ctx context.Context, req *pb.GetReferenceMarkersRequest) (*pb.GetReferenceMarkersResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.GetReferenceMarkers")
	defer span.Finish()

	//从minio中获取PDF元数据
	metadata, err := s.GetPdfMetadata(ctx, req.PdfId)
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		return &pb.GetReferenceMarkersResponse{
			NeedFetch: true,
		}, nil
	}
	//转换对象
	figureAndTableList := s.convertToFigureAndTableMarker(ctx, metadata.FigureAndTableMarkers)
	referenceMarkers := s.convertToReferenceMarker(ctx, metadata.ReferenceMarkers, metadata.References)

	return &pb.GetReferenceMarkersResponse{
		NeedFetch:             false,
		FigureAndTableMarkers: figureAndTableList,
		Markers:               referenceMarkers,
	}, nil
}

// GetFiguresAndTables 获取图表信息   手动分页
func (s *PdfParseService) GetFiguresAndTables(ctx context.Context, req *pb.GetFiguresAndTablesListRequest) (*pb.GetFiguresAndTablesListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.GetFiguresAndTables")
	defer span.Finish()

	//获取PDF元数据
	metadata, err := s.GetPdfMetadata(ctx, req.PdfId)
	if err != nil {
		return nil, err
	}
	if metadata == nil {
		return &pb.GetFiguresAndTablesListResponse{
			NeedFetch: true,
		}, nil
	}
	//获取oss上传记录
	imageRecords, err := s.GetFigureAndTable(ctx, req.PdfId)
	if err != nil {
		return nil, err
	}
	//获取图表信息
	figuresAndTables := metadata.FiguresAndTables

	// 处理分页参数
	pageNum := int(req.PageReq.PageNum)
	pageSize := int(req.PageReq.PageSize)

	// 默认值处理
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // 默认每页10条
	}

	// 计算总数量和分页信息
	totalCount := len(figuresAndTables)
	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize

	// 边界检查
	if startIndex >= totalCount {
		// 如果起始索引超出范围，返回空数组
		return &pb.GetFiguresAndTablesListResponse{
			NeedFetch:          false,
			FigureAndTableList: []*pb.PdfFigureAndTableInfo{},
			PageResp: &commonPb.PageResponse{
				PageNum:     int32(pageNum),
				Total:       int32(totalCount),
				IsFirstPage: pageNum == 1,
				IsLastPage:  true,
			},
		}, nil
	}

	if endIndex > totalCount {
		endIndex = totalCount
	}

	// 获取当前页的图表
	pageFiguresAndTables := figuresAndTables[startIndex:endIndex]
	//转换对象
	figureAndTableList := s.convertToFigureAndTableInfo(ctx, pageFiguresAndTables, imageRecords)

	response := pb.GetFiguresAndTablesListResponse{
		NeedFetch:          false,
		FigureAndTableList: figureAndTableList,
		PageResp: &commonPb.PageResponse{
			PageNum:     int32(pageNum),
			Total:       int32(totalCount),
			IsFirstPage: pageNum == 1,
			IsLastPage:  endIndex >= totalCount,
		},
	}
	return &response, nil
}

// GetReference	获取引用  手动分页
func (s *PdfParseService) GetReference(ctx context.Context, req *pb.GetReferenceRequest) (*pb.GetReferenceResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.GetReference")
	defer span.Finish()

	//获取PDF元数据
	metadata, err := s.GetPdfMetadata(ctx, req.PdfId)
	if err != nil {
		return nil, errors.Biz("pdf.pdf.errors.deserialize_failed")
	}
	if metadata == nil {
		return &pb.GetReferenceResponse{
			NeedFetch: true,
		}, nil
	}

	// 如果引用信息不存在，转换为PdfCatalogueInfo数组
	if metadata.References == nil {
		return &pb.GetReferenceResponse{
			NeedFetch: false,
		}, nil
	}
	// 处理分页参数
	pageNum := int(req.PageReq.PageNum)
	pageSize := int(req.PageReq.PageSize)

	// 默认值处理
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10 // 默认每页10条
	}

	// 计算总数量和分页信息
	totalCount := len(metadata.References)
	startIndex := (pageNum - 1) * pageSize
	endIndex := startIndex + pageSize

	// 边界检查
	if startIndex >= totalCount {
		// 如果起始索引超出范围，返回空数组
		return &pb.GetReferenceResponse{
			NeedFetch:         false,
			ReferenceInfoList: []*pb.ReferenceInfo{},
			ReferenceCount:    int32(totalCount),
			PageResp: &commonPb.PageResponse{
				PageNum:     int32(pageNum),
				Total:       int32(totalCount),
				IsFirstPage: pageNum == 1,
				IsLastPage:  true,
			},
		}, nil
	}

	if endIndex > totalCount {
		endIndex = totalCount
	}

	// 获取当前页的引用
	pageReferences := metadata.References[startIndex:endIndex]

	references := []*pb.ReferenceInfo{}
	for _, reference := range pageReferences {
		refer := &pb.ReferenceInfo{
			RefIdx:    reference.RefIdx,
			SearchKey: reference.Title,
			Title:     reference.ContentText,
			PageNum:   reference.Bbox.PageNumber,
		}
		//Bbox 转换
		if reference.Bbox != nil {
			refer.Bbox = &pb.PdfBBox{
				X0:           float32(reference.Bbox.X0),
				Y0:           float32(reference.Bbox.Y0),
				X1:           float32(reference.Bbox.X1),
				Y1:           float32(reference.Bbox.Y1),
				OriginHeight: float32(reference.Bbox.OriginHeight),
				OriginWidth:  float32(reference.Bbox.OriginWidth),
			}
			//如果页面的pageNum为空或者为0  则使用bbox中的pageNum
			if refer.PageNum == 0 {
				refer.PageNum = int32(reference.Bbox.PageNumber)
			}
		}

		if reference.Authors != nil {
			//作者
			authors := []*pb.Author{}
			for _, author := range reference.Authors {
				authors = append(authors, &pb.Author{
					Name: author.FullName,
				})
			}
			refer.PaperMetaExtend = &pb.PaperMetaExtend{
				AuthorList: authors,
			}
		}
		references = append(references, refer)
	}

	//TODO 这里的结构对不上
	return &pb.GetReferenceResponse{
		NeedFetch:         false,
		ReferenceInfoList: references,
		ReferenceCount:    int32(totalCount),
		PageResp: &commonPb.PageResponse{
			PageNum:     int32(pageNum),
			Total:       int32(totalCount),
			IsFirstPage: pageNum == 1,
			IsLastPage:  endIndex >= totalCount,
		},
	}, nil
}

// GetCatalogue 获取PDF目录
func (s *PdfParseService) GetCatalogue(ctx context.Context, req *pb.GetCatalogueRequest) (*pb.GetCatalogueResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.GetCatalogue")
	defer span.Finish()

	metadata, err := s.GetPdfMetadata(ctx, req.PdfId)
	if err != nil {
		return nil, errors.Biz("pdf.pdf.errors.deserialize_failed")
	}
	catalogueResponse := &pb.GetCatalogueResponse{
		NeedFetch: false,
	}
	if metadata == nil {
		// 需要发送mq消息，告诉parse服务解析pdf，因为这里可能没有解析数据
		needErr := s.needParsedPDF(ctx, req.PdfId)
		if needErr != nil {
			return catalogueResponse, nil
		}
		catalogueResponse.NeedFetch = true
		return catalogueResponse, nil
	}
	catalogueItems := metadata.Catalogue
	if len(catalogueItems) == 0 {
		return catalogueResponse, nil
	}

	// 定义转换函数
	var convertCatalogueItem func(item *parsedPb.CatalogueItem) *pb.PdfCatalogueInfo
	convertCatalogueItem = func(item *parsedPb.CatalogueItem) *pb.PdfCatalogueInfo {
		if item == nil {
			return nil
		}

		info := &pb.PdfCatalogueInfo{
			Title:   item.Title,
			PageNum: item.Bbox.PageNumber,
		}

		// 转换边界框
		if item.Bbox != nil {
			info.Bbox = &pb.PdfBBox{
				X0:           float32(item.Bbox.X0),
				Y0:           float32(item.Bbox.Y0),
				X1:           float32(item.Bbox.X1),
				Y1:           float32(item.Bbox.Y1),
				OriginHeight: float32(item.Bbox.OriginHeight),
				OriginWidth:  float32(item.Bbox.OriginWidth),
			}
			//如果页面的pageNum为空或者为0  则使用bbox中的pageNum
			if info.PageNum == 0 {
				info.PageNum = int32(item.Bbox.PageNumber)
			}
		}
		// 递归处理子目录
		if len(item.Child) > 0 {
			info.Child = make([]*pb.PdfCatalogueInfo, 0, len(item.Child))
			for _, child := range item.Child {
				if childInfo := convertCatalogueItem(child); childInfo != nil {
					info.Child = append(info.Child, childInfo)
				}
			}
		}
		return info
	}
	// 转换所有目录项
	catalogueInfoItems := make([]*pb.PdfCatalogueInfo, 0, len(catalogueItems))
	for _, item := range catalogueItems {
		if info := convertCatalogueItem(item); info != nil {
			catalogueInfoItems = append(catalogueInfoItems, info)
		}
	}
	//构建响应
	catalogueResponse.PdfCatalogue = &pb.PdfCatalogueInfo{
		Child: catalogueInfoItems,
	}
	return catalogueResponse, nil
}

// ReParse 重新解析
func (s *PdfParseService) ReParse(ctx context.Context, req *pb.GetCatalogueRequest) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.ReParse")
	defer span.Finish()

	//验证状态  看解析执行到哪一步了，来决定往哪个topic发送消息
	//根据pdf查询paperPdf记录
	paperPdf, err := s.paperPdfService.GetById(ctx, req.PdfId)
	if err != nil {
		return err
	}
	if paperPdf == nil {
		return errors.Biz("pdf.pdf.errors.not_found")
	}
	fileSHA256 := paperPdf.FileSHA256
	//
	cacheKey := fmt.Sprintf("%s%s", docConstant.ParsePDFStatusKeyPrefix, fileSHA256)
	s.logger.Info("msg", "prepare to parse pdf text", "file sha256", paperPdf.FileSHA256)
	var cacheValue docpb.UserDocParsedStatusEnum
	exists, err := s.cache.GetNotBizPrefix(ctx, cacheKey, &cacheValue)
	if err != nil {
		s.logger.Error("mq.ParseUploadPdfService error, failed to get redis cache!", "error", err)
		return errors.Biz("get parse text pdf redis cache failed")
	}
	if !exists {
		//如果不存在则进行查库
		userDoc, err := s.userDocService.GetByPdId(ctx, req.PdfId)
		if err != nil {
			return err
		}
		if userDoc == nil {
			return errors.Biz("pdf.pdf.errors.not_found")
		}
		cacheValue = docpb.UserDocParsedStatusEnum(userDoc.ParseStatus)
	}
	//判断状态是否到了解析文本数据完成
	if cacheValue == docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSED {
		//发送mq消息，告诉parse服务解析pdf，因为这里可能没有解析数据
		s.needEmbedding(ctx, req.PdfId)
		return nil
	}
	// 需要发送mq消息，告诉parse服务解析pdf，因为这里可能没有解析数据
	needErr := s.needParsedPDF(ctx, req.PdfId)
	if needErr != nil {
		return needErr
	}
	return nil
}

// GetPdfMetadata 获取PDF元数据
func (s *PdfParseService) GetPdfMetadata(ctx context.Context, pdfId string) (*parsedPb.DocumentMetadata, error) {
	paperPdfParsed, err := s.getPdfSource(ctx, pdfId, parseConstant.PdfOssTypeMetadata)
	if err != nil {
		return nil, err
	}
	if paperPdfParsed == nil || paperPdfParsed.ObjectKey == "" {
		return nil, nil
	}
	//下载oss文件并返回byte[]
	bucketType := ossConstant.BucketTypeToEnum(s.config, paperPdfParsed.BucketName)

	data, err := s.ossService.DownloadObjectAsBytes(ctx, bucketType, paperPdfParsed.ObjectKey)
	if err != nil {
		return nil, errors.Biz("pdf.pdf.errors.download_failed")
	}

	// 处理元信息
	return util.DeserializePdfMetadata(data)
}

// GetPdfMarkDown 获取PDF的MarkDown
func (s *PdfParseService) GetPdfMarkDown(ctx context.Context, pdfId string) (*string, error) {
	paperPdfParsed, err := s.getPdfSource(ctx, pdfId, parseConstant.PdfOssTypeMarkdown)
	if err != nil {
		return nil, err
	}
	if paperPdfParsed == nil || paperPdfParsed.ObjectKey == "" {
		return nil, nil
	}
	//下载oss文件并返回byte[]
	bucketType := ossConstant.BucketTypeToEnum(s.config, paperPdfParsed.BucketName)

	data, err := s.ossService.DownloadObjectAsBytes(ctx, bucketType, paperPdfParsed.ObjectKey)
	if err != nil {
		return nil, errors.Biz("pdf.pdf.errors.download_failed")
	}
	markdown := string(data)
	return &markdown, nil
}

// 获取图表信息
func (s *PdfParseService) GetFigureAndTable(ctx context.Context, pdfId string) ([]*parsedPb.ImageRecord, error) {
	paperPdfParsed, err := s.getPdfSource(ctx, pdfId, parseConstant.PdfOssTypeFigureAndTable)
	if err != nil {
		return nil, err
	}
	var imageRecords []*parsedPb.ImageRecord
	err = json.Unmarshal([]byte(paperPdfParsed.RecordsJson), &imageRecords)
	if err != nil {
		return nil, errors.Biz("pdf.pdf.errors.deserialize_failed")
	}
	return imageRecords, nil
}

// GetPdfFullDocument 获取PDF段落信息
func (s *PdfParseService) GetPdfFullDocument(ctx context.Context, pdfId string) (*parsedPb.FullDocument, error) {
	paperPdfParsed, err := s.getPdfSource(ctx, pdfId, parseConstant.PdfOssTypeParagraphs)
	if err != nil {
		return nil, err
	}

	//下载oss文件并返回byte[]
	bucketType := ossConstant.BucketTypeToEnum(s.config, paperPdfParsed.BucketName)
	data, err := s.ossService.DownloadObjectAsBytes(ctx, bucketType, paperPdfParsed.ObjectKey)
	if err != nil {
		return nil, errors.Biz("pdf.pdf.errors.download_failed")
	}

	// 处理段落信息
	return util.DeserializePdfParagraphs(data)
}

// getPdfSource 根据pdfId获取PDF源文件的上传记录
func (s *PdfParseService) getPdfSource(ctx context.Context, pdfId string, fileType string) (*paperModel.PaperPdfParsed, error) {
	//根据pdfId查询
	paperPdf, err := s.paperPdfService.GetById(ctx, pdfId)
	if err != nil {
		s.logger.Error("获取PDF失败", "error", err)
		return nil, errors.Biz("pdf.pdf.errors.get_failed")
	}
	if paperPdf == nil {
		return nil, errors.Biz("pdf.pdf.errors.not_found")
	}
	//根据fileSHA256查询paper_pdf_parsed记录
	return s.paperPdfParsedService.GetBySourcePdfFileSHA256AndTypeAndVersion(ctx, paperPdf.FileSHA256, fileType, parseConstant.ParseVersion)
}

func (s *PdfParseService) convertToFigureAndTableInfo(ctx context.Context, figureAndTables []*parsedPb.FigureTable, imageRecords []*parsedPb.ImageRecord) []*pb.PdfFigureAndTableInfo {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.convertToFigureAndTableInfo")
	defer span.Finish()

	//转map
	imageRecordMap := make(map[string]*parsedPb.ImageRecord)
	for _, imageRecord := range imageRecords {
		imageRecordMap[imageRecord.Id] = imageRecord
	}

	//正则表达式模式
	rePattern := regexp.MustCompile(s.config.PDF.RePattern)

	var figureAndTableList []*pb.PdfFigureAndTableInfo
	for _, figureAndTable := range figureAndTables {
		//转换
		refContent := figureAndTable.RefContent
		var refIndex string
		if colonIndex := strings.Index(refContent, ":"); colonIndex != -1 {
			refIndex = strings.TrimSpace(refContent[:colonIndex])
		} else {
			refIndex = refContent
		}
		// 第二步：对 refIndex 进行正则表达式处理
		// 只有当 refIndex 非空，并且可能匹配 Figure/Table 模式时才进行处理
		// 第二步：对 refIndex 进行正则表达式处理
		// 只有当 refIndex 非空时才进行处理
		if refIndex != "" {
			if matches := rePattern.FindStringSubmatch(refIndex); len(matches) >= 3 {
				// 匹配成功，提取 "Figure N" 或 "Table N" 部分
				// matches[1] 是 "Figure" 或 "Table" (不区分大小写)
				// matches[2] 是数字部分，如 "2" 或 "1" 或 "2.1"
				cleanedRefIndex := strings.TrimSpace(matches[1] + " " + matches[2])
				refIndex = cleanedRefIndex
			} else {
				// 正则表达式不匹配，或者匹配到的内容不完整 (例如只有 "Figure" 没有数字)
				refIndex = ""
			}
		}

		var bbox *pb.PdfBBox
		if figureAndTable.Bbox != nil {
			bbox = &pb.PdfBBox{
				X0:           float32(figureAndTable.Bbox.X0),
				Y0:           float32(figureAndTable.Bbox.Y0),
				X1:           float32(figureAndTable.Bbox.X1),
				Y1:           float32(figureAndTable.Bbox.Y1),
				OriginHeight: float32(figureAndTable.Bbox.OriginHeight),
				OriginWidth:  float32(figureAndTable.Bbox.OriginWidth),
			}
		}
		//获取url
		var url string
		imageRecord := imageRecordMap[figureAndTable.Id]
		if imageRecord != nil {
			fileUrl, err := s.ossService.GetFileTemporaryURL(ctx, ossConstant.BucketTypeToEnum(s.config, imageRecord.BucketName), imageRecord.ObjectKey, 120*60)
			if err != nil {
				s.logger.Error("pdf.pdf.errors.get_file_temporary_url_failed", "error", err)
			}
			url = fileUrl
		}

		figureAndTableInfo := &pb.PdfFigureAndTableInfo{
			Desc:    figureAndTable.RefContent,
			PageNum: figureAndTable.Bbox.PageNumber,
			RefIdx:  refIndex,
			Bbox:    bbox,
			Url:     url,
		}
		figureAndTableList = append(figureAndTableList, figureAndTableInfo)
	}
	return figureAndTableList
}

func (s *PdfParseService) convertToFigureAndTableMarker(ctx context.Context, figureAndTables []*parsedPb.RefMarker) []*pb.FigureAndTableReferenceMarker {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.convertToFigureAndTableMarker")
	defer span.Finish()

	var figureAndTableList []*pb.FigureAndTableReferenceMarker
	for _, figureAndTable := range figureAndTables {
		var pageNum int32

		var bbox *pb.PdfBBox
		if figureAndTable.Bbox != nil {
			bbox = &pb.PdfBBox{
				X0:           float32(figureAndTable.Bbox.X0),
				Y0:           float32(figureAndTable.Bbox.Y0),
				X1:           float32(figureAndTable.Bbox.X1),
				Y1:           float32(figureAndTable.Bbox.Y1),
				OriginHeight: float32(figureAndTable.Bbox.OriginHeight),
				OriginWidth:  float32(figureAndTable.Bbox.OriginWidth),
			}
			pageNum = int32(figureAndTable.Bbox.PageNumber)
		}
		figureAndTableInfo := &pb.FigureAndTableReferenceMarker{
			RefContent: figureAndTable.RefContent,
			PageNum:    pageNum,
			RefIdx:     figureAndTable.RefIdx,
			Bbox:       bbox,
		}
		figureAndTableList = append(figureAndTableList, figureAndTableInfo)
	}
	return figureAndTableList
}

func (s *PdfParseService) convertToReferenceMarker(ctx context.Context, referenceMarkers []*parsedPb.RefMarker, references []*parsedPb.Reference) []*pb.ReferenceMarker {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.convertToReferenceMarker")
	defer span.Finish()

	referenceMap := make(map[string]*parsedPb.Reference)
	for _, reference := range references {
		referenceMap[reference.RefIdx] = reference
	}

	var refMarkers []*pb.ReferenceMarker
	for _, referenceMarker := range referenceMarkers {
		var bbox *pb.PdfBBox
		var pageNum int32
		if referenceMarker.Bbox != nil {
			bbox = &pb.PdfBBox{
				X0:           float32(referenceMarker.Bbox.X0),
				Y0:           float32(referenceMarker.Bbox.Y0),
				X1:           float32(referenceMarker.Bbox.X1),
				Y1:           float32(referenceMarker.Bbox.Y1),
				OriginHeight: float32(referenceMarker.Bbox.OriginHeight),
				OriginWidth:  float32(referenceMarker.Bbox.OriginWidth),
			}
			pageNum = int32(referenceMarker.Bbox.PageNumber)
		}
		// set refRaw
		var refRaw string
		if refMarker, ok := referenceMap[referenceMarker.RefIdx]; ok {
			refRaw = refMarker.ContentText
		}
		referenceMarker := &pb.ReferenceMarker{
			RefContent: referenceMarker.RefContent,
			RefIdx:     referenceMarker.RefIdx,
			RefRaw:     refRaw,
			Bbox:       bbox,
			PageNum:    pageNum,
			PaperId:    "999999999",
		}
		refMarkers = append(refMarkers, referenceMarker)
	}
	return refMarkers
}

func (s *PdfParseService) needParsedPDF(ctx context.Context, pdfId string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.needParsedPDF")
	defer span.Finish()

	//通过pdfId获取对应的Paper  PaperPdf UserDoc
	paperPdf, err := s.paperPdfService.paperPDFDAO.FindById(ctx, pdfId)
	if err != nil {
		return err
	}
	if paperPdf == nil {
		return errors.Biz("paper pdf not found")
	}
	//判断redis的状态  只有在解析失败的时候才重新发送解析消息
	cacheKey := fmt.Sprintf("%s%s", docConstant.ParsePDFStatusKeyPrefix, paperPdf.FileSHA256)
	var cacheValue docpb.UserDocParsedStatusEnum
	exists, err := s.cache.GetNotBizPrefix(ctx, cacheKey, &cacheValue)
	if err != nil {
		return errors.Biz("get redis cache failed")
	}
	if exists {
		//当redis中存在信息时，只要在解析过程中的状态，都不应该进行重复解析
		if cacheValue != docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSE_FAILED &&
			cacheValue != docpb.UserDocParsedStatusEnum_HEADER_DATA_PARSE_FAILED &&
			cacheValue != docpb.UserDocParsedStatusEnum_PARSE_FAILED {
			//如果redis中的状态不是失败状态，则不继续执行，直接返回
			return nil
		}
	}
	userDoc, paperPdf, paper, err := s.userDocService.GetUserUploadBaseDataBySHA256AndUserId(ctx, paperPdf.CreatorId, paperPdf.FileSHA256)
	if err != nil {
		return err
	}
	if userDoc == nil || paperPdf == nil || paper == nil {
		return errors.Biz("user doc not found")
	}
	//如果是已经解析成功了的，则直接返回
	if userDoc.ParseStatus == int(docpb.UserDocParsedStatusEnum_CONTENT_DATA_PARSED) {
		return nil
	}
	//设置状态为重新发送
	cacheValue = docpb.UserDocParsedStatusEnum_REPARSE
	err = s.cache.SetNotBizPrefix(ctx, cacheKey, cacheValue, 30*time.Minute)
	if err != nil {
		return errors.Biz("redis status update failed")
	}

	return nil
}

func (s *PdfParseService) needEmbedding(ctx context.Context, pdfId string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.needEmbedding")
	defer span.Finish()

	//通过pdfId获取对应的UserDoc
	userDoc, err := s.userDocService.GetByPdId(ctx, pdfId)
	if err != nil {
		return err
	}
	if userDoc == nil {
		return errors.Biz("user doc not found")
	}
	// 集成全文embedding - 发送文档创建事件到MQ
	err = s.sendDocCreatedEvent(ctx, userDoc.UserId, userDoc.DocName, userDoc.Id, userDoc.PdfId)
	if err != nil {
		// 集成失败，记录警告但不返回错误
		s.logger.Warn("send doc created event failed", "error", err, "docId", userDoc.Id, "pdfId", pdfId)
	}
	return nil
}

// sendDocCreatedEvent 发送文档创建事件到MQ，用于触发embedding集成
func (s *PdfParseService) sendDocCreatedEvent(ctx context.Context, userId string, docName string, docId, pdfId string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PdfParseService.sendDocCreatedEvent")
	defer span.Finish()

	// 检查producer是否已初始化
	if s.producer == nil {
		return errors.Biz("producer is not initialized")
	}

	// 构建文档创建事件消息
	eventMessage := map[string]interface{}{
		"event_type": "doc.created",
		"data": map[string]string{
			"doc_id":        docId,
			"pdf_id":        pdfId,
			"document_name": docName,
			"user_id":       userId,
		},
	}

	body, err := json.Marshal(eventMessage)
	if err != nil {
		s.logger.Error("marshal doc event message failed", "error", err)
		return err
	}

	message := &rocketmq.Message{
		Topic: s.config.RocketMQ.Topic.Event.Doc2DifyIntegrationEvent.Name,
		Body:  body,
	}

	messageId, err := s.producer.SendSync(ctx, message)
	if err != nil {
		s.logger.Error("send doc created event failed",
			"topic", message.GetTopic(),
			"docId", docId,
			"pdfId", pdfId,
			"error", err.Error())
		return errors.Biz("send doc created event failed")
	}

	s.logger.Info("doc created event sent successfully",
		"topic", message.GetTopic(),
		"messageId", messageId,
		"docId", docId,
		"pdfId", pdfId)
	return nil
}
