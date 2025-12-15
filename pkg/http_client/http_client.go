// pkg/http_client/http_client.go
package http_client

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yb2020/odoc/pkg/logging"
)

// HttpClient HTTP客户端接口
type HttpClient interface {
	// Get 发送GET请求并获取响应内容
	Get(url string, headers map[string]string) ([]byte, error)

	GetWithTimeout(url string, headers map[string]string, timeout time.Duration) ([]byte, error)

	// GetWithParams 发送带查询参数的GET请求并获取响应内容
	GetWithParams(url string, queryParams map[string]string, headers map[string]string) ([]byte, error)

	// Post 发送POST请求并获取响应内容
	Post(url string, body interface{}, headers map[string]string) ([]byte, error)

	PostWithTimeout(url string, body interface{}, headers map[string]string, timeout time.Duration) ([]byte, error)

	// Put 发送PUT请求并获取响应内容
	Put(url string, body interface{}, headers map[string]string) ([]byte, error)

	// Delete 发送DELETE请求并获取响应内容
	Delete(url string, headers map[string]string) ([]byte, error)

	// Patch 发送PATCH请求并获取响应内容
	Patch(url string, body interface{}, headers map[string]string) ([]byte, error)

	// Head 发送HEAD请求并获取响应头
	Head(url string, headers map[string]string) (map[string][]string, error)

	// PostForm 发送表单POST请求并获取响应内容
	PostForm(url string, formData map[string]string, headers map[string]string) ([]byte, error)

	// PostMultipartForm 发送多部分表单POST请求并获取响应内容
	PostMultipartForm(url string, formData map[string]string, files map[string]string, headers map[string]string) ([]byte, error)

	// PostMultipartFormWithFileInput 发送多部分表单POST请求并获取响应内容，支持文件流而不是文件路径
	PostMultipartFormWithFileInput(url string, formData map[string]string, fileInputStream map[string][]byte, fileNames map[string]string, headers map[string]string, timeout time.Duration) ([]byte, error)

	// PostStream 发送POST请求并获取流式响应
	PostStream(url string, body interface{}, headers map[string]string, handler func([]byte) error) error

	// GetJSEngine 获取JavaScript引擎
	GetJSEngine() *JSEngine

	// ExecuteJS 执行JavaScript代码并调用函数
	ExecuteJS(scriptPath string, functionName string, args ...interface{}) (interface{}, error)
}

// DefaultClient 默认HTTP客户端实现
type DefaultClient struct {
	client   *resty.Client
	jsEngine *JSEngine
	logger   logging.Logger
}

// HttpClientOption HTTP客户端选项
type HttpClientOption func(*DefaultClient)

