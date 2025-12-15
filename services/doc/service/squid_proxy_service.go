package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/opentracing/opentracing-go"
	pb "github.com/yb2020/odoc-proto/gen/go/doc"
	"github.com/yb2020/odoc/config"
	userContext "github.com/yb2020/odoc/context"
	"github.com/yb2020/odoc/pkg/logging"
	"github.com/yb2020/odoc/pkg/squid_proxy"
	ossModel "github.com/yb2020/odoc/services/oss/model"
	userService "github.com/yb2020/odoc/services/user/service"
)

// SquidProxyService squid 代理服务实现
type SquidProxyService struct {
	logger      logging.Logger
	tracer      opentracing.Tracer
	cfg         *config.Config
	docService  *UserDocService
	userService *userService.UserService
}

// NewSquidProxyService 创建新的 squid 代理服务
func NewSquidProxyService(
	logger logging.Logger,
	tracer opentracing.Tracer,
	cfg *config.Config,
	docService *UserDocService,
	userService *userService.UserService,
) *SquidProxyService {
	return &SquidProxyService{
		logger:      logger,
		tracer:      tracer,
		cfg:         cfg,
		docService:  docService,
		userService: userService,
	}
}

// NewUserContext 创建用户上下文
func (s *SquidProxyService) NewUserContext(ctx context.Context, userId string) context.Context {
	user, err := s.userService.GetUserByID(ctx, userId)
	if err != nil {
		s.logger.Error("msg", "Get user by id failed", "userId", userId, "error", err)
		return ctx
	}
	ctx = userContext.SetUserContext(ctx, user)
	return ctx
}

// ProxyHttp 通过 squid 代理发送http请求
func (s *SquidProxyService) ProxyHttp(ctx context.Context, url string) ([]byte, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.ProxyHttp")
	defer span.Finish()

	userAgent := squid_proxy.GetRandomUserAgent()

	// 1. 创建代理配置
	config := &squid_proxy.SquidConfig{
		ProxyURL:  s.cfg.Squid.ProxyUrl,
		Username:  s.cfg.Squid.Username,
		Password:  s.cfg.Squid.Password,
		Timeout:   time.Duration(s.cfg.Squid.Timeout) * time.Second,
		UserAgent: userAgent,
	}

	// 2. 创建代理客户端
	client, err := squid_proxy.NewSquidProxyClient(config)
	if err != nil {
		s.logger.Error("msg", "创建squid代理客户端失败", "error", err)
		return nil, err
	}

	// 3. 发送GET请求
	headers := map[string]string{
		"Accept": "application/json",
	}
	result, err := client.Get(url, headers)
	if err != nil {
		s.logger.Error("msg", "通过squid代理发送http请求失败", "url", url, "error", err)
		return nil, err
	}

	return result.Body, nil
}

