package i18n

import (
	translatePb "github.com/yb2020/odoc-proto/gen/go/translate"
)

// 翻译服务语言常量和转换函数
// 用于处理翻译服务特有的语言枚举

// GetTranslateProtoEnUS 获取翻译服务的英语枚举
func GetTranslateProtoEnUS() translatePb.TranslateLanguage {
	return translatePb.TranslateLanguage_EN_US
}

// GetTranslateProtoZhCN 获取翻译服务的中文枚举
func GetTranslateProtoZhCN() translatePb.TranslateLanguage {
	return translatePb.TranslateLanguage_ZH_CN
}

// GetTranslateLanguageEnum 根据 RFC 5646 格式获取翻译服务的语言枚举
func GetTranslateLanguageEnum(rfc5646Lang string) translatePb.TranslateLanguage {
	switch rfc5646Lang {
	case LanguageEnUS:
		return translatePb.TranslateLanguage_EN_US
	case LanguageZhCN:
		return translatePb.TranslateLanguage_ZH_CN
	default:
		return translatePb.TranslateLanguage_EN_US
	}
}

// GetRFC5646FromTranslateProto 从翻译服务枚举转换为 RFC 5646 格式
func GetRFC5646FromTranslateProto(translateEnum translatePb.TranslateLanguage) string {
	switch translateEnum {
	case translatePb.TranslateLanguage_EN_US:
		return LanguageEnUS
	case translatePb.TranslateLanguage_ZH_CN:
		return LanguageZhCN
	default:
		return LanguageEnUS
	}
}

// GetDefaultTranslateSourceLanguage 获取默认的翻译源语言枚举
func GetDefaultTranslateSourceLanguage() translatePb.TranslateLanguage {
	return translatePb.TranslateLanguage_EN_US
}

// GetDefaultTranslateTargetLanguage 获取默认的翻译目标语言枚举
func GetDefaultTranslateTargetLanguage() translatePb.TranslateLanguage {
	return translatePb.TranslateLanguage_ZH_CN
}

// ParseTranslateLanguageFromString 从字符串解析翻译语言枚举
// 支持 EN-US, EN_US, zh-CN, ZH_CN 等各种格式
func ParseTranslateLanguageFromString(lang string) translatePb.TranslateLanguage {
	rfc5646Lang := GlobalConverter.NormalizeToRFC5646(lang)
	return GetTranslateLanguageEnum(rfc5646Lang)
}

// GetTranslateLanguageString 获取翻译服务枚举的字符串表示（RFC 5646 格式）
func GetTranslateLanguageString(translateEnum translatePb.TranslateLanguage) string {
	return GetRFC5646FromTranslateProto(translateEnum)
}