// WithTimeout 设置超时选项
func WithTimeout(timeout time.Duration) HttpClientOption {
	return func(c *DefaultClient) {
		c.client.SetTimeout(timeout)
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger logging.Logger) HttpClientOption {
	return func(c *DefaultClient) {
		c.logger = logger
		c.jsEngine = NewJSEngine(logger)
	}
}

// WithRetry 设置重试选项
func WithRetry(count int, waitTime time.Duration, maxWaitTime time.Duration) HttpClientOption {
	return func(c *DefaultClient) {
		c.client.SetRetryCount(count).
			SetRetryWaitTime(waitTime).
			SetRetryMaxWaitTime(maxWaitTime)
	}
}

// NewHttpClient 创建新的HTTP客户端实例
func NewHttpClient(logger logging.Logger, options ...HttpClientOption) HttpClient {
	// 创建默认HTTP客户端
	client := resty.New().
		SetTimeout(40*time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36")

	// 创建HTTP客户端实例
	httpClient := &DefaultClient{
		client:   client,
		jsEngine: NewJSEngine(logger),
		logger:   logger,
	}

	// 应用选项
	for _, option := range options {
		option(httpClient)
	}

	return httpClient
}

// Get 发送GET请求并获取响应内容
func (c *DefaultClient) Get(url string, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 发送请求
	resp, err := req.Get(url)
	if err != nil {
		c.logger.Error("发送GET请求失败")
		return nil, fmt.Errorf("发送GET请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

func (c *DefaultClient) GetWithTimeout(url string, headers map[string]string, timeout time.Duration) ([]byte, error) {
	// 1. 克隆一个新的 resty 客户端实例
	localClient := c.client.Clone()

	// 2. 在克隆出的客户端上设置本次请求的特定超时时间
	localClient.SetTimeout(timeout)

	// 3. 使用克隆出的客户端来创建请求
	req := localClient.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 发送请求
	resp, err := req.Get(url)
	if err != nil {
		c.logger.Error("发送GET请求失败")
		return nil, fmt.Errorf("发送GET请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// GetWithParams 发送带查询参数的GET请求并获取响应内容
func (c *DefaultClient) GetWithParams(url string, queryParams map[string]string, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 添加查询参数
	if queryParams != nil {
		req.SetQueryParams(queryParams)
	}

	// 发送请求
	resp, err := req.Get(url)
	if err != nil {
		c.logger.Error("发送带参数的GET请求失败", "error", err, "url", url, "params", queryParams)
		return nil, fmt.Errorf("发送带参数的GET请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url, "params", queryParams)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// Post 发送POST请求并获取响应内容
func (c *DefaultClient) Post(url string, body interface{}, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}

	// 发送请求
	resp, err := req.Post(url)
	if err != nil {
		c.logger.Error("发送POST请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送POST请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// PostWithTimeout 发送POST请求并获取响应内容
func (c *DefaultClient) PostWithTimeout(url string, body interface{}, headers map[string]string, timeout time.Duration) ([]byte, error) {
	// 1. 克隆一个新的 resty 客户端实例
	localClient := c.client.Clone()
	// 2. 在克隆出的客户端上设置本次请求的特定超时时间
	localClient.SetTimeout(timeout)
	// 3. 使用克隆出的客户端来创建请求
	req := localClient.R()
	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}
	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}
	// 发送请求
	resp, err := req.Post(url)
	if err != nil {
		c.logger.Error("发送POST请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送POST请求失败: %w", err)
	}
	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}
	return resp.Body(), nil
}

// Put 发送PUT请求并获取响应内容
func (c *DefaultClient) Put(url string, body interface{}, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}

	// 发送请求
	resp, err := req.Put(url)
	if err != nil {
		c.logger.Error("发送PUT请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送PUT请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}
	return resp.Body(), nil
}

// Delete 发送DELETE请求并获取响应内容
func (c *DefaultClient) Delete(url string, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 发送请求
	resp, err := req.Delete(url)
	if err != nil {
		c.logger.Error("发送DELETE请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送DELETE请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// Patch 发送PATCH请求并获取响应内容
func (c *DefaultClient) Patch(url string, body interface{}, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}

	// 发送请求
	resp, err := req.Patch(url)
	if err != nil {
		c.logger.Error("发送PATCH请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送PATCH请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// Head 发送HEAD请求并获取响应头
func (c *DefaultClient) Head(url string, headers map[string]string) (map[string][]string, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 发送请求
	resp, err := req.Head(url)
	if err != nil {
		c.logger.Error("发送HEAD请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送HEAD请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Header(), nil
}

// PostForm 发送表单POST请求并获取响应内容
func (c *DefaultClient) PostForm(url string, formData map[string]string, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 设置表单数据
	if formData != nil {
		req.SetFormData(formData)
	}

	// 发送请求
	resp, err := req.Post(url)
	if err != nil {
		c.logger.Error("发送表单POST请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送表单POST请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// PostMultipartForm 发送多部分表单POST请求并获取响应内容
func (c *DefaultClient) PostMultipartForm(url string, formData map[string]string, files map[string]string, headers map[string]string) ([]byte, error) {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 设置表单数据
	if formData != nil {
		req.SetFormData(formData)
	}

	// 添加文件
	if files != nil {
		for fieldName, filePath := range files {
			req.SetFile(fieldName, filePath)
		}
	}

	// 发送请求
	resp, err := req.Post(url)
	if err != nil {
		c.logger.Error("发送多部分表单POST请求失败", "error", err, "url", url)
		return nil, fmt.Errorf("发送多部分表单POST请求失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode() != 200 {
		c.logger.Error("HTTP请求返回非200状态码", "status", resp.StatusCode(), "url", url)
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d", resp.StatusCode())
	}

	return resp.Body(), nil
}

// PostMultipartFormWithFileInput 发送多部分表单POST请求并获取响应内容，支持文件内容而不是文件路径
func (c *DefaultClient) PostMultipartFormWithFileInput(url string, formData map[string]string, fileByte map[string][]byte, fileNames map[string]string, headers map[string]string, timeout time.Duration) ([]byte, error) {
	// 1. 克隆一个新的 resty 客户端实例
	localClient := c.client.Clone()
	// 2. 在克隆出的客户端上设置本次请求的特定超时时间
	localClient.SetTimeout(timeout)
	// 3. 使用克隆出的客户端来创建请求
	req := localClient.R()
	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}
	// 设置表单数据
	if formData != nil {
		for key, value := range formData {
			// 处理多值参数，如 teiCoordinates=figure&teiCoordinates=biblStruct
			if strings.Contains(value, ",") && (key == "teiCoordinates" || key == "include" || key == "exclude") {
				values := strings.Split(value, ",")
				for _, v := range values {
					req.SetMultipartFormData(map[string]string{key: v})
				}
			} else {
				req.SetMultipartFormData(map[string]string{key: value})
			}
		}
	}

	// 添加文件内容
	if fileByte != nil {
		for fieldName, fileContent := range fileByte {
			fileName := fileNames[fieldName]
			if fileName == "" {
				fileName = fieldName
			}
			req.SetFileReader(fieldName, fileName, bytes.NewReader(fileContent))
		}
	}

	// 发送请求
	resp, err := req.Post(url)
	if err != nil {
		return nil, fmt.Errorf("发送多部分表单POST请求失败: %w", err)
	}
	// 检查状态码
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("HTTP请求返回非200状态码: %d, 响应: %s", resp.StatusCode(), string(resp.Body()))
	}

	return resp.Body(), nil
}

// PostStream 发送POST请求并获取流式响应
func (c *DefaultClient) PostStream(url string, body interface{}, headers map[string]string, handler func([]byte) error) error {
	// 创建请求
	req := c.client.R()

	// 添加请求头
	if headers != nil {
		req.SetHeaders(headers)
	}

	// 设置请求体
	if body != nil {
		req.SetBody(body)
	}

	// 设置流式处理模式 - 这是处理流式响应的关键设置
	req.SetDoNotParseResponse(true)

	// 记录请求信息
	c.logger.Debug("发送流式POST请求", "url", url)

	// 发送请求
	resp, err := req.Post(url)
	if err != nil {
		c.logger.Error("发送POST请求失败", "error", err, "url", url)
		return fmt.Errorf("发送POST请求失败: %w", err)
	}
	defer resp.RawResponse.Body.Close()

	// 检查状态码
	if resp.StatusCode() != 200 {
		// 读取错误响应体
		bodyBytes := make([]byte, 1024)
		n, _ := resp.RawResponse.Body.Read(bodyBytes)

		c.logger.Error("HTTP请求返回非200状态码",
			"status", resp.StatusCode(),
			"url", url,
			"response_body", string(bodyBytes[:n]))

		return fmt.Errorf("HTTP请求返回非200状态码: %d, 响应体: %s",
			resp.StatusCode(), string(bodyBytes[:n]))
	}

	// 处理流式响应
	var dataReceived bool
	buf := make([]byte, 1024)

	for {
		n, err := resp.RawResponse.Body.Read(buf)

		// 处理读取到的数据
		if n > 0 {
			dataReceived = true
			if err := handler(buf[:n]); err != nil {
				return fmt.Errorf("处理流式响应数据失败: %w", err)
			}
		}

		// 处理读取结束或错误
		if err != nil {
			// 正常结束
			if err.Error() == "EOF" {
				return nil
			}

			// 连接被关闭但已接收到数据，视为正常结束
			if err.Error() == "http2: response body closed" && dataReceived {
				c.logger.Warn("服务器关闭了连接，但已接收到数据")
				return nil
			}

			// 其他错误
			c.logger.Error("读取流式响应失败", "error", err)
			return fmt.Errorf("读取流式响应失败: %w", err)
		}
	}
}

// GetJSEngine 获取JavaScript引擎
func (c *DefaultClient) GetJSEngine() *JSEngine {
	return c.jsEngine
}

// ExecuteJS 执行JavaScript代码并调用函数
func (c *DefaultClient) ExecuteJS(scriptPath string, functionName string, args ...interface{}) (interface{}, error) {
	// 加载脚本
	err := c.jsEngine.LoadScriptFromFile(scriptPath)
	if err != nil {
		return nil, err
	}

	// 调用函数
	return c.jsEngine.CallFunction(functionName, args...)
}
