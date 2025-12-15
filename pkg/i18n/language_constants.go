package i18n

import (
	langPb "github.com/yb2020/odoc-proto/gen/go/lang"
)

// RFC 5646 标准语言常量 - 项目统一使用这些常量
// 禁止在代码中直接使用硬编码的语言字符串
const (
	LanguageEnUS = "en-US"
	LanguageZhCN = "zh-CN"
)

// Proto 枚举获取函数 - 禁止直接调用 .String() 方法
// 这些函数提供了统一的 Proto 枚举访问接口

// GetLangProtoEnUS 获取英语的 Proto 枚举字符串
func GetLangProtoEnUS() string {
	return langPb.Language_EN_US.String()
}

// GetLangProtoZhCN 获取中文的 Proto 枚举字符串
func GetLangProtoZhCN() string {
	return langPb.Language_ZH_CN.String()
}

// GetLangProtoEnum 根据 RFC 5646 格式获取对应的 Proto 枚举
func GetLangProtoEnum(rfc5646Lang string) langPb.Language {
	switch rfc5646Lang {
	case LanguageEnUS:
		return langPb.Language_EN_US
	case LanguageZhCN:
		return langPb.Language_ZH_CN
	default:
		return langPb.Language_EN_US
	}
}

// 语言转换便捷函数

// GetRFC5646FromLangProto 从 Lang Proto 枚举字符串转换为 RFC 5646 格式
func GetRFC5646FromLangProto(protoLang string) string {
	return GlobalConverter.ProtoToRFC5646(protoLang)
}

// GetLangProtoFromRFC5646 从 RFC 5646 格式转换为 Lang Proto 枚举字符串
func GetLangProtoFromRFC5646(rfc5646Lang string) string {
	return GlobalConverter.RFC5646ToProto(rfc5646Lang)
}

// NormalizeLanguage 标准化语言格式为 RFC 5646
func NormalizeLanguage(lang string) string {
	return GlobalConverter.NormalizeToRFC5646(lang)
}

// IsLanguageSupported 检查语言是否被支持
func IsLanguageSupported(lang string) bool {
	return GlobalConverter.IsSupported(lang)
}

// GetDefaultLanguage 获取默认语言
func GetDefaultLanguage() string {
	return GlobalConverter.GetDefaultLanguage()
}

// GetSupportedLanguages 获取支持的语言列表
func GetSupportedLanguages() []string {
	return GlobalConverter.GetSupportedLanguages()
}

// GetFallbackLanguage 获取回退语言
func GetFallbackLanguage() string {
	return GlobalConverter.GetFallbackLanguage()
}
