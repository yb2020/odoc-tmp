package service

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/errors"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/idgen"
	"github.com/yb2020/odoc/pkg/logging"
	docpb "github.com/yb2020/odoc/proto/gen/go/doc"
	ossPb "github.com/yb2020/odoc/proto/gen/go/oss"
	parsepb "github.com/yb2020/odoc/proto/gen/go/parsed"
	ossModel "github.com/yb2020/odoc/services/oss/model"
	ossService "github.com/yb2020/odoc/services/oss/service"
	paperModel "github.com/yb2020/odoc/services/paper/model"
	paperService "github.com/yb2020/odoc/services/paper/service"
	"github.com/yb2020/odoc/services/parse/constant"
	"github.com/yb2020/odoc/services/parse/util"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

type ParseOperateService struct {
	config                *config.Config
	tracer                opentracing.Tracer
	logger                logging.Logger
	paperPdfParsedService *paperService.PaperPdfParsedService
	ossService            ossService.OssServiceInterface
	httpClient            http_client.HttpClient
}

func NewParseOperateService(
	config *config.Config,
	tracer opentracing.Tracer,
	logger logging.Logger,
	paperPdfParsedService *paperService.PaperPdfParsedService,
	ossService ossService.OssServiceInterface,
	httpClient http_client.HttpClient,
) *ParseOperateService {
	return &ParseOperateService{
		config:                config,
		tracer:                tracer,
		logger:                logger,
		paperPdfParsedService: paperPdfParsedService,
		ossService:            ossService,
		httpClient:            httpClient,
	}
}

// UploadPDFResponse 上传PDF响应结构
// 实际返回的是一个文件路径的数组
type UploadPDFResponse []string

type TempDataResult struct {
	Msg     string `json:"msg"`
	EventId string `json:"event_id"`
	Output  struct {
		Data []json.RawMessage `json:"data"`
	} `json:"output"`
	Success bool   `json:"success"`
	Title   string `json:"title"`
}

// EventOutput 事件输出结构
type EventOutput struct {
	Data []json.RawMessage `json:"data,omitempty"` // 使用json.RawMessage表示任意JSON数据
}

// ProcessCompletedEventMessage 处理完成事件消息
type ProcessCompletedEventMessage struct {
	Msg     string       `json:"msg,omitempty"`      // 消息
	EventId string       `json:"event_id,omitempty"` // 事件ID
	Output  *EventOutput `json:"output,omitempty"`   // 输出内容
	Success bool         `json:"success,omitempty"`  // 是否成功
	Title   string       `json:"title,omitempty"`    // 标题
}

// 根据parse pb 的解析结果，上传meta和到oss
func (s *ParseOperateService) UploadDocument(
	ctx context.Context,
	documentMetadata *parsepb.DocumentMetadata,
	pdfSHA256 string,
	bucketType ossPb.OSSBucketEnum,
) (metaRecord *ossModel.OssRecord, err error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.uploadDocument")
	defer span.Finish()

	// 2. 将文档转换为JSON
	documentMetaJSON, err := json.Marshal(documentMetadata)
	if err != nil {
		return nil, errors.Biz("mq.ParseUploadPdfService error, failed to marshal document metadata!")
	}
	// 计算JSON数据的SHA256值
	metaJsonSHA256 := fmt.Sprintf("%x", sha256.Sum256(documentMetaJSON))

	// 3. 上传JSON文件到OSS
	uuid := idgen.GenerateUUID()
	// meta文件
	metaFileName := fmt.Sprintf("%s%s", uuid, constant.MetadataJsonSuffix)
	// 上传不包含段落的文档
	metaReader := bytes.NewReader(documentMetaJSON)
	//单位： byte
	metaSize := int64(len(documentMetaJSON))

	metaObjectKeyGen := func(uniqueId, fileName string) string {
		return fmt.Sprintf("%s/%s/%s/%s", pdfSHA256, constant.ParsedPdfCatalog, constant.ParseVersion, metaFileName)
	}
	metaRecord, err = s.ossService.UploadObjectAndSaveRecord(
		ctx,
		bucketType,
		metaFileName,
		metaJsonSHA256,
		metaReader,
		metaSize,
		metaObjectKeyGen,
		"",
		nil,
	)
	if err != nil {
		return nil, errors.Biz("mq.ParseUploadPdfService error, failed to upload meta file!")
	}
	return metaRecord, nil
}

