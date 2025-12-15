package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	parsepb "github.com/yb2020/odoc-proto/gen/go/parsed"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	ossService "github.com/yb2020/odoc/services/oss/service"
	"github.com/yb2020/odoc/services/parse/util/mineru"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

// MineruPDFParseService 实现基于Mineru的PDF解析服务
type MineruPDFParseService struct {
	config              *config.Config
	tracer              opentracing.Tracer
	logger              logging.Logger
	httpClient          http_client.HttpClient
	parseOperateService *ParseOperateService
	ossService          ossService.OssServiceInterface
}

// NewMineruPDFParseService 创建新的Mineru PDF解析服务实例
func NewMineruPDFParseService(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	httpClient http_client.HttpClient,
	parseOperateService *ParseOperateService,
	ossService ossService.OssServiceInterface,
) *MineruPDFParseService {
	return &MineruPDFParseService{
		config:              config,
		logger:              logger,
		tracer:              tracer,
		httpClient:          httpClient,
		parseOperateService: parseOperateService,
		ossService:          ossService,
	}
}

func (s *MineruPDFParseService) ExecuteParsePDF(ctx context.Context, pdfContent []byte, paperPdf *pdfModel.PaperPdf) (*parsepb.DocumentMetadata, *parsepb.FullDocument, []*parsepb.PageBlockData, map[string]*parsepb.ImageRecord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MineruPDFParseService.ExecuteParseWorkflow")
	defer span.Finish()
	//
	fileSHA256 := paperPdf.FileSHA256
	// 生成随机上传ID
	uploadId := idgen.GenerateUUID()
	// 1. 上传PDF   返回一个文件路径
	tempFilePath, err := s.parseOperateService.MineruUploadPDF(ctx, pdfContent, fileSHA256+".pdf", uploadId)
	if err != nil {
		s.logger.Error("Mineru upload pdf failed", "error", err)
		return nil, nil, nil, nil, err
	}
	s.logger.Info("Mineru 上传文件进度检测 ", "tempFilePath", tempFilePath)
	// 2. 循环检查进度直到完成
	isCompleted := false
	maxRetries := 60 // 最大重试次数，防止无限循环
	retryCount := 0
	var uploadProgressResponse *parsepb.UploadProgressResponse
	for !isCompleted && retryCount < maxRetries {
		isCompleted, uploadProgressResponse, err = s.parseOperateService.UploadProgress(ctx, uploadId)
		if err != nil {
			retryCount++
			time.Sleep(5 * time.Second)
			continue
		}
		if !isCompleted {
			// 等待一段时间再检查
			time.Sleep(3 * time.Second)
			retryCount++
		}
	}

	if !isCompleted {
		s.logger.Error("PDF parse timeout or failed")
		return nil, nil, nil, nil, errors.Biz("PDF parse timeout or failed")
	}
	s.logger.Info("上传文件进度检测完成", "uploadProgressResponse", uploadProgressResponse)
	//文件名称
	fileName := uploadProgressResponse.OrigName
	//文件大小
	fileSize := uploadProgressResponse.ChunkSize
	// 生成会话哈希值
	sessionHash := idgen.GenerateUUID()
	// 3. 加入解析队列
	s.logger.Info("Mineru 加入解析队列", "tempFilePath", tempFilePath, "sessionHash", sessionHash)
	resultID, err := s.parseOperateService.Join(ctx, tempFilePath, uploadId, sessionHash, fileName, int64(fileSize))
	if err != nil {
		s.logger.Error("Mineru join parse failed", "error", err)
		return nil, nil, nil, nil, err
	}
	s.logger.Info("Mineru 加入解析队列成功并开始获取数据", "resultID", resultID, "sessionHash", sessionHash)
	// 4. 获取数据
	documentData, err := s.parseOperateService.GetData(ctx, resultID, sessionHash)
	if err != nil {
		s.logger.Error("Mineru get data failed", "error", err)
		return nil, nil, nil, nil, err
	}
	// 5. 解析documentData  解析出zip包的下载地址
	output := documentData.Output
	if output == nil {
		s.logger.Error("parse result is empty")
		return nil, nil, nil, nil, errors.Biz("parse result is empty")
	}
	s.logger.Info("Mineru 解析数据成功")
	// 创建存储字符串和对象的集合
	// var stringValues []string
	var pathInfoValues []*parsepb.PathInfo
	// 这里的结果数组为[string,string,PathInfo,PathInfo] 不确定是否会随着版本进行变化
	// 目前来说，第三条数据就是我们需要的数据 这里面包含了zip包的下载地址
	for _, rawData := range output.Data {
		// 先尝试解析为PathInfo对象
		var pathInfo parsepb.PathInfo
		err := json.Unmarshal([]byte(rawData), &pathInfo)

		if err == nil && pathInfo.Path != "" {
			s.logger.Info("Mineru line pathInfo数据", "rawData", rawData)
			// 成功解析为PathInfo对象
			pathInfoValues = append(pathInfoValues, &pathInfo)
		}
		// else {
		// 	// 尝试解析为字符串
		// 	var strValue string
		// 	err = json.Unmarshal([]byte(rawData), &strValue)

		// 	if err == nil {
		// 		// 成功解析为字符串
		// 		stringValues = append(stringValues, strValue)
		// 	}
		// }
	}
	s.logger.Info("Mineru 解析数据成功", "pathInfoValues size:", len(pathInfoValues))
	// 6. 下载zip包并解压
	for _, pathInfo := range pathInfoValues {
		if pathInfo == nil {
			continue
		}
		s.logger.Info("下载zip包的文件对象", "pathInfo", pathInfo.Path, "pathInfo url", pathInfo.Url)
		if pathInfo.Path == "" {
			continue
		}
		if pathInfo.Size <= 0 {
			continue
		}
		// 判断是否是zip包
		if !strings.HasSuffix(pathInfo.Path, ".zip") {
			continue
		}
		// 下载并解压zip包
		var err error
		_, err = s.parseOperateService.DownloadAndExtractZip(ctx, pathInfo.Url, fileSHA256)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		break
	}
	s.logger.Info("Mineru 下载zip包并解压成功")
	// 7. 解析解压后的文件
	extractedDir := filepath.Join(s.config.PDF.Download.TempDirectory, fileSHA256)
	// 列出目录中的images文件夹
	s.logger.Info("Mineru 解析解压后的文件", "extractedDir", extractedDir)
	imagesDir := filepath.Join(extractedDir, s.config.PDF.Download.MineruImageDirectory)
	//处理images文件
	imageRecords, err := mineru.HandleImageFiles(extractedDir, imagesDir)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	s.logger.Info("Mineru 处理images文件成功")
	//处理content_list.json文件
	contentListJsonBytes, middleJsonBytes, err := mineru.GetMineruContentListAndMiddleJsonBytes(extractedDir)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	s.logger.Info("Mineru 处理content_list.json和middle.json文件开始")
	documentMetadata, fullDocument, pageBlocks, err := s.HandleContentListAndMiddleJsonByte(ctx, contentListJsonBytes, middleJsonBytes, imageRecords)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	s.logger.Info("Mineru 处理content_list.json和middle.json文件成功")
	// 判断这些数据是否为空， 为空则返回错误，让消费者进行重新解析就好
	if documentMetadata == nil || fullDocument == nil {
		return nil, nil, nil, nil, errors.Biz("Mineru parse pdf failed")
	}
	return documentMetadata, fullDocument, pageBlocks, imageRecords, nil
}

