// pkg/squid_proxy/utils.go
package squid_proxy

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"math/rand"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

	unipdf "github.com/unidoc/unipdf/v3/model"
)

// GetPdfPageCount 使用unipdf库获取PDF文件的页数
func GetPdfPageCount(reader *bytes.Reader) (int, error) {
	pdfReader, err := unipdf.NewPdfReader(reader)
	if err != nil {
		return 0, fmt.Errorf("无法创建PDF读取器: %w", err)
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return 0, fmt.Errorf("无法获取PDF页数: %w", err)
	}

	// 在读取后，需要将reader的指针重置到开头，以便后续流程（如上传）可以从头读取文件内容
	_, err = reader.Seek(0, 0)
	if err != nil {
		return 0, fmt.Errorf("重置文件读取器失败: %w", err)
	}

	return numPages, nil
}

// GetFileNameFromResponse 从HTTP响应头或URL中获取文件名
// 优先从 Content-Disposition 头解析，如果失败则回退到从URL解析
func GetFileNameFromResponse(headers http.Header, fallbackURL string) string {
	// 1. 尝试从 Content-Disposition 获取文件名
	disposition := headers.Get("Content-Disposition")
	if disposition != "" {
		// 解析媒体类型，例如 "attachment; filename=example.pdf"
		_, params, err := mime.ParseMediaType(disposition)
		if err == nil {
			if filename, ok := params["filename"]; ok {
				// 如果成功获取到文件名，直接返回
				if filename != "" {
					return filename
				}
			}
		}
	}

	// 2. 如果无法从响应头获取，回退到从URL解析
	fileNameFromURL := GetFileNameFromURL(fallbackURL)

	// 3. 如果URL解析结果不含后缀，且响应类型是PDF，则强制添加.pdf后缀
	if !strings.Contains(fileNameFromURL, ".") {
		contentType := headers.Get("Content-Type")
		if strings.Contains(strings.ToLower(contentType), "application/pdf") {
			return fileNameFromURL + ".pdf"
		}
	}

	return fileNameFromURL
}

// ValidateURL 验证URL有效性
func ValidateURL(targetURL string) error {
	if targetURL == "" {
		return NewSquidError(ErrInvalidURL, "URL不能为空")
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return NewSquidError(ErrInvalidURL, fmt.Sprintf("URL格式无效: %v", err))
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return NewSquidError(ErrInvalidURL, "URL必须使用http或https协议")
	}

	if parsedURL.Host == "" {
		return NewSquidError(ErrInvalidURL, "URL必须包含主机名")
	}

	return nil
}

// IsPDFURL 判断URL是否指向PDF文件
func IsPDFURL(targetURL string) bool {
	if targetURL == "" {
		return false
	}

	// 检查URL路径是否以.pdf结尾
	if strings.HasSuffix(strings.ToLower(targetURL), ".pdf") {
		return true
	}

	// 检查URL是否包含pdf关键字
	lowerURL := strings.ToLower(targetURL)
	return strings.Contains(lowerURL, "pdf") || strings.Contains(lowerURL, "document")
}

// GetContentTypeFromURL 从URL推断内容类型
func GetContentTypeFromURL(targetURL string) string {
	if targetURL == "" {
		return "application/octet-stream"
	}

	lowerURL := strings.ToLower(targetURL)

	if strings.HasSuffix(lowerURL, ".pdf") {
		return "application/pdf"
	}
	if strings.HasSuffix(lowerURL, ".jpg") || strings.HasSuffix(lowerURL, ".jpeg") {
		return "image/jpeg"
	}
	if strings.HasSuffix(lowerURL, ".png") {
		return "image/png"
	}
	if strings.HasSuffix(lowerURL, ".gif") {
		return "image/gif"
	}
	if strings.HasSuffix(lowerURL, ".html") || strings.HasSuffix(lowerURL, ".htm") {
		return "text/html"
	}
	if strings.HasSuffix(lowerURL, ".txt") {
		return "text/plain"
	}
	if strings.HasSuffix(lowerURL, ".json") {
		return "application/json"
	}
	if strings.HasSuffix(lowerURL, ".xml") {
		return "application/xml"
	}

	return "application/octet-stream"
}