// 根据parse pb 的解析结果，上传full document到oss
func (s *ParseOperateService) UploadFullDocument(
	ctx context.Context,
	fullDocument *parsepb.FullDocument,
	pdfSHA256 string,
	bucketType ossPb.OSSBucketEnum,
) (fullDocumentRecord *ossModel.OssRecord, err error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.uploadFullDocument")
	defer span.Finish()

	// 1. 将文档转换为JSON
	fullDocumentJSON, err := json.Marshal(fullDocument)
	if err != nil {
		return nil, errors.Biz("mq.ParseOperateService error, failed to marshal full document!")
	}
	// 计算JSON数据的SHA256值
	fullDocumentJsonSHA256 := fmt.Sprintf("%x", sha256.Sum256(fullDocumentJSON))

	// 2. 上传JSON文件到OSS
	uuid := idgen.GenerateUUID()
	// paragraphs文件
	paragraphsFileName := fmt.Sprintf("%s%s", uuid, constant.ParagraphsJsonSuffix)

	// 上传仅包含段落的文档
	paragraphsReader := bytes.NewReader(fullDocumentJSON)
	paragraphsSize := int64(len(fullDocumentJSON))

	paragraphsObjectKeyGen := func(uniqueId, fileName string) string {
		return fmt.Sprintf("%s/%s/%s/%s", pdfSHA256, constant.ParsedPdfCatalog, constant.ParseVersion, paragraphsFileName)
	}
	fullDocumentRecord, err = s.ossService.UploadObjectAndSaveRecord(
		ctx,
		bucketType,
		paragraphsFileName,
		fullDocumentJsonSHA256,
		paragraphsReader,
		paragraphsSize,
		paragraphsObjectKeyGen,
		"",
		nil,
	)
	if err != nil {
		return nil, errors.Biz("mq.ParseUploadPdfService error, failed to upload paragraphs file!")
	}
	return fullDocumentRecord, nil
}

func (s *ParseOperateService) UploadPageBlocks(
	ctx context.Context,
	pageBlocks []*parsepb.PageBlockData,
	pdfSHA256 string,
	bucketType ossPb.OSSBucketEnum,
) (pageBlocksRecord *ossModel.OssRecord, err error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.uploadPageBlocks")
	defer span.Finish()

	// 1. 将文档转换为JSON
	pageBlocksJSON, err := json.Marshal(pageBlocks)
	if err != nil {
		return nil, errors.Biz("mq.ParseUploadPdfService error, failed to marshal page blocks!")
	}
	// 计算JSON数据的SHA256值
	pageBlocksJsonSHA256 := fmt.Sprintf("%x", sha256.Sum256(pageBlocksJSON))

	// 2. 上传JSON文件到OSS
	uuid := idgen.GenerateUUID()
	// pageBlocks文件
	pageBlocksFileName := fmt.Sprintf("%s%s", uuid, constant.PageBlocksJsonSuffix)

	// 上传仅包含page blocks的文档
	pageBlocksReader := bytes.NewReader(pageBlocksJSON)
	pageBlocksSize := int64(len(pageBlocksJSON))

	pageBlocksObjectKeyGen := func(uniqueId, fileName string) string {
		return fmt.Sprintf("%s/%s/%s/%s", pdfSHA256, constant.ParsedPdfCatalog, constant.ParseVersion, pageBlocksFileName)
	}
	pageBlocksRecord, err = s.ossService.UploadObjectAndSaveRecord(
		ctx,
		bucketType,
		pageBlocksFileName,
		pageBlocksJsonSHA256,
		pageBlocksReader,
		pageBlocksSize,
		pageBlocksObjectKeyGen,
		"",
		nil,
	)
	if err != nil {
		return nil, errors.Biz("mq.ParseUploadPdfService error, failed to upload page blocks file!")
	}
	return pageBlocksRecord, nil
}

func (s *ParseOperateService) UploadMarkdown(ctx context.Context, pdfSHA256 string, bucketType ossPb.OSSBucketEnum) (*ossModel.OssRecord, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.uploadMarkdown")
	defer span.Finish()

	markdownBytes, err := s.GetMineruMarkdownFile(ctx, pdfSHA256)
	if err != nil {
		return nil, errors.BizWrap("read markdown file failed", err)
	}

	markdownSHA256 := fmt.Sprintf("%x", sha256.Sum256(markdownBytes))
	// 上传markdown文件到OSS
	markdownFileName := fmt.Sprintf("%s%s", idgen.GenerateUUID(), constant.MineruMdSuffix)
	markdownObjectKeyGen := func(uniqueId, fileName string) string {
		return fmt.Sprintf("%s/%s/%s/%s", pdfSHA256, constant.ParsedPdfCatalog, constant.ParseVersion, markdownFileName)
	}
	markdownRecord, err := s.ossService.UploadObjectAndSaveRecord(
		ctx,
		bucketType,
		markdownFileName,
		markdownSHA256,
		bytes.NewReader(markdownBytes),
		int64(len(markdownBytes)),
		markdownObjectKeyGen,
		"",
		nil,
	)
	if err != nil {
		return nil, errors.BizWrap("upload markdown file failed", err)
	}
	return markdownRecord, nil
}

