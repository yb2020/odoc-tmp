<template>
  <main class="main max-w-7xl w-full h-full mx-auto text-rp-neutral-10">
    <Locked v-if="disabled" />
    <SummaryEditor
      :note-id="params.noteId"
      @loaded="onLoad"
    />
  </main>
</template>

<script setup lang="ts">
import { useUrlSearchParams, useEventListener } from '@vueuse/core';
import { useSummaryDetachedStatus } from '@common/components/Notes/useSummary';

import SummaryEditor from '@common/components/Notes/components/SummaryEditor.vue';
import Locked from '@/components/Right/TabPanel/Note/common/Locked.vue';
import { GetSummaryResponse } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import { computed } from 'vue';

const params = useUrlSearchParams<{
  noteId: string;
  active: string;
}>();
const { detached } = useSummaryDetachedStatus(params.noteId);
// 笔记管理处打开也需要锁住
const disabled = computed(() => params.active !== '1' && detached.value);

const toggleDetached = (v = false) => {
  if (params.active === '1') {
    detached.value = v;
  }
};

useEventListener(window, 'load', () => toggleDetached(true));

// https://github.com/electron/electron/issues/1929
// window.addEventListener('beforeunload', clearDetached);
useEventListener(window, 'unload', () => toggleDetached(false));

const onLoad = (data: GetSummaryResponse) => {
  document.title = data.docName;
};
</script>

<style lang="postcss">
#app {
  background-color: theme('colors.rp-neutral-1');
  height: 100%;
}

.main {
  .summary-editor {
    .milkdown-menu {
      border-bottom: 1px solid theme('colors.rp-neutral-3');
      background: theme('colors.rp-neutral-2');
      justify-content: flex-start;
      align-items: center;
      gap: 8px;
      height: 48px;
      padding: 0 12px;

      .menu-selector {
        margin: 0;
      }

      .divider {
        margin: 0;
        height: 16px;
        min-height: 16px;
      }

      .button {
        margin: 0;
        width: 16px;
        height: 16px;
        background-color: transparent;

        &:hover {
          background-color: rgba(129, 161, 193, 0.12);
        }
      }

      .material-icons {
        font-size: 16px;
      }
    }

    .milkdown {
      padding: 0 theme('spacing.6');
    }
  }
}
</style>
