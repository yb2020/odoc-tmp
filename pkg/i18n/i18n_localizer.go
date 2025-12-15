package i18n

import (
	"sync"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/yb2020/odoc/pkg/logging"
	"golang.org/x/text/language"
)

// Localizer 本地化器接口
type Localizer interface {
	// Localize 本地化消息
	Localize(messageID string, c *gin.Context) string
	// LocalizeWithData 使用数据本地化消息
	LocalizeWithData(messageID string, data map[string]interface{}, c *gin.Context) string
	// GetDefaultLanguage 获取默认语言
	GetDefaultLanguage() string
	// GetSupportedLanguages 获取支持的语言列表
	GetSupportedLanguages() []string
	//获取当前上下文中的语言类型
	GetLanguage(c *gin.Context) string
}

// Loader i18n资源加载器接口
type Loader interface {
	GetBundle() *i18n.Bundle
	GetDefaultLanguage() string
	GetSupportedLanguages() []string
	GetFallbackLanguage() string
	GetLanguage(c *gin.Context) string
}

// I18nLocalizer 国际化本地化器
// 实现 Localizer 接口
type I18nLocalizer struct {
	bundle           *i18n.Bundle
	defaultLanguage  string
	supportedLangs   []string
	fallbackLanguage string
	logger           logging.Logger
	// 本地化器缓存，按语言缓存
	localizerCache map[string]*i18n.Localizer
	localizerMutex sync.RWMutex
}

// 全局本地化器
var (
	// 全局本地化器实例
	atomicLocalizer atomic.Value
)

// NewLocalizer 创建新的本地化器
func NewLocalizer(loader Loader, logger logging.Logger) *I18nLocalizer {
	// 如果没有提供 logger，创建一个默认的
	if logger == nil {
		logger = logging.NewLogger("info", "logfmt")
	}

	return &I18nLocalizer{
		bundle:           loader.GetBundle(),
		defaultLanguage:  loader.GetDefaultLanguage(),
		supportedLangs:   loader.GetSupportedLanguages(),
		fallbackLanguage: loader.GetFallbackLanguage(),
		logger:           logger,
		localizerCache:   make(map[string]*i18n.Localizer),
	}
}

// Localize 本地化消息
func (l *I18nLocalizer) Localize(messageID string, c *gin.Context) string {
	// 获取语言
	lang := l.getLanguage(c)

	// 使用语言本地化消息
	return l.localizeWithLocale(messageID, nil, lang)
}

// LocalizeWithData 使用数据本地化消息
func (l *I18nLocalizer) LocalizeWithData(messageID string, data map[string]interface{}, c *gin.Context) string {
	// 获取语言
	lang := l.getLanguage(c)

	// 使用语言本地化消息
	return l.localizeWithLocale(messageID, data, lang)
}

func (l *I18nLocalizer) GetLanguage(c *gin.Context) string {
	return l.getLanguage(c)
}

// localizeWithLocale 使用指定语言本地化消息
func (l *I18nLocalizer) localizeWithLocale(messageID string, data map[string]interface{}, locale string) string {
	l.logger.Debug("msg", "localizeWithLocale called", "messageID", messageID, "locale", locale)

	// 从缓存获取或创建本地化器
	localizer := l.getOrCreateLocalizer(locale)

	// 添加调试代码，检查 bundle 中的语言标签
	l.logger.Debug("msg", "Bundle language tags", "tags", l.bundle.LanguageTags())

	// 本地化消息
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
		DefaultMessage: &i18n.Message{
			ID:    messageID,
			Other: messageID, // 使用消息ID作为默认消息
		},
	})

	if err != nil {
		// 打印详细错误信息
		l.logger.Debug("msg", "Localization failed", "error", err, "messageID", messageID, "locale", locale, "errorType", err.Error())

		// 如果本地化失败，且当前语言不是回退语言，尝试使用回退语言
		if locale != l.fallbackLanguage {
			// 使用回退语言
			fallbackLocalizer := l.getOrCreateLocalizer(l.fallbackLanguage)
			fallbackMsg, fallbackErr := fallbackLocalizer.Localize(&i18n.LocalizeConfig{
				MessageID:    messageID,
				TemplateData: data,
				DefaultMessage: &i18n.Message{
					ID:    messageID,
					Other: messageID, // 使用消息ID作为默认消息
				},
			})

			if fallbackErr == nil {
				return fallbackMsg
			}

			// 如果回退语言也失败，记录错误
			l.logger.Debug("msg", "Fallback localization also failed", "error", fallbackErr, "messageID", messageID)
		}

		// 如果回退语言也失败或者当前已经是回退语言，返回消息ID
		return messageID
	}

	return msg
}

