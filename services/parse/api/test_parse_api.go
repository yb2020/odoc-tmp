package api

import (
	"fmt"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/pkg/response"
	"github.com/yb2020/odoc/services/parse/service"
	pdfModel "github.com/yb2020/odoc/services/pdf/model"
)

type TestParseAPI struct {
	tracer                opentracing.Tracer
	grobidPdfParseService *service.GrobidPDFParseService
	mineruPdfParseService *service.MineruPDFParseService
}

// TestParseAPI 测试解析API处理器
func NewTestParseApi(
	tracer opentracing.Tracer,
	grobidPdfParseService *service.GrobidPDFParseService,
	mineruPdfParseService *service.MineruPDFParseService,
) *TestParseAPI {
	return &TestParseAPI{
		tracer:                tracer,
		grobidPdfParseService: grobidPdfParseService,
		mineruPdfParseService: mineruPdfParseService,
	}
}

func (api *TestParseAPI) ParseFulltextV8(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TestParseAPI.ParseFulltext")
	defer span.Finish()

	content, err := api.readLocalPdfFile(c)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}
	// TODO: 实现解析逻辑
	_, fullDocument, err := api.grobidPdfParseService.ParsePDFFulltextWithGrobidEn(c.Request.Context(), content)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}

	response.Success(c, "Success", fullDocument)
}

func (api *TestParseAPI) ParseMetadataV8(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TestParseAPI.ParseFulltext")
	defer span.Finish()

	content, err := api.readLocalPdfFile(c)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}
	// TODO: 实现解析逻辑
	metadata, _, err := api.grobidPdfParseService.ParsePDFFulltextWithGrobidEn(c.Request.Context(), content)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}

	response.Success(c, "Success", metadata)
}
func (api *TestParseAPI) ParseHeaderV8(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TestParseAPI.ParseFulltext")
	defer span.Finish()

	content, err := api.readLocalPdfFile(c)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}
	metadata, err := api.grobidPdfParseService.ParsePDFHeaderWithGrobid(c.Request.Context(), content)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}

	response.Success(c, "Success", metadata)
}

func (api *TestParseAPI) ParseMineru(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TestParseAPI.ParseMineru")
	defer span.Finish()

	content, err := api.readLocalPdfFile(c)
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}
	// TODO: 实现解析逻辑
	metadata, fullDocument, imageRecords, pageBlocks, err := api.mineruPdfParseService.ExecuteParsePDF(c.Request.Context(), content, &pdfModel.PaperPdf{
		FileSHA256: "af91eb0d9382ed649e65f3d80bae456a",
		Language:   "en",
	})
	if err != nil {
		response.ErrorNoData(c, "解析失败")
		return
	}
	fmt.Println("metadata", metadata)
	fmt.Println("fullDocument", fullDocument)
	fmt.Println("imageRecords", imageRecords)
	fmt.Println("pageBlocks", pageBlocks)
	response.Success(c, "Success", metadata)
}

// 读取本地pdf文件
func (api *TestParseAPI) readLocalPdfFile(c *gin.Context) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(c.Request.Context(), api.tracer, "TestParseAPI.readLocalPdfFile")
	defer span.Finish()
	//有问题的文档
	// inputFile := "C:/Users/IDEA/Desktop/parsedPdf/上传测试pdf/readPaper2.pdf"
	//英文
	inputFile := "C:/Users/IDEA/Desktop/parsedPdf/mineruTest.pdf"
	//中文
	// inputFile := "C:/Users/IDEA/Desktop/parsedPdf/marinedrugs-10-01812.pdf"
	//日文
	// inputFile := "C:/Users/IDEA/Desktop/parsedPdf/Pre1_78.pdf"
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
