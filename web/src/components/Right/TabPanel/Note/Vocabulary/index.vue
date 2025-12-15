<template>
  <section class="h-full pt-1 pb-3 flex flex-col">
    <div class="flex items-center px-4 text-xs">
      <a-button
        class="!p-0 !text-xs"
        type="text"
        size="small"
        @click="handleEdit"
      >
        <PlusOutlined class="-mr-2" />{{ $t('note.vocabulary.add') }}
      </a-button>
      <a-tooltip>
        <template #title>
          {{ $t('note.vocabulary.addTip') }}
        </template>
        <InfoCircleOutlined class="ml-1" />
      </a-tooltip>
      <p class="flex-1" />
      <a-dropdown>
        <a
          class="mr-5 leading-[1]"
          @click.prevent
        >
          <LoadingOutlined v-if="isConfiging" />
          <DownOutlined v-else />
          <span class="ml-2">{{ $t(modeTxt) }}</span>
        </a>
        <template #overlay>
          <a-menu>
            <a-menu-item v-for="(option, value) in DISPLAY_OPTIONS">
              <a
                href="javascript:;"
                @click="handleModeChange(+value)"
              >{{
                $t(option.i18n)
              }}</a>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
      <span>{{ $t('note.vocabulary.color') }}</span>
      <a-spin :spinning="isConfiging">
        <DotSelect
          editable
          class="w-5 h-5 flex items-center justify-center cursor-pointer dot-select-hover"
          :color="gcolorHex"
          :color-map="vocabularyStyleMap"
          @change="handleColorChange"
        />
      </a-spin>
    </div>
    <div
      v-if="isEditing"
      class="py-3 px-4"
    >
      <div class="p-4" :style="{ backgroundColor: 'var(--site-theme-pdf-panel-secondary)' }">
        <a-input
          ref="input"
          v-model:value="text"
          :style="{ 
            backgroundColor: 'var(--site-theme-bg-light) !important',
            color: 'var(--site-theme-text-primary-inverse) !important'
          }"
          class="focus:border-rp-blue-6"
          @pressEnter="handleAdd"
          @blur="handleAdd"
        />
      </div>
    </div>
    <ScrollList
      class="flex-1 flex flex-col gap-3 pt-1 px-4 overflow-auto"
      :total="total"
      :hasmore="words.length < total"
      @scroll-load="onLoadMore"
    >
      <template #empty>
        <Empty :message="$t('message.noDataTip')" />
      </template>

      <div
        v-for="item in words"
        class="volcabulary"
        :style="{ backgroundColor: 'var(--site-theme-pdf-panel-secondary)' }"
      >
        <WordCard
          :data="item"
          @deleted="remove(item.id!)"
          @submited="mutate(item.id!, $event)"
        />
      </div>
    </ScrollList>
  </section>
</template>

<script setup lang="ts">
import _ from 'lodash';
import { computed, onMounted, onUnmounted, ref, watch, nextTick } from 'vue';
import {
  PlusOutlined,
  InfoCircleOutlined,
  DownOutlined,
  LoadingOutlined,
} from '@ant-design/icons-vue';
import { RectOptions } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';

import type { PromiseType, ValuesType } from 'utility-types';
import type PDFPageView from '@idea/pdf-annotate-viewer/typing/PDFPageView';
import type WordMatchController from '@idea/pdf-annotate-viewer/typing/WordMatchController';
import { ViewerEvent } from '@idea/pdf-annotate-viewer';
import {
  PDF_ANNOTATIONLAYER,
  PDF_ANNOTATIONLAYER_GROUP_VOCABULARY,
  renderAnnotationSvg,
  scaleUpRaw,
  ToolBarType,
  ANNOTATION_MOUSEOVER,
} from '@idea/pdf-annotate-core';
import { Annotation } from '@idea/pdf-annotate-core/render/appendChild';
import { vocabularyStyleMap, VocabularyColorKey } from '@/style/select';

import { selfNoteInfo } from '@/store';
import { useAnnotationStore } from '@/stores/annotationStore';
import {
  useWordNotes,
  DisplayMode,
} from '@common/components/Notes/useWordNote';
import {
  ElementClick,
  getPdfIdFromUrl,
  PageType,
  reportElementClick,
} from '@/api/report';
import { getPDFCurPages } from '@/dom/pdf';

import ScrollList from '@/components/Common/ScrollList.vue';
import Empty from '../common/Empty.vue';
import DotSelect from '../common/DotSelect.vue';
import WordCard from '@common/components/Notes/components/WordCard.vue';
import { usePdfStore } from '~/src/stores/pdfStore';
import { WordColor } from 'go-sea-proto/gen/ts/note/NoteWord';

