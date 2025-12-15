<template>
  <Drawer
    v-model:visible="sideTabSettings.shown"
    dir="right"
    class="right-drawer"
    :min-width="280"
    :initial-width="sideTabSettings.width"
    :shortcut="shortcut"
    @visible-change="onVisibleChange"
    @width-change="onWidthChange"
  >
    <TabPanelVue
      v-if="ownNoteOrVisitSharedNote"
      :paper-id="paperId"
      :pdf-id="pdfId"
      :is-private-paper="isPrivatePaper"
      :clip-selecting="clipSelecting"
      :clip-action="clipAction"
      @visible-change="handleVisibleChange"
    />
    <ForbiddenVue v-else />
  </Drawer>
</template>
<script setup lang="ts">
import TabPanelVue from './TabPanel/index.vue';
import Drawer from '../Common/Drawer.vue';
import ForbiddenVue from './Forbidden.vue';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { BaseActionTypes } from '~/src/store/base';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { computed, watch } from 'vue';
import { useDocStore } from '~/src/stores/docStore';
import { selfNoteInfo } from '~/src/store';
import { PAGE_ROUTE_NAME } from '../../routes/type';
import useShortcuts from '../../hooks/useShortcuts';
import { getPlatformKey } from '../../store/shortcuts';
import { clearArrow } from '../../dom/arrow';

import { store, ownNoteOrVisitSharedNote } from '~/src/store';
import { GetPdfStatusInfoResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import { ElementClick, reportClick } from '~/src/api/report';
import { useI18n } from 'vue-i18n';
import { useClip } from '~/src/hooks/useHeaderScreenShot';

const props = defineProps<{
  isOwner: boolean;
  paperId: string;
  pdfId: string;
  isPrivatePaper: boolean;
  pdfStatusInfo: GetPdfStatusInfoResponse;
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const { t } = useI18n();

const handleVisibleChange = (visible?: boolean) => {
  const shown = visible !== undefined ? visible : !sideTabSettings.value.shown;
  setSideTabSetting({
    shown,
  });
};

const { sideTabSettings, setSideTabSetting } = useRightSideTabSettings();

const onVisibleChange = (visible: boolean) => {
  setSideTabSetting({ shown: visible });
  reportClick(ElementClick.note_bar, visible ? 'on' : 'off');
};

const onWidthChange = (width: number) => {
  setSideTabSetting({ width });
};

if (!props.isOwner) {
  store.dispatch(`base/${BaseActionTypes.SWITCH_TO_GROUP}`, {
    groupId: SELF_NOTEINFO_GROUPID,
    t,
  });
}

const shortcutsConfig = computed(
  () => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {}
);
const platformKey = getPlatformKey();
const shortcut = computed(
  () => shortcutsConfig.value.shortcuts.TOGGLE_NOTE.value[platformKey]
);
const opts = computed(() => ({ scope: shortcutsConfig.value.scope || 'all' }));
const handler = () => {
  clearArrow();
  handleVisibleChange();
};
useShortcuts(shortcut, handler, opts);

const docStore = useDocStore();
watch(
  () => selfNoteInfo.value?.noteId,
  (noteId) => {
    if (noteId) {
      docStore.fetchDocInfo(noteId);
    }
  },
  { immediate: true }
);
</script>
<style lang="less" scoped>
.right {
  height: 100%;
  position: relative;
  transition: width 1.3s cubic-bezier(0.7, 0.3, 0.1, 1);
  overflow: hidden;

  
}

.right-drawer {
  background-color: var(--site-theme-pdf-panel);
  z-index: 10;
}


</style>
