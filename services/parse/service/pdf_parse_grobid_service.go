package service

import (
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/config"
	"github.com/yb2020/odoc/pkg/http_client"
	"github.com/yb2020/odoc/pkg/logging"
	pb "github.com/yb2020/odoc/proto/gen/go/parsed"
	grobidUtil "github.com/yb2020/odoc/services/parse/util/grobid"
)

// GrobidPDFParseService 实现PDF解析服务
type GrobidPDFParseService struct {
	config              *config.Config
	tracer              opentracing.Tracer
	logger              logging.Logger
	httpClient          http_client.HttpClient
	parseOperateService *ParseOperateService
}

// NewGrobidPDFParseService 创建新的PDF解析服务实例
func NewGrobidPDFParseService(
	config *config.Config,
	logger logging.Logger,
	tracer opentracing.Tracer,
	httpClient http_client.HttpClient,
	parseOperateService *ParseOperateService,
) *GrobidPDFParseService {
	return &GrobidPDFParseService{
		config:              config,
		logger:              logger,
		tracer:              tracer,
		httpClient:          httpClient,
		parseOperateService: parseOperateService,
	}
}

// ParsePDFHeaderWithGrobid 使用Grobid解析PDF的Header信息
func (s *GrobidPDFParseService) ParsePDFHeaderWithGrobid(ctx context.Context, pdfContent []byte) (*pb.DocumentHeader, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.ParsePDFHeaderWithGrobid")
	defer span.Finish()

	responseData, err := s.parseOperateService.GetGrobidParsedHeaderByte(ctx, pdfContent)
	if err != nil {
		s.logger.Error("parse document header failed", "error", err)
		return nil, err
	}
	parser := grobidUtil.NewGrobidParser(grobidUtil.GrobidParserV8)
	parseDocumentHeader, err := parser.ParseDocumentHeaderXml(responseData)
	if err != nil {
		s.logger.Error("parse document header failed", "error", err)
		return nil, err
	}
	return parseDocumentHeader, nil
}

// ParsePDFFulltextWithGrobidEn 使用Grobid解析英文PDF的全文信息
func (s *GrobidPDFParseService) ParsePDFFulltextWithGrobidEn(ctx context.Context, pdfContent []byte) (*pb.DocumentMetadata, *pb.FullDocument, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.ParsePDFFulltextWithGrobidEn")
	defer span.Finish()

	responseData, err := s.parseOperateService.GetGrobidParsedDocumentByte(ctx, pdfContent)
	if err != nil {
		return nil, nil, err
	}
	// // TODO：本地测试 写入grobid.xml
	// err = s.writeGrobidXmlToLocalFile(ctx, responseData)
	// if err != nil {
	// 	return nil, nil, err
	// }
	// 创建Grobid解析器
	parser := grobidUtil.NewGrobidParser(grobidUtil.GrobidParserV8)
	// 解析Grobid响应数据
	metadataDocument, fullDocument, err := parser.ParseDocument(responseData)
	if err != nil {
		return nil, nil, err
	}
	return metadataDocument, fullDocument, nil
}

// =============================== mock 接口 ===============================

// 这个是自己部署的grobid 0.8返回的fulltext xml
func (s *GrobidPDFParseService) MockParsePDFFulltextXmlWithGrobid(ctx context.Context) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.MockParsePDFFulltextXmlWithGrobid")
	defer span.Finish()

	inputFile := "C:/Users/IDEA/Desktop/parsedPdf/fulltextXml.txt"
	// 打开输入文件
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	// 读取文件内容到字节数组
	xmlContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return xmlContent, nil
}

func (s *GrobidPDFParseService) ParsePDFFulltextWithGrobidEnMock(ctx context.Context, pdfContent []byte) (*pb.DocumentMetadata, *pb.FullDocument, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "PDFParseService.ParsePDFFulltextWithGrobidEn")
	defer span.Finish()

	//TODO: 由于这个接口请求响应过慢，这里先mock
	responseData, err := s.MockParsePDFFulltextXmlWithGrobid(ctx)
	if err != nil {
		return nil, nil, err
	}
	// 创建Grobid解析器
	parser := grobidUtil.NewGrobidParser(grobidUtil.GrobidParserV8)
	// 解析Grobid响应数据
	metadataDocument, fullDocument, err := parser.ParseDocument(responseData)
	if err != nil {
		return nil, nil, err
	}
	return metadataDocument, fullDocument, nil
}

func (s *GrobidPDFParseService) writeGrobidXmlToLocalFile(ctx context.Context, xmlContent []byte) error {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "GrobidPDFParseService.writeGrobidXmlToLocalFile")
	defer span.Finish()

	//写入grobid.xml
	filePath := filepath.Join(s.config.PDF.Download.TempDirectory, "grobid.xml")

	os.WriteFile(filePath, xmlContent, 0644)

	return nil
}
