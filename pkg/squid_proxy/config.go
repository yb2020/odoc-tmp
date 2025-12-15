// pkg/squid_proxy/config.go
package squid_proxy

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// SquidConfig Squid代理配置
type SquidConfig struct {
	ProxyURL  string        // 代理服务器URL，如: "http://proxy.example.com:3128"
	Username  string        // 代理认证用户名（可选）
	Password  string        // 代理认证密码（可选）
	Timeout   time.Duration // 请求超时时间
	UserAgent string        // 用户代理字符串
}

// DownloadResult 下载结果
type DownloadResult struct {
	StatusCode    int           // HTTP状态码
	Headers       http.Header   // 响应头
	Body          []byte        // 响应体
	ContentType   string        // 内容类型
	ContentLength int64         // 内容长度
	Duration      time.Duration // 请求耗时
}

// DownloadedPdfInfo 封装了从URL下载的PDF文件的核心信息
type DownloadedPdfInfo struct {
	FileReader *bytes.Reader // 文件内容读取器
	FileName   string        // 文件名
	SHA256     string        // SHA256值(16进制)
	Size       int64         // 文件大小(字节)
	PageCount  int           // PDF页数
}

// DownloadOptions 下载选项
type DownloadOptions struct {
	Headers    map[string]string // 自定义请求头
	MaxSize    int64             // 最大文件大小限制（字节）
	Timeout    time.Duration     // 单次请求超时
	RetryCount int               // 重试次数
	RetryDelay time.Duration     // 重试间隔
}

// Validate 验证配置有效性
func (c *SquidConfig) Validate() error {
	if c.ProxyURL == "" {
		return NewSquidError(ErrInvalidConfig, "代理URL不能为空")
	}

	// 验证URL格式
	proxyURL, err := url.Parse(c.ProxyURL)
	if err != nil {
		return NewSquidError(ErrInvalidConfig, fmt.Sprintf("无效的代理URL格式: %v", err))
	}

	// 检查协议
	if proxyURL.Scheme != "http" && proxyURL.Scheme != "https" {
		return NewSquidError(ErrInvalidConfig, "代理URL必须使用http或https协议")
	}

	// 检查主机名
	if proxyURL.Host == "" {
		return NewSquidError(ErrInvalidConfig, "代理URL必须包含主机名和端口")
	}

	// 设置默认超时时间
	if c.Timeout <= 0 {
		c.Timeout = 30 * time.Second
	}

	// 设置默认User-Agent
	if c.UserAgent == "" {
		c.UserAgent = "SquidProxy/1.0"
	}

	return nil
}

// DefaultConfig 创建默认配置
func DefaultConfig(proxyURL string) *SquidConfig {
	return &SquidConfig{
		ProxyURL:  proxyURL,
		Timeout:   30 * time.Second,
		UserAgent: "SquidProxy/1.0",
	}
}

// DefaultDownloadOptions 创建默认下载选项
func DefaultDownloadOptions() *DownloadOptions {
	return &DownloadOptions{
		MaxSize:    100 * 1024 * 1024, // 100MB
		Timeout:    60 * time.Second,
		RetryCount: 3,
		RetryDelay: 2 * time.Second,
		Headers:    make(map[string]string),
	}
}

// PDFDownloadOptions 创建PDF下载优化选项
func PDFDownloadOptions() *DownloadOptions {
	options := DefaultDownloadOptions()
	options.Headers["Accept"] = "application/pdf,*/*"
	options.Headers["Accept-Language"] = "en-US,en;q=0.9"
	options.MaxSize = 50 * 1024 * 1024 // 50MB，PDF文件通常较小
	return options
}

// LargeFileDownloadOptions 创建大文件下载选项
func LargeFileDownloadOptions() *DownloadOptions {
	options := DefaultDownloadOptions()
	options.MaxSize = 500 * 1024 * 1024 // 500MB
	options.Timeout = 300 * time.Second // 5分钟
	options.RetryCount = 2              // 大文件重试次数少一些
	options.RetryDelay = 5 * time.Second
	return options
}
