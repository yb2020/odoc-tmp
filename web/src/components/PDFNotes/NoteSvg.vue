<template>
  <teleport
    v-if="annotation"
    :to="toEl"
  >
    <div
      ref="noteEl"
      class="note-popup js-note-popup js-resizable"
      :style="styles"
      @pointerenter.stop="keep"
      @pointerleave.stop="hide()"
      @pointermove.stop
    >
      <NoteCard
        ref="card"
        md-tip-visible
        :note="annotation"
        :note-width="
          typeof noteWidth === 'number' ? `${noteWidth}px` : noteWidth
        "
        :locked="note.locked"
        @editing="onEditing"
        @edited="onEdited"
      />
    </div>
  </teleport>
</template>

<script setup lang="ts">
import NoteCard from '@/components/Right/TabPanel/Note/ItemCard.vue';
import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
// import { onClickOutside } from '@vueuse/core';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import tippy, { Instance } from 'tippy.js';
import enableTippyResizable from '@/dom/enableTippyResizable';
import { getTextRectByAspectRadio } from '@/util/popup';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { NoteSlice } from '~/src/stores/annotationStore/types';
import { onClickOutside } from '@vueuse/core';
import isMobile from 'is-mobile';

const props = defineProps<{
  note: NoteSlice;
  pdfViewInstance: ViewerController;
}>();

const card = ref();
const noteEl = ref<HTMLDivElement>();
const noteWidth = ref<number | string>();
const customWidth = ref<number>();
const tippyInstance = ref<Instance>();

const annotationStore = useAnnotationStore();
const isInputing = computed(
  () => annotationStore.inputingAnnotationId === props.note.uuid
);
const isFocus = computed(
  () => annotationStore.currentAnnotationId === props.note.uuid
);
const isBlurOutside = ref(false);
const mountTs = ref(0);

const checkEdit = () => {
  if (isFocus.value) {
    card.value?.editNote();
  }
};
watch(isFocus, checkEdit);
onMounted(checkEdit);

const annotation = computed(() => {
  const { page, uuid } = props.note;
  const arr = annotationStore.crossPageMap[page];

  return arr?.find((x) => x.uuid === uuid);
});

const toEl = computed(() => {
  const { page, pageEl } = props.note;

  return (
    pageEl ??
    (() => {
      const container = props.pdfViewInstance?.getDocumentViewer().container!;
      return container.querySelector(
        `.page[data-page-number="${page}"]`
      )! as HTMLElement;
    })()
  );
});

const baseEl = computed(() => {
  const { uuid } = props.note;

  const el = toEl.value?.querySelector(`[uuid="${uuid}"]`);
  // const svgEl = el?.closest('svg');

  return el;
});

const MIN_WIDTH_OTHERS = 14 + 20 + 8 * 2 + 12 * 2;
const MIN_WIDTH_NO_LABEL = 94 + MIN_WIDTH_OTHERS;
// 最长标签 + MD提示 + 高亮颜色 + gap * 2 + padding * 2
const MIN_WIDTH_ELLIPSIS_LABEL = 168 + MIN_WIDTH_OTHERS;
const styles = computed(() => {
  if (noteWidth.value || !baseEl.value || !toEl.value) {
    return;
  }

  const { width: baseW } = baseEl.value.getBoundingClientRect();
  const { width: pageW } = toEl.value.getBoundingClientRect();

  let minWidth = Math.max(
    annotation.value?.tags?.length
      ? MIN_WIDTH_ELLIPSIS_LABEL
      : MIN_WIDTH_NO_LABEL,
    baseW
  );
  let width = customWidth.value;
  if (!width && annotation.value?.idea) {
    width = getTextRectByAspectRadio(
      annotation.value.idea,
      4 / 3,
      {
        'font-size': '14px',
        'line-height': '24px',
      },
      minWidth,
      pageW
    );
  }

  return {
    width: width ? `${width}px` : undefined,
    // 用最小宽度以避免标签增加时，长度不够导致标签溢出浮窗
    'min-width': `${minWidth}px`,
    'max-width': `${pageW}px`,
  };
});

