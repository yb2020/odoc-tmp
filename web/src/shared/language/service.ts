import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import Cookies from 'js-cookie';

const I18N_COOKIE_KEY = 'i18n';

// Proto 枚举值到标准格式的映射
const PROTO_TO_STANDARD: Record<string, string> = {
  'EN_US': 'en-US',
  'ZH_CN': 'zh-CN'
};

// 标准格式到 Proto 枚举值的映射
const STANDARD_TO_PROTO: Record<string, string> = {
  'en-US': 'EN_US',
  'zh-CN': 'ZH_CN'
};

// 标准语言类型定义
export type StandardLanguage = 'en-US' | 'zh-CN';
export const SUPPORTED_LANGUAGES: StandardLanguage[] = ['en-US', 'zh-CN'];

/**
 * Proto 枚举转标准格式
 * @param lang Proto 语言枚举
 * @returns 标准格式语言代码
 */
export function languageEnumToStandard(lang: Language): StandardLanguage {
  const protoString = Language[lang]; // 获取 'EN_US' 或 'ZH_CN'
  return PROTO_TO_STANDARD[protoString] as StandardLanguage || 'en-US';
}

/**
 * 标准格式转 Proto 枚举
 * @param standard 标准格式语言代码
 * @returns Proto 语言枚举或 null
 */
export function standardToLanguageEnum(standard: string): Language | null {
  const protoString = STANDARD_TO_PROTO[standard];
  if (!protoString) return null;
  
  // 通过枚举键名获取枚举值
  return Language[protoString as keyof typeof Language] ?? null;
}

/**
 * 验证是否为支持的标准格式
 * @param lang 语言代码
 * @returns 是否为有效的标准语言格式
 */
export function isValidStandardLanguage(lang: string): lang is StandardLanguage {
  return SUPPORTED_LANGUAGES.includes(lang as StandardLanguage);
}

/**
 * 获取语言 Cookie（返回标准格式）
 * @returns 标准格式语言代码或 null
 */
export function getLanguageCookie(): StandardLanguage | null {
  const value = Cookies.get(I18N_COOKIE_KEY);
  
  // 如果是标准格式，直接返回
  if (value && isValidStandardLanguage(value)) {
    return value;
  }
  
  // 如果是旧格式（EN_US, ZH_CN），转换为标准格式
  if (value && PROTO_TO_STANDARD[value]) {
    const standardValue = PROTO_TO_STANDARD[value] as StandardLanguage;
    // 立即更新 Cookie 为标准格式
    setLanguageCookie(standardValue);
    return standardValue;
  }
  
  return null;
}

/**
 * 设置语言 Cookie（存储标准格式）
 * @param standardLang 标准格式语言代码
 */
export function setLanguageCookie(standardLang: StandardLanguage): void {
  // 设置 Cookie 并确保立即生效
  const date = new Date();
  date.setTime(date.getTime() + 365 * 24 * 60 * 60 * 1000); // 1年过期
  
  Cookies.set(I18N_COOKIE_KEY, standardLang, { 
    expires: date,
    path: '/',
    sameSite: 'lax'
  });
  
  // 立即验证 Cookie 是否设置成功
  const verifyValue = Cookies.get(I18N_COOKIE_KEY);
  if (verifyValue !== standardLang) {
    console.warn(`Cookie 设置可能失败: 期望 ${standardLang}, 实际 ${verifyValue}`);
  }
}

/**
 * 统一的语言切换服务
 * @param lang Proto 语言枚举
 * @returns 标准格式语言代码
 */
export function switchLanguage(lang: Language): StandardLanguage {
  const standardLang = languageEnumToStandard(lang); // 转换为标准格式
  setLanguageCookie(standardLang); // 存储标准格式
  return standardLang;
}

/**
 * 判断当前语言是否为指定的 proto 枚举语言
 * @param targetLang 目标 proto 语言枚举
 * @param currentStandardLang 当前标准格式语言（可选，不传则从Cookie读取）
 * @returns 是否匹配
 */
export function isCurrentLanguage(targetLang: Language, currentStandardLang?: string): boolean {
  const current = currentStandardLang || getLanguageCookie();
  if (!current) return false;
  
  const targetStandard = languageEnumToStandard(targetLang);
  return current === targetStandard;
}

/**
 * 获取当前语言对应的 proto 枚举
 * @param currentStandardLang 当前标准格式语言（可选，不传则从Cookie读取）
 * @returns proto 语言枚举或 null
 */
export function getCurrentLanguageEnum(currentStandardLang?: string): Language | null {
  const current = currentStandardLang || getLanguageCookie();
  if (!current) return null;
  
  return standardToLanguageEnum(current);
}

/**
 * 获取默认语言（标准格式）
 * @returns 默认的标准格式语言代码
 */
export function getDefaultLanguage(): StandardLanguage {
  // 统一默认语言为英文，与 core.ts 中的 getCurrentLanguage() 保持一致
  return 'en-US'; // 默认英文
}
