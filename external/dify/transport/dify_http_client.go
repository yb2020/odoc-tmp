package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/yb2020/odoc/external/dify/exception"
	"github.com/yb2020/odoc/pkg/logging"
)

const (
	// ContentTypeJson JSON内容类型
	ContentTypeJson = "application/json"
	// ContentTypeForm 表单内容类型
	ContentTypeForm = "application/x-www-form-urlencoded"
	// ContentTypeMultipart 多部分表单内容类型
	ContentTypeMultipart = "multipart/form-data"
)

// HttpClient HTTP客户端接口
type HttpClient interface {
	// Get 发送GET请求
	Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	// Post 发送POST请求
	Post(ctx context.Context, url string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	// Put 发送PUT请求
	Put(ctx context.Context, url string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	// Patch 发送PATCH请求
	Patch(ctx context.Context, url string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error)
	// Delete 发送DELETE请求
	Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error)
	// PostForm 发送表单POST请求
	PostForm(ctx context.Context, url string, data url.Values, headers map[string]string) (*http.Response, error)
	// PostMultipart 发送多部分表单POST请求
	PostMultipart(ctx context.Context, url string, params map[string]string, fileParams map[string]string, headers map[string]string) (*http.Response, error)
	// PostStream 发送POST请求并处理流式响应
	PostStream(ctx context.Context, url string, body interface{}, headers map[string]string, handler func([]byte) error) error
}

// DefaultHttpClient 默认HTTP客户端实现
type DefaultHttpClient struct {
	client *http.Client
	logger logging.Logger
}

// NewDefaultHttpClient 创建新的默认HTTP客户端
func NewDefaultHttpClient(logger logging.Logger, timeout int, responseHeaderTimeout int) *DefaultHttpClient {
	// 自定义 transport 以实现更精细的超时控制
	transport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second, // 连接超时: 5秒
			KeepAlive: 30 * time.Second,                     // TCP Keep-Alive
		}).DialContext,
		ResponseHeaderTimeout: time.Duration(responseHeaderTimeout) * time.Second, // 响应头超时: 120秒
		TLSHandshakeTimeout:   10 * time.Second,                                   // TLS握手超时: 10秒
		MaxIdleConns:          100,                                                // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second,                                   // 空闲连接超时: 90秒
		ExpectContinueTimeout: 1 * time.Second,                                    // 100-continue超时: 1秒
	}

	return &DefaultHttpClient{
		client: &http.Client{
			Transport: transport,
		},
		logger: logger,
	}
}

// Get 发送GET请求
func (c *DefaultHttpClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	return c.client.Do(req)
}

// Post 发送POST请求
func (c *DefaultHttpClient) Post(ctx context.Context, url string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}

	// 设置内容类型
	req.Header.Set("Content-Type", contentType)

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	return c.client.Do(req)
}

// Put 发送PUT请求
func (c *DefaultHttpClient) Put(ctx context.Context, url string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}

	// 设置内容类型
	req.Header.Set("Content-Type", contentType)

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	return c.client.Do(req)
}

// Patch 发送PATCH请求
func (c *DefaultHttpClient) Patch(ctx context.Context, url string, body io.Reader, contentType string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}

	// 设置内容类型
	req.Header.Set("Content-Type", contentType)

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	return c.client.Do(req)
}

// Delete 发送DELETE请求
func (c *DefaultHttpClient) Delete(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}

	// 添加请求头
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// 发送请求
	return c.client.Do(req)
}

// PostForm 发送表单POST请求
func (c *DefaultHttpClient) PostForm(ctx context.Context, url string, data url.Values, headers map[string]string) (*http.Response, error) {
	body := strings.NewReader(data.Encode())

	// 添加表单内容类型
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = ContentTypeForm

	return c.Post(ctx, url, body, ContentTypeForm, headers)
}

// PostMultipart 发送多部分表单POST请求
func (c *DefaultHttpClient) PostMultipart(ctx context.Context, url string, params map[string]string, fileParams map[string]string, headers map[string]string) (*http.Response, error) {
	// 创建一个buffer用于存储multipart数据
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加普通参数
	for key, value := range params {
		if err := writer.WriteField(key, value); err != nil {
			return nil, err
		}
	}

	// 添加文件参数
	for key, filePath := range fileParams {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		part, err := writer.CreateFormFile(key, filepath.Base(filePath))
		if err != nil {
			return nil, err
		}

		if _, err = io.Copy(part, file); err != nil {
			return nil, err
		}
	}

	// 关闭writer以完成multipart消息
	if err := writer.Close(); err != nil {
		return nil, err
	}

	// 添加multipart内容类型
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = writer.FormDataContentType()

	return c.Post(ctx, url, body, writer.FormDataContentType(), headers)
}

// PostStream 发送POST请求并处理流式响应
func (c *DefaultHttpClient) PostStream(ctx context.Context, url string, body interface{}, headers map[string]string, handler func([]byte) error) error {
	// 这是一个简单的实现，实际上需要根据body类型进行处理
	var bodyReader io.Reader
	if body != nil {
		// 假设body是一个可以序列化为JSON的对象
		jsonBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = bytes.NewReader(jsonBytes)
	}

	// 设置内容类型
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = ContentTypeJson

	// 发送请求
	resp, err := c.Post(ctx, url, bodyReader, ContentTypeJson, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// 读取错误响应
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode)
		}
		errorCode, errorMessage, _ := ExtractErrorInfo(string(body))
		return exception.NewDifyApiException(resp.StatusCode, errorCode, errorMessage)
	}

	// 处理流式响应
	buffer := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			if err := handler(buffer[:n]); err != nil {
				return err
			}
		}
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}

// HandleResponse 处理HTTP响应
func HandleResponse(resp *http.Response) ([]byte, error) {
	if resp == nil {
		return nil, fmt.Errorf("响应为空")
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 检查状态码
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return body, nil
	}

	// 处理错误响应
	errorCode, errorMessage, _ := ExtractErrorInfo(string(body))
	return nil, exception.NewDifyApiException(resp.StatusCode, errorCode, errorMessage)
}

// BuildUrl 构建URL
func BuildUrl(baseUrl string, path string) string {
	if !strings.HasSuffix(baseUrl, "/") && !strings.HasPrefix(path, "/") {
		return baseUrl + "/" + path
	} else if strings.HasSuffix(baseUrl, "/") && strings.HasPrefix(path, "/") {
		return baseUrl + path[1:]
	} else {
		return baseUrl + path
	}
}