const emit = defineEmits<{
  (e: 'counted', v: number): void;
}>();

const SIZE_MAX_HIGHLIGHT = 200;

const DISPLAY_OPTIONS = {
  [DisplayMode.DEFAULT]: {
    text: '单处高亮',
    i18n: 'note.vocabulary.singleHighlight',
  },
  [DisplayMode.GLOBAL]: {
    text: '全局高亮',
    i18n: 'note.vocabulary.allHighlight',
  },
  [DisplayMode.HIDDEN]: {
    text: '不高亮',
    i18n: 'note.vocabulary.noHighlight',
  },
} as const;

const pdfId = computed(() => selfNoteInfo.value.pdfId);
const noteId = computed(() => selfNoteInfo.value.noteId);
const pdfStore = usePdfStore();
const annotationStore = useAnnotationStore();
const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(pdfId.value);
});
const pdfAnnotaterRef = computed(() => {
  return pdfStore.getAnnotater(noteId.value);
});

const card = ref<(typeof WordCard)[]>();
const input = ref();
const text = ref('');
const isEditing = ref(false);

const {
  run,
  data: words,
  add,
  mutate,
  remove,
  isAdding,
  gmode,
  gcolor,
  gcolorHex,
  config,
  isConfiging,
  total,
  isLastPage,
  next,
} = useWordNotes(pdfId, noteId, {
  pageSize: SIZE_MAX_HIGHLIGHT,
});

run();

const modeTxt = computed(
  // @ts-ignore
  () => DISPLAY_OPTIONS[gmode.value]?.i18n
);

watch(total, () => {
  emit('counted', total.value);
});

type Matches = {
  [k: string]: ValuesType<
    PromiseType<ReturnType<WordMatchController['getWordMatchRects']>>
  >;
};
const renderMatchesCache: {
  [k: number]:
    | undefined
    | {
        flag: string;
        matches: Matches;
      };
} = {};
const renderHighlights = async (page?: number) => {
  const viewer = pdfViewerRef.value?.getDocumentViewer()?.getPdfViewer();
  const matcher = pdfViewerRef.value?.getWordMatchController();
  const pageNumber = page ?? (viewer?.currentPageNumber || 1);
  const source: undefined | PDFPageView = viewer?.getPageView(pageNumber - 1);
  const instance = pdfAnnotaterRef.value;
  if (!instance || !source) {
    // 延迟执行
    requestIdleCallback(() => renderHighlights(page), {
      timeout: 100,
    });
    return;
  }

  const strs = words.value.map((x) => x.word);
  const { scale, rotation } = source.viewport;
  const flag = `${gcolorHex.value}-${gmode.value}-[${strs}]-${scale}-${rotation}`;

  // console.debug('[Note Vocabulary]rendering page', pageNumber, flag);

  let viewport = Object.create(source.viewport);
  let word2MatchRes: Matches = {};
  if (gmode.value === DisplayMode.GLOBAL) {
    // 这里获取的rects信息是当前viewport下的，渲染svg时无需再次处理
    viewport.scale = 1;
    viewport.rotation = 0;

    const cache = renderMatchesCache[pageNumber];
    if (cache?.flag === flag && cache?.matches) {
      word2MatchRes = cache.matches;
    } else {
      const result = await matcher?.getWordMatchRects(strs, pageNumber - 1);

      word2MatchRes = _.keyBy(result, 'word');
      // console.log('[Note Vocabulary]matchRes', word2MatchRes);
      renderMatchesCache[pageNumber] = { flag, matches: word2MatchRes };
    }
  }

  const annotationsVocabulary = words.value.reduce(
    (arr, { id, word, rectangle }) => {
      let matchedRectangles = rectangle;
      if (word2MatchRes[word]) {
        const { matches } = word2MatchRes[word];
        const pageMatches = matches[pageNumber - 1];

        matchedRectangles = (pageMatches?.rects.map((x) => {
          return {
            ...x,
            pageNumber,
          };
        }) ?? rectangle) as unknown as RectOptions[];
      }

      const rectS = scaleUpRaw(source.viewport.scale, {
        x: rectangle[0]?.x || 0,
        y: rectangle[0]?.y || 0,
      });
      if (matchedRectangles?.length && word) {
        // 必须要为每一处单独渲染g元素，不然浮窗定位比较难搞
        matchedRectangles.forEach((rect: RectOptions, i) => {
          if (
            // gmode.value === DisplayMode.DEFAULT &&
            rect.pageNumber !== pageNumber
          ) {
            return;
          }
          const x: Annotation & { source: string } = {
            uuid: id,
            source:
              Math.abs(rect.x - rectS.x) < 2 && Math.abs(rect.y - rectS.y) < 2
                ? '1'
                : '0',
            type: ToolBarType.Vocabulary,
            rectangles: [rect],
            rectStr: word,
            idea: word,
            color: gcolorHex.value,
            pageNumber,
            tags: [],
          };

          arr.push(x);
        });
      }

      return arr;
    },
    [] as Annotation[]
  );

  renderAnnotationSvg({
    documentId: noteId.value,
    pageNumber,
    source,
    viewport,
    instance,
    annotationsVocabulary,
  });
};