// 保存解析之后的原始记录到PaperPdfParsed
func (s *ParseOperateService) SavePaperPdfParsed(
	ctx context.Context,
	paperPdf *pdfModel.PaperPdf,
	metaRecord *ossModel.OssRecord,
	fullDocumentRecord *ossModel.OssRecord,
	imageRecords []parsepb.ImageRecord,
	markdownRecord *ossModel.OssRecord,
	pageBlocksRecord *ossModel.OssRecord,
) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.savePaperPdfParsedParagraphs")
	defer span.Finish()

	// 创建元数据记录
	pdfMetaRecord := &paperModel.PaperPdfParsed{
		SourcePdfSHA256: paperPdf.FileSHA256,
		FileSHA256:      metaRecord.FileSHA256,
		FileType:        constant.PdfOssTypeMetadata,
		FileName:        metaRecord.FileName,
		FileSize:        metaRecord.FileSize,
		ObjectKey:       metaRecord.ObjectKey,
		BucketName:      metaRecord.BucketName,
	}
	pdfMetaRecord.BaseModel.Id = idgen.GenerateUUID()
	pdfMetaRecord.BaseModel.CreatorId = paperPdf.BaseModel.CreatorId
	pdfMetaRecord.BaseModel.ModifierId = paperPdf.BaseModel.CreatorId
	pdfMetaRecord.BaseModel.CreatedAt = time.Now()
	pdfMetaRecord.Version = constant.ParseVersion

	//创建全文记录
	pdfFullDocumentRecord := &paperModel.PaperPdfParsed{
		SourcePdfSHA256: paperPdf.FileSHA256,
		FileSHA256:      fullDocumentRecord.FileSHA256,
		FileType:        constant.PdfOssTypeParagraphs,
		FileName:        fullDocumentRecord.FileName,
		FileSize:        fullDocumentRecord.FileSize,
		ObjectKey:       fullDocumentRecord.ObjectKey,
		BucketName:      fullDocumentRecord.BucketName,
	}
	pdfFullDocumentRecord.BaseModel.Id = idgen.GenerateUUID()
	pdfFullDocumentRecord.BaseModel.CreatorId = paperPdf.BaseModel.CreatorId
	pdfFullDocumentRecord.BaseModel.ModifierId = paperPdf.BaseModel.CreatorId
	pdfFullDocumentRecord.BaseModel.CreatedAt = time.Now()
	pdfFullDocumentRecord.Version = constant.ParseVersion

	//创建图表记录
	// 文件大小
	var imagesFileSize int64
	for i := range imageRecords {
		imagesFileSize += imageRecords[i].FileSize
	}
	jsonStr, jsonParseErr := json.Marshal(imageRecords)
	if jsonParseErr != nil {
		return errors.Biz("PDF解析失败")
	}
	imageRecord := &paperModel.PaperPdfParsed{
		SourcePdfSHA256: paperPdf.FileSHA256,
		FileType:        constant.PdfOssTypeFigureAndTable,
		FileSize:        imagesFileSize,
		RecordsJson:     string(jsonStr),
	}
	imageRecord.BaseModel.Id = idgen.GenerateUUID()
	imageRecord.BaseModel.CreatorId = paperPdf.BaseModel.CreatorId
	imageRecord.BaseModel.ModifierId = paperPdf.BaseModel.CreatorId
	imageRecord.BaseModel.CreatedAt = time.Now()
	imageRecord.Version = constant.ParseVersion

	//创建markdown记录
	markdownParsedRecord := &paperModel.PaperPdfParsed{
		SourcePdfSHA256: paperPdf.FileSHA256,
		FileSHA256:      markdownRecord.FileSHA256,
		FileType:        constant.PdfOssTypeMarkdown,
		FileName:        markdownRecord.FileName,
		FileSize:        markdownRecord.FileSize,
		ObjectKey:       markdownRecord.ObjectKey,
		BucketName:      markdownRecord.BucketName,
	}
	markdownParsedRecord.BaseModel.Id = idgen.GenerateUUID()
	markdownParsedRecord.BaseModel.CreatorId = paperPdf.BaseModel.CreatorId
	markdownParsedRecord.BaseModel.ModifierId = paperPdf.BaseModel.CreatorId
	markdownParsedRecord.BaseModel.CreatedAt = time.Now()
	markdownParsedRecord.Version = constant.ParseVersion

	//创建page blocks记录
	pageBlocksParsedRecord := &paperModel.PaperPdfParsed{
		SourcePdfSHA256: paperPdf.FileSHA256,
		FileType:        constant.PdfOssTypePageBlocks,
		FileSize:        pageBlocksRecord.FileSize,
		ObjectKey:       pageBlocksRecord.ObjectKey,
		BucketName:      pageBlocksRecord.BucketName,
	}
	pageBlocksParsedRecord.BaseModel.Id = idgen.GenerateUUID()
	pageBlocksParsedRecord.BaseModel.CreatorId = paperPdf.BaseModel.CreatorId
	pageBlocksParsedRecord.BaseModel.ModifierId = paperPdf.BaseModel.CreatorId
	pageBlocksParsedRecord.BaseModel.CreatedAt = time.Now()
	pageBlocksParsedRecord.Version = constant.ParseVersion

	err := s.paperPdfParsedService.BatchSave(ctx, []*paperModel.PaperPdfParsed{pdfMetaRecord, pdfFullDocumentRecord, imageRecord, markdownParsedRecord, pageBlocksParsedRecord})
	if err != nil {
		return errors.Biz("mq.ParseUploadPdfService error, failed to save pdf oss records!")
	}
	return nil
}

