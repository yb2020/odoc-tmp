<template>
  <div class="header-wrap">
    <div class="header">
      <!-- Logo -->
      <div class="logo-container relative">
        <a
          v-if="envStore.viewerConfig.logoGoHomeAnchor !== false"
          :href="homeUrl"
          :title="$t('viewer.goHomePage')"
          target="_blank"
          @click="reportClick(ElementClick.go_home)"
          class="logo-text inline-block"
        >
          <img :src="faviconIco" alt="ODOC.AI logo" class="site-logo" width="20" height="20" />
          <span class="site-title">ODOC.AI</span>
        </a>
        <span class="site-badge">
          BETA
        </span>
      </div>
      
      <slot name="title">
        <Title v-if="!IS_MOBILE && !IS_ELECTRON_MODE" />
      </slot>
      
      <!-- <div class="spacer"></div> -->
      
      <!-- 右侧操作区域 -->
      <div class="actions-container">
        <!-- 导航菜单 -->
        <div class="nav-menu">
          <a :href="docsHomeLink" class="nav-item" :class="{ active: isActivePath('/docs/') || isActivePath('/docs/zh/') }">{{ t('nav.home') }}</a>
          <a :href="docsWorkbenchLink" class="nav-item" :class="{ active: isActivePath('/workbench') }">{{ t('nav.workbench') }}</a>
          <a :href="docsGuideLink" class="nav-item" :class="{ active: isActivePath('/docs/guide') || isActivePath('/docs/zh/guide') }">{{ t('nav.guide') }}</a>
          <a :href="docsPricingLink" class="nav-item" :class="{ active: isActivePath('/docs/pricing') || isActivePath('/docs/zh/pricing') }">{{ t('nav.pricing') }}</a>
        </div>
        
        <slot name="actions"></slot>
        
        <!-- 分隔线 -->
        <div class="divider"></div>
        
        <!-- 语言切换 -->
        <div class="i18n-container">
          <a-dropdown 
            :trigger="['hover']" 
            placement="bottomRight" 
            overlay-class-name="i18n-dropdown"
            @visible-change="handleDropdownVisibleChange"
          >
            <div class="i18n-icon globe-icon">
              <img :src="globeIcon" alt="language" width="20" height="20" />
              <span class="vpi-chevron-down" :class="{ 'rotated': isDropdownVisible }"></span>
            </div>
            <template #overlay>
              <a-menu class="custom-language-menu">
                <a-menu-item key="en-US" @click="switchLanguage(Language.EN_US)">
                  English
                </a-menu-item>
                <a-menu-item key="zh-CN" @click="switchLanguage(Language.ZH_CN)">
                  简体中文
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>

        <!-- 分隔线 -->
        <div class="divider"></div>
        
        <!-- 主题切换开关 -->
        <div class="theme-switch">
          <a-switch
            :checked="currentTheme === 'dark'"
            class="theme-toggle"
            @change="toggleTheme"
            size="default"
          >
            <template #checkedChildren>
              <span class="icon-wrapper sun">
                <SunOutlined />
              </span>
            </template>
            <template #unCheckedChildren>
              <span class="icon-wrapper moon">
                <MoonOutlined />
              </span>
            </template>
          </a-switch>
        </div>

        <!-- 分隔线 -->
        <div class="divider"></div>
        
        <!-- GitHub 图标 -->
        <!-- <a href="https://github.com/your-repo" target="_blank" class="github-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
          </svg>
        </a> -->
        
        <!-- 用户头像 -->
        <div class="avatar-container">
          <a-dropdown :trigger="['click']" placement="bottomRight">
            <a-avatar :size="32" class="user-avatar">
              <template #icon><UserOutlined /></template>
            </a-avatar>
            <template #overlay>
              <a-menu>
                <!--<a-menu-item key="profile" @click="handleProfile">
                  <ProfileOutlined />
                  <span>{{ $t('user.profile') }}</span>
                </a-menu-item>
                <a-menu-item key="settings" @click="handleSettings">
                  <SettingOutlined />
                  <span>{{ $t('user.settings') }}</span>
                </a-menu-item>
                <a-menu-divider /> -->
                <a-menu-item key="logout" @click="handleLogout">
                  <LogoutOutlined />
                  <span>{{ $t('user.logout') }}</span>
                </a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
  </div>
  <GlobalAlert />
