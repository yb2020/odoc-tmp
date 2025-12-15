import { getCurrentLanguage } from './core';

/**
 * 创建主项目的自定义翻译函数
 */
export function createMainProjectTranslator(originalT: any) {
  return function customTranslator(key: string, ...args: any[]) {
    try {
      const result = originalT(key, ...args);
      
      // 检查是否为缺失的翻译键
      if (result === key) {
        console.warn(`翻译键缺失: ${key}`);
      }
      
      return result;
    } catch (error) {
      console.error(`翻译过程中发生错误，键: ${key}`, error);
      return key; // 降级处理，返回原始键
    }
  };
}

/**
 * 创建主项目的自定义数字格式化函数
 */
export function createMainProjectNumberFormatter(originalN: any) {
  return function customNumberFormatter(value: number, format: string) {
    try {
      // 获取当前 locale，getCurrentLanguage() 现在直接返回标准格式
      const standardLocale = getCurrentLanguage();
      
      // 使用原生 Intl.NumberFormat 进行格式化
      if (format === 'integer') {
        return new Intl.NumberFormat(standardLocale, {
          style: 'decimal',
          useGrouping: true,
          maximumFractionDigits: 0,
        }).format(value);
      } else if (format === 'percent') {
        return new Intl.NumberFormat(standardLocale, {
          style: 'percent',
          useGrouping: false,
        }).format(value);
      }
      
      // 对于其他格式，尝试使用原始函数
      return originalN(value, format);
    } catch (error) {
      console.error(`数字格式化失败，值: ${value}, 格式: ${format}`, error);
      // 降级处理，返回原始数字的字符串形式
      return value.toString();
    }
  };
}

/**
 * 主项目的自定义钩子配置
 */
export const mainProjectHook = {
  name: 'mainProject',
  applyT: createMainProjectTranslator,
  applyN: createMainProjectNumberFormatter,
};
