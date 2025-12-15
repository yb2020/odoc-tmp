import { ref, watch } from 'vue'
import { useData } from 'vitepress'
import { createI18n } from 'vue-i18n'
import { Language } from 'go-sea-proto/gen/ts/lang/Language'; // 直接导入Language枚举

// 导入共享的i18n核心模块
import globalI18nInstance, { 
  registerCustomHook,
  convertEnumToStandardLocale,
  convertStandardToEnumLocale,
  docsProjectMessages, // 导入docs项目的语言包
  LanguageCodeMap,
  getCurrentLanguage,
} from '../../src/shared/i18n/core';
import { docsProjectHook } from '../../src/shared/i18n/docsProjectHooks';
import { useLanguage } from '../../src/hooks/useLanguage'

// 获取当前语言并转换为标准格式
const currentEnumLanguage = getCurrentLanguage();
const currentStandardLanguage = convertEnumToStandardLocale(currentEnumLanguage);

// 创建docs项目专用的i18n实例，使用标准格式的locale
const docsI18nInstance = createI18n({
  locale: currentStandardLanguage, // 使用转换后的标准格式
  legacy: false,
  globalInjection: true,
  messages: docsProjectMessages,
  // 复用主项目的数字和日期时间格式
  numberFormats: globalI18nInstance.global.numberFormats.value,
  datetimeFormats: globalI18nInstance.global.datetimeFormats.value,
});

// 注册docs项目的自定义钩子
registerCustomHook(docsProjectHook);

// 导出docs项目专用的i18n实例
export const i18n = docsI18nInstance;

// ==================== VitePress集成钩子 ====================

// const currentLocale = ref(docsI18nInstance.global.locale.value);

/**
 * VitePress特定的i18n钩子函数
 * 使用docs项目专用的i18n实例
 */
export function useVitePressI18n() {
  const { lang } = useData() // 这是来自 VitePress 的语言源头, e.g., 'en-US'
  const { changeLanguage } = useLanguage()

  // 防止循环更新的标志
  let isUpdating = false;

  // 监听 VitePress 语言变化 (这是唯一的语言变更入口)
  watch(
    () => lang.value,
    (newLang) => { // newLang 是标准格式，如 'en-US', 'zh-CN'
      if (isUpdating || newLang === docsI18nInstance.global.locale.value) {
        return;
      }
      
      isUpdating = true;

      // 1. 【核心】更新 docs 实例的 locale。这是驱动组件更新的关键！
      docsI18nInstance.global.locale.value = newLang;
      // 2. (可选，但良好实践) 同步另一个 i18n 实例
      const newEnumLang = convertStandardToEnumLocale(newLang); 
      if (newEnumLang !== undefined) {
        globalI18nInstance.global.locale.value = newEnumLang;
      }
      
      // 3. 同步语言切换
      if (newLang === LanguageCodeMap[Language.EN_US]) {
        changeLanguage(Language.EN_US);
      } else {
        changeLanguage(Language.ZH_CN);
      }

      isUpdating = false;
    },
    { immediate: true }
  )

  // 创建安全的翻译函数
  const safeT = (key: string, values?: any, options?: any) => {
    try {
      // 这个函数现在很纯粹，只负责翻译
      return docsI18nInstance.global.t(key, values, options);
    } catch (error) {
      console.error(`[Docs] 翻译错误 for key "${key}":`, error);
      return key;
    }
  };

  return {
    t: safeT,
    n: docsI18nInstance.global.n,
    // 【修改】直接返回 i18n 实例内部的响应式 locale
    locale: docsI18nInstance.global.locale,
    // 为保持兼容性，也提供 currentLocale
    currentLocale: docsI18nInstance.global.locale, 
  }
}

// ==================== 向后兼容的导出函数 ====================

/**
 * 导出安装函数，用于在 enhanceApp 中注册 i18n
 */
export function setupI18n(app: any) {
  app.use(docsI18nInstance);
}

// ==================== 转换函数导出 ====================

// 保持原有的转换函数导出，供其他地方使用
export { convertEnumToStandardLocale, convertStandardToEnumLocale };

// 默认导出全局i18n实例
export default docsI18nInstance;