onMounted(() => {
  if (!baseEl.value || !toEl.value) {
    return;
  }
  // @ts-ignore
  tippyInstance.value = tippy(baseEl.value, {
    theme: 'ref-paper',
    content: noteEl.value,
    appendTo: toEl.value,
    maxWidth: 'none',
    trigger: 'manual',
    arrow: false,
    hideOnClick: false,
    interactive: true,
    showOnCreate: true,
    // 默认是10px
    // offset: [0, 10],
  });

  if (!noteEl.value) {
    return;
  }
  let isNeedRestoreEdit = false;
  enableTippyResizable(noteEl.value, {
    edges: {
      left: true,
      right: true,
      top: false,
      bottom: false,
    },
    hold: 1,
    listeners: {
      // 需要down事件才能比blur先触发
      start: () => {
        annotationStore.mutHoverNote(props.note.uuid, {
          locked: true,
        });
      },
      move: (event: Interact.ResizeEvent) => {
        const target = event.target as HTMLElement;
        customWidth.value = target.offsetWidth;
        if (isNeedRestoreEdit) {
          noteWidth.value = '100%';
        }
      },
      end: (event: Interact.ResizeEvent) => {
        const target = event.target as HTMLElement;
        customWidth.value = target.offsetWidth;
        keep();
        if (isNeedRestoreEdit) {
          card.value?.editNote();
          isNeedRestoreEdit = false;
        }
        annotationStore.mutHoverNote(props.note.uuid, {
          locked: false,
        });
      },
    },
  }).on('down', () => {
    isNeedRestoreEdit = isInputing.value;
  });
});

onMounted(() => {
  mountTs.value = Date.now();
  onClickOutside(noteEl, () => {
    if (isInputing.value) {
      // 还没触发edited
      isBlurOutside.value = true;
    } else if (
      !isMobile({ tablet: true }) ||
      Date.now() - mountTs.value >= 500
    ) {
      // 避免移动端点击显示后立即触发hide
      hide(0);
    }
  });
});

onUnmounted(() => {
  tippyInstance.value?.destroy();
});

const keep = () => {
  annotationStore.showHoverNote(props.note);
};
const hide = (delay?: number) => {
  if (
    isInputing.value ||
    // windows默认输入法输入时未失焦输入框但会引起mouseleave事件...
    document.activeElement?.closest('.js-note-popup') === noteEl.value
  ) {
    return;
  }
  annotationStore.delHoverNote(props.note.uuid, delay);
  if (typeof delay === 'number' && delay <= 0) {
    annotationStore.currentAnnotationId = '';
  }
};
const onEditing = () => {
  const viewEl = noteEl.value?.querySelector<HTMLElement>(
    '.idea-markdown-view-container'
  );

  if (viewEl) {
    noteWidth.value = viewEl.offsetWidth;
  }
};
const onEdited = () => {
  noteWidth.value = undefined;
  if (isBlurOutside.value) {
    hide(0);
    isBlurOutside.value = false;
  }
};
</script>

<style lang="postcss">
.note-popup {
  background-color: #fff;

  &::before,
  &::after {
    content: '';
    position: absolute;
    width: 100%;
    height: 10px;
    left: 0;
  }
  &::before {
    bottom: 100%;
  }
  &::after {
    top: 100%;
  }

  .note {
    padding: theme('spacing.3');
  }

  .note-tags {
    padding: 0;
    gap: 4px;
  }

  .note-extra {
    justify-content: space-between;
  }

  .note-dot {
    padding: 5px;

    &:hover {
      background: theme('colors.rp-neutral-1');
    }
  }

  .idea-tag-add,
  .idea-tag-entry-edit,
  .idea-tag-item-container {
    visibility: visible;
    margin-bottom: 0;
    color: theme('colors.rp-neutral-6');
    background-color: theme('colors.rp-neutral-2');
  }

  .idea-tag-entry-edit {
    top: 0;
  }

  .idea-tag-input {
    border: 1px solid theme('colors.rp-blue-6');
  }

  .markdown-tip {
    width: 20px;
    height: 20px;

    & > i {
      color: theme('colors.rp-neutral-6');
    }
    &:hover {
      background-color: #00000014;
    }
  }

  .idea-markdown-view-container {
    h1,
    h2,
    h3,
    h4,
    h5,
    h6 {
      color: rgba(0, 0, 0, 0.85);
    }
  }

  .idea-markdown-view-container,
  .idea-markdown-edit-container {
    color: theme('colors.rp-neutral-10');
    border-width: 1px;
    border-style: solid;
    border-color: transparent;
    padding: 3px 0;
    min-height: 32px;
    line-height: 24px;
    margin-bottom: 7px !important;
    /* 保持一致避免失焦重排抖动 */

    &:hover {
      background-color: theme('colors.rp-neutral-1') !important;
      border-color: transparent;
    }

    &:focus {
      background-color: transparent !important;
      border-color: theme('colors.rp-blue-6');
      box-shadow: none;
    }
  }
}
</style>
