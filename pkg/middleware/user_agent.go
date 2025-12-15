package middleware

import (
	"context"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

// UserAgentInfo 存储用户代理解析的信息
type UserAgentInfo struct {
	OsName             string
	OsVersion          string
	BrowserName        string
	BrowserVersion     string
	DeviceType         string
	DeviceManufacturer string
	RawUserAgent       string
}

// UserAgentMiddleware 创建一个中间件，用于解析用户代理信息并存储在 context 中
func UserAgentMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userAgent := c.Request.UserAgent()

		// 解析用户代理信息
		info := parseUserAgent(userAgent)

		// 将信息存储在 Gin context 中
		c.Set("user_agent_info", info)

		// 将信息存储在请求的 context.Context 中
		ctx := context.WithValue(c.Request.Context(), userAgentInfoKey, info)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

// 定义 context key
type contextKey string

const userAgentInfoKey contextKey = "user_agent_info"

// GetUserAgentInfo 从 context 中获取用户代理信息
func GetUserAgentInfo(ctx context.Context) (UserAgentInfo, bool) {
	info, ok := ctx.Value(userAgentInfoKey).(UserAgentInfo)
	return info, ok
}

// GetUserAgentInfoFromGin 从 Gin context 中获取用户代理信息
func GetUserAgentInfoFromGin(c *gin.Context) (UserAgentInfo, bool) {
	info, ok := c.Get("user_agent_info")
	if !ok {
		return UserAgentInfo{}, false
	}
	return info.(UserAgentInfo), true
}

func parseUserAgent(userAgentString string) UserAgentInfo {
	ua := user_agent.New(userAgentString)

	// 获取操作系统信息
	osInfo := ua.OS()

	// 获取浏览器信息
	browserName, browserVersion := ua.Browser()

	// 获取设备信息
	deviceType := "desktop"
	if ua.Mobile() {
		deviceType = "mobile"
	}

	// 尝试获取设备制造商（这需要更复杂的解析）
	deviceManufacturer := detectDeviceManufacturer(userAgentString)

	return UserAgentInfo{
		OsName:             osInfo,
		OsVersion:          detectOSVersion(userAgentString, osInfo),
		BrowserName:        browserName,
		BrowserVersion:     browserVersion,
		DeviceType:         deviceType,
		DeviceManufacturer: deviceManufacturer,
		RawUserAgent:       userAgentString,
	}
}

// 检测操作系统版本
func detectOSVersion(userAgent, osName string) string {
	// 根据不同的操作系统实现版本检测逻辑
	// 这里只是一个简化的示例
	if strings.Contains(osName, "Windows") {
		re := regexp.MustCompile(`Windows NT (\d+\.\d+)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if strings.Contains(osName, "Mac OS X") {
		re := regexp.MustCompile(`Mac OS X (\d+[._]\d+[._]?\d*)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			return strings.Replace(matches[1], "_", ".", -1)
		}
	} else if strings.Contains(osName, "Android") {
		re := regexp.MustCompile(`Android (\d+\.\d+)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			return matches[1]
		}
	} else if strings.Contains(osName, "iOS") {
		re := regexp.MustCompile(`OS (\d+[._]\d+[._]?\d*)`)
		matches := re.FindStringSubmatch(userAgent)
		if len(matches) > 1 {
			return strings.Replace(matches[1], "_", ".", -1)
		}
	}
	return ""
}

// 检测设备制造商
func detectDeviceManufacturer(userAgent string) string {
	// 这里只是一个简化的示例，实际上需要更复杂的逻辑
	if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") || strings.Contains(userAgent, "iPod") {
		return "Apple"
	} else if strings.Contains(userAgent, "Samsung") {
		return "Samsung"
	} else if strings.Contains(userAgent, "Huawei") {
		return "Huawei"
	} else if strings.Contains(userAgent, "Xiaomi") {
		return "Xiaomi"
	} else if strings.Contains(userAgent, "OPPO") {
		return "OPPO"
	} else if strings.Contains(userAgent, "vivo") {
		return "vivo"
	}
	return "Unknown"
}
