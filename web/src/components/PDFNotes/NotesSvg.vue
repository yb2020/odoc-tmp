<template>
  <NoteSvg
    v-for="item in annotations"
    :key="item.uuid"
    :note="item"
    :pdfViewInstance="pdfViewInstance"
  />
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue';
import {
  ToolBarType,
  PDF_ANNOTATE_TYPE,
  PDFJSAnnotate,
  ANNOTATION_MOUSEOVER,
} from '@idea/pdf-annotate-core';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { IS_MOBILE } from '~/src/util/env';
import { ownNoteOrVisitSharedNote, currentGroupId } from '@/store';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import NoteSvg from './NoteSvg.vue';

const props = defineProps<{
  pdfViewInstance: ViewerController;
  pdfAnnotateInstance: PDFJSAnnotate;
}>();

const annotationStore = useAnnotationStore();
const annotations = computed(() => annotationStore.activeHoverNotes);

const showAnnotations = (nodes: SVGGElement[], event: MouseEvent) => {
  // 浮窗编辑时不触发其他浮窗
  if (!annotationStore.inputingAnnotationId && nodes.length) {
    showAnnotation(nodes[nodes.length - 1], event);
  }
};

const showAnnotation = (target: SVGGElement, event: MouseEvent) => {
  if (
    IS_MOBILE ||
    target.style.display === 'none' ||
    currentGroupId.value !== SELF_NOTEINFO_GROUPID ||
    !ownNoteOrVisitSharedNote.value ||
    props.pdfViewInstance?.getDocumentViewer().isSelecting?.()
  ) {
    return;
  }

  const annotateType = parseInt(
    target.getAttribute(PDF_ANNOTATE_TYPE) as string
  );
  const pageEl = target.closest<HTMLDivElement>('.page');
  const uuid = target.getAttribute('uuid') as string;
  const page = target.getAttribute('page-number');

  if (annotateType === ToolBarType.hot || !pageEl || !page) {
    return;
  }

  const annotation = annotationStore.crossPageMap[page]?.find(
    (anno) => anno.uuid === uuid
  )!;

  if (!annotation || !(annotation.idea || annotation.tags?.length)) {
    return;
  }

  // 最多显示一个浮窗
  annotationStore.showHoverNote(
    {
      page: parseInt(page, 10),
      pageEl: pageEl,
      uuid,
    },
    true
  );
};

const hideAnnotations = (nodes: SVGGElement[], event: MouseEvent) => {
  nodes.forEach((el) => {
    hideAnnotation(el, event);
  });
};

const hideAnnotation = (target: SVGGElement, event: MouseEvent) => {
  const uuid = target.getAttribute('uuid') as string;

  if (annotationStore.inputingAnnotationId !== uuid) {
    annotationStore.delHoverNote(uuid);
  }
};

onMounted(() => {
  const UI = props.pdfAnnotateInstance.UI;
  UI.on(ANNOTATION_MOUSEOVER, showAnnotations);
  UI.on('annotation:mouseout', hideAnnotations);
});

onUnmounted(() => {
  const UI = props.pdfAnnotateInstance.UI;
  UI.off(ANNOTATION_MOUSEOVER, showAnnotations);
  UI.off('annotation:mouseout', hideAnnotations);
});
</script>
