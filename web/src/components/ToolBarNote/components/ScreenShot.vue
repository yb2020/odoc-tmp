<template>
  <div
    :class="['screenshot-container', clipSelecting ? 'active' : '']"
    @click="handleScreenShot"
  >
    <a-tooltip
      placement="bottom"
      :overlay-style="{
        whiteSpace: 'pre-line',
      }"
    >
      <template #title>
        {{ $t('viewer.screenshot') }}<component :is="getShortcutTxt(shortcut)" />
      </template>
      <i
        class="aiknowledge-icon icon-crop"
        aria-hidden="true"
      />
    </a-tooltip>
  </div>
</template>

<script lang="ts" setup>
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { useStore } from '@/store';
import { PAGE_ROUTE_NAME } from '../../../routes/type';
import useShortcuts from '../../../hooks/useShortcuts';

import { getPlatformKey, getShortcutTxt } from '../../../store/shortcuts';

const store = useStore();

const props = defineProps({
  clipSelecting: {
    type: Boolean,
    required: true,
  },
  cancelClip: {
    type: Function,
    required: true,
  },
  handleScreenShot: {
    type: Function as any,
    required: true,
  },
});

const callback = (e: MouseEvent) => {
  if (props.clipSelecting) {
    e.preventDefault();
    props.cancelClip();
  }
};

onMounted(() => {
  document.addEventListener('contextmenu', callback);
});

onUnmounted(() => {
  document.removeEventListener('contextmenu', callback);
});

const shortcutsConfig = computed(
  () => store.state.shortcuts[PAGE_ROUTE_NAME.NOTE] || {}
);
const platformKey = getPlatformKey();
const shortcut = computed(
  () => shortcutsConfig.value.shortcuts.NOTE_SCREENSHOT.value[platformKey]
);
const opts = computed(() => ({ scope: shortcutsConfig.value.scope || 'all' }));
useShortcuts(shortcut, props.handleScreenShot, opts);
</script>

<style lang="less" scoped>
.screenshot-container {
  width: 40px;
  height: 40px;
  border-radius: 2px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;

  &:hover,
  &:active,
  &.active {
    background-color: rgba(0, 0, 0, 0.2);
  }
  .aiknowledge-icon {
    font-size: 20px;
  }
}
</style>
