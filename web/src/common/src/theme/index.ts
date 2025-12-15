/**
 * 统一的主题管理工具
 * 供主项目和 VitePress 共同使用
 */

import { THEME } from '../constants/storage-keys'
import { localStore } from '../services/storage'

// 导出主题类型
export enum ThemeType {
  beige = 'beige',
  green = 'green',
  dark = 'dark'
}

export type ThemeValue = ThemeType | 'default' | 'light'

/**
 * 获取当前主题
 * 按优先级从不同存储位置获取
 */
export const getCurrentTheme = (): ThemeValue => {
  // 优先级：站点主题 > 通用主题键 > 默认
  let theme = localStore.get(THEME.SITE) || 
              localStore.get('theme') || 
              'default'
              
  return theme as ThemeValue
}

/**
 * 设置主题
 * 同时更新到多个存储位置以确保兼容性
 */
export const setTheme = (theme: ThemeValue) => {
  // 设置到主项目的标准存储键
  localStore.set(THEME.SITE, theme)
  
  // 为了兼容性，也设置到通用主题键
  localStore.set('theme', theme)
  
  // 应用主题到 DOM
  applyThemeToDOM(theme)
  
  console.log(`主题已设置: ${theme}`)
}

/**
 * 应用主题到 DOM
 */
export const applyThemeToDOM = (theme: ThemeValue) => {
  const html = document.documentElement
  const body = document.body

  // 清除所有主题相关的类和属性
  html.classList.remove('dark')
  html.removeAttribute('data-theme')
  body.removeAttribute('data-theme')

  // 应用新主题
  if (theme === 'dark') {
    html.classList.add('dark')
    html.setAttribute('data-theme', 'dark')
    body.setAttribute('data-theme', 'dark')
  } else {
    const themeValue = theme === 'default' ? 'light' : theme
    html.setAttribute('data-theme', themeValue)
    body.setAttribute('data-theme', themeValue)
  }
}

/**
 * 监听主题变化
 */
export const onThemeChange = (callback: (theme: ThemeValue) => void) => {
  const handleStorageChange = (event: StorageEvent) => {
    if (event.key === THEME.SITE || event.key === 'theme') {
      const newTheme = getCurrentTheme()
      callback(newTheme)
    }
  }

  window.addEventListener('storage', handleStorageChange)
  
  // 返回清理函数
  return () => {
    window.removeEventListener('storage', handleStorageChange)
  }
}

/**
 * 初始化主题
 * 从存储中读取主题并应用到 DOM
 */
export const initTheme = () => {
  if (typeof window !== 'undefined') {
    const theme = getCurrentTheme()
    applyThemeToDOM(theme)
    console.log(`主题已初始化: ${theme}`)
    return theme
  }
  return 'default'
}

/**
 * 检查是否为深色主题
 */
export const isDarkTheme = (theme?: ThemeValue): boolean => {
  const currentTheme = theme || getCurrentTheme()
  return currentTheme === 'dark'
}