<template>
  <Align
    v-for="(item, i) in actives"
    :key="item.data.content"
    visible
    alignClass="aihighlight-popuop w-fit overflow-visible z-[9999]"
    :to="item.to"
    :target="item.target"
    :alignProps="{
      points: ['bc', 'tc'],
      offset: [0, 0],
      overflow: { adjustX: true, adjustY: true },
    }"
  >
    <template #align>
      <HighlightCard
        ref="popups"
        :hidden="item.hidden"
        :to="item.to"
        :target="item.target"
        :data="item.data"
        :pdfViewInstance="pdfViewInstance"
        @pointerenter.stop="doKeepAnnotations"
        @pointermove.stop="doKeepAnnotations"
        @pointerleave.stop="doHideAnnotation(item, i)"
        @hide="doHideAnnotation(item, i)"
        @noted="doHideAnnotations(0)"
      />
    </template>
  </Align>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, triggerRef } from 'vue';
import { onClickOutside } from '@vueuse/core';
import { SCIMItem } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/scim/AiSCIMInfo';
import {
  PDFJSAnnotate,
  PDF_ANNOTATE_TYPE,
  ToolBarType,
} from '@idea/pdf-annotate-core';
import { ViewerController } from '@idea/pdf-annotate-viewer';
import { IS_MOBILE } from '@/util/env';
import {
  ownNoteOrVisitSharedNote,
  selfNoteInfo,
  isSelfNoteInfo,
} from '@/store';
import { useAIHighlights } from '@/hooks/useAIHighlighter';
import Align from '@common/components/Align/index.vue';
import HighlightCard from './AIHighlightSvg.vue';

const props = defineProps<{
  pdfViewInstance: ViewerController;
  pdfAnnotateInstance: PDFJSAnnotate;
}>();

interface Active {
  data: SCIMItem & {
    translation?: string;
  };
  to: HTMLElement;
  target: SVGGElement;
  hidden?: boolean;
}

let timer: ReturnType<typeof setTimeout>;

const onShowAnnotations = (nodes: SVGGElement[]) => {
  if (
    IS_MOBILE ||
    !nodes.length ||
    !isSelfNoteInfo ||
    !ownNoteOrVisitSharedNote.value ||
    props.pdfViewInstance?.getDocumentViewer().isSelecting?.()
  ) {
    return;
  }

  let arr: Active[] = [];
  // 需要取最后渲染的node，确保优先级比摘录笔记低
  nodes.slice(-1).forEach((svgEl) => {
    if (svgEl.style.display === 'none') {
      return;
    }

    const type = svgEl.getAttribute(PDF_ANNOTATE_TYPE);
    const uuid = svgEl.getAttribute('uuid') as string;
    const data = highlights.value?.find(
      (x) => `scim-${x.scimType}-${x.taskId}` === uuid
    );
    const to = svgEl.closest<HTMLDivElement>('.page');

    if (
      type !== ToolBarType.AIHighlight ||
      !to ||
      !data ||
      actives.value?.some((x) => x.data === data)
    ) {
      return;
    }

    arr.push({
      data,
      to,
      target: svgEl,
    });
  });

  if (arr.length) {
    // console.debug('>>> enter', arr);
    actives.value = arr;
  }
};

const onHideAnnotations = (nodes: SVGElement[]) => {
  nodes.slice(-1).forEach((svgEl) => {
    const uuid = svgEl.getAttribute('uuid') as string;
    const idx =
      actives.value?.findIndex(
        (x) => `scim-${x.data.scimType}-${x.data.taskId}` === uuid
      ) ?? -1;

    if (idx >= 0) {
      // console.debug('>>> leave', actives.value![idx]);
      doHideAnnotation(actives.value![idx], idx);
    }
  });
};

const doKeepAnnotations = () => {
  // console.debug('>>> keep svg');
  clearTimeout(timer);
};

const doHideAnnotations = (delay = 300) => {
  const fn = () =>
    (actives.value = actives.value?.filter((x) => x.hidden) ?? []);
  if (delay <= 0) {
    fn();
  } else {
    timer = setTimeout(fn, delay);
  }
};

const doHideAnnotation = (item: Active, idx: number, delay = 300) => {
  if (popups.value?.[idx]?.translating) {
    // console.debug('>>> hide svg');
    item.hidden = true;
    triggerRef(actives);
  } else {
    // console.debug('>>> unmount svg');
    timer = setTimeout(() => {
      actives.value = actives.value?.filter((x) => x.data !== item.data) ?? [];
    }, delay);
  }
};

onMounted(() => {
  const UI = props.pdfAnnotateInstance.UI;
  UI.on('annotation:mouseover', onShowAnnotations);
  UI.on('annotation:mouseout', onHideAnnotations);
});

onUnmounted(() => {
  const UI = props.pdfAnnotateInstance.UI;
  UI.off('annotation:mouseover', onShowAnnotations);
  UI.off('annotation:mouseout', onHideAnnotations);
});

const pdfId = computed(() => selfNoteInfo.value.pdfId);
const noteId = computed(() => selfNoteInfo.value.noteId);
const { data } = useAIHighlights(pdfId, noteId);
const highlights = computed(() => data.value?.items ?? []);

const popups = ref<any[]>();
const actives = ref<Active[]>();

onClickOutside(
  computed(() => popups.value?.[0]?.rootEl),
  () => {
    doHideAnnotations(0);
  }
);
</script>

<!-- 伪元素无法触发mouseenter/over -->
<!-- <style lang="less">
.aihighlight-popuop {
  .aihighlight-card {
    &::before,
    &::after {
      content: '';
      position: absolute;
      left: 0;
      width: 100%;
      height: 30px;
    }
    &::before {
      bottom: 100%;
    }
    &::after {
      top: 100%;
    }
  }
}
</style> -->