// 根据pdf的sha256判断是否存在解析记录
func (s *ParseOperateService) IsPaperPdfParsedExists(ctx context.Context, sha256 string) (bool, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.isPaperPdfParsedExists")
	defer span.Finish()

	return s.paperPdfParsedService.HasExistBySourcePdfFileSHA256AndVersion(ctx, sha256, constant.ParseVersion)
}

// 根据pdf的sha256查询解析记录
func (s *ParseOperateService) GetPaperPdfParsedBySHA256(ctx context.Context, sha256 string) ([]*paperModel.PaperPdfParsed, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.getPaperPdfParsedBySHA256")
	defer span.Finish()

	return s.paperPdfParsedService.GetBySourcePdfFileSHA256AndVersion(ctx, sha256, constant.ParseVersion)
}

// DownloadPdfContent 从URL下载PDF内容  TODO: 这里存在性能问题，所以可能需要改造consumer，增加程序是否需要启动消费某个topic的配置
func (s *ParseOperateService) DownloadPdfContent(ctx context.Context, pdfUrl string) ([]byte, error) {
	// 使用HTTP客户端下载PDF内容
	headers := map[string]string{}
	pdfContent, err := s.httpClient.Get(pdfUrl, headers)
	if err != nil {
		return nil, errors.BizWrap("download pdf content failed", err)
	}
	if len(pdfContent) == 0 {
		return nil, errors.Biz("download pdf content is empty")
	}
	return pdfContent, nil
}

// ============================================= 文件操作相关方法 =============================================

// 写入pdf源文件到临时目录
func (s *ParseOperateService) WritePdfContentToTempFile(ctx context.Context, pdfContent []byte, pdfSHA256 string) error {
	// 确保目录存在
	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, pdfSHA256)
	if err := os.MkdirAll(tempDirectory, 0755); err != nil {
		return errors.BizWrap("create temp dir failed", err)
	}
	// 保存文件到临时目录
	tempFilePath := filepath.Join(tempDirectory, pdfSHA256+".pdf")
	if err := os.WriteFile(tempFilePath, pdfContent, 0644); err != nil {
		return errors.BizWrap("save pdf to temp dir failed", err)
	}
	return nil
}

// 从临时目录读取出pdf源文件
func (s *ParseOperateService) ReadPdfContentFromTempFile(ctx context.Context, pdfSHA256 string) ([]byte, error) {
	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, pdfSHA256)
	pdfTempFilePath := filepath.Join(tempDirectory, pdfSHA256+".pdf")
	fileContent, err := os.ReadFile(pdfTempFilePath)
	if err != nil {
		return nil, errors.Biz("mq.ParseUploadPdfService error, failed to read temp pdf!")
	}
	return fileContent, nil
}

// 判断临时目录中是否存在源pdf文件
func (s *ParseOperateService) FileExistsInTempDirectory(ctx context.Context, pdfSHA256 string) bool {
	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, pdfSHA256)
	pdfTempFilePath := filepath.Join(tempDirectory, pdfSHA256+".pdf")
	_, err := os.Stat(pdfTempFilePath)
	return err == nil
}

// 删除临时目录中的源pdf文件
func (s *ParseOperateService) DeletePdfContentFromTempFile(ctx context.Context, pdfSHA256 string) error {
	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, pdfSHA256)
	pdfTempFilePath := filepath.Join(tempDirectory, pdfSHA256+".pdf")
	if err := os.Remove(pdfTempFilePath); err != nil {
		return errors.Biz("mq.ParseUploadPdfService error, failed to delete temp pdf!")
	}
	return nil
}

