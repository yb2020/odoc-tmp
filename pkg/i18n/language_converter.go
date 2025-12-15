package i18n

import (
	langPb "github.com/yb2020/odoc-proto/gen/go/lang"
	"strings"
)

// LanguageConverter 语言格式转换器
// 负责在 Proto 枚举格式和 RFC 5646 标准格式之间进行转换
type LanguageConverter struct{}

// ProtoToRFC5646 将 Proto 枚举转换为 RFC 5646 标准格式
func (lc *LanguageConverter) ProtoToRFC5646(protoLang string) string {
	switch protoLang {
	case langPb.Language_EN_US.String():
		return "en-US"
	case langPb.Language_ZH_CN.String():
		return "zh-CN"
	default:
		return "en-US" // 默认回退
	}
}

// RFC5646ToProto 将 RFC 5646 格式转换为 Proto 枚举字符串
func (lc *LanguageConverter) RFC5646ToProto(rfc5646Lang string) string {
	normalized := lc.NormalizeToRFC5646(rfc5646Lang)
	switch normalized {
	case "en-US":
		return langPb.Language_EN_US.String()
	case "zh-CN":
		return langPb.Language_ZH_CN.String()
	default:
		return langPb.Language_EN_US.String()
	}
}

// NormalizeToRFC5646 将各种格式统一转换为 RFC 5646 标准
// 支持的输入格式：EN_US, EN-US, en_us, en-us, en, EN, ZH_CN, ZH-CN, zh_cn, zh-cn, zh, ZH
func (lc *LanguageConverter) NormalizeToRFC5646(lang string) string {
	// 去除空格并转换为小写进行比较
	normalized := strings.ToLower(strings.TrimSpace(lang))
	
	switch normalized {
	// 英语各种格式
	case "en_us", "en-us", "en", "enus":
		return "en-US"
	// 中文各种格式  
	case "zh_cn", "zh-cn", "zh", "zhcn":
		return "zh-CN"
	default:
		// 尝试大写格式
		upper := strings.ToUpper(lang)
		switch upper {
		case "EN_US", "EN-US", "EN", "ENUS":
			return "en-US"
		case "ZH_CN", "ZH-CN", "ZH", "ZHCN":
			return "zh-CN"
		default:
			return "en-US" // 默认回退到英语
		}
	}
}

// GetSupportedLanguages 获取支持的语言列表（RFC 5646 格式）
func (lc *LanguageConverter) GetSupportedLanguages() []string {
	return []string{"en-US", "zh-CN"}
}

// GetDefaultLanguage 获取默认语言（RFC 5646 格式）
func (lc *LanguageConverter) GetDefaultLanguage() string {
	return "en-US"
}

// GetFallbackLanguage 获取回退语言（RFC 5646 格式）
func (lc *LanguageConverter) GetFallbackLanguage() string {
	return "en-US"
}

// IsSupported 检查语言是否被支持
func (lc *LanguageConverter) IsSupported(lang string) bool {
	normalized := lc.NormalizeToRFC5646(lang)
	for _, supported := range lc.GetSupportedLanguages() {
		if normalized == supported {
			return true
		}
	}
	return false
}

// GetGrobidLanguage 为 Grobid 服务提供专门的格式转换
// Grobid 要求小写格式：zh-cn, en-us
func (lc *LanguageConverter) GetGrobidLanguage(rfc5646Lang string) string {
	switch rfc5646Lang {
	case "zh-CN":
		return "zh-cn"
	case "en-US":
		return "en-us"
	default:
		return "en-us"
	}
}

// GetDifyLanguageKey 为 Dify 服务提供专门的格式转换
// Dify 要求下划线格式：zh_cn, en_us
func (lc *LanguageConverter) GetDifyLanguageKey(rfc5646Lang string) string {
	switch rfc5646Lang {
	case "en-US":
		return "en_us"
	case "zh-CN":
		return "zh_cn"
	default:
		return "en_us"
	}
}

// 全局转换器实例
var GlobalConverter = &LanguageConverter{}