// getOrCreateLocalizer 获取或创建本地化器
func (l *I18nLocalizer) getOrCreateLocalizer(locale string) *i18n.Localizer {
	// 先尝试从缓存中读取
	l.localizerMutex.RLock()
	localizer, exists := l.localizerCache[locale]
	l.localizerMutex.RUnlock()

	if exists {
		return localizer
	}

	// 如果缓存中不存在，创建新的本地化器
	l.localizerMutex.Lock()
	defer l.localizerMutex.Unlock()

	// 双重检查，避免在获取写锁的过程中其他协程已经创建了
	localizer, exists = l.localizerCache[locale]
	if exists {
		return localizer
	}

	// 创建新的本地化器并缓存
	localizer = i18n.NewLocalizer(l.bundle, locale)
	l.localizerCache[locale] = localizer

	return localizer
}

// GetDefaultLanguage 获取默认语言
func (l *I18nLocalizer) GetDefaultLanguage() string {
	return l.defaultLanguage
}

// GetSupportedLanguages 获取支持的语言列表
func (l *I18nLocalizer) GetSupportedLanguages() []string {
	return l.supportedLangs
}

// getLanguage 获取语言
func (l *I18nLocalizer) getLanguage(c *gin.Context) string {
	// 如果上下文为nil，直接返回默认语言
	if c == nil {
		if l.logger != nil {
			l.logger.Debug("msg", "Context is nil, using default language", "defaultLanguage", l.defaultLanguage)
		}
		return l.defaultLanguage
	}

	// 从请求头获取语言
	lang := c.GetHeader("Accept-Language")
	if l.logger != nil {
		l.logger.Debug("msg", "Got language from Accept-Language header", "language", lang)
	}

	if lang == "" {
		// 从查询参数获取语言
		lang = c.Query("lang")
		if l.logger != nil {
			l.logger.Debug("msg", "Got language from query parameter", "language", lang)
		}
	}

	// 如果没有指定语言，使用默认语言
	if lang == "" {
		if l.logger != nil {
			l.logger.Debug("msg", "No language specified, using default language", "defaultLanguage", l.defaultLanguage)
		}
		return l.defaultLanguage
	}

	// 使用转换层标准化语言格式为 RFC 5646
	normalizedLang := GlobalConverter.NormalizeToRFC5646(lang)
	if l.logger != nil {
		l.logger.Debug("msg", "Normalized language", "original", lang, "normalized", normalizedLang)
	}

	// 检查标准化后的语言是否被支持
	if GlobalConverter.IsSupported(normalizedLang) {
		return normalizedLang
	}

	// 如果不是支持的语言，使用默认语言
	if l.logger != nil {
		l.logger.Debug("msg", "Unsupported language, using default language", "language", lang, "normalized", normalizedLang, "defaultLanguage", l.defaultLanguage)
	}
	return l.defaultLanguage
}

// SimpleLoader 简单的 Loader 实现
type SimpleLoader struct {
	bundle           *i18n.Bundle
	defaultLanguage  string
	supportedLangs   []string
	fallbackLanguage string
}

// NewSimpleLoader 创建一个简单的 Loader
func NewSimpleLoader(defaultLang string, supportedLangs []string, fallbackLang string) *SimpleLoader {
	// 创建语言包
	bundle := i18n.NewBundle(language.MustParse(defaultLang))

	return &SimpleLoader{
		bundle:           bundle,
		defaultLanguage:  defaultLang,
		supportedLangs:   supportedLangs,
		fallbackLanguage: fallbackLang,
	}
}

// GetBundle 获取语言包
func (l *SimpleLoader) GetBundle() *i18n.Bundle {
	return l.bundle
}

// GetDefaultLanguage 获取默认语言
func (l *SimpleLoader) GetDefaultLanguage() string {
	return l.defaultLanguage
}

// GetSupportedLanguages 获取支持的语言列表
func (l *SimpleLoader) GetSupportedLanguages() []string {
	return l.supportedLangs
}

// GetFallbackLanguage 获取回退语言
func (l *SimpleLoader) GetFallbackLanguage() string {
	return l.fallbackLanguage
}

// SetLocalizer 设置全局本地化器
func SetLocalizer(localizer Localizer) {
	// 如果提供了有效的本地化器，则更新它
	if localizer != nil {
		atomicLocalizer.Store(localizer)
	}
}

// GetLocalizer 获取全局本地化器
func GetLocalizer() Localizer {
	// 从原子变量中获取本地化器
	localizer := atomicLocalizer.Load()
	if localizer == nil {
		// 如果没有设置，则返回 nil
		return nil
	}

	return localizer.(Localizer)
}
