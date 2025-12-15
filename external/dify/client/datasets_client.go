package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/external/dify/constant"
	"github.com/yb2020/odoc/external/dify/proto"
	"github.com/yb2020/odoc/external/dify/transport"
	"github.com/yb2020/odoc/pkg/logging"
)

// DifyDatasetsClient Dify数据集客户端接口
type DifyDatasetsClient interface {
	DifyBaseClient

	// CreateDataset 创建数据集
	CreateDataset(ctx context.Context, request *proto.DatasetRequest) (*proto.DatasetResponse, error)

	// GetDatasets 获取数据集列表
	GetDatasets(ctx context.Context, page int, limit int) (*proto.DatasetsResponse, error)

	// GetDataset 获取数据集详情
	GetDataset(ctx context.Context, datasetID string) (*proto.DatasetResponse, error)

	// UpdateDataset 更新数据集
	UpdateDataset(ctx context.Context, datasetID string, request *proto.DatasetRequest) (*proto.DatasetResponse, error)

	// DeleteDataset 删除数据集
	DeleteDataset(ctx context.Context, datasetID string) (*proto.SimpleResponse, error)

	// 查询dataset元数据列表
	GetDatasetsMetadataList(ctx context.Context, datasetID string) (*proto.DocumentMetadataListResponse, error)

	// GetDocuments 获取文档列表
	GetDocuments(ctx context.Context, datasetID string, keyword string, page int, limit int) (*proto.DocumentsResponse, error)

	// GetDocument 获取文档详情
	GetDocument(ctx context.Context, datasetID string, documentID string) (*proto.DocumentResponse, error)

	// GetDocumentMetadata 获取文档元数据
	GetDocumentMetadata(ctx context.Context, datasetID string, documentID string) (*proto.DocumentMetadataListResponse, error)

	// BuiltInDatasetMetadata 操作内置元数据
	BuiltInDatasetMetadata(ctx context.Context, datasetID string, action string) error

	// UpdateMetadata 更新数据集元数据
	UpdateDatasetMetadata(ctx context.Context, datasetID, metadataID string, request *proto.UpdateDatasetMetadataRequest) (*proto.MetadataResponse, error)

	// CreateMetadata 创建数据集元数据
	CreateDatasetMetadata(ctx context.Context, datasetID string, request *proto.CreateMetadataRequest) (*proto.MetadataResponse, error)

	// DeleteMetadata 删除数据集元数据
	DeleteDatasetMetadata(ctx context.Context, datasetID string, metadataID string) (*proto.SimpleResponse, error)
	// DeleteDocument 删除文档
	DeleteDocument(ctx context.Context, datasetID string, documentID string) (*proto.SimpleResponse, error)

	// CreateDocumentByText 通过文本创建文档
	CreateDocumentByText(ctx context.Context, datasetID string, request *proto.DocumentCreateByTextRequest) (*proto.DocumentResponse, error)

	// CreateDocumentByFile 通过文件创建文档
	CreateDocumentByFile(ctx context.Context, datasetID string, request *proto.DocumentCreateByFileRequest, filePath string) (*proto.DocumentResponse, error)

	// UpdateDocumentByText 通过文本更新文档
	UpdateDocumentByText(ctx context.Context, datasetID string, documentID string, request *proto.DocumentUpdateByTextRequest) (*proto.DocumentResponse, error)

	// UpdateDocumentByFile 通过文件更新文档
	UpdateDocumentByFile(ctx context.Context, datasetID string, documentID string, request *proto.DocumentUpdateByFileRequest, filePath string) (*proto.DocumentResponse, error)

	// UpdateDocumentMetadata 更新文档元数据【批量赋值】
	UpdateDocumentMetadata(ctx context.Context, datasetID string, request *proto.UpdateDocumentMetadataRequest) error

	// DeleteDocumentMetadata 删除文档元数据
	DeleteDocumentMetadata(ctx context.Context, datasetID string, documentID string) (*proto.SimpleResponse, error)

	// GetIndexingStatus 获取文档嵌入状态
	GetIndexingStatus(ctx context.Context, datasetID string, batch string) (*proto.IndexingStatusResponse, error)

	// CreateSegments 新增文档分段
	CreateSegments(ctx context.Context, datasetID string, documentID string, request *proto.CreateSegmentsRequest) (*proto.SegmentListResponse, error)

	// GetSegments 查询文档分段
	GetSegments(ctx context.Context, datasetID string, documentID string, keyword string, status string, page int, limit int) (*proto.SegmentListResponse, error)

	// DeleteSegment 删除文档分段
	DeleteSegment(ctx context.Context, datasetID string, documentID string, segmentID string) (*proto.SimpleResponse, error)

	// UpdateSegment 更新文档分段
	UpdateSegment(ctx context.Context, datasetID string, documentID string, segmentID string, request *proto.UpdateSegmentRequest) (*proto.SegmentResponse, error)

	// CreateChildChunk 创建子分段
	CreateChildChunk(ctx context.Context, datasetID string, documentID string, segmentID string, request *proto.SaveChildChunkRequest) (*proto.ChildChunkResponse, error)

	// GetChildChunks 获取子分段
	GetChildChunks(ctx context.Context, datasetID string, documentID string, segmentID string, keyword string, page int, limit int) (*proto.ChildChunkListResponse, error)

	// DeleteChildChunk 删除子分段
	DeleteChildChunk(ctx context.Context, datasetID string, documentID string, segmentID string, childChunkID string) (*proto.SimpleResponse, error)

	// UpdateChildChunk 更新子分段
	UpdateChildChunk(ctx context.Context, datasetID string, documentID string, segmentID string, childChunkID string, request *proto.SaveChildChunkRequest) (*proto.ChildChunkResponse, error)

	//getEmbeddingModelList 获取Embedding模型列表
	GetEmbeddingModelList(ctx context.Context) (*proto.EmbeddingModelListResponse, error)
}

