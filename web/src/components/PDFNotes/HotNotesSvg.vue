<template>
  <teleport to="body">
    <div
      v-show="hotTipVisible"
      ref="hotNoteTipsRef"
      class="hot-note-tips"
    >
      {{ $t('message.peopleAnnotateTip', { num: hotNoteScore }) }}
    </div>
  </teleport>
</template>

<script setup lang="ts">
import {
  PDFJSAnnotate,
  PDF_ANNOTATE_TYPE,
  ANNOTATION_MOUSEOVER,
  ToolBarType,
} from '@idea/pdf-annotate-core';
import { computed, onMounted, onUnmounted, ref } from 'vue';
import { currentGroupId, ownNoteOrVisitSharedNote } from '~/src/store';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { IS_MOBILE } from '~/src/util/env';

const props = defineProps<{
  pdfAnnotateInstance: PDFJSAnnotate;
}>();

const annotationStore = useAnnotationStore();
const hotTipVisible = computed({
  get: () => annotationStore.hotTipVisible,
  set: (v: boolean) => (annotationStore.hotTipVisible = v),
});
const hotNoteTipsRef = ref();
const hotNoteScore = ref();

const showHotNoteTips = (annotateScore: string, event: MouseEvent) => {
  const { pageX, pageY } = event;

  hotNoteTipsRef.value.style.left = pageX + 'px';
  hotNoteTipsRef.value.style.top = pageY + 10 + 'px';
  hotNoteScore.value = annotateScore;
  hotTipVisible.value = true;
};

const annotationMouseOverCb = (nodes: SVGGElement[], event: MouseEvent) => {
  nodes.forEach((el) => {
    annotationMouseOver(el, event);
  });
};

const annotationMouseOver = async (target: SVGGElement, event: MouseEvent) => {
  if (
    IS_MOBILE ||
    target.style.display === 'none' ||
    currentGroupId.value !== SELF_NOTEINFO_GROUPID ||
    !ownNoteOrVisitSharedNote.value
  ) {
    return;
  }

  const annotateType = parseInt(
    target.getAttribute(PDF_ANNOTATE_TYPE) as string
  );
  const annotateScore = target.getAttribute(
    'data-pdf-annotate-score'
  ) as string;

  if (annotateType === ToolBarType.hot) {
    showHotNoteTips(annotateScore, event);
  }
};

const annotationMouseOutCb = () => {
  hotTipVisible.value = false;
};

const UI = props.pdfAnnotateInstance.UI;
onMounted(() => {
  UI.on(ANNOTATION_MOUSEOVER, annotationMouseOverCb);
  UI.on('annotation:mouseout', annotationMouseOutCb);
});
onUnmounted(() => {
  UI.off(ANNOTATION_MOUSEOVER, annotationMouseOverCb);
  UI.off('annotation:mouseout', annotationMouseOutCb);
});
</script>

<style scoped>
.hot-note-tips {
  position: fixed;
  background-color: rgba(29, 34, 41, 0.85);
  padding: 3px 8px;
  color: #fff;
}
</style>