</template>

<script lang="ts" setup>
import { computed, onMounted, ref, unref, watch } from 'vue';
import { message } from 'ant-design-vue';
import {
  UserOutlined,
  SettingOutlined,
  LogoutOutlined,
  ProfileOutlined,
} from '@ant-design/icons-vue';
import SunOutlined from '../icons/SunOutlined.vue';
import MoonOutlined from '../icons/MoonOutlined.vue';
import { useI18n } from 'vue-i18n';
import { store } from '~/src/store';
import { useEnvStore } from '~/src/stores/envStore';
import { useLanguage } from '~/src/hooks/useLanguage';
import { ElementClick, reportClick } from '~/src/api/report';
import { bridgeAdaptor } from '~/src/adaptor/bridge';
import {
  getDomainOrigin,
  getHostname,
  IS_ELECTRON_MODE,
  IS_MOBILE,
  isInElectron,
} from '~/src/util/env';
import { goPathPage } from '~/src/common/src/utils/url';
import { doLogout } from '~/src/api/user';

import Title from '../Head/components/Title.vue';
import GlobalAlert from '../Head/components/GlobalAlert.vue';
import Cookies from 'js-cookie';
import { COOKIE_REFRESH_TOKEN } from '~/src/api/const';
import { THEME } from '@/common/src/constants/storage-keys';
import {changeSiteTheme, ThemeType} from '@/theme';
import { localStore } from '~/src/common/src/services/storage';
import I18nOutlined from '../icons/I18nOutlined.vue';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import globeIcon from '@/common/assets/images/language/language.svg';
import faviconIco from '@/common/assets/images/favicon_no_bg.ico';

const i18n = useI18n();
const { t } = i18n;
const { isZhCN, changeLanguage } = useLanguage();
const userInfo = computed(() => store.state.user.userInfo);
const homeUrl = isInElectron() ? `https://${getHostname()}/` : '/';

// 动态文档链接
const isZhLocale = isZhCN;

const docsGuideLink = computed(() => {
  const result = isZhLocale.value ? '/docs/zh/guide' : '/docs/guide';
  return result;
});

const docsPricingLink = computed(() => {
  const result = isZhLocale.value ? '/docs/zh/pricing' : '/docs/pricing';
  return result;
});

const docsHomeLink = computed(() => {
  const result = isZhLocale.value ? '/docs/zh/' : '/docs/';
  return result;
});

const docsWorkbenchLink = computed(() => {
  const result = '/workbench';
  return result;
});

// 当前路径是否激活
const isActivePath = (path: string) => {
  if (typeof window !== 'undefined') {
    return window.location.pathname.startsWith(path);
  }
  return false;
};

// 主题切换
const currentTheme = ref('light');

// 在组件挂载时初始化主题
onMounted(() => {
  // 从本地存储获取当前主题
  const savedTheme = localStore.get(THEME.SITE);
  currentTheme.value = savedTheme || 'light';
  
  // 确保DOM主题属性与当前主题一致
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', currentTheme.value);
  }
});

const toggleTheme = () => {
  const newTheme = currentTheme.value === 'light' ? 'dark' : 'light';
  currentTheme.value = newTheme;
  changeSiteTheme(newTheme as ThemeType);
  
  // 如果有全局主题切换函数，也调用它
  if (typeof window !== 'undefined' && (window as any).changeSiteTheme) {
    (window as any).changeSiteTheme(newTheme);
  }
};

const handleProfile = () => {
  if (!unref(userInfo)?.id) {
    bridgeAdaptor.login();
    return;
  }
  goPathPage(`${getDomainOrigin()}/user/${unref(userInfo)?.id}`);
};

const handleSettings = () => {
  goPathPage(`${getDomainOrigin()}/settings`);
};

