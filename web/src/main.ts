import { createApp } from 'vue';
import App from './App.vue';
import Antd, { message } from 'ant-design-vue';
import router from './routes';
import { useStore } from 'vuex';
import { setupInterceptors } from './api/interceptors';

import PerfectScrollbar from 'vue3-perfect-scrollbar';
import 'vue3-perfect-scrollbar/dist/vue3-perfect-scrollbar.css';

import './assets/less/tailwind.css';
import './assets/less/font.less'

import 'tippy.js/dist/tippy.css';

// 替代阿里 iconfont
import '@idea/aiknowledge-icon/dist/css/iconfont.css';

import './assets/less/antd.less';
import './assets/less/pdf.less';
import './assets/less/style.less';
import { isInElectron } from './util/env';
import { usePinaStore } from './stores';

// 导入主题相关
import './theme'
import './assets/less/theme.less'
import i18n from './locals/i18n';
import { resetLSCurrentTranslateTabKey } from './stores/translateStore';
import { useLanguage } from './hooks/useLanguage';

// 导入 Microsoft Clarity 服务
import { initClarity } from './utils/clarity';

if (!isInElectron()) {
  message.config({
    top: '60px',
  });
}

// 初始化存储服务，包括迁移旧键名数据到新键名
const initStorage = () => {
  console.log('Storage service initialized');
};

// 在应用初始化之前就启动 Clarity
if (!isInElectron() && typeof window !== 'undefined') {
  initClarity();
}

const initApp = (container: string | Element = '#app') => {
  // 初始化存储服务
  initStorage();
  
  // 初始化 Axios 拦截器
  setupInterceptors();
  
  // initSentry();

  const app = createApp(App);
  
  // 在 Vue 应用创建后初始化语言设置
  console.log('[main.ts] 开始初始化语言设置');
  // 直接使用统一语言服务，不依赖 Vue Hook
  import('./shared/language/service').then(({ getLanguageCookie, setLanguageCookie, getDefaultLanguage }) => {
    const cookieLang = getLanguageCookie();
    if (!cookieLang) {
      // 没有 Cookie 时设置默认语言
      const defaultLang = getDefaultLanguage();
      setLanguageCookie(defaultLang);
      console.log(`[main.ts] 设置默认语言: ${defaultLang}`);
    } else {
      console.log(`[main.ts] 读取到语言Cookie: ${cookieLang}`);
    }
    console.log('[main.ts] 语言初始化完成');
  }).catch(error => {
    console.error('[main.ts] 语言初始化失败:', error);
  });

  const pinaStore = usePinaStore()

  app.use(pinaStore);

  // 不要直接使用 useStore 作为插件，这会导致 Vue 警告
  
  app.use(router);

  app.use(Antd);

  app.use(PerfectScrollbar);

  app.use(i18n);

  app.mount(container);

  resetLSCurrentTranslateTabKey();

  return app;
};

import { useUserStore } from './common/src/stores/user';
import { getCurrentLanguage } from './locals/i18n';
import { log } from 'console';

const startup = async () => {
  // 仅在浏览器根路径执行此检查
  if (window.location.pathname === '/') {
        const getRedirectPath = () => {
      // 从 i18n store 获取当前语言设置
      const lang = getCurrentLanguage();
      console.log('lang', lang);
      return lang.toLowerCase().startsWith('zh') ? '/docs/zh/' : '/docs/';
    };

    const pinaStore = usePinaStore();
    const userStore = useUserStore(pinaStore);
    try {
      await userStore.getUserInfo();
      if (!userStore.isLogin()) {
        // 未登录，重定向到对应语言的文档页并停止执行
        window.location.replace(getRedirectPath());
        return;
      }
    } catch (error) { 
      // 获取用户信息失败也视为未登录
      window.location.replace(getRedirectPath());
      return;
    }
  }

  // 检查通过或无需检查，正常初始化应用
  initApp();
};

// 启动应用
startup();

// 同时导出，以便其他地方可以使用
export default initApp;