const toggleHightlights = async (v = false) => {
  const viewer = pdfViewerRef.value?.getDocumentViewer()?.getPdfViewer();
  const layers = viewer?.container.querySelectorAll<SVGElement>(
    `.${PDF_ANNOTATIONLAYER} .${PDF_ANNOTATIONLAYER_GROUP_VOCABULARY}`
  );

  if (layers?.length) {
    layers.forEach((layer) => (layer.style.display = v ? 'block' : 'none'));
  }
};

const onHighlightPage = (pageNumber?: number) => {
  const highlightWords = words.value.slice(0, SIZE_MAX_HIGHLIGHT);
  if (highlightWords.length && gmode.value !== DisplayMode.HIDDEN) {
    const renderedPages = pageNumber
      ? [pageNumber]
      : getPDFCurPages(pdfViewerRef.value);

    renderedPages.forEach((page) =>
      requestIdleCallback(() => renderHighlights(page))
    );
    toggleHightlights(true);
  } else {
    toggleHightlights();
  }
};

watch([gcolor, gmode, words], (_, [, , _prevWords]) => {
  onHighlightPage();
});

const onPageRendered = ({
  pageNumber,
}: {
  pageNumber: number;
  source: any;
}) => {
  onHighlightPage(pageNumber);
};

const onLoadMore = (isRetry?: boolean) => {
  if (isLastPage.value) {
    return;
  }

  isRetry ? run() : next();
};

let isEventInited = false;
const initEvents = () => {
  if (isEventInited || !pdfAnnotaterRef.value?.UI || !pdfViewerRef.value) {
    return;
  }
  isEventInited = true;
  // 必须监听TEXT_LAYER_RENDERED，否则match由于textBounds未更新坐标不准
  pdfViewerRef.value?.addEventListener(
    ViewerEvent.TEXT_LAYER_RENDERED,
    onPageRendered
  );
};

watch(() => annotationStore.instantiated, initEvents);

onMounted(() => {
  initEvents();
});

onUnmounted(() => {
  pdfViewerRef.value?.removeEventListener(
    ViewerEvent.TEXT_LAYER_RENDERED,
    onHighlightPage
  );
});

const handleModeChange = async (
  v: Exclude<DisplayMode, DisplayMode.UNRECOGNIZED>
) => {
  if (isConfiging.value) {
    return;
  }

  await config({
    displayMode: +v,
    color: gcolor.value,
  });

  reportElementClick({
    page_type: PageType.note,
    type_parameter: getPdfIdFromUrl(),
    element_name: {
      [DisplayMode.DEFAULT]: ElementClick.note_word_single_present,
      [DisplayMode.GLOBAL]: ElementClick.note_word_global_present,
      [DisplayMode.HIDDEN]: ElementClick.note_word_not_present,
    }[v],
  });
};

const handleColorChange = async (v: number | string) => {
  if (isConfiging.value) {
    return;
  }

  // const c = v as VocabularyColorKey;
  const c = Number(v);
  await config({
    displayMode: gmode.value,
    color: c as WordColor,
  });
};

const handleEdit = async () => {
  isEditing.value = true;
  await nextTick();
  input.value?.focus();
};

const handleAdd = async () => {
  if (isAdding.value) {
    return;
  }

  await add({
    wordInfo: {
      word: text.value,
      rectangle: [],
    },
  });
  isEditing.value = false;
  text.value = '';
};
</script>

<style lang="less" scoped>
.dot-select-hover:hover {
  background-color: var(--site-theme-background-hover);
}

.volcabulary {
  :deep(.word-tt) {
    color: var(--site-theme-text-primary);

    @apply text-base;
    @apply font-bold;
    @apply pt-4;
    @apply bg-transparent;
  }
  :deep(.word-ct) {
    padding-top: 2px;
  }
  :deep(.pronunciation-item .content) {
    color: var(--site-theme-text-secondary);
  }
  :deep(.word-note .idea-markdown-view-container) {
    @apply text-base;
    color: var(--site-theme-text-secondary);

    &:hover {
      background-color: var(--site-theme-background-hover);
    }
  }
}
</style>
