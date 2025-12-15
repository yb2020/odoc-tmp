// pkg/squid_proxy/squid_proxy.go
package squid_proxy

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// SquidProxyClient Squid代理客户端
type SquidProxyClient struct {
	config *SquidConfig
	client *http.Client
}

// NewSquidProxyClient 创建新的Squid代理客户端
func NewSquidProxyClient(config *SquidConfig) (*SquidProxyClient, error) {
	if config == nil {
		return nil, NewSquidError(ErrInvalidConfig, "配置不能为空")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	// 解析代理URL
	proxyURL, err := url.Parse(config.ProxyURL)
	if err != nil {
		return nil, NewSquidError(ErrInvalidProxy, fmt.Sprintf("无效的代理URL: %v", err))
	}

	// 如果配置了用户名和密码，添加到代理URL中
	// 这样Go的http.Transport会自动处理代理认证
	if config.Username != "" && config.Password != "" {
		proxyURL.User = url.UserPassword(config.Username, config.Password)
	}

	// 创建HTTP传输配置
	var transport *http.Transport
	if config.Username != "" && config.Password != "" {
		// 构造 Basic 认证头，确保在 CONNECT 请求中也携带认证
		basic := "Basic " + base64.StdEncoding.EncodeToString([]byte(config.Username+":"+config.Password))
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			ProxyConnectHeader: http.Header{
				"Proxy-Authorization": []string{basic},
			},
		}
	} else {
		transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}
	}

	// 创建HTTP客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   config.Timeout,
	}

	return &SquidProxyClient{
		config: config,
		client: client,
	}, nil
}

// DownloadFile 通过代理下载文件（主要用于PDF）
func (c *SquidProxyClient) DownloadFile(targetURL string, options *DownloadOptions) (*DownloadResult, error) {
	if targetURL == "" {
		return nil, NewSquidError(ErrInvalidURL, "目标URL不能为空")
	}

	// 验证URL格式和协议
	if err := ValidateURL(targetURL); err != nil {
		return nil, err
	}

	// 使用默认选项
	if options == nil {
		options = &DownloadOptions{
			Timeout:    c.config.Timeout,
			RetryCount: 3,
			RetryDelay: 2 * time.Second,
			MaxSize:    100 * 1024 * 1024, // 100MB默认限制
		}
	}

	var lastErr error
	retryCount := options.RetryCount
	if retryCount <= 0 {
		retryCount = 1
	}

	// 重试逻辑
	for i := 0; i < retryCount; i++ {
		result, err := c.downloadFileOnce(targetURL, options)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// 如果不是网络错误，不重试
		if !isRetryableError(err) {
			break
		}

		// 最后一次重试不需要等待
		if i < retryCount-1 && options.RetryDelay > 0 {
			time.Sleep(options.RetryDelay)
		}
	}

	return nil, lastErr
}

// downloadFileOnce 执行一次下载尝试
func (c *SquidProxyClient) downloadFileOnce(targetURL string, options *DownloadOptions) (*DownloadResult, error) {
	startTime := time.Now()

	// 创建请求
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, NewSquidError(ErrRequestFailed, fmt.Sprintf("创建请求失败: %v", err))
	}

	// 设置请求头
	c.setRequestHeaders(req, options.Headers)

	// 创建临时客户端（如果需要自定义超时）
	client := c.client
	if options.Timeout > 0 && options.Timeout != c.config.Timeout {
		transport := c.client.Transport.(*http.Transport).Clone()
		client = &http.Client{
			Transport: transport,
			Timeout:   options.Timeout,
		}
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return nil, NewSquidError(ErrRequestFailed, fmt.Sprintf("请求失败: %v", err))
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, NewSquidError(ErrHTTPError, fmt.Sprintf("HTTP错误: %d %s", resp.StatusCode, resp.Status))
	}

	// 检查内容长度
	contentLength := resp.ContentLength
	if options.MaxSize > 0 && contentLength > 0 && contentLength > options.MaxSize {
		return nil, NewSquidError(ErrFileTooLarge, fmt.Sprintf("文件过大: %d bytes，限制: %d bytes", contentLength, options.MaxSize))
	}

	// 读取响应体
	var body []byte
	if options.MaxSize > 0 {
		// 使用LimitReader防止内存溢出
		limitReader := io.LimitReader(resp.Body, options.MaxSize+1)
		body, err = io.ReadAll(limitReader)
		if err != nil {
			return nil, NewSquidError(ErrReadFailed, fmt.Sprintf("读取响应失败: %v", err))
		}

		// 检查是否超过大小限制
		if int64(len(body)) > options.MaxSize {
			return nil, NewSquidError(ErrFileTooLarge, fmt.Sprintf("文件过大，超过限制: %d bytes", options.MaxSize))
		}
	} else {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, NewSquidError(ErrReadFailed, fmt.Sprintf("读取响应失败: %v", err))
		}
	}

	duration := time.Since(startTime)

	return &DownloadResult{
		StatusCode:    resp.StatusCode,
		Headers:       resp.Header,
		Body:          body,
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: int64(len(body)),
		Duration:      duration,
	}, nil
}

// DownloadFileToPath 通过代理下载文件并保存到指定路径
func (c *SquidProxyClient) DownloadFileToPath(targetURL, savePath string, options *DownloadOptions) error {
	result, err := c.DownloadFile(targetURL, options)
	if err != nil {
		return err
	}

	// 保存文件
	err = os.WriteFile(savePath, result.Body, 0644)
	if err != nil {
		return NewSquidError(ErrWriteFailed, fmt.Sprintf("保存文件失败: %v", err))
	}

	return nil
}

