<template>
  <a-dropdown
    v-if="rights?.canUse"
    v-model:visible="visible"
    :trigger="['hover']"
  >
    <template #overlay>
      <AIHighlighterPopover
        :status="status"
        :progress="data?.progress ?? 0"
        :finished="finished"
        :times="rights.remainingTimes"
        :traceid="traceid"
        :highlights="highlights"
        @start="onStart"
        @toggle="onToggle"
      />
    </template>
    <button
      class="relative h-8 !p-2 rounded-sm bg-transparent border border-solid cursor-pointer ai-highlighter-btn"
      :class="[
        finished ? 'finished' : 'unfinished',
        {
          'active': !closed,
        },
      ]"
      @click="onToggleAll"
    >
      {{ $t('aiHighlighter.name') }}
      <span
        v-if="status === Status.LOADING"
        class="absolute h-0.5 bottom-0 left-0 progress-bar"
        :style="{
          width: `${data?.progress ?? 1}%`,
        }"
      />
    </button>
  </a-dropdown>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, triggerRef, watch } from 'vue';
import { ViewerEvent } from '@idea/pdf-annotate-viewer';
import type PDFPageView from '@idea/pdf-annotate-viewer/typing/PDFPageView';
import {
  renderAnnotationSvg,
  scaleDownRaw,
  ToolBarType,
} from '@idea/pdf-annotate-core';
import { Annotation } from '@idea/pdf-annotate-core/render/appendChild';
import { RectOptions } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { Status } from '@/api/aiHighlighter';
import { ElementClick, reportClick } from '@/api/report';
import { currentNoteInfo } from '@/store';
import { getPDFCurPages } from '@/dom/pdf';
import { usePdfStore } from '@/stores/pdfStore';
import { useAnnotationStore } from '@/stores/annotationStore';
import { useAIHighligter, useAIHighlights } from '@/hooks/useAIHighlighter';
import AIHighlighterPopover from './AIHighlighterPopover.vue';

const pdfId = computed(() => currentNoteInfo.value?.pdfId);
const noteId = computed(() => currentNoteInfo.value?.noteId);
const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();
const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(pdfId.value);
});
const pdfAnnotaterRef = computed(() => {
  return pdfStore.getAnnotater(noteId.value);
});
const { data: rights, refresh } = useAIHighligter();
const { data, start } = useAIHighlights(pdfId, noteId);
const onStart = () => {
  start().catch(() => {
    rights.value = {
      ...rights.value!,
      remainingTimes: rights.value!.remainingTimes + 1,
    };
    refresh();
  });

  reportClick(ElementClick.scim_generate);
};
const onToggle = (type: string, checked: boolean) => {
  if (data.value && finished.value) {
    const { map } = data.value;
    map[type].checked = checked;

    data.value = {
      ...data.value,
      map: { ...map },
    };

    // @ts-ignore
    reportClick(type, checked ? 'on' : 'off');
  }
};
const onToggleAll = () => {
  if (data.value && finished.value) {
    const { map } = data.value;
    const all = Object.values(map);
    const checked = !all.every((x) => x.checked);
    all.forEach((x) => {
      x.checked = checked;
    });

    data.value = {
      ...data.value,
      map: { ...map },
    };

    reportClick(ElementClick.scim_keypoint, checked ? 'on' : 'off');
  }
};

const visible = ref(false);

const status = computed(() => {
  const { progress = 0 } = data.value || {};
  if (progress === 0) {
    return Status.INITIAL;
  }

  if (progress === 100) {
    return Status.SUCCESS;
  }

  return Status.LOADING;
});
const traceid = computed(() => data.value?.traceid);
const highlights = computed(() => data.value?.map || {});

const finished = computed(() => status.value === Status.SUCCESS);
const closed = computed(() =>
  Object.values(highlights.value).every((x) => !x.checked)
);

const renderHighlights = async (page: number) => {
  const viewer = pdfViewerRef.value?.getDocumentViewer()?.getPdfViewer();
  const matcher = pdfViewerRef.value?.getWordMatchController();
  const pageNumber = page ?? (viewer?.currentPageNumber || 1);
  const source: undefined | PDFPageView = viewer?.getPageView(pageNumber - 1);
  const viewport = source?.viewport;
  const instance = pdfAnnotaterRef.value;
  if (!instance || !source || !viewport) {
    // 延迟执行
    requestIdleCallback(() => renderHighlights(page), {
      timeout: 100,
    });
    return;
  }

  const annotations: Annotation[] = [];
  await Promise.all(
    Object.values(highlights.value)
      .filter(
        (x) => x.checked
        // 后台暂时不支持page
        // && x.page === pageNumber
      )
      .map(async (x) => {
        const wordMatches = await matcher?.getWordMatchRects(
          x.items.map((item) => item.content),
          pageNumber - 1
        );

        wordMatches?.forEach((wordMatch, i) => {
          const matches = wordMatch.matches?.[pageNumber - 1];
          const highlight = x.items.find((y) => y.content === wordMatch.word);

          const rects = matches?.rects.map(
            ({ x, y, width, height, ...rest }) => {
              return {
                ...rest,
                ...scaleDownRaw(viewport.scale, {
                  x,
                  y,
                  width,
                  height,
                }),
                pageNumber,
              };
            }
          ) as unknown as RectOptions[];

          if (rects?.length) {
            const item: Annotation = {
              uuid: `scim-${x.type}-${highlight?.taskId}`,
              type: ToolBarType.AIHighlight,
              rectangles: rects,
              rectStr: wordMatch.word,
              idea: wordMatch.word,
              color: x.color,
              pageNumber,
              tags: [],
            };
            annotations.push(item);
          }
        });
      })
  );

  renderAnnotationSvg({
    documentId: noteId.value,
    pageNumber,
    source,
    viewport,
    instance,
    annotationsAI: annotations,
  });
};

const onRenderPageHighlights = (params?: { pageNumber: number }) => {
  const pages = params?.pageNumber
    ? [params.pageNumber]
    : getPDFCurPages(pdfViewerRef.value);

  pages.forEach((page) => {
    requestIdleCallback(() => renderHighlights(page));
  });
};

let isEventInited = false;
const initEvents = () => {
  if (!isEventInited && pdfViewerRef.value) {
    isEventInited = true;
    // 必须监听TEXT_LAYER_RENDERED，否则match由于textBounds未更新坐标不准
    pdfViewerRef.value.addEventListener(
      ViewerEvent.TEXT_LAYER_RENDERED,
      onRenderPageHighlights
    );
  }
};

watch(() => annotationStore.instantiated, initEvents);
watch(highlights, () => onRenderPageHighlights());

onMounted(() => {
  initEvents();
});

onUnmounted(() => {
  pdfViewerRef.value?.removeEventListener(
    ViewerEvent.TEXT_LAYER_RENDERED,
    onRenderPageHighlights
  );
});
</script>

<style scoped>
.ai-highlighter-btn {
  color: var(--site-theme-text-color, #000000);
  border-color: var(--site-theme-divider, #d9d9d9);
}

.ai-highlighter-btn:hover {
  background-color: var(--site-theme-background-hover, rgba(0, 0, 0, 0.05)) !important;
}

.ai-highlighter-btn.finished {
  border-color: var(--site-theme-primary-color, #52c41a) !important;
}

.ai-highlighter-btn.unfinished {
  border-color: var(--site-theme-divider, #d9d9d9) !important;
}

.ai-highlighter-btn.active {
  background-color: var(--site-theme-background-hover, rgba(0, 0, 0, 0.05)) !important;
}

.progress-bar {
  background-color: var(--site-theme-primary-color, #1890ff);
}
</style>