// DownloadAndExtractZip 下载并解压ZIP文件
// return: 解压后的目录
func (s *ParseOperateService) DownloadAndExtractZip(ctx context.Context, url string, sha256 string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.DownloadAndExtractZip")
	defer span.Finish()

	s.logger.Info("Mineru 下载zip包并解压", "url", url, "sha256", sha256)
	// 确保目录存在
	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, sha256)
	err := os.MkdirAll(tempDirectory, 0755)
	if err != nil {
		s.logger.Error("create mineru pdf temp dir failed", "error", err)
		return "", errors.BizWrap("create mineru pdf temp dir failed", err)
	}

	//由于zip包中包含image目录，需要提前创建
	//打印这里的创建目录
	s.logger.Info("Mineru 下载zip包并解压", "tempDirectory", tempDirectory, "sha256", sha256)
	os.MkdirAll(filepath.Join(tempDirectory, s.config.PDF.Download.MineruImageDirectory), 0755)

	// 生成临时文件路径
	zipFileName := filepath.Base(url)
	zipFilePath := filepath.Join(tempDirectory, zipFileName)

	// 下载ZIP文件
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.BizWrap("download zip failed", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Biz(fmt.Sprintf("download zip failed, status code: %d", resp.StatusCode))
	}

	// 创建目标文件
	out, err := os.Create(zipFilePath)
	if err != nil {
		return "", errors.BizWrap("create zip file failed", err)
	}
	defer out.Close()

	// 将响应体写入文件
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", errors.BizWrap("write zip file failed", err)
	}
	s.logger.Info("Mineru 下载zip包成功 准备打开并解压ZIP包", "zipFilePath", zipFilePath)

	// 打开ZIP文件
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return "", errors.BizWrap("open zip file failed", err)
	}
	defer reader.Close()

	// 解压文件
	for _, file := range reader.File {
		path := filepath.Join(tempDirectory, file.Name)

		// 确保文件路径不会超出目标目录
		if !strings.HasPrefix(path, filepath.Clean(tempDirectory)+string(os.PathSeparator)) {
			s.logger.Error("Mineru 解压zip包失败 文件路径超出目标目录", "path", path)
			return "", errors.Biz(fmt.Sprintf("illegal file path: %s", path))
		}
		// 打开目标文件
		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			s.logger.Error("Mineru 解压zip包失败 创建目标文件失败", "path", path, "error", err.Error())
			return "", errors.BizWrap("create target file failed", err)
		}
		// 打开源文件
		srcFile, err := file.Open()
		if err != nil {
			dstFile.Close()
			s.logger.Error("Mineru 解压zip包失败 打开源文件失败", "path", path, "error", err.Error())
			return "", errors.BizWrap("open source file failed", err)
		}

		// 复制文件内容
		_, err = io.Copy(dstFile, srcFile)
		dstFile.Close()
		srcFile.Close()
		if err != nil {
			s.logger.Error("Mineru 解压zip包失败 复制文件内容失败", "path", path, "error", err.Error())
			return "", errors.BizWrap("copy file content failed", err)
		}
	}
	s.logger.Info("Mineru 解压zip包成功")

	return tempDirectory, nil
}

// 获取临时目录中的图片数据
func (s *ParseOperateService) GetTempDirectoryImage(ctx context.Context, sha256 string) (map[string][]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.GetTempDirectoryImage")
	defer span.Finish()

	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, sha256)
	imagePath := filepath.Join(tempDirectory, s.config.PDF.Download.MineruImageDirectory)
	imagesFiles, err := os.ReadDir(imagePath)
	if err != nil {
		return nil, errors.BizWrap("read images directory failed", err)
	}
	var imageRecords = make(map[string][]byte)
	for _, file := range imagesFiles {
		if file.IsDir() {
			continue
		}
		// 读取文件内容
		imageFile, err := os.ReadFile(filepath.Join(imagePath, file.Name()))
		if err != nil {
			return nil, errors.BizWrap("read image file failed", err)
		}
		imageRecords[file.Name()] = imageFile
	}
	return imageRecords, nil
}

// 删除临时目录
func (s *ParseOperateService) DeleteTempDirectory(ctx context.Context, sha256 string) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.DeleteTempDirectory")
	defer span.Finish()

	tempDirectory := filepath.Join(s.config.PDF.Download.TempDirectory, sha256)
	os.RemoveAll(tempDirectory)
	return nil
}

// 获取mineru中的.md格式的markdown文件
func (s *ParseOperateService) GetMineruMarkdownFile(ctx context.Context, pdfSHA256 string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.GetMineruMarkdownFile")
	defer span.Finish()

	// 1. 获取mineru中的.md文件
	extractedDir := filepath.Join(s.config.PDF.Download.TempDirectory, pdfSHA256)
	// 列出解压目录中的文件
	files, err := os.ReadDir(extractedDir)
	if err != nil {
		return nil, errors.BizWrap("read extracted directory failed", err)
	}
	var markdownFile os.DirEntry
	//找到.md后缀的文件
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), constant.MineruMdSuffix) {
			continue
		}
		markdownFile = file
		break
	}
	if markdownFile == nil {
		return nil, errors.Biz("read markdown file failed")
	}
	// 读取文件内容
	markdownBytes, err := os.ReadFile(filepath.Join(extractedDir, markdownFile.Name()))
	if err != nil {
		return nil, errors.BizWrap("read markdown file failed", err)
	}
	return markdownBytes, nil
}

// ============================================= http接口调用操作 =============================================