// FormatFileSize 格式化文件大小显示
func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDuration 格式化持续时间显示
func FormatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%.2fμs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	}
	if d < time.Minute {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return d.String()
}

// GetFileNameFromURL 从URL中提取文件名
func GetFileNameFromURL(targetURL string) string {
	if targetURL == "" {
		return "download"
	}

	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return "download"
	}

	path := parsedURL.Path
	if path == "" || path == "/" {
		return "download"
	}

	// 获取路径的最后一部分
	parts := strings.Split(path, "/")
	fileName := parts[len(parts)-1]

	// 如果文件名为空，使用默认名称
	if fileName == "" {
		fileName = "download"
	}

	// 如果没有扩展名且是PDF URL，添加.pdf扩展名
	if !strings.Contains(fileName, ".") && IsPDFURL(targetURL) {
		fileName += ".pdf"
	}

	return fileName
}

// MergeHeaders 合并HTTP头
func MergeHeaders(base, custom map[string]string) map[string]string {
	result := make(map[string]string)

	// 复制基础头
	for k, v := range base {
		result[k] = v
	}

	// 覆盖自定义头
	for k, v := range custom {
		result[k] = v
	}

	return result
}

// BuildProxyAuthHeader 构建代理认证头
func BuildProxyAuthHeader(username, password string) string {
	if username == "" || password == "" {
		return ""
	}
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(username+":"+password))
}

// IsSuccessStatusCode 判断HTTP状态码是否表示成功
func IsSuccessStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 300
}

// GetRetryDelay 计算重试延迟（指数退避）
func GetRetryDelay(attempt int, baseDelay time.Duration) time.Duration {
	if attempt <= 0 {
		return baseDelay
	}

	// 指数退避：baseDelay * 2^attempt，最大不超过30秒
	delay := baseDelay
	for i := 0; i < attempt && delay < 30*time.Second; i++ {
		delay *= 2
	}

	if delay > 30*time.Second {
		delay = 30 * time.Second
	}

	return delay
}

// CloneRequest 克隆HTTP请求
func CloneRequest(req *http.Request) *http.Request {
	if req == nil {
		return nil
	}

	// 创建新请求
	newReq := &http.Request{
		Method:        req.Method,
		URL:           req.URL,
		Proto:         req.Proto,
		ProtoMajor:    req.ProtoMajor,
		ProtoMinor:    req.ProtoMinor,
		Header:        make(http.Header),
		Body:          req.Body,
		ContentLength: req.ContentLength,
		Host:          req.Host,
	}

	// 复制头部
	for k, v := range req.Header {
		newReq.Header[k] = v
	}

	return newReq
}

// SanitizeURL 清理URL，移除敏感信息
func SanitizeURL(targetURL string) string {
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		return targetURL
	}

	// 移除用户信息
	parsedURL.User = nil

	// 移除查询参数中的敏感信息
	query := parsedURL.Query()
	sensitiveParams := []string{"token", "key", "password", "secret", "auth"}

	for _, param := range sensitiveParams {
		if query.Has(param) {
			query.Set(param, "***")
		}
	}

	parsedURL.RawQuery = query.Encode()
	return parsedURL.String()
}

// userAgents 存储了一组常见的 User-Agent 字符串
// 实际使用时，你可以扩展这个列表，使其包含更多真实的、多样化的 UA
var userAgents = []string{
	// Chrome 桌面
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",
	"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.0.0 Safari/537.36",

	// Firefox 桌面
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:110.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:110.0) Gecko/20100101 Firefox/110.0",
	"Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:110.0) Gecko/20100101 Firefox/110.0",

	// Safari 桌面
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Safari/605.1.15",

	// Edge 桌面
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Edge/110.0.1587.63 Safari/537.36",

	// Chrome Android 手机
	"Mozilla/5.0 (Linux; Android 10) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/110.0.5481.65 Mobile Safari/537.36",

	// Safari iOS 手机
	"Mozilla/5.0 (iPhone; CPU iPhone OS 16_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/16.3 Mobile/15E148 Safari/604.1",
}

