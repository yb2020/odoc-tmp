<template>
  <div
    class="switcher-container"
    @click="handleLeftSideVisible"
  >
    <a-tooltip>
      <template #title>
        {{ sideTabSettings.shown ? $t('viewer.collapseLeftbar') : $t('viewer.expandLeftbar')
        }}<component :is="getShortcutTxt(shortcut)" />
      </template>
      <i
        class="aiknowledge-icon icon-thumbnail-switcher"
        aria-hidden="true"
      />
    </a-tooltip>
  </div>
</template>

<script lang="ts" setup>
export type ToolbarLeftSideBtnEvent = {
  type: 'toolbar:leftside';
};

import { computed } from 'vue';
import { SideTabCommonSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { useStore } from '@/store';
import { PAGE_ROUTE_NAME } from '../../../routes/type';
import { getPlatformKey, getShortcutTxt } from '../../../store/shortcuts';

const props = defineProps<{
  sideTabSettings: SideTabCommonSettings;
  // setSideTabSetting: (values: Partial<SideTabCommonSettings>) => void;
}>();

const store = useStore();

const emit = defineEmits<{
  (event: 'toggleLeftSide', payload: ToolbarLeftSideBtnEvent): void;
}>();

const handleLeftSideVisible = () => {
  emit('toggleLeftSide', { type: 'toolbar:leftside' });
  // props.setSideTabSetting({
  //   shown: !props.sideTabSettings.shown
  // })
};

const shortcutsConfig = computed(() => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {});
const platformKey = getPlatformKey();
const shortcut = computed(() => shortcutsConfig.value.shortcuts.TOGGLE_CATALOG.value[platformKey]);
</script>

<style lang="less" scoped>
.switcher-container {
  width: 36px;
  height: 36px;
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  &:hover {
    background: #f5f5f5;
  }

  &.active {
    background: #f5f5f5;
  }

  .aiknowledge-icon {
    margin-top: 2px;
    font-size: 16px;
    color: rgba(0, 0, 0, 64%);
  }
}
</style>
