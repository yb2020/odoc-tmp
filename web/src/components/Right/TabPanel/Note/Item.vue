<template>
  <div
    class="note-container"
    :class="[isOwner ? 'note-owner' : 'note-visitor']"
    @dragover="dragOver"
    @dragenter="dragEnterNote"
    @dragleave="dragLeave"
  >
    <div class="note-self group flex">
      <div
        class="note-bar basis-1"
        :style="{
          background:
            annotationDragging === note
              ? 'var(--site-theme-brand)'
              : isFocus
                ? note.color
                : 'var(--site-theme-pdf-panel-blockquote)',
        }"
      />
      <!-- data-*属性被用于连接线 -->
      <ItemCard
        ref="card"
        class="flex-1"
        :data-pdf-annotate-id="note.uuid"
        refVisible
        :note="note"
        :style="dragStyle"
        :draggable="isOwner && !annotationDragFetching && !isInputing"
        @click="clickNote"
        @edited="handleBlur"
        @dblclick.prevent.stop
        @dragstart="dragStart"
        @dragend="dragEnd"
        @dragover="dragOver"
      >
        <template #inner>
          <div class="cursor-grab" />
        </template>
        <template #prefix>
          <div
            v-if="droppable"
            class="absolute z-[1] inset-0 flex flex-col justify-between items-stretch"
          >
            <div
              class="basis-1/2"
              :class="{
                'drag-over':
                  annotationDragHovering === note.uuid &&
                  annotationDragPosition === 0,
              }"
              @dragenter="dragEnterPosition(0)"
              @dragleave.stop.prevent
              @drop.stop.prevent="drop(0)"
            />
            <div
              class="basis-1/2"
              :class="{
                'drag-over':
                  annotationDragHovering === note.uuid &&
                  annotationDragPosition === 1,
              }"
              @dragenter="dragEnterPosition(1)"
              @dragleave.stop.prevent
              @drop.stop.prevent="drop(1)"
            />
          </div>
        </template>
        <template #suffix>
          <div
            v-if="isOwner && !annotationDragging"
            class="note-btn-del w-5 h-5 absolute z-[2] right-0 top-0 flex items-center justify-center opacity-0 group-hover:opacity-100 translate-x-1/2 -translate-y-1/2"
            @click.stop="handleDelete"
          >
            <div
              class="rounded-full flex items-center justify-center"
              :style="{ backgroundColor: 'var(--site-theme-pdf-panel-secondary)' }"
            >
              <close-outlined
                class="text-sm scale-50"
                :style="{ color: 'var(--site-theme-text-secondary)' }"
              />
            </div>
          </div>
        </template>
      </ItemCard>
    </div>

    <div
      :title="$t('message.doubleClickToCreateNote')"
      class="note-next"
      :style="{
        height: bottomHeight,
      }"
      @dblclick="insertNoReferenceAnnotation(note.pageNumber, index + 1)"
      @dragover="dragOver"
      @dragenter="dragEnterPosition(1)"
      @dragleave.stop.prevent
      @drop.stop.prevent="drop(1)"
    >
      <div
        :style="{
          background: isAppendToNext ? 'var(--site-theme-brand)' : 'transparent',
        }"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { Modal } from 'ant-design-vue';
import { CloseOutlined } from '@ant-design/icons-vue';
import { computed, watch, nextTick, StyleValue, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import '@idea/aiknowledge-markdown/dist/style.css';
import { delay } from '@idea/aiknowledge-special-util/delay';

import uuid from '~/src/util/uuid';
import { clearArrow, connect } from '~/src/dom/arrow';
import { scrollMark } from '~/src/util/scroll';
import { AnnotationAll } from '~/src/stores/annotationStore/BaseAnnotationController';
import { currentNoteInfo } from '~/src/store';
import { usePdfStore } from '~/src/stores/pdfStore';
import {
  useAnnotationStore,
  personAnnotationController,
} from '~/src/stores/annotationStore';

import {
  insertNoReferenceAnnotation,
  annotationDragging,
  annotationDragHovering,
  annotationDragPosition,
  annotationDragStyled,
  annotationDragFetching,
  checkNoReferenceAnnotation,
} from './annotation-state';
import ItemCard from './ItemCard.vue';

const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(currentNoteInfo.value?.pdfId);
});

const props = defineProps<{
  note: AnnotationAll;
  index: number;
  noteList: AnnotationAll[];
  isOwner: boolean;
}>();

const connectLeftRight = async () => {
  if (isNoRef.value) {
    return;
  }

  await delay(300);
  connect(props.note.uuid!);
};

