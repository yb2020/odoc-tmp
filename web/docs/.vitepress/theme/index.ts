// docs/.vitepress/theme/index.ts
import DefaultTheme from 'vitepress/theme'
import Antd from 'ant-design-vue'
import type { EnhanceAppContext } from 'vitepress'
import { createPinia } from 'pinia'
// 导入国际化配置
import { setupI18n } from '../../i18n/index'
// 导入统一的 Clarity 服务
import { initClarity } from '../../../src/utils/clarity'
// 导入语言设置函数
import { setLanguageCookie } from '../../../src/shared/language/service'

// 导入主项目的样式
import 'ant-design-vue/dist/antd.css'
import '@/assets/less/theme.less' //导入主项目主题

// 导入自定义CSS文件
import './css/index.css'
// 导入暗黑模式CSS文件
import './css/dark-mode.css'
// 导入全局样式文件
import '../style.css'

// 创建 Pinia 实例
const pinia = createPinia()


import MyLayout from '../../components/MyLayout.vue'

export default {
  ...DefaultTheme,
  Layout: MyLayout,
  enhanceApp({ app,router }: EnhanceAppContext) {
    // 初始化 Microsoft Clarity (仅在浏览器环境)
    if (typeof window !== 'undefined') {
      initClarity();
      
      // 根据当前URL路径同步Cookie语言设置
      const syncLanguageCookie = () => {
        const currentPath = window.location.pathname;
        if (currentPath.startsWith('/docs/zh/')) {
          // 访问中文路径，同步Cookie为中文
          setLanguageCookie('zh-CN');
        } else if (currentPath.startsWith('/docs/')) {
          // 访问英文路径，同步Cookie为英文
          setLanguageCookie('en-US');
        }
      };
      
      // 初始化时同步一次
      syncLanguageCookie();
    }

    // 集成 Ant Design Vue
    app.use(Antd)

    app.use(pinia)

    // 集成国际化
    setupI18n(app)

        // 监听路由变化，同步语言设置
    if (router) {
      // 监听路由变化
      router.onBeforeRouteChange = (to) => {
        // 检查是否是指向主应用的路由
        if (!to.startsWith("/docs")) {
          // 检查是否在浏览器环境中
          if (typeof window !== 'undefined') {
            // 重定向到主应用
            window.location.href = to
            return false // 阻止 VitePress 的默认路由行为
          }
        }
      }
      router.onAfterRouteChange = (to) => {
        // 检查是否在浏览器环境中
        if (typeof window !== 'undefined') {
          // 路由变化时同步Cookie语言设置
          if (to.startsWith('/docs/zh/')) {
            setLanguageCookie('zh-CN');
          } else if (to.startsWith('/docs/')) {
            setLanguageCookie('en-US');
          }
        }
      }
    }
  }
}