// DifyDatasetClientImpl Dify数据集客户端实现
type DifyDatasetClientImpl struct {
	DifyBaseClientImpl
	logger logging.Logger
	tracer opentracing.Tracer
}

// NewDifyDatasetsClient 创建新的Dify数据集客户端
func NewDifyDatasetsClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, timeout int, responseHeaderTimeout int) DifyDatasetsClient {
	return &DifyDatasetClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: transport.NewDefaultHttpClient(logger, timeout, responseHeaderTimeout),
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// NewDifyDatasetsClientWithHttpClient 使用自定义HTTP客户端创建新的Dify数据集客户端
func NewDifyDatasetsClientWithHttpClient(logger logging.Logger, tracer opentracing.Tracer, apiBaseUrl, apiKey string, httpClient transport.HttpClient) DifyDatasetsClient {
	return &DifyDatasetClientImpl{
		DifyBaseClientImpl: DifyBaseClientImpl{
			apiKey:     apiKey,
			apiBaseUrl: apiBaseUrl,
			httpClient: httpClient,
			tracer:     tracer,
		},
		logger: logger,
		tracer: tracer,
	}
}

// CreateDataset 创建数据集
func (c *DifyDatasetClientImpl) CreateDataset(ctx context.Context, request *proto.DatasetRequest) (*proto.DatasetResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateDataset")
	defer span.Finish()

	if request == nil {
		c.logger.Error("请求不能为空")
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH)

	// 执行请求并解析响应
	var response proto.DatasetResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDatasets 获取数据集列表
func (c *DifyDatasetClientImpl) GetDatasets(ctx context.Context, page int, limit int) (*proto.DatasetsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetDatasets")
	defer span.Finish()

	// 构建查询参数
	queryParams := make(map[string]string)
	if page > 0 {
		queryParams["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		queryParams["limit"] = strconv.Itoa(limit)
	}

	// 执行请求并解析响应
	var response proto.DatasetsResponse
	err := c.getWithParams(ctx, constant.DATASETS_PATH, queryParams, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDataset 获取数据集详情
func (c *DifyDatasetClientImpl) GetDataset(ctx context.Context, datasetID string) (*proto.DatasetResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetDataset")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID)

	// 执行请求并解析响应
	var response proto.DatasetResponse
	err := c.doJSONRequest(ctx, http.MethodGet, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateDataset 更新数据集
func (c *DifyDatasetClientImpl) UpdateDataset(ctx context.Context, datasetID string, request *proto.DatasetRequest) (*proto.DatasetResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateDataset")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID)

	// 执行请求并解析响应
	var response proto.DatasetResponse
	err := c.doJSONRequest(ctx, http.MethodPatch, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteDataset 删除数据集
func (c *DifyDatasetClientImpl) DeleteDataset(ctx context.Context, datasetID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.DeleteDataset")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID)

	// 执行请求并解析响应
	var response proto.SimpleResponse
	err := c.doJSONRequest(ctx, http.MethodDelete, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDatasetsMetadataList 获取数据集元数据列表
func (c *DifyDatasetClientImpl) GetDatasetsMetadataList(ctx context.Context, datasetID string) (*proto.DocumentMetadataListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetDatasetsMetadataList")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.METADATA_PATH)

	// 执行请求并解析响应
	var response proto.DocumentMetadataListResponse
	err := c.doJSONRequest(ctx, http.MethodGet, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateDocument 创建文档
func (c *DifyDatasetClientImpl) CreateDocument(ctx context.Context, datasetID string, request *proto.DocumentCreateRequest) (*proto.DocumentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateDocument")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH)

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// 执行请求并解析响应
	var response proto.DocumentResponse
	err = c.doJSONRequest(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData), &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDocuments 获取文档列表
func (c *DifyDatasetClientImpl) GetDocuments(ctx context.Context, datasetID string, keyword string, page int, limit int) (*proto.DocumentsResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetDocuments")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}

	// 构建URL
	baseURL := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH)

	// 添加查询参数
	params := url.Values{}
	if keyword != "" {
		params.Add("keyword", keyword)
	}
	if page > 0 {
		params.Add("page", strconv.Itoa(page))
	}
	if limit > 0 {
		params.Add("limit", strconv.Itoa(limit))
	}

	requestURL := baseURL
	if len(params) > 0 {
		requestURL += "?" + params.Encode()
	}

	// 执行请求并解析响应
	var response proto.DocumentsResponse
	err := c.getWithParams(ctx, requestURL, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDocumentMetadata 获取文档元数据
func (c *DifyDatasetClientImpl) GetDocumentMetadata(ctx context.Context, datasetID string, documentID string) (*proto.DocumentMetadataListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetDocumentMetadata")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.METADATA_PATH)

	// 执行请求并解析响应
	var response proto.DocumentMetadataListResponse
	err := c.getWithParams(ctx, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetDocument 获取上传文档的详情
func (c *DifyDatasetClientImpl) GetDocument(ctx context.Context, datasetID string, documentID string) (*proto.DocumentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetDocument")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.UPLOAD_FILE_PATH)

	// 执行请求并解析响应
	var response proto.DocumentResponse
	err := c.getWithParams(ctx, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteDocument 删除文档
func (c *DifyDatasetClientImpl) DeleteDocument(ctx context.Context, datasetID string, documentID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.DeleteDocument")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID)

	// 执行请求并解析响应
	err := c.doJSONRequest(ctx, http.MethodDelete, url, nil, nil)
	if err != nil {
		return nil, err
	}

	return &proto.SimpleResponse{}, nil
}

// CreateDocumentByText 通过文本创建文档
func (c *DifyDatasetClientImpl) CreateDocumentByText(ctx context.Context, datasetID string, request *proto.DocumentCreateByTextRequest) (*proto.DocumentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateDocumentByText")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENT_CREATE_BY_TEXT_PATH)

	// 执行请求并解析响应
	var response proto.DocumentResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateDocumentByFile 通过文件创建文档
func (c *DifyDatasetClientImpl) CreateDocumentByFile(ctx context.Context, datasetID string, request *proto.DocumentCreateByFileRequest, filePath string) (*proto.DocumentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateDocumentByFile")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}
	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENT_CREATE_BY_FILE_PATH)

	// 准备表单数据
	formData := map[string]interface{}{
		"data": request,
	}
	fileParams := map[string]string{
		"file": filePath,
	}

	// 使用公共方法执行多部分表单请求
	var response proto.DocumentResponse
	err := c.doMultipartRequest(ctx, http.MethodPost, url, formData, fileParams, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateDocumentByText 通过文本更新文档
func (c *DifyDatasetClientImpl) UpdateDocumentByText(ctx context.Context, datasetID string, documentID string, request *proto.DocumentUpdateByTextRequest) (*proto.DocumentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateDocumentByText")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.UPDATE_BY_TEXT_PATH)

	// 将请求转换为JSON
	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	// 执行请求
	respBody, err := c.doRequest(ctx, http.MethodPut, url, bytes.NewBuffer(jsonData), transport.ContentTypeJson, nil)
	if err != nil {
		return nil, err
	}

	// 准备响应对象
	var response proto.DocumentResponse

	// 解析响应
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateDocumentByFile 通过文件更新文档
func (c *DifyDatasetClientImpl) UpdateDocumentByFile(ctx context.Context, datasetID string, documentID string, request *proto.DocumentUpdateByFileRequest, filePath string) (*proto.DocumentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateDocumentByFile")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}
	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.UPDATE_BY_FILE_PATH)

	// 准备表单数据
	formData := map[string]interface{}{
		"data": request,
	}
	fileParams := map[string]string{
		"file": filePath,
	}

	// 使用公共方法执行多部分表单请求
	var response proto.DocumentResponse
	err := c.doMultipartRequest(ctx, http.MethodPost, url, formData, fileParams, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetIndexingStatus 获取文档索引状态
func (c *DifyDatasetClientImpl) GetIndexingStatus(ctx context.Context, datasetID string, documentID string) (*proto.IndexingStatusResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetIndexingStatus")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.INDEXING_STATUS_PATH)

	// 执行请求
	respBody, err := c.doRequest(ctx, http.MethodGet, url, nil, "", nil)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response proto.IndexingStatusResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateSegments 创建段落
func (c *DifyDatasetClientImpl) CreateSegments(ctx context.Context, datasetID string, documentID string, request *proto.CreateSegmentsRequest) (*proto.SegmentListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateSegments")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH)

	// 执行请求并解析响应
	var response proto.SegmentListResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetSegments 获取段落列表
func (c *DifyDatasetClientImpl) GetSegments(ctx context.Context, datasetID string, documentID string, keyword string, status string, page int, limit int) (*proto.SegmentListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetSegments")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}

	// 构建查询参数
	queryParams := make(map[string]string)
	if page > 0 {
		queryParams["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		queryParams["limit"] = strconv.Itoa(limit)
	}
	if keyword != "" {
		queryParams["keyword"] = keyword
	}
	if status != "" {
		queryParams["status"] = status
	}

	// 构建URL路径
	path := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH)

	// 执行请求并解析响应
	var response proto.SegmentListResponse
	err := c.getWithParams(ctx, path, queryParams, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteSegment 删除段落
func (c *DifyDatasetClientImpl) DeleteSegment(ctx context.Context, datasetID string, documentID string, segmentID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.DeleteSegment")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if segmentID == "" {
		return nil, fmt.Errorf("段落ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH, segmentID)

	// 执行请求并解析响应
	var response proto.SimpleResponse
	err := c.doJSONRequest(ctx, http.MethodDelete, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateSegment 更新段落
func (c *DifyDatasetClientImpl) UpdateSegment(ctx context.Context, datasetID string, documentID string, segmentID string, request *proto.UpdateSegmentRequest) (*proto.SegmentResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateSegment")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if segmentID == "" {
		return nil, fmt.Errorf("段落ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH, segmentID)

	// 执行请求并解析响应
	var response proto.SegmentResponse
	err := c.doJSONRequest(ctx, http.MethodPatch, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateChildChunk 创建子块
func (c *DifyDatasetClientImpl) CreateChildChunk(ctx context.Context, datasetID string, documentID string, segmentID string, request *proto.SaveChildChunkRequest) (*proto.ChildChunkResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateChildChunk")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if segmentID == "" {
		return nil, fmt.Errorf("段落ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH, segmentID, constant.CHILD_CHUNKS_PATH)

	// 执行请求并解析响应
	var response proto.ChildChunkResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetChildChunks 获取子块列表
func (c *DifyDatasetClientImpl) GetChildChunks(ctx context.Context, datasetID string, documentID string, segmentID string, keyword string, page int, limit int) (*proto.ChildChunkListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetChildChunks")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if segmentID == "" {
		return nil, fmt.Errorf("段落ID不能为空")
	}

	// 构建查询参数
	queryParams := make(map[string]string)
	if page > 0 {
		queryParams["page"] = strconv.Itoa(page)
	}
	if limit > 0 {
		queryParams["limit"] = strconv.Itoa(limit)
	}
	if keyword != "" {
		queryParams["keyword"] = keyword
	}

	// 构建URL路径
	path := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH, segmentID, constant.CHILD_CHUNKS_PATH)

	// 执行请求并解析响应
	var response proto.ChildChunkListResponse
	err := c.getWithParams(ctx, path, queryParams, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// DeleteChildChunk 删除子块
func (c *DifyDatasetClientImpl) DeleteChildChunk(ctx context.Context, datasetID string, documentID string, segmentID string, childChunkID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.DeleteChildChunk")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if segmentID == "" {
		return nil, fmt.Errorf("段落ID不能为空")
	}
	if childChunkID == "" {
		return nil, fmt.Errorf("子块ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH, segmentID, constant.CHILD_CHUNKS_PATH, childChunkID)

	// 执行请求并解析响应
	var response proto.SimpleResponse
	err := c.doJSONRequest(ctx, http.MethodDelete, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// UpdateChildChunk 更新子块
func (c *DifyDatasetClientImpl) UpdateChildChunk(ctx context.Context, datasetID string, documentID string, segmentID string, childChunkID string, request *proto.SaveChildChunkRequest) (*proto.ChildChunkResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateChildChunk")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}
	if segmentID == "" {
		return nil, fmt.Errorf("段落ID不能为空")
	}
	if childChunkID == "" {
		return nil, fmt.Errorf("子块ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.SEGMENTS_PATH, segmentID, constant.CHILD_CHUNKS_PATH, childChunkID)

	// 执行请求并解析响应
	var response proto.ChildChunkResponse
	err := c.doJSONRequest(ctx, http.MethodPut, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetEmbeddingModelList 获取Embedding模型列表
func (c *DifyDatasetClientImpl) GetEmbeddingModelList(ctx context.Context) (*proto.EmbeddingModelListResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.GetEmbeddingModelList")
	defer span.Finish()

	// 构建URL
	url := c.buildURL(constant.EMBEDDING_MODEL_TYPES_PATH)

	// 执行请求并解析响应
	var response proto.EmbeddingModelListResponse
	err := c.doJSONRequest(ctx, http.MethodGet, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *DifyDatasetClientImpl) UpdateDatasetMetadata(ctx context.Context, datasetID string, metadataID string, request *proto.UpdateDatasetMetadataRequest) (*proto.MetadataResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateDatasetMetadata")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, metadataID, constant.METADATA_PATH)

	// 执行请求并解析响应
	var response proto.MetadataResponse
	err := c.doJSONRequest(ctx, http.MethodPut, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateMetadata 创建元数据
func (c *DifyDatasetClientImpl) CreateDatasetMetadata(ctx context.Context, datasetID string, request *proto.CreateMetadataRequest) (*proto.MetadataResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.CreateDatasetMetadata")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return nil, fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.METADATA_PATH)

	// 执行请求并解析响应
	var response proto.MetadataResponse
	err := c.doJSONRequest(ctx, http.MethodPost, url, request, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *DifyDatasetClientImpl) DeleteDatasetMetadata(ctx context.Context, datasetID, metadataID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.DeleteDatasetMetadata")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.METADATA_PATH, metadataID)

	// 执行请求并解析响应
	var response proto.SimpleResponse
	err := c.doJSONRequest(ctx, http.MethodDelete, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func (c *DifyDatasetClientImpl) BuiltInDatasetMetadata(ctx context.Context, datasetID string, action string) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.BuiltInDatasetMetadata")
	defer span.Finish()

	if datasetID == "" {
		return fmt.Errorf("数据集ID不能为空")
	}
	if action == "" {
		return fmt.Errorf("操作类型不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.METADATA_PATH, "builtin")

	// 构建请求参数
	queryParams := map[string]string{"action": action}

	// 执行请求并解析响应
	var response proto.MetadataResponse
	err := c.getWithParams(ctx, url, queryParams, &response)
	if err != nil {
		return err
	}

	return nil
}

// UpdateDocumentMetadata 更新文档元数据【赋值】
func (c *DifyDatasetClientImpl) UpdateDocumentMetadata(ctx context.Context, datasetID string, request *proto.UpdateDocumentMetadataRequest) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.UpdateDocumentMetadata")
	defer span.Finish()

	if datasetID == "" {
		return fmt.Errorf("数据集ID不能为空")
	}
	if request == nil {
		return fmt.Errorf("请求不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, constant.METADATA_PATH)

	// 执行请求并解析响应
	return c.doJSONRequest(ctx, http.MethodPost, url, request, nil)
}

func (c *DifyDatasetClientImpl) DeleteDocumentMetadata(ctx context.Context, datasetID string, documentID string) (*proto.SimpleResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyDatasetClientImpl.DeleteDocumentMetadata")
	defer span.Finish()

	if datasetID == "" {
		return nil, fmt.Errorf("数据集ID不能为空")
	}
	if documentID == "" {
		return nil, fmt.Errorf("文档ID不能为空")
	}

	// 构建URL
	url := c.buildURL(constant.DATASETS_PATH, datasetID, constant.DOCUMENTS_PATH, documentID, constant.METADATA_PATH)

	// 执行请求并解析响应
	var response proto.SimpleResponse
	err := c.doJSONRequest(ctx, http.MethodDelete, url, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