// 调用grobid接口 解析pdf的全文信息
func (s *ParseOperateService) GetGrobidParsedDocumentByte(ctx context.Context, pdfContent []byte) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.getGrobidParsedDocumentByte")
	defer span.Finish()

	url := s.config.PDF.Parse.Grobid.URL + s.config.PDF.Parse.Grobid.FulltextDocument
	// 构建表单参数
	formData := map[string]string{
		"consolidateHeader":      "1",
		"consolidateCitations":   "0", //这个参数设置为1的情况下会很慢
		"consolidateFunders":     "1",
		"includeRawAffiliations": "1",
		"includeRawCitations":    "1",
		"includeRawCopyrights":   "1",
		"segmentSentences":       "1",
		"generateIDs":            "1",
		//这里必须都加上，否则会导致解析的部分数据不存在坐标信息figure formula这两个可以不需要
		"teiCoordinates": "persName,figure,ref,biblStruct,formula,head,note,title,affiliation,s,p",
	}
	// 设置文件内容
	fileContents := map[string][]byte{
		"input": pdfContent,
	}
	// 设置文件名
	fileNameUUid := idgen.GenerateUUID()
	fileNames := map[string]string{
		"input": fileNameUUid + ".pdf",
	}
	// 设置请求头
	headers := map[string]string{
		"Accept": "application/xml",
	}
	// 发送请求
	timeout := time.Duration(s.config.PDF.Parse.Grobid.Timeout) * time.Minute
	responseData, err := s.httpClient.PostMultipartFormWithFileInput(url, formData, fileContents, fileNames, headers, timeout)
	if err != nil {
		s.logger.Error("grobid request failed", "error", err)
		return nil, err
	}
	if len(responseData) == 0 {
		s.logger.Error("grobid response is empty")
		return nil, nil
	}
	return responseData, nil
}

// 调用grobid接口  获取header信息
func (s *ParseOperateService) GetGrobidParsedHeaderByte(ctx context.Context, pdfContent []byte) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.ParsePDFHeaderWithGrobid")
	defer span.Finish()

	url := s.config.PDF.Parse.Grobid.URL + s.config.PDF.Parse.Grobid.HeaderDocument
	// 设置请求头
	headers := map[string]string{
		"Accept": "application/xml",
	}
	// 设置文件内容
	fileContents := map[string][]byte{
		"input": pdfContent, // 你已经下载的 PDF 内容
	}
	// 设置文件名
	fileNames := map[string]string{
		"input": "document.pdf",
	}
	// 发送请求
	timeout := time.Duration(s.config.PDF.Parse.Grobid.Timeout) * time.Minute
	responseData, err := s.httpClient.PostMultipartFormWithFileInput(url, nil, fileContents, fileNames, headers, timeout)
	if err != nil {
		s.logger.Error("grobid request failed", "error", err)
		return nil, err
	}
	if len(responseData) == 0 {
		s.logger.Error("grobid response is empty")
		return nil, nil
	}
	return responseData, nil
}

// ============================================= mineru api 调用 =============================================

// 调用mineru接口  上传文件
// 返回值为文件在mineru服务器上的临时路径
func (s *ParseOperateService) MineruUploadPDF(ctx context.Context, pdfContent []byte, fileName string, uploadId string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.mineruUploadPDF")
	defer span.Finish()

	// 从配置文件获取URL
	uploadConfigURL := s.config.PDF.Parse.Mineru.URL + s.config.PDF.Parse.Mineru.UploadURL
	// 构建上传URL
	uploadURL := fmt.Sprintf("%s?upload_id=%s", uploadConfigURL, uploadId)

	// 设置请求头
	headers := map[string]string{
		"Accept":       "*/*",
		"Content-Type": "multipart/form-data",
	}
	// 设置文件内容，注意表单字段名为'files'
	fileContents := map[string][]byte{
		"files": pdfContent,
	}
	// 设置文件名
	fileNames := map[string]string{
		"files": fileName,
	}
	timeout := time.Duration(s.config.PDF.Parse.Mineru.Timeout) * time.Minute
	// 发送上传请求
	responseData, err := s.httpClient.PostMultipartFormWithFileInput(uploadURL, nil, fileContents, fileNames, headers, timeout)
	if err != nil {
		s.logger.Error("Mineru upload failed , uploadUrl:", uploadURL, "error", err)
		return "", err
	}
	// 解析响应获取文件路径
	var response UploadPDFResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		s.logger.Error("parse upload response failed", "error", err, "responseData", string(responseData))
		return "", errors.BizWrap("parse upload response failed", err)
	}
	// 检查响应是否为空
	if len(response) == 0 {
		s.logger.Error("upload pdf failed, response is empty array")
		return "", errors.Biz("upload pdf failed, response is empty array")
	}
	// 返回第一个文件路径作为任务ID
	filePath := response[0]
	s.logger.Info("文件上传成功", "filePath", filePath)
	return filePath, nil
}

