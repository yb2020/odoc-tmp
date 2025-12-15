<template>
  <CollapsedBar
    class="toolbar-bt"
    :collapsed="!toolbarSettings.toolBarVisible"
    :style="toolbarStyle"
    :icon-style="toolbarStyle"
    :report-props="{ element: ElementClick.bottom_tool_bar }"
    @toggled="
      (v) =>
        setToolbarSettings({
          toolBarVisible: !v,
        })
    "
  >
    <!-- <LeftSideBtn
      v-if="leftBtnProps"
      :side-tab-settings="leftBtnProps.sideTabSettings"
      @toggle-left-side="toolbarEventHandler"
    /> -->
    <Page
      :num-pages="pageProps.numPages"
      :current-page="pageProps.currentPage"
      @go-to-page="handleEvent"
    />
    <div class="split" />
    <Theme :onChangeTheme="handleChangeTheme" />
    <template v-if="leftBtnProps">
      <GroupEye
        v-if="activeTab === RightSideBarType.Group"
        :pdfId="pdfId"
        :noteId="noteId"
      />
      <Eye
        v-else
        :pdfId="pdfId"
        :noteId="noteId"
        :current-page="pageProps.currentPage"
      />
    </template>
    <FullPage />
    <Scale
      :scale="scaleProps.scale"
      :scale-preset-value="scaleProps.scalePresetValue"
      @change-scale="handleEvent"
    />
  </CollapsedBar>
</template>

<script lang="ts" setup>
import { computed, StyleValue } from 'vue';
import CollapsedBar from '@/components/Common/CollapsedBar.vue';
import Page, { ToolbarPageEvent } from './components/Page.vue';
import Scale, { ToolbarScaleEvent } from './components/Scale.vue';
import Theme from './components/Theme.vue';
import FullPage from './components/FullPage.vue';
import Eye from './components/Eye.vue';
import GroupEye from './components/GroupEye.vue';
import { ToolbarLeftSideBtnEvent } from './components/LeftSideBtn.vue';
import { RightSideBarType } from '~/src/components/Right/TabPanel/type';
import {
  SideTabCommonSettings,
  useRightSideTabSettings,
} from '~/src/hooks/UserSettings/useSideTabSettings';
import { useToolBarSettings } from '~/src/hooks/UserSettings/useToolBarSettings';
import { ElementClick, PageType, reportElementClick } from '~/src/api/report';
import { selfNoteInfo } from '~/src/store';

const { activeTab } = useRightSideTabSettings();
const { toolbarSettings, setToolbarSettings } = useToolBarSettings();

const props = defineProps<{
  pdfId: string;
  noteId?: string;
  pageProps: {
    numPages: number;
    currentPage: number;
  };
  scaleProps: {
    scale: number;
    scalePresetValue: string;
  };
  leftBtnProps?: {
    sideTabSettings: SideTabCommonSettings;
  };
  isLoginUser: boolean;
  bottomDistance: number;
}>();

const emit = defineEmits<{
  (
    event: 'toolbar-event',
    payload: ToolbarPageEvent | ToolbarLeftSideBtnEvent | ToolbarScaleEvent
  ): void;
}>();

const handleEvent = (
  payload: ToolbarPageEvent | ToolbarLeftSideBtnEvent | ToolbarScaleEvent
) => {
  emit('toolbar-event', payload);
};

const handleChangeTheme = (theme: string) =>
  reportElementClick({
    page_type: PageType.note,
    type_parameter: selfNoteInfo.value?.pdfId || '',
    element_name: `color_${theme}` as any,
  });

const toolbarStyle = computed(() => {
  const style: Partial<StyleValue> = {};

  if (props.bottomDistance) {
    style.bottom = `${props.bottomDistance}px`;
  }

  return style;
});
</script>

<style lang="postcss">
.toolbar-bt.collapsed-icon,
.toolbar-bt.collapsed-bar {
  bottom: 30px;
  background: var(--site-theme-bg-primary);
  color: var(--site-theme-text-secondary);
}

.toolbar-bt.collapsed-icon {
  left: 0;
}

.toolbar-bt.collapsed-bar {
  width: max-content;
  transform: translateX(-50%);

  & > *:not(.split) {
    flex: 1 0 auto;
    display: flex;
    align-items: center;
    justify-content: center;
    height: theme('spacing.8');
    padding: theme('spacing.2') 10px;
    line-height: 1;
    border-radius: 2px;

    &:hover {
      background-color: var(--site-theme-bg-hover);
    }
  }

  .btn-icon {
    color: inherit;
  }
  .split {
    flex: 1 0 auto;
    width: 1px;
    height: 12px;
    background: var(--site-theme-divider-light);
    margin: 0 10px;
  }
}
</style>
