import { PDFViewerColorToneAddon } from '@idea/pdf-annotate-viewer'
import { message } from 'ant-design-vue'
import { debounce } from 'lodash-es'
import { localStore } from '../common/src/services/storage'
import { THEME } from '../common/src/constants/storage-keys'

const PDFViewerThemeAddonInstance = new PDFViewerColorToneAddon()

export enum ThemeType {
  beige = 'beige',
  green = 'green',
  dark = 'dark'
}

/**
 * 切换PDF阅读器主题
 * 只更新PDF阅读器相关的主题设置，不影响站点全局主题
 */
export const changePageTheme = (theme: ThemeType | 'default') => {
  PDFViewerThemeAddonInstance.changeReaderColorTone(theme)
  document.documentElement.setAttribute('data-theme', theme)
  // 只更新PDF阅读器的主题设置
  localStore.set(THEME.PDF_READER, theme)
}

/**
 * 切换全局主题
 * 更新全局主题设置，但不影响PDF阅读器主题
 */
export const changeSiteTheme = (theme: ThemeType | 'default') => {
  document.documentElement.setAttribute('data-theme', theme)
  // 使用THEME.SITE作为全站统一主题键名
  localStore.set(THEME.SITE, theme)
}

/**
 * 同时切换PDF阅读器和站点主题
 */
export const changeAllTheme = (theme: ThemeType | 'default') => {
  PDFViewerThemeAddonInstance.changeReaderColorTone(theme)
  document.documentElement.setAttribute('data-theme', theme)
  // 更新PDF阅读器主题
  localStore.set(THEME.PDF_READER, theme)
  // 更新全局主题设置，使用THEME.SITE作为统一键名
  localStore.set(THEME.SITE, theme)
}

/**
 * 同时切换PDF阅读器主题
 */
export const changeReaderTheme = (theme: ThemeType | 'default') => {
  PDFViewerThemeAddonInstance.changeReaderColorTone(theme)
  // 更新PDF阅读器主题
  localStore.set(THEME.PDF_READER, theme)
}

export const debounceChangePageTheme = debounce(changePageTheme, 200, {
  leading: true,
})

const initPageTheme = () => {
  let theme: ThemeType | 'default' | undefined
  try {
    // 使用新的键名获取主题
    theme = localStore.get(THEME.PDF_READER) as ThemeType
    if (!theme) {
      // 如果没有找到PDF阅读器主题，尝试使用全局主题
      theme = localStore.get(THEME.SITE) as ThemeType
    }
  } catch (error) {
    //
  }
  
  if (theme) {
    document.documentElement.setAttribute('data-theme', theme)
  }

  PDFViewerThemeAddonInstance.initial(theme === 'default' ? undefined : theme);
}

initPageTheme()

export const switchThemeColorToDefault = () => {
  let theme: ThemeType | 'default' | undefined
  try {
    // 使用新的键名获取主题
    theme = localStore.get(THEME.PDF_READER) as ThemeType
  } catch (error) {
    //
  }
  if (theme === ThemeType.dark) {
    message.error('全文翻译暂不支持深色模式')
    debounceChangePageTheme('default')
  }
}

export default PDFViewerThemeAddonInstance