package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/yb2020/odoc/external/dify/constant"
	"github.com/yb2020/odoc/external/dify/proto"
	"github.com/yb2020/odoc/external/dify/transport"
	"github.com/yb2020/odoc/pkg/logging"
)

// DifyBaseClient Dify基础客户端接口
type DifyBaseClient interface {
	// GetAppInfo 获取应用信息
	GetAppInfo(ctx context.Context) (*proto.AppInfoResponse, error)
}

// DifyBaseClientImpl Dify基础客户端实现
type DifyBaseClientImpl struct {
	apiKey     string
	apiBaseUrl string
	httpClient transport.HttpClient
	tracer     opentracing.Tracer
}

// Generic HTTP request helper methods

// doRequest 执行HTTP请求并返回响应体
func (c *DifyBaseClientImpl) doRequest(ctx context.Context, method, url string, body io.Reader, contentType string, headers map[string]string) ([]byte, error) {
	// 设置通用请求头
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Authorization"] = "Bearer " + c.apiKey

	// 如果提供了内容类型，则设置Content-Type头
	if contentType != "" && body != nil {
		headers["Content-Type"] = contentType
	}

	// 根据HTTP方法发送请求
	var resp *http.Response
	var err error

	switch method {
	case http.MethodGet:
		resp, err = c.httpClient.Get(ctx, url, headers)
	case http.MethodPost:
		resp, err = c.httpClient.Post(ctx, url, body, contentType, headers)
	case http.MethodPut:
		resp, err = c.httpClient.Put(ctx, url, body, contentType, headers)
	case http.MethodDelete:
		resp, err = c.httpClient.Delete(ctx, url, headers)
	case http.MethodPatch:
		resp, err = c.httpClient.Patch(ctx, url, body, contentType, headers)
	default:
		return nil, fmt.Errorf("不支持的HTTP方法: %s", method)
	}

	if err != nil {
		return nil, err
	}

	// 处理响应
	return transport.HandleResponse(resp)
}

// doJSONRequest 执行JSON请求并解析响应
func (c *DifyBaseClientImpl) doJSONRequest(ctx context.Context, method, url string, requestBody interface{}, responseObj interface{}) error {
	var body io.Reader

	// 如果提供了请求体，则处理为JSON
	if requestBody != nil {
		switch data := requestBody.(type) {
		case []byte:
			// 如果是已经序列化的JSON数据，直接使用
			body = bytes.NewBuffer(data)
		default:
			// 否则进行序列化
			jsonData, err := json.Marshal(requestBody)
			if err != nil {
				return err
			}
			body = bytes.NewBuffer(jsonData)
		}
	}

	// 执行请求
	respBody, err := c.doRequest(ctx, method, url, body, transport.ContentTypeJson, nil)
	if err != nil {
		return err
	}

	// 如果提供了响应对象，则将响应体反序列化为该对象
	if responseObj != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, responseObj); err != nil {
			return err
		}
	}

	return nil
}

// buildURL 构建完整的URL
func (c *DifyBaseClientImpl) buildURL(path string, params ...string) string {
	// 构建基本URL
	url := transport.BuildUrl(c.apiBaseUrl, path)

	// 如果提供了额外的路径参数，则添加到URL中
	if len(params) > 0 {
		for _, param := range params {
			if param != "" {
				// 处理参数中的斜杠，避免出现双斜杠
				cleanParam := strings.TrimPrefix(param, "/")

				if !strings.HasSuffix(url, "/") {
					url += "/"
				}
				url += cleanParam
			}
		}
	}

	return url
}

// getWithParams 执行带查询参数的GET请求
func (c *DifyBaseClientImpl) getWithParams(ctx context.Context, path string, queryParams map[string]string, responseObj interface{}) error {
	// 构建URL
	url := c.buildURL(path)

	// 添加查询参数
	if len(queryParams) > 0 {
		query := url
		if !strings.Contains(query, "?") {
			query += "?"
		} else if !strings.HasSuffix(query, "?") && !strings.HasSuffix(query, "&") {
			query += "&"
		}

		params := make([]string, 0, len(queryParams))
		for key, value := range queryParams {
			params = append(params, fmt.Sprintf("%s=%s", neturl.QueryEscape(key), neturl.QueryEscape(value)))
		}
		query += strings.Join(params, "&")
		url = query
	}

	// 执行请求
	respBody, err := c.doRequest(ctx, http.MethodGet, url, nil, "", nil)
	if err != nil {
		return err
	}

	// 解析响应
	if responseObj != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, responseObj); err != nil {
			return err
		}
	}

	return nil
}

// NewDifyBaseClient 创建新的Dify基础客户端
func NewDifyBaseClient(logger logging.Logger, tracer opentracing.Tracer, apiKey, apiBaseUrl string) DifyBaseClient {
	return &DifyBaseClientImpl{
		apiKey:     apiKey,
		apiBaseUrl: apiBaseUrl,
		httpClient: transport.NewDefaultHttpClient(logger, 5, 120),
		tracer:     tracer,
	}
}

// NewDifyBaseClientWithHttpClient 使用自定义HTTP客户端创建新的Dify基础客户端
func NewDifyBaseClientWithHttpClient(apiKey, apiBaseUrl string, httpClient transport.HttpClient, tracer opentracing.Tracer) DifyBaseClient {
	return &DifyBaseClientImpl{
		apiKey:     apiKey,
		apiBaseUrl: apiBaseUrl,
		httpClient: httpClient,
		tracer:     tracer,
	}
}

// GetAppInfo 获取应用信息
func (c *DifyBaseClientImpl) GetAppInfo(ctx context.Context) (*proto.AppInfoResponse, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyBaseClientImpl.GetAppInfo")
	defer span.Finish()

	// 构建URL
	url := c.buildURL(constant.PARAMETERS_PATH)

	// 执行请求
	respBody, err := c.doRequest(ctx, http.MethodGet, url, nil, "", nil)
	if err != nil {
		return nil, err
	}

	// 解析响应
	var response proto.AppInfoResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// doMultipartRequest 执行多部分表单请求
func (c *DifyBaseClientImpl) doMultipartRequest(ctx context.Context, method, url string, formData map[string]interface{}, fileParams map[string]string, responseObj interface{}) error {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, c.tracer, "DifyBaseClientImpl.doMultipartRequest")
	defer span.Finish()

	// 创建multipart表单
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 处理数据字段
	if data, ok := formData["data"]; ok {
		// 将data序列化为JSON
		jsonData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		// 写入data字段
		if err := writer.WriteField("data", string(jsonData)); err != nil {
			return err
		}
	}

	// 处理其他表单字段
	for key, value := range formData {
		if key != "data" { // 跳过已处理的data字段
			strValue := fmt.Sprintf("%v", value)
			if err := writer.WriteField(key, strValue); err != nil {
				return err
			}
		}
	}

	// 添加文件
	for fieldName, filePath := range fileParams {
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 创建文件表单字段
		part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			return err
		}

		// 复制文件内容到表单字段
		if _, err = io.Copy(part, file); err != nil {
			return err
		}
	}

	// 关闭writer以完成multipart消息
	if err := writer.Close(); err != nil {
		return err
	}

	// 执行请求
	respBody, err := c.doRequest(ctx, method, url, body, writer.FormDataContentType(), nil)
	if err != nil {
		return err
	}

	// 解析响应
	if responseObj != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, responseObj); err != nil {
			return err
		}
	}

	return nil
}
