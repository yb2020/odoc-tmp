/**
 * 统一的存储键名常量
 * 按照功能模块分类，便于管理和查找
 */

// 命名空间前缀
const APP_PREFIX = 'odoc'

// 主题相关
export const THEME = {
  // 站点主题 (当前使用: 全局主题，目前只有 light/dark)
  SITE: `${APP_PREFIX}:site:theme`,
  // 主题设置 (当前使用: theme)
  PREFERENCE: `${APP_PREFIX}:preference:theme`,
  // PDF阅读器主题 (当前使用: pdf-annotate/2.0/theme)
  PDF_READER: `${APP_PREFIX}:pdf-reader:theme`,
}

// 用户界面相关
export const UI = {
  // 复制标签位置 (当前使用: copilotTabPositioned)
  COPILOT_TAB_POSITION: `${APP_PREFIX}:ui:copilot-tab-position`,
  // 纸张列表头部 (当前使用: paperHeadList)
  PAPER_HEAD_LIST: `${APP_PREFIX}:ui:paper-head-list`,
  // 纸张列表页面大小 (当前使用: paperListPageSize)
  PAPER_LIST_PAGE_SIZE: `${APP_PREFIX}:ui:paper-list-page-size`,
  // 纸张列表总数 (当前使用: paperListTotal)
  PAPER_LIST_TOTAL: `${APP_PREFIX}:ui:paper-list-total`,
  // 纸张列表排序 (当前使用: paperListSort)
  PAPER_LIST_SORT: `${APP_PREFIX}:ui:paper-list-sort`,
  // 纸张列表排序方向 (当前使用: paperListSortDirection)
  PAPER_LIST_SORT_DIRECTION: `${APP_PREFIX}:ui:paper-list-sort-direction`,
}

// PDF注释相关
export const PDF_READER = {
  // 术语表 (当前使用: pdf-reader/glossary)
  GLOSSARY: `${APP_PREFIX}:pdf-reader:glossary`,
  // 设置 (当前使用: pdf-reader/settings)
  SETTINGS: `${APP_PREFIX}:pdf-reader:settings`,
  // 提示 (当前使用: pdf-reader/tippy)
  TIPPY: `${APP_PREFIX}:pdf-reader:tippy`,
  // 翻译字体大小 (当前使用: pdf-reader/translate-font-size)
  TRANSLATE_FONT_SIZE: `${APP_PREFIX}:pdf-reader:translate-font-size`,
  // 翻译锁定 (当前使用: pdf-reader/translate-lock)
  TRANSLATE_LOCK: `${APP_PREFIX}:pdf-reader:translate-lock`,
  // 翻译标签 (当前使用: pdf-reader/translate-tab)
  TRANSLATE_TAB: `${APP_PREFIX}:pdf-reader:translate-tab`,
  // 翻译标签重置持续时间 (当前使用: pdf-reader/translate-tab-reset-duration)
  TRANSLATE_TAB_RESET_DURATION: `${APP_PREFIX}:pdf-reader:translate-tab-reset-duration`,
  // 翻译标签重置时间 (当前使用: pdf-reader/translate-tab-reset-time)
  TRANSLATE_TAB_RESET_TIME: `${APP_PREFIX}:pdf-reader:translate-tab-reset-time`,
  // 注释 (当前使用: pdf-reader/comment)
  COMMENT: `${APP_PREFIX}:pdf-reader:comment`,
  // 显示图片组选择 (当前使用: pdf-reader/show-image-group-select)
  SHOW_IMAGE_GROUP_SELECT: `${APP_PREFIX}:pdf-reader:show-image-group-select`,
  // 复制建议标签 (当前使用: pdf-reader/copilot-suggestion-tab)
  COPILOT_SUGGESTION_TAB: `${APP_PREFIX}:pdf-reader:copilot-suggestion-tab`,
}


// 用户中心相关
export const USER_CENTER = {
  // 文献排序 (当前使用: user-center/literature/sort)
  LITERATURE_SORT: `${APP_PREFIX}:user-center:literature-sort`,
  // 文献排序方向 (当前使用: user-center/literature/sort-direction)
  LITERATURE_SORT_DIRECTION: `${APP_PREFIX}:user-center:literature-sort-direction`,
}

// 国际化相关
export const I18N = {
  // 语言设置 (当前使用: i18n)
  LANGUAGE: `${APP_PREFIX}:i18n:language`,
}

// 开发工具相关
export const DEV_TOOLS = {
  // ESM插件设置 (当前使用: __VUE_DEVTOOLS_NEXT_PLUGIN_SETTINGS__dev.esm.pinia__)
  PINIA_SETTINGS: `${APP_PREFIX}:dev-tools:pinia-settings`,
}

// 系统相关
export const SYSTEM = {
  // 设备ID (当前使用: REPORT_DEVICE_ID)
  DEVICE_ID: `${APP_PREFIX}:system:device-id`,
  // 笔记词汇量大小 (当前使用: NOTE_VOCABULARY_SIZE)
  NOTE_VOCABULARY_SIZE: `${APP_PREFIX}:system:note-vocabulary-size`,
}

// 用户认证相关
export const AUTH = {
  // 本地登录类型 (当前使用: localLoginTpye)
  LOGIN_TYPE: `${APP_PREFIX}:auth:login-type`,
  // 本地图标 (当前使用: localI18nIcon)
  I18N_ICON: `${APP_PREFIX}:auth:i18n-icon`,
}

// 全局消息相关
export const MESSAGE = {
  // 全局警告消息 (当前使用: global-alert-message)
  GLOBAL_ALERT: `${APP_PREFIX}:message:global-alert`,
}