const clickReference = () => {
  scrollMark(props.note.pageNumber, props.note.uuid || '', pdfViewerRef.value!);
  connectLeftRight();
};

const clickNote = () => {
  clickReference();
};

const card = ref();

const isInputing = computed(
  () => annotationStore.inputingAnnotationId === props.note.uuid
);
const isFocus = computed(
  () => annotationStore.currentAnnotationId === props.note.uuid
);

const isNoRef = computed(() => checkNoReferenceAnnotation(props.note));

const isEditingNewNoRefAnnotation = computed(
  () =>
    isNoRef.value &&
    !props.note.idea &&
    annotationStore.currentAnnotationId === props.note.uuid
);

watch(isEditingNewNoRefAnnotation, () => {
  if (isEditingNewNoRefAnnotation.value) {
    card.value?.editNote();
  }
});

const { t } = useI18n();

const handleBlur = async () => {
  if (
    isNoRef.value &&
    annotationStore.currentAnnotationId === props.note.uuid
  ) {
    if (!props.note.idea) {
      await deleteNote();
    }

    annotationStore.currentAnnotationId = '';
  }
};

const handleDelete = () => {
  if (isEditingNewNoRefAnnotation.value) {
    triggerBlur();
    return;
  }

  Modal.confirm({
    title: t('message.confirmToDeleteNoteTip'),
    onOk: deleteNote,
    okButtonProps: {
      danger: true,
    },
    cancelButtonProps: { type: 'primary' },
    // cancelText: '取消',
    okText: t('viewer.delete'),
  });
};

const deleteNote = async () => {
  await annotationStore.controller.deleteAnnotation(
    props.note.uuid,
    props.note.pageNumber
  );
  await annotationStore.controller.loadAnnotationMap();
  annotationStore.currentAnnotationId = '';

  clearArrow();
  annotationStore.delHoverNote(props.note.uuid!);
};

const dragStyle = computed(() => {
  const style: StyleValue = {};

  if (props.isOwner) {
    style.cursor = 'pointer';
    if (annotationDragging.value === props.note) {
      if (annotationDragStyled.value) {
        style.zIndex = -1;
      } else {
        style.border = '1px solid var(--site-theme-brand)';
      }
    }
  }

  return style;
});

const dragStart = () => {
  annotationDragging.value = props.note;
  setTimeout(async () => {
    await nextTick();
    annotationDragStyled.value = true;
  }, 10);
};

const clearDragging = () => {
  annotationDragging.value = null;
  annotationDragStyled.value = false;
  clearDragHovering();
};

const clearDragHovering = () => {
  annotationDragHovering.value = '';
  annotationDragPosition.value = NaN;
};

const dragEnd = async () => {
  await nextTick();
  clearDragging();
};

const droppable = computed(() => {
  return (
    annotationDragging.value &&
    annotationDragging.value !== props.note &&
    (annotationDragging.value.pageNumber === props.note.pageNumber ||
      annotationDragging.value.rectangles.length === 0)
  );
});

const dragOver = (event: Event) => {
  if (droppable.value) {
    event.preventDefault();
  }
};

const dragEnterNote = () => {
  if (!annotationDragging.value) {
    return;
  }

  annotationDragHovering.value = props.note.uuid;
};

const dragEnterPosition = (position: number) => {
  if (!annotationDragging.value) {
    return;
  }

  annotationDragPosition.value = position;
};

const dragLeave = () => {
  if (!annotationDragging.value) {
    return;
  }

  if (annotationDragHovering.value !== props.note.uuid) {
    return;
  }

  clearDragHovering();
};

const drop = async (position: number) => {
  const getId = (item: AnnotationAll) => item.uuid;
  const fromPageNumber = annotationDragging.value!.pageNumber;
  const fromNoteListOld = annotationStore.pageMap[fromPageNumber].map(getId);
  const fromNoteListNew: AnnotationAll['uuid'][] = [...fromNoteListOld];
  const fromIndex = fromNoteListOld.findIndex(
    (id) => id === annotationDragging.value!.uuid
  );
  const toPageNumber = props.note.pageNumber;
  const toNoteListOld = props.noteList.map(getId);
  const toNoteListNew =
    fromPageNumber === toPageNumber ? fromNoteListNew : [...toNoteListOld];
  const movePageNoteId =
    fromPageNumber === toPageNumber ? '' : annotationDragging.value!.uuid;
  const toIndex = props.index + position;
  const deleteSymbol = uuid();
  fromNoteListNew.splice(fromIndex, 1, deleteSymbol);
  toNoteListNew.splice(toIndex, 0, annotationDragging.value!.uuid);
  const deleteIndex = fromNoteListNew.findIndex((id) => id === deleteSymbol);
  fromNoteListNew.splice(deleteIndex, 1);
  annotationDragFetching.value = true;

  await personAnnotationController.sortAnnotation(
    movePageNoteId,
    toPageNumber,
    {
      [fromPageNumber]: fromNoteListNew,
      [toPageNumber]: toNoteListNew,
    }
  );
  annotationDragFetching.value = false;
  clearDragging();
};

