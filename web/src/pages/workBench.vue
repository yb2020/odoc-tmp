<template>
  <div class="workbench-container">
    <!-- 左侧导航 -->
    <div class="sidebar" :class="{ 'collapsed': isSidebarCollapsed }">
      <div class="sidebar-item" :class="{ 'active': currentPath === '/workbench/recent' }" @click="router.push('/workbench/recent')">
        <read-outlined />
        <span v-show="!isSidebarCollapsed">{{ t('workbench.recentlyRead') }}</span>
      </div>
      <div class="sidebar-item" :class="{ 'active': currentPath === '/workbench/library' }" @click="router.push('/workbench/library')">
        <folder-outlined />
        <span v-show="!isSidebarCollapsed">{{ t('workbench.myDocuments') }}</span>
      </div>
      <div class="sidebar-item" :class="{ 'active': currentPath === '/workbench/notes' }" @click="router.push('/workbench/notes')">
        <file-outlined />
        <span v-show="!isSidebarCollapsed">{{ t('workbench.myNotes') }}</span>
      </div>
      <!-- <div class="sidebar-item">
        <translation-outlined />
        <span v-show="!isSidebarCollapsed">{{ t('workbench.myTranslations') }}</span>
      </div>
      <div class="sidebar-item">
        <setting-outlined />
        <span v-show="!isSidebarCollapsed">{{ t('workbench.tagManagement') }}</span>
      </div>
      <div class="sidebar-item">
        <delete-outlined />
        <span v-show="!isSidebarCollapsed">{{ t('workbench.trash') }}</span>
      </div> -->
      
      <div class="sidebar-divider" v-show="!isSidebarCollapsed"></div>
      
      <!-- 专业文档类型 -->
      <NavWebsiteBar v-show="!isSidebarCollapsed"/>
    </div>

    <!-- 收缩按钮 -->
    <div class="collapse-button" @click="toggleSidebar">
      <left-outlined v-if="!isSidebarCollapsed" />
      <right-outlined v-else />
    </div>

    <!-- 右侧内容区 -->
    <div class="content">
      <router-view v-slot="{ Component }">
        <transition name="slide-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, onMounted, defineAsyncComponent } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { message } from 'ant-design-vue';
import { useUserStore } from '@common/stores/user';
import { useI18n } from 'vue-i18n';
// 使用defineAsyncComponent正确处理懒加载组件
const NavWebsiteBar = defineAsyncComponent(() => import('../components/NavBar/NavWebsiteBar.vue'));
import { ThemeType } from '@/theme';
import {
  ReadOutlined,
  FolderOutlined,
  TranslationOutlined,
  FileOutlined,
  SettingOutlined,
  DeleteOutlined,
  PlusOutlined,
  AppstoreOutlined,
  SortAscendingOutlined,
  FilePdfOutlined,
  LeftOutlined,
  RightOutlined
} from '@ant-design/icons-vue';

const router = useRouter();
const route = useRoute();
const userStore = useUserStore();
const { t } = useI18n();

// 侧边栏收缩状态
const isSidebarCollapsed = ref(false);

// 当前路由路径
const currentPath = computed(() => route.path);

// 获取当前主题
const currentTheme = computed(() => {
  if (typeof localStorage !== 'undefined') {
    return localStorage.getItem('theme') || ThemeType.LIGHT;
  }
  return ThemeType.LIGHT;
});

// 判断是否为暗黑模式
const isDarkMode = computed(() => currentTheme.value === ThemeType.DARK);

// 切换侧边栏收缩状态
const toggleSidebar = () => {
  isSidebarCollapsed.value = !isSidebarCollapsed.value;
};
</script>

<style lang="less" scoped>
.workbench-container {
  display: flex;
  height: 100vh;
  background-color: var(--site-theme-background);
  color: var(--site-theme-text-color);
  position: relative;
}

/* 添加过渡动画样式 */
.slide-fade-enter-active {
  transition: all 0.3s ease;
}

.slide-fade-leave-active {
  transition: all 0.3s cubic-bezier(1, 0.5, 0.8, 1);
}

.slide-fade-enter-from {
  transform: translateX(20px);
  opacity: 0;
}

.slide-fade-leave-to {
  transform: translateX(-20px);
  opacity: 0;
}

