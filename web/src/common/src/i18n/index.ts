/**
 * 统一的国际化管理工具
 * 供主项目和 VitePress 共同使用
 */

import Cookies from 'js-cookie'
import { Language } from 'go-sea-proto/gen/ts/lang/Language'
import { LanguageCodeMap, CodeToLanguageMap } from '../locals/i18n'
import { getLanguageCookie, languageEnumToStandard } from '../../../shared/language/service'

const I18N_COOKIE_KEY = 'i18n'

/**
 * 获取当前语言
 */
export const getCurrentLocale = (): string => {
  // 使用统一的语言管理服务获取标准格式语言
  const standardLocale = getLanguageCookie()
  
  if (standardLocale) {
    return standardLocale // 返回标准格式 'en-US' 或 'zh-CN'
  }
  
  // 默认返回英文标准格式
  return languageEnumToStandard(Language.EN_US)
}

/**
 * 设置语言
 */
export const setLocale = (locale: string) => {
  // 设置到 Cookie（主项目标准存储位置）
  Cookies.set(I18N_COOKIE_KEY, locale)
  
  console.log(`语言已设置: ${locale}`)
}

/**
 * 从路径检测语言 
 */
export const detectLocaleFromPath = (path: string): string => {
  if (path.startsWith('/zh/') || path.startsWith('/docs/zh/')) {
    return 'zh-CN'
  }
  return 'en-US'
}

/**
 * 监听语言变化   TODO： 这里必须优化，放到下个版本解决
 * 由于 Cookie 没有原生的变化事件，我们使用轮询机制
 */
export const onLocaleChange = (callback: (locale: string) => void) => {
  let currentLocale = getCurrentLocale()
  
  // 每秒检查一次
  const intervalId = setInterval(() => {
    const newLocale = getCurrentLocale()
    if (newLocale !== currentLocale) {
      currentLocale = newLocale
      callback(currentLocale)
    }
  }, 1000)
  
  // 返回清理函数
  return () => {
    clearInterval(intervalId)
  }
}

/**
 * 初始化语言
 * 根据 URL 路径和 Cookie 设置初始语言
 */
export const initLocale = (): string => {
  const path = window.location.pathname
  const detectedLocale = detectLocaleFromPath(path)
  setLocale(detectedLocale)
  return detectedLocale
}

// 重新导出 Language 以便外部使用
export { Language, LanguageCodeMap, CodeToLanguageMap }