const handleLogout = async () => {
  try {
    const result = await doLogout();
    message.success(t('user.logoutSuccess'));
    if (result) {
      Cookies.remove(COOKIE_REFRESH_TOKEN);
      window.location.href = getDomainOrigin();
    }
    // 重新加载页面以清除用户状态
    window.location.reload();
  } catch (error) {
    message.error(t('user.logoutFailed'));
    console.error('Logout error:', error);
  }
};

const envStore = useEnvStore();

// 语言切换函数
const switchLanguage = (lang: Language) => {
  console.log('Switching language to:', lang);
  // 使用统一的语言切换服务
  changeLanguage(lang);
};

const isDropdownVisible = ref(false);

const handleDropdownVisibleChange = (visible: boolean) => {
  isDropdownVisible.value = visible;
};
</script>

<style lang="less" scoped>
.header-wrap {
  position: sticky;
  top: 0;
  z-index: 100;
  width: 100%;
  background-color: var(--site-theme-bg-primary);
  border-bottom: 1px solid var(--site-theme-divider);
}

.header {
  display: flex;
  align-items: center;
  height: 60px;
  padding: 0 16px;
  margin: 0 auto;
  width: 100%;
}

.logo-container {
  display: flex;
  align-items: center;
  
  .logo-text {
    font-size: 18px;
    font-weight: 600;
    color: var(--site-theme-text-primary);
    text-decoration: none;
    display: flex;
    align-items: center;
  }
  
  .site-logo {
    width: 20px;
    height: 20px;
    margin-right: 4px;
    display: inline-block;
  }
  
  .site-badge {
    display: inline-block;
    margin-left: 4px;
    position: relative;
    top: -4px;
    padding: 2px 2px; /* py-0.5 px-1.5 */
    font-size: 10px;
    line-height: 1;
    font-weight: 500;
    color: #ffffff;
    background-color: #6d6d6d; /* Tailwind blue-400 */
    border-radius: 2px;
    user-select: none;
    pointer-events: none;
    z-index: 10;
  }
}

.spacer {
  flex: 1;
}

.actions-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

.nav-menu {
  display: flex;
  align-items: center;
  margin-right: 8px;
  
  .nav-item {
    padding: 0 12px;
    font-size: 14px;
    color: var(--site-theme-text-primary);
    text-decoration: none;
    height: 60px;
    line-height: 60px;
    transition: color 0.3s;
    
    &:hover {
      color: var(--site-theme-brand);
    }
    
    &.active {
      color: var(--site-theme-brand);
    }
  }
}

.divider {
  height: 20px;
  width: 1px;
  background-color: var(--site-theme-divider);
}

.i18n-container {
  cursor: pointer;
  
  .i18n-icon {
    display: flex;
    align-items: center;
    color: var(--site-theme-text-primary);
    
    .down-icon {
      display: inline-flex;
      margin-left: 4px;
    }
    
    &:hover {
      color: var(--site-theme-brand);
    }
  }
  
  .language-text {
    font-size: 14px;
    margin-right: 4px;
  }
  
  .globe-icon {
    width: 20px;
    height: 20px;
    margin: 0 8px;
    display: flex;
    align-items: center;
    gap: 4px;
    
    /* 语言图标在暗黑模式下的适配 */
    img {
      filter: var(--site-theme-icon-filter, none);
    }
    
    /* 暗黑模式下直接设置图标反转 */
    html[data-theme='dark'] & img {
      filter: invert(1);
    }
    
    /* 鼠标悬停时的语言图标颜色变化 */
    &:hover img {
      filter: brightness(0) saturate(100%) invert(35%) sepia(85%) saturate(1800%) hue-rotate(215deg) brightness(95%) contrast(130%);
    }
    
    /* 暗黑模式下鼠标悬停时的语言图标颜色变化 */
    html[data-theme='dark'] &:hover img {
      filter: brightness(0) saturate(100%) invert(35%) sepia(85%) saturate(1800%) hue-rotate(215deg) brightness(95%) contrast(130%);
    }
  }
  
  /* 模拟 VitePress 的下拉箭头样式 */
  .vpi-chevron-down {
    margin-left: 4px;
    width: 12px;
    height: 12px;
    display: inline-block;
    position: relative;
    transition: transform 0.2s ease;
  }
  
  .vpi-chevron-down.rotated {
    transform: rotate(180deg);
  }
  
  .vpi-chevron-down::before {
    content: "";
    position: absolute;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%) rotate(45deg);
    width: 6px;
    height: 6px;
    border-right: 1px solid var(--site-theme-text-primary);
    border-bottom: 1px solid var(--site-theme-text-primary);
  }
  
  /* 暗黑模式下的箭头样式优化 */
  html[data-theme='dark'] & .vpi-chevron-down::before {
    border-right: 1px solid #ffffff;
    border-bottom: 1px solid #ffffff;
  }
}

