import { createI18n, I18nOptions } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { isInOverseaseElectron } from '../../util/env';
import Cookies from 'js-cookie';
import { standardToLanguageEnum, StandardLanguage } from '../language/service';

// 导入主项目语言包
import mainZhJSON from '../../locals/zh-CN.json';
import mainEnJSON from '../../locals/en-US.json';
import zhCommonJSON from '../../common/src/locals/files/zh-CN.json';
import enCommonJSON from '../../common/src/locals/files/en-US.json';

// 导入docs项目语言包
import docsZhCN from '../../../docs/i18n/zh-CN';
import docsEnUS from '../../../docs/i18n/en-US';

// ==================== 语言标识符转换配置 ====================

/**
 * 项目枚举格式 → vue-i18n 标准格式转换映射
 */
export const ENUM_TO_STANDARD_LOCALE_MAP: Record<string, string> = {
  'EN_US': 'en-US',  // 'EN_US' → 'en-US'
  'ZH_CN': 'zh-CN',  // 'ZH_CN' → 'zh-CN'
};

/**
 * vue-i18n 标准格式 → 项目枚举格式转换映射
 */
export const STANDARD_TO_ENUM_LOCALE_MAP: Record<string, string> = {
  'en-US': 'EN_US',  // 'en-US' → 'EN_US'
  'zh-CN': 'ZH_CN',  // 'zh-CN' → 'ZH_CN'
};

/**
 * 将项目枚举格式转换为 vue-i18n 标准格式
 */
export function convertEnumToStandardLocale(enumLocale: string): string {
  const standardLocale = ENUM_TO_STANDARD_LOCALE_MAP[enumLocale];
  if (!standardLocale) {
    console.warn(`未知的枚举语言标识符: ${enumLocale}, 使用默认值 'en-US'`);
    return 'en-US';
  }
  return standardLocale;
}

/**
 * 将 vue-i18n 标准格式转换为项目枚举格式
 */
export function convertStandardToEnumLocale(standardLocale: string): string {
  const enumLocale = STANDARD_TO_ENUM_LOCALE_MAP[standardLocale];
  if (!enumLocale) {
    console.warn(`未知的标准语言标识符: ${standardLocale}, 使用默认值 'EN_US'`);
    return 'EN_US';
  }
  return enumLocale;
}

// ==================== 语言映射配置 ====================

// 更新为标准格式映射（RFC 5646）
export const LanguageCodeMap = {
  [Language.EN_US]: 'en-US',  // 0 → 'en-US'
  [Language.ZH_CN]: 'zh-CN',  // 1 → 'zh-CN'
};

export const CodeToLanguageMap = {
  'en-US': Language.EN_US,  // 'en-US' → 0
  'zh-CN': Language.ZH_CN,  // 'zh-CN' → 1
};

// ==================== 获取当前语言 ====================

export function getCurrentLanguage(): string {
  // 优先从 Cookie 获取
  const cookieLanguage = Cookies.get('i18n');
  if (cookieLanguage) {
    // 现在 Cookie 中存储的是标准格式（'en-US', 'zh-CN'）
    // 检查是否为有效的标准格式语言代码
    try {
      standardToLanguageEnum(cookieLanguage);
      console.log(`[getCurrentLanguage] 从Cookie读取到标准格式语言: ${cookieLanguage}`);
      return cookieLanguage as StandardLanguage;
    } catch (error) {
      // 不是有效的标准格式，继续检查旧格式
    }
    
    // 兼容旧格式（'EN_US', 'ZH_CN'）
    if (cookieLanguage === 'EN_US' || cookieLanguage === 'ZH_CN') {
      console.log(`[getCurrentLanguage] 从Cookie读取到旧格式语言: ${cookieLanguage}，转换为标准格式`);
      // 转换为标准格式
      return cookieLanguage === 'EN_US' ? 'en-US' : 'zh-CN';
    }
  }

  // 统一默认语言为 en-US（所有环境）
  const defaultLang = 'en-US';
  console.log(`[getCurrentLanguage] 使用默认语言: ${defaultLang}`);
  return defaultLang;
}

// ==================== 合并所有语言包 ====================

// 主项目语言包
const mainProjectMessages = {
  // 标准格式（RFC 5646）
  'zh-CN': {
    // 主项目语言包
    ...mainZhJSON,
    ...zhCommonJSON,
    // 添加顶级键以解决 Not found 问题
    library: '文献管理',
    notes: '笔记管理',
    zh: {
      library: '文献管理',
      notes: '笔记管理',
    },
  },
  'en-US': {
    // 主项目语言包
    ...mainEnJSON,
    ...enCommonJSON,
    // 添加顶级键以解决 Not found 问题
    library: 'Library',
    notes: 'Notes',
    en: {
      library: 'Library',
      notes: 'Notes',
    },
  },
};

