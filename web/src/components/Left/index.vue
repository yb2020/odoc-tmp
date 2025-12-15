<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { LeftSideBarType } from './type';
import Catalog from './Catalog/index.vue';
import Thumbnail from './Thumbnails.vue';
import Drawer from '../Common/Drawer.vue';
import { CloseOutlined } from '@ant-design/icons-vue';
import { computed } from 'vue';
import { SideTabCommonSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { useStore } from '@/store';
import { getPlatformKey } from '../../store/shortcuts';
import { PAGE_ROUTE_NAME } from '../../routes/type';
import { useI18n } from 'vue-i18n';
import useShortcuts from '~/src/hooks/useShortcuts';

const props = defineProps<{
  pdfId: string;
  sideTabSettings: SideTabCommonSettings;
  setSideTabSetting: (values: Partial<SideTabCommonSettings>) => void;
  pdfViewInstance: ViewerController;
}>();

const { t } = useI18n();

const tabs = computed(() => {
  return [
    {
      key: LeftSideBarType.Thumbnail,
      icon: 'iconsuolvetupailie',
      title: t('viewer.thumbnails'),
    },
    {
      key: LeftSideBarType.Catalog,
      icon: 'iconmulu',
      title: t('viewer.catalogue'),
    },
  ];
});

const changeTab = (tab: (typeof tabs.value)[0]) => {
  props.setSideTabSetting({ tab: tab.key });
};

const drawerRef = ref();

const handleClose = () => {
  drawerRef.value.handleVisibleChange(false);
};

const activeTab = computed(() => {
  return props.sideTabSettings.tab === LeftSideBarType.Catalog
    ? LeftSideBarType.Catalog
    : LeftSideBarType.Thumbnail;
});

/**
 * vue3 watch监听props内属性的值的变化 无响应情况分析
 * https://blog.csdn.net/wuyxinu/article/details/124477647?spm=1001.2101.3001.6661.1&utm_medium=distribute.pc_relevant_t0.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-124477647-blog-121246917.t0_edu_mlt&depth_1-utm_source=distribute.pc_relevant_t0.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-124477647-blog-121246917.t0_edu_mlt&utm_relevant_index=1
 */
watch(props.sideTabSettings, (newVal) => {
  if (newVal.shown && newVal.tab === LeftSideBarType.Thumbnail) {
    setTimeout(() => {
      props.pdfViewInstance.getThumbnailViewer()?.forceRendering();
    }, 1000);
  }
});

const onVisibleChange = (visible: boolean) => {
  props.setSideTabSetting({ shown: visible });
};

const onWidthChange = (width: number) => {
  props.setSideTabSetting({ width });
};

const store = useStore();
const shortcutsConfig = computed(
  () => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {}
);
const platformKey = getPlatformKey();
const shortcut = computed(
  () => shortcutsConfig.value.shortcuts.TOGGLE_CATALOG.value[platformKey]
);
const opts = computed(() => ({ scope: shortcutsConfig.value.scope || 'all' }));
useShortcuts(
  shortcut,
  () => onVisibleChange(!props.sideTabSettings.shown),
  opts
);
</script>

<template>
  <Drawer
    ref="drawerRef"
    v-model:visible="sideTabSettings.shown"
    dir="left"
    :min-width="160"
    :shortcut="shortcut"
    :initial-width="sideTabSettings.width"
    class="left-drawer"
    @visible-change="onVisibleChange"
    @width-change="onWidthChange"
  >
    <div class="tabs">
      <div
        v-for="tab in tabs"
        :key="tab.icon"
        :class="['tab', activeTab === tab.key ? 'active' : '']"
        @click="changeTab(tab)"
      >
        {{ tab.title }}
      </div>
      <span class="close-icon">
        <close-outlined
          style="cursor: pointer"
          @click="handleClose"
        />
      </span>
    </div>
    <keep-alive>
      <Catalog
        v-if="activeTab === LeftSideBarType.Catalog"
        class="container"
        :pdf-id="pdfId"
        :pdfViewInstance="pdfViewInstance"
      />
      <Thumbnail
        v-else
        class="container"
        :pdf-id="pdfId"
        :pdfViewInstance="pdfViewInstance"
      />
    </keep-alive>
  </Drawer>
</template>

<style scoped lang="less">
.left-drawer {
  background-color: var(--site-theme-pdf-panel);
}

.tabs {
  display: flex;
  align-items: center;
  position: relative;
  padding: 0 5px 12px;
  overflow: hidden;

  .tab {
    font-size: 13px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    color: var(--site-theme-text-secondary);
    padding: 8px 10px;
    margin: 0 0px;
    transition: all 0.3s;

    &.active {
      color: var(--site-theme-pdf-tab-active-text);
      background-color: var(--site-theme-pdf-tab-active-bg);
      font-weight: 500;

      .iconfont {
        color: var(--site-theme-pdf-tab-active-text);
      }
    }
    
    &:hover:not(.active) {
      background-color: var(--site-theme-bg-hover);
      opacity: 0.7;
    }
  }

  .close-icon {
    position: absolute;
    right: 0;
    padding: 12px;
    color: var(--site-theme-pdf-panel-text);
    
    &:hover {
      color: var(--site-theme-text-primary);
    }
  }
}

.container {
  height: calc(100% - 60px);
}
</style>