// 调用mineru接口  检查上传进度
// 返回值为true表示上传完成
func (s *ParseOperateService) UploadProgress(ctx context.Context, uploadId string) (bool, *parsepb.UploadProgressResponse, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MineruPDFParseService.UploadProgress")
	defer span.Finish()

	MineruProgressURL := s.config.PDF.Parse.Mineru.URL + s.config.PDF.Parse.Mineru.ProgressURL
	progressURL := fmt.Sprintf("%s?upload_id=%s", MineruProgressURL, uploadId)
	// 设置请求头
	headers := map[string]string{
		"Accept": "text/event-stream", // 事件流格式
	}
	// 发送进度检查请求
	responseData, err := s.httpClient.Get(progressURL, headers)
	if err != nil {
		s.logger.Error("check parse progress failed", "error", err, "uploadId", uploadId)
		return false, nil, errors.BizWrap("check parse progress failed", err)
	}
	// 响应是事件流格式，需要按行解析
	lines := strings.Split(string(responseData), "\n")
	// 默认不完成
	isCompleted := false

	var response parsepb.UploadProgressResponse
	// 处理每一行数据
	for _, line := range lines {
		// 跳过空行
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 去除前缀"data: "
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		} else {
			continue // 不是数据行，跳过
		}
		// 解析JSON   这里的响应结构和mineru中的一致，所以这里直接使用minerU的协议就好
		if err := json.Unmarshal([]byte(line), &response); err != nil {
			continue // 跳过解析失败的行
		}
		// 检查是否完成
		if response.Msg == "done" {
			isCompleted = true
			break
		}
	}
	return isCompleted, &response, nil
}

// Join 加入解析队列   返回一个eventID 事件ID  这个事件ID用于对比后续的解析结果，判断解析结果是否是当前事件的
// TODO 这里暂时可以不保存这个数据，因为这里的处理结果将是异步且单线程的
func (s *ParseOperateService) Join(ctx context.Context, filePath string, uploadId string, sessionHash string, fileName string, fileSize int64) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.Join")
	defer span.Finish()

	// 获取配置的joinURL
	joinURL := s.config.PDF.Parse.Mineru.URL + s.config.PDF.Parse.Mineru.JoinURL

	// 构造文件URL
	fileURL := fmt.Sprintf("%s%s", s.config.PDF.Parse.Mineru.URL+s.config.PDF.Parse.Mineru.FileURL, filePath)

	//获取最大解析页数
	maxPage := s.config.PDF.Parse.Mineru.MaxPage
	// 根据API要求构造请求体
	requestBody := map[string]interface{}{
		"data": []interface{}{
			map[string]interface{}{
				"path":      filePath,
				"url":       fileURL,
				"orig_name": fileName,
				"size":      fileSize,
				"mime_type": "application/pdf",
				"meta": map[string]string{
					"_type": "gradio.FileData",
				},
			},
			maxPage, // 解析的最大页数
			false,
			true,
			true,
			"en", // 语言
		},
		"event_data":   nil,
		"fn_index":     2,
		"trigger_id":   19,
		"session_hash": sessionHash,
	}

	// 设置请求头
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// 发送请求
	timeout := time.Duration(s.config.PDF.Parse.Mineru.Timeout) * time.Minute
	responseData, err := s.httpClient.PostWithTimeout(joinURL, requestBody, headers, timeout)
	if err != nil {
		s.logger.Error("merge parse result failed", "error", err, "filePath", filePath)
		return "", errors.BizWrap("merge parse result failed", err)
	}

	// 解析响应获取事件ID
	var response parsepb.JoinResponse
	if err := json.Unmarshal(responseData, &response); err != nil {
		s.logger.Error("parse merge response failed", "error", err, "responseData", string(responseData))
		return "", errors.BizWrap("parse merge response failed", err)
	}

	// 检查事件ID是否存在
	if response.EventId == "" {
		s.logger.Error("merge parse result failed", "response", string(responseData))
		return "", errors.Biz("merge parse result failed: no event id")
	}
	return response.EventId, nil
}

