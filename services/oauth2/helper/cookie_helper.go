package helper

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"

	"strings"
	"time"

	"github.com/yb2020/odoc/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/yb2020/odoc/config"
)

// 常量定义
const (
	// Cookie 相关常量
	DefaultAppID = "default" // 默认应用 ID
)

// WriteSuccessLoginCookie 设置登录成功的 cookie
func WriteSuccessLoginCookie(c *gin.Context, userId string, accessToken string, refreshExpiresAt uint64, appId string) {
	// 获取全局配置
	config := config.GetConfig()
	if appId == "" {
		appId = DefaultAppID
	}

	// 创建 token cookie 标签
	cookieLabel := fmt.Sprintf(config.OAuth2.TokenStorage.CookieTokenLabel, appId)
	cookieUIDLabel := fmt.Sprintf(config.OAuth2.TokenStorage.CookieTokenUIDLabel, appId)

	// 先删除要加入的 cookie
	DestroyCookieByAppId(c, appId)

	// 计算相对过期时间（秒）
	// refreshExpiresAt 是一个 Unix 时间戳，表示过期的绝对时间
	now := time.Now().Unix()
	var maxAge int
	if refreshExpiresAt > uint64(now) {
		// 如果过期时间在未来，计算相对时间
		maxAge = int(refreshExpiresAt - uint64(now))
	} else {
		// 如果过期时间无效或已过期，使用默认值（24小时）
		maxAge = 86400 // 24小时
	}

	// 设置 token cookie
	tokenValue := base64.StdEncoding.EncodeToString([]byte(accessToken))
	loginCookie := &http.Cookie{
		Name:     cookieLabel,
		Value:    tokenValue,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: true,
	}
	// 设置 domain
	SetDomain(c, loginCookie)
	http.SetCookie(c.Writer, loginCookie)

	// 设置 uid cookie
	uidCookie := &http.Cookie{
		Name:     cookieUIDLabel,
		Value:    userId,
		MaxAge:   maxAge,
		Path:     "/",
		HttpOnly: false,
	}
	// 设置 domain
	SetDomain(c, uidCookie)
	http.SetCookie(c.Writer, uidCookie)
}

// DestroyCookieByAppId 根据 appId 删除 cookie
func DestroyCookieByAppId(c *gin.Context, appId string) {
	// 获取全局配置
	config := config.GetConfig()
	if appId == "" {
		appId = DefaultAppID
	}

	// 创建 token cookie 标签
	cookieLabel := fmt.Sprintf(config.OAuth2.TokenStorage.CookieTokenLabel, appId)
	cookieUIDLabel := fmt.Sprintf(config.OAuth2.TokenStorage.CookieTokenUIDLabel, appId)

	// 删除 token cookie
	expiredCookie := &http.Cookie{
		Name:     cookieLabel,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: true,
	}
	SetDomain(c, expiredCookie)
	http.SetCookie(c.Writer, expiredCookie)

	// 删除 uid cookie
	expiredUIDCookie := &http.Cookie{
		Name:     cookieUIDLabel,
		Value:    "",
		MaxAge:   -1,
		Path:     "/",
		HttpOnly: false,
	}
	SetDomain(c, expiredUIDCookie)
	http.SetCookie(c.Writer, expiredUIDCookie)
}

// SetDomain 设置 cookie 的 domain
func SetDomain(c *gin.Context, cookie *http.Cookie) {
	// 1. 优先从 X-Forwarded-Host Header 获取原始 Host。这是反向代理传递客户端原始请求的标准方式。
	domain := c.Request.Header.Get("X-Forwarded-Host")

	// 2. 如果 X-Forwarded-Host 不存在，则回退到使用 Request.Host。
	//    这保证了在没有反向代理的直连环境下代码也能正常工作。
	if domain == "" {
		domain = c.Request.Host
	}
	// 去掉端口号
	if strings.Contains(domain, ":") {
		domain = strings.Split(domain, ":")[0]
	}

	// 检查是否是 IP 地址
	if !IsIPv4(domain) {
		mainDomain := utils.GetMainDomain(domain)
		// 添加点前缀，设置到主域
		cookie.Domain = fmt.Sprintf(".%s", mainDomain)
	}
}

// IsIPv4 检查是否是 IPv4 地址
func IsIPv4(ip string) bool {
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return false
		}
	}

	return true
}
