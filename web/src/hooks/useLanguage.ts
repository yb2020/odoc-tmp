import { computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { 
  switchLanguage as serviceSwitchLanguage, 
  getLanguageCookie, 
  setLanguageCookie,
  languageEnumToStandard, 
  standardToLanguageEnum,
  StandardLanguage,
  getDefaultLanguage,
  isCurrentLanguage,
  getCurrentLanguageEnum
} from '@/shared/language/service';

/**
 * 统一的语言管理 Hook
 * 提供语言切换、获取当前语言、语言判断等功能
 */
export function useLanguage() {
  const { locale } = useI18n();

  /**
   * 获取当前语言（Proto 枚举格式）
   */
  const getCurrentLanguage = (): Language => {
    const currentStandard = locale.value as StandardLanguage;
    const protoEnum = standardToLanguageEnum(currentStandard);
    return protoEnum || Language.EN_US; // 默认英文
  };

  /**
   * 获取当前语言（标准格式）
   */
  const getCurrentStandardLanguage = (): StandardLanguage => {
    return locale.value as StandardLanguage;
  };

  /**
   * 判断当前是否为指定语言
   */
  const isCurrentLanguage = (lang: Language): boolean => {
    const standardLang = languageEnumToStandard(lang);
    return locale.value === standardLang;
  };

  /**
   * 判断当前是否为中文
   */
  const isZhCN = computed(() => isCurrentLanguage(Language.ZH_CN));

  /**
   * 判断当前是否为英文
   */
  const isEnUS = computed(() => isCurrentLanguage(Language.EN_US));

  /**
   * 切换语言并同步所有状态
   * @param protoLang Proto 枚举语言
   */
  const switchLanguage = (protoLang: Language): void => {
    const standardLang = languageEnumToStandard(protoLang);
    if (!standardLang) {
      console.error(`[useLanguage] 无效的语言枚举:`, protoLang);
      return;
    }

    // 更新状态
    locale.value = standardLang;
    
    // 只在浏览器环境中设置 Cookie 和进行同步验证
    if (typeof window !== 'undefined') {
      // 设置 Cookie
      setLanguageCookie(standardLang);
      
      // 开发环境下输出详细日志
      if (process.env.NODE_ENV === 'development') {
        console.log(`[useLanguage] 语言切换完成:`, {
          Proto: Language[protoLang],
          Standard: standardLang,
          Cookie: getLanguageCookie(),
          Locale: locale.value,
        });
      }
      
      // 立即验证同步状态（仅在浏览器环境）
      setTimeout(() => {
        const cookieValue = getLanguageCookie();
        const localeValue = locale.value;
        
        // 开发环境下输出同步验证日志
        if (process.env.NODE_ENV === 'development') {
          console.log(`[useLanguage] 同步验证:`, {
            期望语言: standardLang,
            Cookie值: cookieValue,
            Locale值: localeValue,
            是否同步: cookieValue === standardLang && localeValue === standardLang
          });
        }
        
        // 只有在真正不同步时才输出错误日志
        if (cookieValue !== standardLang || localeValue !== standardLang) {
          console.warn(`[useLanguage] 语言切换同步异常:`, {
            期望: standardLang,
            Cookie: cookieValue,
            Locale: localeValue
          });
        }
      }, 100); // 100ms 后验证
    } else {
      // 服务端环境下只输出简单日志
      if (process.env.NODE_ENV === 'development') {
        console.log(`[useLanguage] 服务端语言设置:`, {
          Proto: Language[protoLang],
          Standard: standardLang
        });
      }
    }
  };

  /**
   * 切换语言（统一入口）
   * @param lang Proto 语言枚举
   */
  const changeLanguage = (lang: Language): void => {
    if (process.env.NODE_ENV === 'development') {
      console.log(`[useLanguage] 开始切换语言到: ${Language[lang]} (${lang})`);
    }
    
    // 使用内部的 switchLanguage 函数
    switchLanguage(lang);
  };

  /**
   * 初始化语言设置
   * 从 Cookie 读取语言设置，如果没有则使用默认语言
   */
  const initializeLanguage = (): void => {
    const cookieLang = getLanguageCookie();
    if (cookieLang) {
      locale.value = cookieLang;
    } else {
      // 没有 Cookie 时使用默认语言并设置 Cookie
      const defaultLang = getDefaultLanguage();
      locale.value = defaultLang;
      
      // 将默认语言转换为 Proto 枚举并设置
      const protoEnum = standardToLanguageEnum(defaultLang);
      if (protoEnum) {
        switchLanguage(protoEnum);
      }
    }
  };

  return {
    // 状态
    locale,
    isZhCN,
    isEnUS,
    
    // 方法
    getCurrentLanguage,
    getCurrentStandardLanguage,
    isCurrentLanguage,
    changeLanguage,
    initializeLanguage
  };
}