// GetData 获取解析数据
func (s *ParseOperateService) GetData(ctx context.Context, eventID string, sessionHash string) (*parsepb.ProcessCompletedEventMessage, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "MineruPDFParseService.GetData")
	defer span.Finish()

	// 获取配置的获取数据URL
	MineruDataURL := s.config.PDF.Parse.Mineru.URL + s.config.PDF.Parse.Mineru.DataURL

	// 使用固定的获取数据URL，并添加session_hash参数
	dataURL := fmt.Sprintf("%s?session_hash=%s", MineruDataURL, sessionHash)

	// 设置请求头
	headers := map[string]string{
		"Accept": "text/event-stream", // 事件流格式
	}

	// 发送获取数据请求
	timeout := time.Duration(s.config.PDF.Parse.Mineru.Timeout) * time.Minute
	responseData, err := s.httpClient.GetWithTimeout(dataURL, headers, timeout)
	if err != nil {
		s.logger.Error("get parse data failed", "error", err, "eventID", eventID)
		return nil, errors.BizWrap("get parse data failed", err)
	}

	// 响应是事件流格式，需要按行解析
	lines := strings.Split(string(responseData), "\n")

	// 存储process_completed事件消息
	var result *parsepb.ProcessCompletedEventMessage
	// 处理每一行数据
	for _, line := range lines {
		// 跳过空行
		if line == "" {
			continue
		}
		// 去除前缀"data: "
		if strings.HasPrefix(line, "data: ") {
			line = strings.TrimPrefix(line, "data: ")
		} else {
			continue // 不是数据行，跳过
		}
		if len(line) > 300 {
			s.logger.Info("Mineru 数据行", "line", line[:300])
		} else {
			s.logger.Info("Mineru 数据行", "line", line)
		}
		// 先解析消息类型
		var msgObj struct {
			Msg string `json:"msg"`
		}
		if err := json.Unmarshal([]byte(line), &msgObj); err != nil {
			continue // 跳过解析失败的行
		}
		// 这里的消息类型有四种 estimation,process_starts,process_completed,close_stream
		// 在处理完 process_completed 消息后
		if msgObj.Msg == "process_completed" {
			// 打印出来line的前100个字符
			s.logger.Info("Mineru 获取解析数据成功状态《process_completed》")
			// 先创建一个临时结构体来解析复杂的 JSON 结构
			var tempResult TempDataResult
			if err := json.Unmarshal([]byte(line), &tempResult); err != nil {
				s.logger.Error("解析 process_completed 消息失败", "error", err)
				continue
			}

			// 创建最终的 ProcessCompletedEventMessage 结构
			result = &parsepb.ProcessCompletedEventMessage{
				Msg:     tempResult.Msg,
				EventId: tempResult.EventId,
				Success: tempResult.Success,
				Title:   tempResult.Title,
			}
			s.logger.Info("解析数据长度", "length", len(tempResult.Output.Data))
			// 处理 Output 字段
			if len(tempResult.Output.Data) > 0 {
				// 创建字符串数组
				stringData := make([]string, 0, len(tempResult.Output.Data))
				// 将每个 json.RawMessage 转换为字符串
				for _, rawMsg := range tempResult.Output.Data {
					// 将 json.RawMessage 转换为字符串
					stringData = append(stringData, string(rawMsg))
				}
				// 设置到 result 中
				result.Output = &parsepb.EventOutput{
					Data: stringData,
				}
			}
		}
		// 检查是否是完成消息
		if msgObj.Msg == "close_stream" {
			break
		}
	}

	return result, nil
}

// 检测markdown文件的语言
func (s *ParseOperateService) DetectMarkdownLanguage(ctx context.Context, pdfSHA256 string) string {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "ParseOperateService.DetectMarkdownLanguage")
	defer span.Finish()
	// 读取文件内容
	content, err := s.GetMineruMarkdownFile(ctx, pdfSHA256)
	if err != nil {
		s.logger.Warn(ctx, "Failed to get markdown file: %v", err)
		return constant.LanguageEnUS
	}

	// 将内容转换为字符串
	fullText := string(content)

	// 匹配带有数字编号的章节标题（如 # 1 或 # 1.1）
	numberedSectionRegex := regexp.MustCompile(`(?m)^#\s+(\d+(?:\.\d+)*)[.\s]+(.+)$`)
	matches := numberedSectionRegex.FindAllStringIndex(fullText, -1)

	var extractedText strings.Builder
	var titleText strings.Builder

	// 如果找到了编号的章节标题
	if len(matches) > 0 {
		// 查找所有章节标题（包括非编号的，如"参考文献"）
		allSectionRegex := regexp.MustCompile(`(?m)^#\s+([^0-9].+)$`)
		nonNumberedMatches := allSectionRegex.FindAllStringIndex(fullText, -1)
		// 确定提取内容的结束位置
		endPos := len(fullText)
		// 找到第一个非编号章节标题的位置作为结束位置
		for _, match := range nonNumberedMatches {
			// 确保这个非编号章节在第一个编号章节之后
			if match[0] > matches[0][0] {
				endPos = match[0]
				break
			}
		}
		// 提取从第一个编号标题开始到第一个非编号标题之前的内容
		extractedText.WriteString(fullText[matches[0][0]:endPos])

		// 单独提取所有标题行
		for _, match := range matches {
			titleLine := strings.TrimSpace(fullText[match[0]:match[1]])
			titleText.WriteString(titleLine)
		}
	}
	if titleText.Len() == 0 && extractedText.Len() == 0 {
		return util.DetectLanguageFallback(fullText)
	}

	return util.DetectLanguage(titleText.String(), extractedText.String())
}

func (s *ParseOperateService) AuthorListToAuthorEditeds(authorList []*parsepb.Author) []*docpb.AuthorInfo {
	if authorList == nil {
		return nil
	}
	var authorEditeds []*docpb.AuthorInfo
	for _, author := range authorList {
		authorEdited := &docpb.AuthorInfo{
			Literal:          author.FullName,
			Family:           &author.GivenName,
			Given:            &author.Surname,
			IsAuthentication: false,
		}
		authorEditeds = append(authorEditeds, authorEdited)
	}
	return authorEditeds
}
