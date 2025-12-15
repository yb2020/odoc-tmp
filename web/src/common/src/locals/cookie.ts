// 重新导出统一语言服务，保持向后兼容
// @deprecated 请直接使用 '../../../shared/language/service' 中的函数，这些别名将在未来版本中移除
export { 
  getLanguageCookie as getI18nCookie,
  setLanguageCookie as setI18nCookie,
  type StandardLanguage
} from '../../../shared/language/service';