// docs项目语言包 - 保存但不合并到主项目
export const docsProjectMessages = {
  'zh-CN': docsZhCN,
  'en-US': docsEnUS,
};

// 使用主项目的语言包作为i18n实例的messages
const messages = mainProjectMessages;

// ==================== 数字格式化配置 ====================

const numberFormats: I18nOptions['numberFormats'] = {
  // 标准格式（RFC 5646）
  'en-US': {
    integer: {
      style: 'decimal',
      useGrouping: true,
    },
    percent: {
      style: 'percent',
      useGrouping: false,
    },
  },
  'zh-CN': {
    integer: {
      style: 'decimal',
      useGrouping: true,
    },
    percent: {
      style: 'percent',
      useGrouping: false,
    },
  },
};

// ==================== 日期时间格式化配置 ====================

const datetimeFormats: I18nOptions['datetimeFormats'] = {
  // 标准格式（RFC 5646）
  'en-US': {
    short: {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    },
    long: {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      weekday: 'short',
      hour: 'numeric',
      minute: 'numeric',
    },
    time: {
      hour: 'numeric',
      minute: 'numeric',
    },
    date: {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    },
  },
  'zh-CN': {
    short: {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    },
    long: {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      weekday: 'short',
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    },
    time: {
      hour: 'numeric',
      minute: 'numeric',
      hour12: false,
    },
    date: {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    },
  },
} as const;

// ==================== 创建全局唯一的i18n实例 ====================

const currentLocal = getCurrentLanguage();

const globalI18nInstance = createI18n({
  locale: currentLocal,
  legacy: false,
  globalInjection: true,
  fallbackLocale: 'en-US',
  messages,
  numberFormats,
  datetimeFormats,
});

// ==================== 自定义钩子注册系统 ====================

type CustomHook = {
  name: string;
  applyT?: (originalT: any) => any;
  applyN?: (originalN: any) => any;
};

const registeredHooks: CustomHook[] = [];

/**
 * 注册自定义钩子
 */
export function registerCustomHook(hook: CustomHook) {
  registeredHooks.push(hook);
  applyAllHooks();
}

/**
 * 应用所有注册的自定义钩子
 */
function applyAllHooks() {
  const originalT = globalI18nInstance.global.t;
  const originalN = globalI18nInstance.global.n;
  
  let finalT: any = originalT;
  let finalN: any = originalN;
  
  registeredHooks.forEach(hook => {
    if (hook.applyT) {
      finalT = hook.applyT(finalT);
    }
    if (hook.applyN) {
      finalN = hook.applyN(finalN);
    }
  });
  
  // 使用类型断言来避免严格的类型检查
  (globalI18nInstance.global as any).t = finalT;
  (globalI18nInstance.global as any).n = finalN;
  
  // 重写Vue应用的全局属性，确保$t和$n也使用增强的翻译函数
  if (globalI18nInstance.global.d) {
    // 保存原始函数引用，用于调试
    (globalI18nInstance.global.d as any).$originalT = (globalI18nInstance.global.d as any).$t;
    (globalI18nInstance.global.d as any).$originalN = (globalI18nInstance.global.d as any).$n;
    
    // 重写全局注入的$t和$n属性
    (globalI18nInstance.global.d as any).$t = finalT;
    (globalI18nInstance.global.d as any).$n = finalN;
  }
}

/**
 * 统一的语言设置函数
 * @deprecated 请使用统一语言服务 setLanguageCookie 替代
 */
export function setGlobalLocale(locale: string) {
  console.log(`[setGlobalLocale] 设置全局语言: ${locale}`);
  globalI18nInstance.global.locale.value = locale;
  
  // 使用统一语言服务存储标准格式到 Cookie
  // 注意：这里直接存储标准格式，不再转换为枚举格式
  const date = new Date();
  date.setTime(date.getTime() + 365 * 24 * 60 * 60 * 1000); // 1年过期
  
  Cookies.set('i18n', locale, { 
    expires: date,
    path: '/',
    sameSite: 'lax'
  });
  
  console.log(`[setGlobalLocale] Cookie已设置为标准格式: ${locale}`);
}

export default globalI18nInstance;