.theme-switch {
  display: flex;
  align-items: center;
  
  .theme-toggle {
    width: 48px;
    height: 24px;
    min-width: 48px;
    border-radius: 24px;
    background-color: #f0f0f0;
    border: 1px solid #aaa;
    
    &.ant-switch-checked {
      background-color: #222;
    }
    
    .ant-switch-handle {
      width: 20px;
      height: 20px;
      top: 1px;
      left: 2px;
    }
    
    &.ant-switch-checked .ant-switch-handle {
      left: calc(100% - 22px);
    }
  }
  
  .icon-wrapper {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 16px;
    height: 16px;
    position: relative;
    
    &.sun {
      color: #aaa;
      top: 2px;
    }
    
    &.moon {
      color: #333;
      top: 2px;
    }
    
    :deep(svg) {
      width: 14px;
      height: 14px;
      stroke: currentColor;
      stroke-width: 2.5;
    }
  }
}

.github-icon {
  display: flex;
  align-items: center;
  color: var(--site-theme-text-primary);
  margin-left: 16px;
  
  &:hover {
    color: var(--site-theme-brand);
  }
}

.avatar-container {
  cursor: pointer;
  margin-left: 16px;
  
  .user-avatar {
    background-color: var(--site-theme-brand);
    color: #fff;
  }
}

html[data-theme="dark"] {
  .header-wrap {
    background-color: var(--bg-color, #1f1f1f);
    border-bottom-color: var(--border-color, #303030);
  }
  
  .logo-text,
  .nav-menu .nav-item,
  .i18n-icon,
  .github-icon {
    color: var(--site-theme-text-primary);
  }
}

:global(.i18n-dropdown) {
  min-width: 150px;
}

:global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu) {
  padding: 8px 0;
  border-radius: 12px;
  box-shadow: var(--site-theme-shadow-2);
  background-color: var(--site-theme-bg-primary) !important;
  border: none;
}

:global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu-item) {
  padding: 13px 20px;
  margin: 0;
  color: var(--site-theme-text-primary) !important;
  font-size: 14px;
  line-height: 1.5;
  transition: background-color 0.2s ease;
  background-color: transparent !important;
}

:global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu-item:hover) {
  background-color: var(--site-theme-bg-soft) !important;
  color: var(--site-theme-brand) !important;
}

:global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu-item:first-child) {
  border-top-left-radius: 8px;
  border-top-right-radius: 8px;
}

:global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu-item:last-child) {
  border-bottom-left-radius: 8px;
  border-bottom-right-radius: 8px;
}

/* 暗黑模式下的样式调整 */
html[data-theme="dark"] {
  :global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu) {
    background-color: var(--site-theme-bg-secondary) !important;
  }
  
  :global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu-item) {
    color: var(--site-theme-text-primary) !important;
  }
  
  :global(.ant-dropdown.i18n-dropdown .ant-dropdown-menu-item:hover) {
    background-color: var(--site-theme-bg-mute) !important;
  }
  
  .header-wrap {
    background-color: var(--site-theme-bg-primary);
    border-bottom-color: var(--site-theme-divider);
  }
  
  .logo-text {
    color: var(--site-theme-text-primary);
  }
  
  .nav-menu .nav-item {
    color: var(--site-theme-text-primary);
    
    &:hover, &.active {
      color: var(--site-theme-brand);
    }
  }
  
  .github-icon {
    color: var(--site-theme-text-primary);
    
    &:hover {
      color: var(--site-theme-brand);
    }
  }
}
</style>