// GetRandomUserAgent 从 userAgents 列表中随机选择一个 User-Agent 字符串并返回
func GetRandomUserAgent() string {
	if len(userAgents) == 0 {
		return "" // 如果列表为空，返回空字符串或一个默认值
	}
	// rand.Intn(n) 返回一个 [0, n) 范围的随机整数
	randomIndex := rand.Intn(len(userAgents))
	return userAgents[randomIndex]
}

// ValidatePDFMagicNumber 验证文件是否以PDF魔数开头
// PDF文件必须以 "%PDF-" 开头（字节序列：0x25 0x50 0x44 0x46 0x2D）
func ValidatePDFMagicNumber(data []byte) error {
	if len(data) < 5 {
		return NewSquidError(ErrInvalidResponse, "文件太小，不是有效的PDF文件")
	}

	// PDF魔数：%PDF-
	pdfMagic := []byte{0x25, 0x50, 0x44, 0x46, 0x2D}

	for i := 0; i < 5; i++ {
		if data[i] != pdfMagic[i] {
			return NewSquidError(ErrInvalidResponse, fmt.Sprintf("文件头不是有效的PDF魔数，期望 %%PDF-，实际: %s", string(data[:min(10, len(data))])))
		}
	}

	return nil
}

// ValidatePDFContentType 验证HTTP响应的Content-Type是否为PDF
func ValidatePDFContentType(contentType string) error {
	if contentType == "" {
		return NewSquidError(ErrInvalidResponse, "响应头缺少Content-Type")
	}

	// 标准化处理，移除参数（如 charset）
	ct := strings.ToLower(strings.Split(contentType, ";")[0])
	ct = strings.TrimSpace(ct)

	// 接受的PDF Content-Type
	validTypes := []string{
		"application/pdf",
		"application/x-pdf",
		"application/x-bzpdf",
		"application/x-gzpdf",
	}

	for _, validType := range validTypes {
		if ct == validType {
			return nil
		}
	}

	return NewSquidError(ErrInvalidResponse, fmt.Sprintf("Content-Type不是PDF类型，实际: %s", contentType))
}

// ValidatePDFStructure 验证PDF文件的基本结构完整性
// 检查是否包含 %%EOF 结尾标记
func ValidatePDFStructure(data []byte) error {
	if len(data) < 100 {
		return NewSquidError(ErrInvalidResponse, "文件太小，不是有效的PDF文件")
	}

	// 检查文件末尾是否有 %%EOF 标记
	// PDF规范要求文件以 %%EOF 结束，但可能有尾随空白
	eofMarker := []byte("%%EOF")
	lastBytes := data
	if len(data) > 1024 {
		// 只检查最后1KB，提高效率
		lastBytes = data[len(data)-1024:]
	}

	if !bytes.Contains(lastBytes, eofMarker) {
		return NewSquidError(ErrInvalidResponse, "PDF文件缺少%%EOF结尾标记，文件可能不完整")
	}

	return nil
}

// ValidatePDFFile 综合验证PDF文件的合法性
// 包括：Content-Type验证、魔数验证、结构完整性验证
func ValidatePDFFile(data []byte, contentType string) error {
	// 1. 验证Content-Type（如果是 application/octet-stream，跳过此验证，因为很多服务器会返回通用二进制流）
	ct := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	if ct != "application/octet-stream" && ct != "" {
		// 只有当 Content-Type 不是 octet-stream 且不为空时，才进行严格验证
		if err := ValidatePDFContentType(contentType); err != nil {
			return err
		}
	}

	// 2. 验证PDF魔数（这是最可靠的验证方式）
	if err := ValidatePDFMagicNumber(data); err != nil {
		return err
	}

	// 3. 验证PDF结构完整性
	if err := ValidatePDFStructure(data); err != nil {
		return err
	}

	return nil
}

// min 返回两个整数中的较小值
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
