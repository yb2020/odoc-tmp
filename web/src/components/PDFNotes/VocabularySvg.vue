<template>
  <Align
    v-for="item in activeWords"
    visible
    alignClass="volcabulary-popuop w-fit max-w-[200px] overflow-hidden z-[9999]"
    :to="item.to"
    :target="item.target"
    :alignProps="{
      points: ['bc', 'tc'],
      offset: [0, 0],
      overflow: { adjustX: true, adjustY: true },
    }"
    @align="onAlign"
  >
    <template #align>
      <WordCard
        ref="card"
        disabledPhonetic
        class="border border-solid border-rp-neutral-3 rounded-[6px]"
        :cardStyle="
          item.width
            ? {
              width: `${item.width}px`,
            }
            : undefined
        "
        :data="item.word"
        @deleted="remove(item.word.id!)"
        @submited="mutate(item.word.id!, $event)"
        @editing="onEditing(item.word.id!)"
        @edited="onEdited(item.word.id!)"
        @pointerenter.stop="doKeepAnnotations"
        @pointermove.stop="doKeepAnnotations"
        @pointerleave.stop="onHideAnnotations"
      />
      <aside class="h-2.5" />
    </template>
  </Align>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, triggerRef } from 'vue';
import { onClickOutside } from '@vueuse/core';
import isElVisible from 'element-visible';
import { WordInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/NoteModule';
import {
  PDF_ANNOTATE_TYPE,
  PDF_ANNOTATIONLAYER,
  ANNOTATION_MOUSEOVER,
  ANNOTATION_PAGESVG_RENDERED,
  ToolBarType,
  PDFJSAnnotate,
} from '@idea/pdf-annotate-core';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { IS_MOBILE } from '@/util/env';
import {
  currentGroupId,
  ownNoteOrVisitSharedNote,
  selfNoteInfo,
} from '@/store';
import { SELF_NOTEINFO_GROUPID } from '@/store/base/type';
import { useWordNotes } from '@common/components/Notes/useWordNote';
import Align from '@common/components/Align/index.vue';
import WordCard from '@common/components/Notes/components/WordCard.vue';
import isMobile from 'is-mobile';

const props = defineProps<{
  pdfViewInstance: ViewerController;
  pdfAnnotateInstance: PDFJSAnnotate;
}>();

interface ActiveWord {
  word: WordInfo;
  to: HTMLElement;
  target: SVGGElement;
  width?: number;
}

let timer: ReturnType<typeof setTimeout>;

const onShowAnnotations = (nodes: SVGGElement[]) => {
  if (
    !nodes.length ||
    IS_MOBILE ||
    currentGroupId.value !== SELF_NOTEINFO_GROUPID ||
    !ownNoteOrVisitSharedNote.value ||
    props.pdfViewInstance?.getDocumentViewer().isSelecting?.() ||
    card.value?.[0]?.isEdit
  ) {
    return;
  }

  let arr: ActiveWord[] = [];
  nodes.slice(-1).forEach((svgEl) => {
    if (!isElVisible(svgEl)) {
      return;
    }

    const type = svgEl.getAttribute(PDF_ANNOTATE_TYPE);
    const uuid = svgEl.getAttribute('uuid') as string;
    const word = words.value.find((x) => x.id === uuid);
    const to = svgEl.closest<HTMLDivElement>('.page');

    if (type !== ToolBarType.Vocabulary || !to || !word) {
      return;
    }

    arr.push({
      word,
      to,
      target: svgEl,
    });
  });

  activeWords.value = arr;
  mountTs.value = Date.now();
};

const onHideAnnotations = (nodes: SVGGElement[]) => {
  doHideAnnotations();
};

const doKeepAnnotations = () => {
  clearTimeout(timer);
};

const doHideAnnotations = (delay = 300) => {
  if (delay <= 0) {
    activeWords.value = [];
  }

  timer = setTimeout(() => {
    activeWords.value =
      activeWords.value?.filter((word, i) => {
        return !!card.value?.[i]?.isEdit;
      }) || [];
  }, delay);
};

onMounted(() => {
  const UI = props.pdfAnnotateInstance.UI;
  UI.on(ANNOTATION_MOUSEOVER, onShowAnnotations);
  UI.on('annotation:mouseout', onHideAnnotations);
  UI.on(ANNOTATION_PAGESVG_RENDERED, checkIfActiveNewWord);
});

onUnmounted(() => {
  const UI = props.pdfAnnotateInstance.UI;
  UI.off(ANNOTATION_MOUSEOVER, onShowAnnotations);
  UI.off('annotation:mouseout', onHideAnnotations);
  UI.off(ANNOTATION_PAGESVG_RENDERED, checkIfActiveNewWord);
});

const pdfId = computed(() => selfNoteInfo.value.pdfId);
const noteId = computed(() => selfNoteInfo.value.noteId);
const { data: words, added, mutate, remove } = useWordNotes(pdfId, noteId);

const card = ref<(typeof WordCard)[]>();
const activeWords = ref<ActiveWord[]>();
const mountTs = ref(0);

onMounted(() => {
  mountTs.value = Date.now();
  onClickOutside(
    computed(() => card.value?.[0]?.rootEl),
    () => {
      if (!isMobile({ tablet: true }) || Date.now() - mountTs.value >= 500) {
        doHideAnnotations();
      }
    }
  );
});

const checkIfActiveNewWord = () => {
  if (added.value) {
    const uuid = added.value.id;
    // 显示浮窗并编辑
    const { container } = props.pdfViewInstance.getDocumentViewer() || {};
    const g = container?.querySelector<SVGGElement>(
      `.${PDF_ANNOTATIONLAYER} g[uuid='${uuid}'][source='1']`
    );
    if (g) {
      onShowAnnotations([g]);
    }
  }
};

const onAlign = () => {
  if (added.value) {
    added.value = undefined;
    card.value?.[0]?.handleEdit();
  }
};

const onEditing = (id: string) => {
  const i = activeWords.value?.findIndex((x) => x.word.id === id) ?? 0;
  const item = activeWords.value?.[i];
  const cardEl: null | HTMLElement = card.value?.[i]?.rootEl;

  if (item && cardEl) {
    item.width = cardEl.offsetWidth;
  }
};

const onEdited = (id: string) => {
  const item = activeWords.value?.find((x) => x.word.id === id);
  if (item) {
    delete item?.width;
    triggerRef(activeWords);
    doHideAnnotations(0);
  }
};
</script>

<style lang="less">
.volcabulary-popuop {
  filter: drop-shadow(2px 6px 16px rgba(0, 0, 0, 0.12))
    drop-shadow(4px 10px 28px rgba(0, 0, 0, 0.08));

  .word-card {
    background: white;
  }

  .word-card > .btn-delete {
    display: block;
    top: 6px;
    @apply right-3;
    padding: 1.5px 0;
    @apply border-0;
    @apply text-rp-neutral-6;
    @apply text-sm;

    &:hover {
      @apply text-rp-neutral-8;
    }

    span {
      @apply ml-1;
    }
  }

  .word-tt {
    @apply text-sm;
    @apply overflow-hidden;
    @apply text-ellipsis;
    padding: 6px theme('spacing.7') 6px 12px;
  }

  .word-ct {
    padding: 7px 12px 12px;
  }

  .word-note {
    .idea-markdown-edit-container,
    .idea-markdown-view-container {
      line-height: 22px;
      min-height: 22px;
      max-height: calc(5lh + 2px);
      @apply !mb-0;

      &:hover {
        @apply bg-rp-neutral-1;
      }
    }
  }
}
</style>