.sidebar {
  width: 200px;
  background-color: var(--site-theme-background-secondary);
  padding: 20px 0;
  border-right: 1px solid var(--site-theme-border-color);
  overflow-y: auto;
  transition: width 0.3s ease;
  
  &.collapsed {
    width: 60px;
  }

  &-item {
    display: flex;
    align-items: center;
    padding: 10px 20px;
    cursor: pointer;
    transition: all 0.3s;
    font-size: 14px;
    white-space: nowrap;
    height: 44px; /* 固定高度确保图标位置一致 */

    .anticon, i {
      margin-right: 10px;
      font-size: 16px;
      flex-shrink: 0; /* 防止图标被压缩 */
    }

    &:hover {
      background-color: var(--site-theme-background-hover);
    }

    &.active {
      background-color: var(--site-theme-primary-color-fade);
      color: var(--site-theme-primary-color);
      border-left: 3px solid var(--site-theme-primary-color);
    }
  }

  &.collapsed &-item {
    justify-content: center;
    padding: 10px 0;

    .anticon, i {
      margin-right: 0;
    }
  }

  &-divider {
    height: 1px;
    background-color: var(--site-theme-border-color);
    margin: 15px 0;
  }

  &-category {
    padding: 10px 20px;
    font-size: 12px;
    color: var(--site-theme-text-secondary-color);
    text-transform: uppercase;
  }
}

.collapse-button {
  position: absolute;
  left: 190px;
  top: 30px;
  transform: translateY(0);
  width: 20px;
  height: 20px;
  background-color: var(--site-theme-collapse-button-bg);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  z-index: 10;
  transition: left 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  
  .anticon {
    color: var(--site-theme-collapse-button-icon);
    font-size: 12px;
  }
  
  &:hover {
    background-color: var(--site-theme-collapse-button-hover-bg);
  }
}

.sidebar.collapsed + .collapse-button {
  left: 50px;
}

.content {
  flex: 1;
  overflow-y: auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;

  &-left {
    display: flex;
    align-items: center;
    gap: 15px;
  }
}

.upload-btn {
  background-color: var(--site-theme-primary-color);
  border-color: var(--site-theme-primary-color);

  &:hover {
    background-color: var(--site-theme-primary-color-hover);
    border-color: var(--site-theme-primary-color-hover);
  }
}

.document-section {
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;

    h2 {
      font-size: 18px;
      margin: 0;
      color: var(--site-theme-text-color);
    }

    .section-actions {
      display: flex;
      gap: 10px;
    }
  }
}

.time-groups {
  .time-group {
    margin-bottom: 30px;

    .time-label {
      font-size: 14px;
      color: var(--site-theme-text-secondary-color);
      margin-bottom: 10px;
      display: flex;
      align-items: center;

      &::after {
        content: '';
        flex: 1;
        height: 1px;
        background-color: var(--site-theme-border-color);
        margin-left: 10px;
      }
    }
  }
}

.document-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.document-item {
  display: flex;
  align-items: center;
  padding: 12px 15px;
  background-color: var(--site-theme-background-secondary);
  border-radius: 4px;
  transition: all 0.3s;
  cursor: pointer;

  &:hover {
    background-color: var(--site-theme-background-hover);
  }

  .document-icon {
    font-size: 24px;
    margin-right: 15px;
  }

  .document-info {
    flex: 1;

    .document-title {
      font-size: 14px;
      margin-bottom: 5px;
      color: var(--site-theme-text-color);
    }

    .document-meta {
      font-size: 12px;
      color: var(--site-theme-text-secondary-color);
    }
  }

  .document-date {
    font-size: 12px;
    color: var(--site-theme-text-secondary-color);
  }
}

// 自定义 Ant Design 组件样式
:deep(.ant-input-search) {
  .ant-input {
    background-color: var(--site-theme-background-secondary);
    border-color: var(--site-theme-border-color);
    color: var(--site-theme-text-color);

    &:focus {
      border-color: var(--site-theme-primary-color);
    }
  }

  .ant-input-search-button {
    background-color: var(--site-theme-primary-color);
    border-color: var(--site-theme-primary-color);

    &:hover {
      background-color: var(--site-theme-primary-color-hover);
      border-color: var(--site-theme-primary-color-hover);
    }
  }
}

:deep(.ant-checkbox-wrapper) {
  color: var(--site-theme-text-color);

  .ant-checkbox-checked .ant-checkbox-inner {
    background-color: var(--site-theme-primary-color);
    border-color: var(--site-theme-primary-color);
  }
}

:deep(.ant-btn-text) {
  color: var(--site-theme-text-color);

  &:hover {
    background-color: var(--site-theme-background-hover);
    color: var(--site-theme-primary-color);
  }
}
</style>
