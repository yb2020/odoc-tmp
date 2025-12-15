import { createApp } from 'vue';
import App from './pages/noteManage/index.vue';
import Antd from 'ant-design-vue';
import './assets/less/tailwind.css';
import './assets/less/font.less';

import 'tippy.js/dist/tippy.css';

// 替代阿里 iconfont
import '@idea/aiknowledge-icon/dist/css/iconfont.css';

import './assets/less/antd.light.less';
import { usePinaStore } from './stores';
import i18n from './locals/i18n';
import { initClarity } from './utils/clarity';

// 在应用初始化之前启动 Clarity
if (typeof window !== 'undefined') {
  initClarity();
}

const initApp = (container: string | Element = '#app') => {
  const app = createApp(App);

  const pinaStore = usePinaStore();

  app.use(pinaStore);

  app.use(Antd);

  app.use(i18n);

  app.mount(container);

  return app;
};

initApp();