const bottomHeight = computed(() => {
  const itemHeight = '12px';
  const lastItemHeight = '18px';
  const lastItemOfLastPageHeight = '160px';
  const isLastPage =
    annotationStore.headTailPageNumber.tail === props.note.pageNumber;
  const isLastItem = props.index === props.noteList.length - 1;
  if (isLastPage) {
    if (isLastItem) {
      return lastItemOfLastPageHeight;
    } else {
      return itemHeight;
    }
  } else if (isLastItem) {
    return lastItemHeight;
  } else {
    return itemHeight;
  }
});

const isAppendToNext = computed(() => {
  return (
    (annotationDragHovering.value === props.note.uuid &&
      annotationDragPosition.value === 1) ||
    (props.noteList[props.index + 1] &&
      annotationDragHovering.value === props.noteList[props.index + 1].uuid &&
      annotationDragPosition.value === 0)
  );
});

const triggerBlur = () => {
  document.body.click();
};
</script>

<style lang="less" scoped>
.note-container {
  .note-next {
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    color: var(--site-theme-text-primary);

    div {
      flex: 0 0 96%;
      height: 3px;
      pointer-events: none;
    }
  }

  &.note-visitor {
    .note-next {
      pointer-events: none;
    }
  }

  &.note-owner {
    .cursor-grab {
      cursor: grab;
    }
  }

  .note-btn-del {
    transform: translate(50%, -50%);

    .anticon {
      transform: scale(0.5);
    }
  }

  .cursor-grab {
    height: 16px;
  }

  :deep(.note) {
    background: var(--site-theme-pdf-panel-secondary);
    border-radius: 2px;
    display: flex;
    position: relative;
  }

  :deep(.note-inner) {
    width: 100%;
    padding: 0 12px 6px 17px;

    .ref-content {
      background: var(--site-theme-pdf-panel-ref-content);
      mix-blend-mode: normal;

      font-family: 'Lato';
      font-style: italic;
      font-weight: 400;
      font-size: 13px;
      line-height: 20px;

      color: var(--site-theme-text-secondary);
      word-break: break-word;
      padding: 8px;
      border-radius: 2px;
      margin-bottom: 8px;
    }

    .ref-picture {
      display: flex;
      justify-content: center;
      margin-bottom: 8px;

      img {
        opacity: 0.4;
        // max-height: 360px;
        max-width: 100%;
      }
    }
  }
}
</style>
<style lang="less">
.note_edit_container .w-e-text-container {
  padding-top: 4px !important;
  padding-bottom: 5px !important;
}

.note-container {
  .idea-markdown-edit-container {
    margin-top: 8px !important;
  }
  .idea-markdown-view-container {
    padding-top: 8px;
    color: var(--site-theme-pdf-panel-text);
    blockquote {
      border-left-color: var(--site-theme-pdf-panel-blockquote);
    }
  }
}

.idea-markdown-view-container {
  font-family: 'Lato';
  font-style: normal;
  font-weight: 400;
  font-size: 14px;
  line-height: 22px;

  min-height: 8px;

  cursor: text;
  word-break: break-word;
  display: block;
  padding-bottom: 6px;
}

.idea-markdown-edit-container {
  margin-bottom: 6px !important;
  transition: none !important;
}

.idea-tag-entry-container {
  padding-top: 2px;
  outline: 0;
}

.idea-tag-input {
  border-color: var(--site-theme-brand-light);
  outline: 0;
}

.idea-markdown-edit-container,
.idea-tag-input {
  color: var(--site-theme-text-primary) !important;
  background-color: var(--site-theme-bg-light) !important;
}

.idea-tag-select-container {
  background-color: var(--site-theme-pdf-panel-secondary);
  color: var(--site-theme-text-secondary);
}

.idea-tag-select-cursor {
  background-color: var(--site-theme-bg-hover);
  color: var(--site-theme-text-primary);
}

.idea-tag-add {
  background-color: var(--site-theme-bg-hover);
}

.idea-tag-item-container {
  background-color: var(--site-theme-bg-hover);
  outline: 0;
}
</style>