// Get 通过代理发送GET请求
func (c *SquidProxyClient) Get(targetURL string, headers map[string]string) (*DownloadResult, error) {
	options := &DownloadOptions{
		Headers:    headers,
		Timeout:    c.config.Timeout,
		RetryCount: 1,
		MaxSize:    10 * 1024 * 1024, // 10MB默认限制
	}

	return c.DownloadFile(targetURL, options)
}

// GetWithStream 通过代理获取流式响应（用于大文件）
func (c *SquidProxyClient) GetWithStream(targetURL string, writer io.Writer, options *DownloadOptions) error {
	if targetURL == "" {
		return NewSquidError(ErrInvalidURL, "目标URL不能为空")
	}

	if writer == nil {
		return NewSquidError(ErrInvalidConfig, "writer不能为空")
	}

	// 使用默认选项
	if options == nil {
		options = &DownloadOptions{
			Timeout:    c.config.Timeout,
			RetryCount: 1,
		}
	}

	// 创建请求
	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return NewSquidError(ErrRequestFailed, fmt.Sprintf("创建请求失败: %v", err))
	}

	// 设置请求头
	c.setRequestHeaders(req, options.Headers)

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return NewSquidError(ErrRequestFailed, fmt.Sprintf("请求失败: %v", err))
	}
	defer resp.Body.Close()

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return NewSquidError(ErrHTTPError, fmt.Sprintf("HTTP错误: %d %s", resp.StatusCode, resp.Status))
	}

	// 流式复制数据
	var reader io.Reader = resp.Body
	if options.MaxSize > 0 {
		reader = io.LimitReader(resp.Body, options.MaxSize)
	}

	_, err = io.Copy(writer, reader)
	if err != nil {
		return NewSquidError(ErrWriteFailed, fmt.Sprintf("写入数据失败: %v", err))
	}

	return nil
}

// TestConnection 测试代理连接
func (c *SquidProxyClient) TestConnection() error {
	// 使用一个简单的HTTP请求测试连接
	testURL := "http://httpbin.org/ip"

	req, err := http.NewRequest("GET", testURL, nil)
	if err != nil {
		return NewSquidError(ErrRequestFailed, fmt.Sprintf("创建测试请求失败: %v", err))
	}

	// 设置较短的超时时间用于测试
	client := &http.Client{
		Transport: c.client.Transport,
		Timeout:   10 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return NewSquidError(ErrProxyConnection, fmt.Sprintf("代理连接测试失败: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return NewSquidError(ErrProxyConnection, fmt.Sprintf("代理连接测试返回错误状态: %d", resp.StatusCode))
	}

	return nil
}

// setRequestHeaders 设置请求头
func (c *SquidProxyClient) setRequestHeaders(req *http.Request, customHeaders map[string]string) {
	// 设置默认User-Agent
	if c.config.UserAgent != "" {
		req.Header.Set("User-Agent", c.config.UserAgent)
	}

	// 注意：代理认证已经在Transport层面通过proxyURL.User处理
	// 不需要手动设置Proxy-Authorization头

	// 设置默认Accept头（针对PDF优化）
	if req.Header.Get("Accept") == "" {
		if strings.Contains(req.URL.Path, ".pdf") || strings.Contains(req.URL.String(), "pdf") {
			req.Header.Set("Accept", "application/pdf,*/*")
		} else {
			req.Header.Set("Accept", "*/*")
		}
	}

	// 设置自定义请求头
	for key, value := range customHeaders {
		req.Header.Set(key, value)
	}
}

// // basicAuth 生成基本认证头
// func basicAuth(username, password string) string {
// 	auth := username + ":" + password
// 	return "Basic " + base64Encode(auth)
// }

// // base64Encode 简单的base64编码
// func base64Encode(data string) string {
// 	const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

// 	input := []byte(data)
// 	output := make([]byte, ((len(input)+2)/3)*4)

// 	for i, j := 0, 0; i < len(input); i, j = i+3, j+4 {
// 		a, b, c := input[i], byte(0), byte(0)

// 		if i+1 < len(input) {
// 			b = input[i+1]
// 		}
// 		if i+2 < len(input) {
// 			c = input[i+2]
// 		}

// 		bitmap := (uint32(a) << 16) | (uint32(b) << 8) | uint32(c)

// 		output[j] = base64Table[(bitmap>>18)&63]
// 		output[j+1] = base64Table[(bitmap>>12)&63]

// 		if i+1 < len(input) {
// 			output[j+2] = base64Table[(bitmap>>6)&63]
// 		} else {
// 			output[j+2] = '='
// 		}

// 		if i+2 < len(input) {
// 			output[j+3] = base64Table[bitmap&63]
// 		} else {
// 			output[j+3] = '='
// 		}
// 	}

// 	return string(output)
// }

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	squidErr, ok := err.(*SquidError)
	if !ok {
		return true // 未知错误，尝试重试
	}

	switch squidErr.Code {
	case ErrProxyConnection, ErrRequestFailed:
		return true
	case ErrInvalidURL, ErrInvalidConfig, ErrFileTooLarge:
		return false
	default:
		return true
	}
}
