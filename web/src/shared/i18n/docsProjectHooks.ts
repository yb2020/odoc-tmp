/**
 * docs项目特定的i18n钩子
 * 提供VitePress特定的翻译和格式化逻辑
 */

// 已通过 setLanguageCookie 设置，无需额外导入
import { setLanguageCookie, languageEnumToStandard, standardToLanguageEnum } from '../language/service';

/**
 * 自定义i18n钩子接口
 */
interface I18nCustomHook {
  name: string;
  applyT?: (originalT: any) => any;
  applyN?: (originalN: any) => any;
  setup?: (i18n: any) => any;
}

/**
 * 创建docs项目的自定义翻译函数
 */
export function createDocsProjectTranslator(originalT: any) {
  return function docsCustomTranslator(key: string, ...args: any[]) {
    try {
      const result = originalT(key, ...args);
      
      // docs项目特定的翻译处理逻辑
      if (result === key) {
        console.warn(`[Docs] 翻译键缺失: ${key}`);
      }
      
      return result;
    } catch (error) {
      console.error(`[Docs] 翻译过程中发生错误，键: ${key}`, error);
      return key; // 降级处理，返回原始键
    }
  };
}

/**
 * 同步语言设置到主项目
 */
export function syncLanguageToMainProject(projectEnumLocale: string) {
  // 使用统一语言服务转换为标准格式
  const standardLocale = languageEnumToStandard(projectEnumLocale as any);
  
  // 使用统一的语言管理服务设置Cookie
  // 通过转换验证语言格式的有效性
  try {
    standardToLanguageEnum(standardLocale);
    setLanguageCookie(standardLocale as 'en-US' | 'zh-CN');
  } catch (error) {
    console.warn('[docsProjectHooks] Invalid language format:', standardLocale);
  }
  
  // 通过统一语言管理服务设置全局语言
  // 注意：这里已经通过上面的 setLanguageCookie 设置了 Cookie
  // changeLanguage 会自动同步 i18n 实例，无需重复设置
}

/**
 * docs项目的自定义钩子配置
 */
export const docsProjectHook: I18nCustomHook = {
  name: 'docsProject',
  applyT: createDocsProjectTranslator,
  // docs项目使用主项目的数字格式化逻辑
};
