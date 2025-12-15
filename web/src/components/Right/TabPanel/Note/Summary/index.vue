<template>
  <section
    class="summary relative px-4 py-3 h-full text-rp-neutral-10"
    :class="{
      'summary-client': detachable,
    }"
  >
    <Locked v-if="detached" />
    <SummaryEditor
      ref="editor"
      :note-id="selfNoteInfo.noteId"
    />
    <div
      v-if="detachable"
      class="absolute top-3 right-4 p-2 my-[1px] h-9 bg-white flex items-center justify-end"
    >
      <ArrowsAltOutlined
        class="text-base px-1"
        @click="onDetach"
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { ArrowsAltOutlined } from '@ant-design/icons-vue';
import { gteElectronVersion, IS_ELECTRON_MODE } from '@/util/env';
import { ROOT_PATH } from '@/routes/index';
import { selfNoteInfo } from '@/store';
import { useSummaryDetachedStatus } from '@common/components/Notes/useSummary';

import SummaryEditor from '@common/components/Notes/components/SummaryEditor.vue';
import Locked from '../common/Locked.vue';
import {
  ElementClick,
  getPdfIdFromUrl,
  PageType,
  reportElementClick,
} from '~/src/api/report';

const editor = ref<typeof SummaryEditor>();
const detachable = IS_ELECTRON_MODE && gteElectronVersion('1.24.0');

const { detached } = useSummaryDetachedStatus(selfNoteInfo.value.noteId);

watch(detached, (curr, prev) => {
  if (!curr && prev) {
    editor.value?.refresh();
  }
});

const onDetach = () => {
  if (IS_ELECTRON_MODE) {
    detached.value = true;

    window.open(
      `${ROOT_PATH}/summary.html?noteId=${selfNoteInfo.value.noteId}`
    );

    reportElementClick({
      page_type: PageType.note,
      type_parameter: getPdfIdFromUrl(),
      element_name: ElementClick.note_independent_window,
    });
  }
};
</script>

<style lang="postcss" scoped>
.summary-client {
  :deep(.milkdown-menu) {
    padding-right: 40px;
  }
}

.summary {
  :deep(.milkdown),
  :deep(.milkdown-wrapper),
  :deep(.milkdown-menu-wrapper) {
    height: 100%;
  }

  :deep(.milkdown-menu-wrapper) {
    display: flex;
    flex-direction: column;
  }

  :deep(.milkdown-inner) {
    flex: 1;
  }

  :deep(.milkdown-menu) {
    justify-content: flex-start;
    align-items: center;
    gap: 8px;
    /* overflow-x: hidden; */

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
    }

    .material-icons {
      font-size: 16px;
    }
  }

  :deep(.milkdown) {
    padding: 6px theme('spacing.4');

    .editor {
      padding: 0;
    }
  }
}
</style>
