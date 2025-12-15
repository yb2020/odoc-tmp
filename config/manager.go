package config

import (
	"sync"
)

var (
	globalConfig *Config
	configMutex  sync.RWMutex
)

// SetConfig 设置全局配置
func SetConfig(cfg *Config) {
	configMutex.Lock()
	defer configMutex.Unlock()
	globalConfig = cfg
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// 如果配置为空，返回默认配置
	if globalConfig == nil {
		return GlobalConfig()
	}

	return globalConfig
}

// Sub 获取子配置
// 简单实现，实际项目中应该使用更复杂的配置管理库，如viper
func (c *Config) Sub(key string) *SubConfig {
	// 这里简化处理，实际应该根据key路径解析配置
	return &SubConfig{
		parent: c,
		path:   key,
	}
}

// SubConfig 子配置
type SubConfig struct {
	parent *Config
	path   string
}

// GetString 获取字符串配置
func (s *SubConfig) GetString(key string) string {
	// 简化实现，实际应该根据path和key获取配置
	// 这里返回一些默认值用于测试
	if s.path == "oauth2.jwt" {
		switch key {
		case "secretKey":
			return "your-secret-key"
		case "issuer":
			return "go-sea"
		}
	}
	return ""
}

// GetInt 获取整数配置
func (s *SubConfig) GetInt(key string) int {
	// 简化实现，实际应该根据path和key获取配置
	if s.path == "oauth2.jwt" {
		switch key {
		case "expiration":
			return 3600 // 1小时
		case "refreshExpiration":
			return 604800 // 7天
		}
	}
	return 0
}

// GetStringSlice 获取字符串切片配置
func (s *SubConfig) GetStringSlice(key string) []string {
	// 简化实现，实际应该根据path和key获取配置
	if s.path == "oauth2.resourceProtection" {
		switch key {
		case "publicPaths":
			return []string{"/api/public/", "/api/oauth2/token"}
		case "adminPaths":
			return []string{"/api/admin/"}
		}
	}
	return nil
}
