// 导入共享的i18n核心模块
import globalI18nInstance, { 
  registerCustomHook,
  LanguageCodeMap,
  CodeToLanguageMap,
  convertEnumToStandardLocale,
  convertStandardToEnumLocale,
  setGlobalLocale
} from '../shared/i18n/core';
import { mainProjectHook } from '../shared/i18n/mainProjectHooks';

// 注册主项目的自定义钩子
registerCustomHook(mainProjectHook);

// ==================== 向后兼容的导出 ====================

// 保持原有的导出方式，完全向后兼容
export default globalI18nInstance;

// 导出自定义的翻译和数字格式化函数
export const { t, n } = globalI18nInstance.global;

// 导出语言映射配置
export { LanguageCodeMap, CodeToLanguageMap };

// 导出转换函数
export { convertEnumToStandardLocale, convertStandardToEnumLocale };

// 导出语言设置函数
export { setGlobalLocale };

// ==================== 兼容性函数 ====================

/**
 * 获取当前语言（兼容原有代码）
 */
export function getCurrentLanguage(): string {
  return globalI18nInstance.global.locale.value;
}

/**
 * Vue应用安装插件（兼容原有代码）
 */
export function install(app: any) {
  app.use(globalI18nInstance);
  
  // 保存原始函数引用，用于调试
  app.config.globalProperties.$originalT = globalI18nInstance.global.t;
  app.config.globalProperties.$originalN = globalI18nInstance.global.n;
}

// ==================== 类型定义 ====================

// 保持原有的类型导出
export type I18nMessageType = any;
