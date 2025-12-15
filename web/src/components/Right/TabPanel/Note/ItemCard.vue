<template>
  <div
    ref="noteRef"
    class="note"
    :class="{
      'note--disabled': !isOwner,
    }"
    @click="clickNode"
  >
    <slot name="prefix" />

    <div class="note-inner">
      <slot name="inner" />

      <div
        v-if="
          refVisible &&
            note.type === ToolBarType.select &&
            annotationStore.activeColorMap.ref
        "
        class="ref-content"
      >
        {{ (note as AnnotationSelect).rectStr }}
      </div>

      <div
        v-if="
          refVisible &&
            note.type === ToolBarType.rect &&
            annotationStore.activeColorMap.ref
        "
        class="ref-picture"
      >
        <img
          :src="(note as AnnotationRect).picUrl"
          alt=""
        >
      </div>

      <IdeaMarkdown
        :raw="note.idea || ''"
        :width="noteWidth"
        :uniq-id="note.uuid"
        :editing="isEditing"
        :blur-delay="0"
        :upload="upload"
        @change="handleChange($event)"
        @blur="handleBlur($event)"
        @submit="handleBlur($event, true)"
        @click-view="editNote()"
        @click.stop
        @mousedown.stop
      />

      <div
        class="note-extra flex justify-end"
        @click.stop
        @mousedown.stop
      >
        <Tags
          class="flex-1"
          :note="note"
          @selecting="toggleSelecting"
          @selectend="toggleSelecting(false)"
        />
        <MarkdownTip v-if="mdTipVisible && isOwner" />
        <Dot
          class="note-dot"
          :item="note"
        />
      </div>
    </div>

    <slot name="suffix" />
  </div>
</template>

<script setup lang="ts">
import { IdeaMarkdown } from '@idea/aiknowledge-markdown';
import { ToolBarType } from '@idea/pdf-annotate-core';
import '@idea/aiknowledge-markdown/dist/style.css';

import { message } from 'ant-design-vue';
import { ref } from 'vue';
import { uploadImage, ImageStorageType } from '~/src/api/upload';

import { isOwner } from '~/src/store';
import {
  AnnotationAll,
  AnnotationSelect,
  AnnotationRect,
} from '~/src/stores/annotationStore/BaseAnnotationController';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import MarkdownTip from '@/components/Common/MarkdownTip.vue';
import Tags from './Tags.vue';
import Dot from './Dot.vue';
// import { ResponseError } from '~/src/api/type';

const props = defineProps<{
  note: AnnotationAll;
  noteWidth?: string;
  locked?: boolean;
  refVisible?: boolean;
  mdTipVisible?: boolean;
}>();

const emit = defineEmits(['click', 'editing', 'edited']);

const annotationStore = useAnnotationStore();
const noteRef = ref<HTMLElement>();
const noteIdea = ref('');
const isEditing = ref(false);
const isSubmiting = ref(false);
const isSelectingTag = ref(false);

const upload = async (src: File | string) => {
  return uploadImage(src, ImageStorageType.markdown);
};

const showHover = () => {
  annotationStore.showHoverNote({
    uuid: props.note.uuid,
    page: props.note.pageNumber,
  });
};

const clickNode = () => {
  emit('click');
};

const editNote = () => {
  if (!isOwner.value) {
    return;
  }

  noteIdea.value = props.note.idea;
  annotationStore.inputingAnnotationId = props.note.uuid;
  showHover();
  isEditing.value = true;

  emit('editing');
};

const editFocus = () => {
  const $md = noteRef.value?.querySelector<HTMLElement>(
    '.idea-markdown-edit-container'
  );
  $md?.focus();
};

const toggleSelecting = (v = true) => {
  isSelectingTag.value = v;
};

const handleChange = (idea: string) => {
  const { uuid, pageNumber } = props.note;
  annotationStore.controller.localPatchAnnotation(uuid, pageNumber, {
    idea,
  });
};

const handleBlur = async (idea: string, isSubmit = false) => {
  if (props.locked) {
    editFocus();
    return;
  }

  if (
    idea.length > 2000 ||
    ((props.note as AnnotationSelect).rectStr || '').length > 2000
  ) {
    message.error('单个笔记字数最多2000字！');
    editFocus();
    return;
  }

  const { uuid } = props.note;
  if (noteIdea.value !== idea && !isSubmiting.value) {
    isSubmiting.value = true;
    annotationStore.controller
      .patchAnnotation(uuid, {
        idea,
      })
      // 已有默认处理
      // .catch((e: ResponseError) => {
      //   message.error(`笔记保存失败(${e.code ?? -1})：${e.message || e.cause}`);
      // })
      .finally(() => {
        isSubmiting.value = false;
      });
    annotationStore.refreshTagList();
  }

  isEditing.value = false;
  noteIdea.value = '';
  if (annotationStore.inputingAnnotationId === props.note.uuid) {
    annotationStore.inputingAnnotationId = '';
  }
  if (isSubmit) {
    annotationStore.delHoverNote(uuid, 0);
    annotationStore.currentAnnotationId = '';
  }
  emit('edited');
};

defineExpose({
  editNote,
});
</script>

<style lang="postcss">
.note--disabled {
  .note-tags {
    .idea-tag-add {
      display: none;
    }
    .idea-tag-item-container {
      padding-right: 8px !important;
    }
    .idea-tag-item-remove {
      display: none;
    }
  }
  .note-dot {
    &:hover {
      background: transparent;
    }
  }

  .note-inner {
    .idea-markdown-view-container {
      &:hover {
        background-color: transparent !important;
      }
    }
  }
}
</style>