// DownloadPdfByUrl 通过 squid 代理下载pdf，参数 url 为 pdf 文件的 url
// 返回值：封装了文件信息的结构体和错误
func (s *SquidProxyService) DownloadPdfByUrl(ctx context.Context, url string) (*squid_proxy.DownloadedPdfInfo, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.DownloadPdfByUrl")
	defer span.Finish()

	userAgent := squid_proxy.GetRandomUserAgent()

	// 1. 创建代理配置
	config := &squid_proxy.SquidConfig{
		ProxyURL:  s.cfg.Squid.ProxyUrl,                             // 使用指定的Squid代理服务器
		Username:  s.cfg.Squid.Username,                             // 如果需要认证，填写用户名
		Password:  s.cfg.Squid.Password,                             // 如果需要认证，填写密码
		Timeout:   time.Duration(s.cfg.Squid.Timeout) * time.Second, // 增加超时时间，arXiv文件可能较大
		UserAgent: userAgent,
	}

	// 2. 创建代理客户端
	s.logger.Info("msg", "代理配置", "config ", config)
	client, err := squid_proxy.NewSquidProxyClient(config)
	if err != nil {
		s.logger.Error("msg", "创建代理客户端失败: ", err)
		return nil, err
	}

	// 3. 测试代理连接（可选）
	s.logger.Info("msg", "测试代理连接...")
	if err := client.TestConnection(); err != nil {
		s.logger.Error("msg", "代理连接失败", "error", err)
		return nil, err
	}
	s.logger.Info("msg", "代理连接成功!")

	// 4. 下载PDF文件到内存
	s.logger.Info("msg", "正在下载论文: ", url)

	result, err := client.DownloadFile(url, squid_proxy.PDFDownloadOptions())
	if err != nil {
		s.logger.Error("msg", "下载失败: ", err)
		return nil, err
	}

	s.logger.Info("msg", "下载成功!")

	// ===== PDF文件验证 =====
	contentType := result.Headers.Get("Content-Type")
	s.logger.Info("msg", "开始验证PDF文件", "Content-Type", contentType, "文件大小", squid_proxy.FormatFileSize(result.ContentLength))

	// 1. 综合验证PDF文件（Content-Type + 魔数 + 结构完整性）
	if err := squid_proxy.ValidatePDFFile(result.Body, contentType); err != nil {
		s.logger.Error("msg", "PDF文件验证失败", "error", err, "url", url)
		return nil, fmt.Errorf("PDF文件验证失败: %w", err)
	}

	s.logger.Info("msg", "PDF文件验证通过")

	// 提取文件名（优先从响应头获取）
	fileName := squid_proxy.GetFileNameFromResponse(result.Headers, url)
	s.logger.Info("msg", "文件名:", fileName)

	// 计算SHA256值
	hash := sha256.Sum256(result.Body)
	sha256Hash := fmt.Sprintf("%x", hash)
	s.logger.Info("msg", "SHA256:", sha256Hash)

	s.logger.Info("msg", "文件大小:", squid_proxy.FormatFileSize(result.ContentLength))
	s.logger.Info("msg", "下载耗时:", squid_proxy.FormatDuration(result.Duration))

	// 创建一个可重用的文件读取器
	fileReader := bytes.NewReader(result.Body)

	// 获取PDF页数
	var pageCount int
	pageCount, err = squid_proxy.GetPdfPageCount(fileReader) // 使用同一个reader
	if err != nil {
		s.logger.Warn("msg", "获取PDF页数失败，将使用默认值0", "error", err)
		pageCount = 0 // 出错时，页数给默认值0
	} else {
		s.logger.Info("msg", "PDF页数:", pageCount)
	}

	// 封装返回结果
	info := &squid_proxy.DownloadedPdfInfo{
		FileReader: fileReader, // 复用reader
		FileName:   fileName,
		SHA256:     sha256Hash,
		Size:       result.ContentLength,
		PageCount:  pageCount,
	}

	return info, nil
}

// DownloadAndUploadPdfByUrl 通过URL下载PDF并上传到OSS的完整流程
func (s *SquidProxyService) DownloadAndUploadPdfByUrl(ctx context.Context, url string, token string, userID string) (*ossModel.OssRecord, error) {
	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.DownloadAndUploadPdfByUrl")
	defer span.Finish()

	// 1. 设置状态为下载中
	if err := s.docService.ChangeUploadTokenStatus(ctx, token, pb.UserDocParsedStatusEnum_DOWNLOADING); err != nil {
		s.logger.Error("msg", "change upload token status to downloading failed", "error", err)
		return nil, err
	}

	// 2. 下载PDF文件
	pdfInfo, err := s.DownloadPdfByUrl(ctx, url)
	if err != nil {
		s.logger.Error("msg", "download pdf by url failed", "error", err)
		// 下载失败时设置状态为下载失败
		if statusErr := s.docService.ChangeUploadTokenStatus(ctx, token, pb.UserDocParsedStatusEnum_DOWNLOAD_FAILED); statusErr != nil {
			s.logger.Error("msg", "change upload token status to download failed error", "error", statusErr)
		}
		return nil, err
	}

	// 3. 设置状态为下载完成
	if err := s.docService.ChangeUploadTokenStatus(ctx, token, pb.UserDocParsedStatusEnum_DOWNLOADED); err != nil {
		s.logger.Error("msg", "change upload token status to downloaded failed", "error", err)
		return nil, err
	}

	s.logger.Info("msg", "download pdf by url success", "pdfInfo", pdfInfo)

	// 4. 上传PDF到OSS
	ossRecord, err := s.docService.UploadPdf(ctx, pdfInfo.FileReader, pdfInfo.FileName, pdfInfo.SHA256, pdfInfo.Size, int64(pdfInfo.PageCount), userID, token)
	if err != nil {
		s.logger.Error("msg", "upload pdf failed", "error", err)
		return nil, err
	}

	s.logger.Info("msg", "upload pdf by url success", "ossRecord", ossRecord)
	return ossRecord, nil
}