func (s *MineruPDFParseService) HandleContentListAndMiddleJsonByte(ctx context.Context, contentListJsonBytes []byte, middleJsonBytes []byte, imageRecords map[string]*parsepb.ImageRecord) (*parsepb.DocumentMetadata, *parsepb.FullDocument, []*parsepb.PageBlockData, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MineruPDFParseService.handleContentListAndMiddleJsonByte")
	defer span.Finish()

	// 记录公式的位置信息和隶属的标题
	formulaRecords := make(map[string]*parsepb.FormulaRecord)

	//处理content_list.json   contentTitles这个是该文件标记是有文字级别的 text_level
	metadata, contentTitles, err := mineru.HandleContentListJsonByte(contentListJsonBytes, &imageRecords, &formulaRecords)
	if err != nil {
		return nil, nil, nil, err
	}
	//处理middle.json   这种时候imageRecords和formulaRecords已经填充了。 并且imageRecords中填充了bbox信息，这个bbox坐标需要进行反填
	paragraphs, documentMetadata, pageBlocks, err := mineru.HandleMiddleJsonByte(middleJsonBytes, &imageRecords, &formulaRecords, contentTitles)
	if err != nil {
		s.logger.Error("handle middle json failed", "error", err)
		return nil, nil, nil, err
	}
	//反填充bbox - 图片和表格
	for _, ft := range metadata.FiguresAndTables {
		if ft == nil {
			continue
		}
		// 通过Id查找对应的imageRecord获取Bbox
		for _, imageRecord := range imageRecords {
			if imageRecord != nil && imageRecord.Id == ft.Id {
				ft.Bbox = imageRecord.Bbox
				break
			}
		}
	}
	//反填充bbox - 公式
	for _, formula := range metadata.Formulas {
		if formula == nil {
			continue
		}
		// 通过RefContent查找对应的formulaRecord获取Bbox
		if formulaRecord, exists := formulaRecords[formula.RefContent]; exists {
			formula.Bbox = formulaRecord.Bbox
		}
	}

	metadata.Catalogue = documentMetadata.Catalogue
	metadata.Abstract = documentMetadata.Abstract
	metadata.Acknowledgment = documentMetadata.Acknowledgment
	metadata.Authors = documentMetadata.Authors
	metadata.Title = documentMetadata.Title
	metadata.Pages = documentMetadata.Pages

	fullDocument := &parsepb.FullDocument{
		Paragraphs: paragraphs,
	}
	// TODO：本地测试 写入paragraphs.json和metadata.json
	// err = s.writeParagraphsToLocalFile(ctx, paragraphs)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// err = s.writeMetadataToLocalFile(ctx, metadata)
	// if err != nil {
	// 	return nil, nil, err
	// }
	return metadata, fullDocument, pageBlocks, nil
}

func (s *MineruPDFParseService) writeParagraphsToLocalFile(ctx context.Context, paragraphs []*parsepb.Paragraph) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MineruPDFParseService.writeParagraphsToLocalFile")
	defer span.Finish()

	//写入paragraphs.json
	filePath := filepath.Join(s.config.PDF.Download.TempDirectory, "paragraphs.json")

	paragraphsBytes, err := json.Marshal(paragraphs)
	if err != nil {
		return err
	}

	os.WriteFile(filePath, paragraphsBytes, 0644)

	return nil
}

func (s *MineruPDFParseService) writeMetadataToLocalFile(ctx context.Context, metadata *parsepb.DocumentMetadata) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MineruPDFParseService.writeMetadataToLocalFile")
	defer span.Finish()

	//写入metadata.json
	filePath := filepath.Join(s.config.PDF.Download.TempDirectory, "metadata.json")

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return err
	}

	os.WriteFile(filePath, metadataBytes, 0644)

	return nil
}