// DownloadAndUploadPdfByUrlAsync 异步版本：通过URL下载PDF并上传到OSS
func (s *SquidProxyService) DownloadAndUploadPdfByUrlAsync(ctx context.Context, url string, token string, userID string) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.DownloadAndUploadPdfByUrlAsync")
	defer span.Finish()

	s.logger.Info("msg", "start async pdf download and upload", "url", url, "token", token, "userID", userID)

	// 异步执行下载和上传流程
	go func() {
		// 创建带有用户上下文的context
		asyncCtx := s.NewUserContext(context.Background(), userID)

		// 调用原有的同步方法
		_, err := s.DownloadAndUploadPdfByUrl(asyncCtx, url, token, userID)
		if err != nil {
			s.logger.Error("msg", "async pdf download and upload failed", "error", err, "token", token)
		} else {
			s.logger.Info("msg", "async pdf download and upload completed successfully", "token", token)
		}
	}()
}

// ValidateIsArxivUrl 验证是否为 arxiv url
func (s *SquidProxyService) ValidateIsArxivUrl(ctx context.Context, urlStr string) (bool, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.ValidateArxivUrl")
	defer span.Finish()

	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false, errors.New("invalid url format:" + urlStr)
	}

	hostname := parsedUrl.Hostname()
	if hostname != "arxiv.org" && !strings.HasSuffix(hostname, ".arxiv.org") {
		return false, nil
	}

	return true, nil
}

// 将 arxiv 详情页面 url 替换为 pdf url，例如：
// https://arxiv.org/abs/2206.00001 -> https://arxiv.org/pdf/2206.00001
func (s *SquidProxyService) ReplaceArxivUrlWithPdfUrl(ctx context.Context, urlStr string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.ReplaceArxivUrlWithPdfUrl")
	defer span.Finish()

	// 检查 URL 是否包含 /abs/ 路径
	if strings.Contains(urlStr, "/abs/") {
		// 替换 /abs/ 为 /pdf/
		pdfUrl := strings.Replace(urlStr, "/abs/", "/pdf/", 1)
		return pdfUrl, nil
	}

	// 如果不包含 /abs/，则假定它已经是 PDF URL 或其他格式，直接返回原 URL
	return urlStr, nil
}

// ValidateIsGithubUrl 验证是否为 github url
func (s *SquidProxyService) ValidateIsGithubUrl(ctx context.Context, urlStr string) (bool, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.ValidateGithubUrl")
	defer span.Finish()

	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false, errors.New("invalid url format:" + urlStr)
	}

	hostname := parsedUrl.Hostname()
	if hostname != "github.com" && !strings.HasSuffix(hostname, ".github.com") {
		return false, nil
	}

	return true, nil
}

// 将 github 详情页面 url 替换为 pdf url，例如：
// https://github.com/deepseek-ai/DeepSeek-OCR/blob/main/DeepSeek_OCR_paper.pdf -> https://github.com/deepseek-ai/DeepSeek-OCR/raw/main/DeepSeek_OCR_paper.pdf
func (s *SquidProxyService) ReplaceGithubUrlWithPdfUrl(ctx context.Context, urlStr string) (string, error) {
	span, _ := opentracing.StartSpanFromContextWithTracer(ctx, s.tracer, "SquidProxyService.ReplaceGithubUrlWithPdfUrl")
	defer span.Finish()

	// 检查 URL 是否包含 /blob/ 路径
	if strings.Contains(urlStr, "/blob/") {
		// 替换 /blob/ 为 /raw/
		pdfUrl := strings.Replace(urlStr, "/blob/", "/raw/", 1)
		return pdfUrl, nil
	}

	// 如果不包含 /blob/，则假定它已经是 PDF URL 或其他格式，直接返回原 URL
	return urlStr, nil
}
