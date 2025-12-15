<template>
  <teleport to="body">
    <div
      ref="toolTip"
      :class="[TOOLTIP_CLASSNAME, JS_IGNORE_MOUSE_OUTSIDE]"
      @click.stop
      @mousedown.stop
    >
      <div class="highlight-comment-options">
        <AIBtn
          v-if="
            showTool.ai && isSelfNoteInfo && toolTipType === ToolBarType.rect
          "
          @ask-copilot="askCopilot"
        />

        <TranslateBtn
          v-if="showTool.translate"
          :toolTipType="toolTipType"
          :annotateId="annotateId"
          :pdfId="pdfId"
          :rects="rectList"
          :pageTexts="pageSelectTexts"
          :get-ocr-image="getOcrImage"
          :add-ocr-note="addOcrNote"
          :alive="annotationStore.hotTipVisible"
        />

        <Highlight
          v-if="isOwner && showTool.highlight && !multiSegment"
          :pdfId="pdfId"
          :noteId="noteId"
          :toolTipType="toolTipType"
          :annotateId="annotateId"
          :rects="rectList"
          :rect-str="rectStr"
          :page-number="pageNumber"
          :color="annotationColor"
          :is-highlight-annotation="annotationIsHighlight"
          :create-rect="createRect"
          @create-highlight="createHighlight"
          @hide-tool-tip="hideToolTip"
        />

        <Comment
          v-if="isOwner && showTool.comment && !multiSegment"
          :toolTipType="toolTipType"
          :annotateId="annotateId"
          :rects="rectList"
          :rect-str="rectStr"
          :page-number="pageNumber"
          :is-highlight-annotation="annotationIsHighlight"
          :create-rect="createRect"
          :pdfId="pdfId"
          @hide-tool-tip="hideToolTip"
        />
      </div>
    </div>
  </teleport>
</template>

<script lang="ts" setup>
import { ref, onUnmounted, computed, nextTick, onMounted } from 'vue';
import {
  ToolBarType,
  ANNOTATION_CLICK,
  attrSelector,
  PDF_ANNOTATE_ID,
  PDFJSAnnotate,
} from '@idea/pdf-annotate-core';
// import setAttributes from '@idea/pdf-annotate-core/utils/setAttributes';
import setAttributes from '~/src/pdf-annotate-core/utils/setAttributes';
import {
  PageSelectText,
  ViewerEvent,
  ViewerController,
  JS_IGNORE_MOUSE_OUTSIDE,
} from '@idea/pdf-annotate-viewer';
import Highlight from './Highlight.vue';
import Comment from './Comment.vue';
import {
  store,
  isOwner,
  isSelfNoteInfo,
  ownNoteOrVisitSharedNote,
  currentGroupId,
} from '@/store';
import { TOOLTIP_CLASSNAME } from '~/src/dom/tooltip';
import {
  createEditOverlay,
  destroyEditOverlay,
  appendPopperElement,
} from '@/dom/tooltip';
import TranslateBtn from './TranslateBtn.vue';
import AIBtn from './AIBtn.vue';
import { clearArrow, connect } from '@/dom/arrow';
import scrollIntoView from 'scroll-into-view-if-needed';
import { RightSideBarType } from '../Right/TabPanel/type';
import {
  createTranslateTippyVue,
  enableTranslateTippyVueToAddNote,
  checkMultiSegment,
} from '~/src/dom/tippy';
import { useRightSideTabSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import { IS_MOBILE } from '~/src/util/env';
import { emitter, COPILOT_ASK } from '~/src/util/eventbus';
import { useAnnotationStore } from '~/src/stores/annotationStore';
import { RectOptions } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/notev2/Common';
import { useEnvStore } from '~/src/stores/envStore';
import { delay } from '@idea/aiknowledge-special-util/delay';
import {
  ANNOTATION_SCREENSHOT,
  ClipPayload,
  SUBMIT_CLIP_PAYLOAD_TO_COPILTO,
  ScreenShotPayload,
  saveRectAnnotation,
  uploadImageAndUpdateAnnotation,
  useClip,
} from '~/src/hooks/useHeaderScreenShot';
import { AskImagePayload } from '~/src/hooks/useCopilot';
import { AnnotationRect } from '~/src/stores/annotationStore/BaseAnnotationController';
import { rectStyleMap } from '~/src/style/select';
import { useCommentGlobalState } from '~/src/hooks/useNoteState';
import { ImageMimeType, compressImageBase64 } from '~/src/util/image';

const props = defineProps<{
  pdfId: string;
  noteId: string;
  pdfViewInstance: ViewerController;
  pdfAnnotateInstance: PDFJSAnnotate;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const toolTipType = ref<ToolBarType>(ToolBarType.select);

const annotationStore = useAnnotationStore();
const envStore = useEnvStore();

const showTool = ref({
  ai: true,
  highlight: true,
  translate: envStore.viewerConfig.toolTipTranslateIcon !== false,
  comment: true,
});

const toolTip = ref<HTMLElement>();

const annotateId = ref<string>();
const anntationItem = computed(() => {
  if (!annotateId.value) {
    return null;
  }

  return annotationStore.controller.findAnnotation(annotateId.value)
    ?.annotation;
});
const annotationColor = computed(() => {
  return anntationItem.value?.color ?? '';
});
const annotationIsHighlight = computed(() => {
  if (currentGroupId.value !== SELF_NOTEINFO_GROUPID) {
    return false;
  }

  return anntationItem.value?.isHighlight ?? false;
});

const pdfViewer = props.pdfViewInstance;
const UI = props.pdfAnnotateInstance!.UI;

const { activeTab, switchAndShowTab, checkNoteTabVisible } =
  useRightSideTabSettings();

const connectComment = (pdfAnnotateId: string) => {
  setTimeout(() => {
    const comment = document.querySelector(
      `.note${attrSelector(PDF_ANNOTATE_ID, pdfAnnotateId)}`
    );

    if (!comment) return;

    scrollIntoView(comment, {
      block: 'center',
      inline: 'center',
    });

    connect(pdfAnnotateId);
  }, 300);
};

async function annotationClick(target: SVGGElement, event: MouseEvent) {
  console.log('annotationClick triggered!', { target, event });
  
  clear();

  const commentator = target.getAttribute('commentator-info-view')!;

  // if (commentator) {
  //   const { userId } = JSON.parse(commentator);

  //   if (userId !== store.state.user.userInfo?.id) {
  //     showTool.value.comment = false;
  //     showTool.value.highlight = false;
  //   }
  // }
  if (commentator) {
    const commentatorInfo = JSON.parse(commentator);
    if (commentatorInfo) {
      const { userId } = commentatorInfo;
      if (userId !== store.state.user.userInfo?.id) {
        showTool.value.comment = false;
        showTool.value.highlight = false;
      }
    }
  }

  const page = +(target.getAttribute('page-number') || 1);
  const path = getMouseEventPath(event);
  
  // 安全检查 PDF viewer 是否已初始化
  const pdfViewerInstance = pdfViewer.getDocumentViewer().getPdfViewer();
  const pageElement = pdfViewerInstance?.viewer?.querySelector(`.page[data-page-number="${page}"]`);
  
  if (
    path.includes(toolTip.value) ||
    !pageElement ||
    !path.includes(pageElement)
  ) {
    return;
  }

  if (!ownNoteOrVisitSharedNote.value) {
    return;
  }

  if (IS_MOBILE) {
    return;
  }

  const selectType = +(target.getAttribute('type') as string);
  if (![ToolBarType.select, ToolBarType.rect].includes(selectType)) {
    return;
  }

  const pdfAnnotateId = target.dataset.pdfAnnotateId as string;

  const _rectStr = target.getAttribute('rect-str') || '';

  toolTipType.value = selectType;

  annotateId.value = pdfAnnotateId;

  pageNumber.value = page;

  const overlay = createEditOverlay(
    annotateId.value,
    page,
    selectType === ToolBarType.rect ? { border: 'none' } : {},
    props.pdfViewInstance
  );

  if (!toolTip.value) {
    return;
  }

  rectStr.value = _rectStr;

  pdfViewer?.getDocumentViewer().setSelectedText(_rectStr);

  const result = annotationStore.controller.findAnnotation(annotateId.value);
  rectList.value = result.annotation?.rectangles ?? [];

  multiSegment.value = false;

  const pageNum = result.pageNumber || page;
  // const viewport = pdfViewer
  //   ?.getDocumentViewer()
  //   .getPdfViewer()
  //   .getPageView(pageNum - 1).viewport;
  const pageView = pdfViewer
    ?.getDocumentViewer()
    .getPdfViewer()
    .getPageView(pageNum - 1);

  if (!pageView) {
    return;
  }
  const viewport = pageView.viewport;

  const selectTexts: PageSelectText[] = [];
  rectList.value.forEach((rect) => {
    const equal = (text: PageSelectText) => text.pageNum === rect.pageNumber;
    if (!selectTexts.some(equal)) {
      selectTexts.push({
        pageNum: rect.pageNumber,
        text: '',
        viewport,
        rects: [],
        multiSegment: false,
        crossing: false,
      });
    }

    selectTexts.find(equal)?.rects.push({
      x: rect.x * viewport.scale,
      y: rect.y * viewport.scale,
      width: rect.width * viewport.scale,
      height: rect.height * viewport.scale,
      text: '',
      shouldScaleText: true,
      pageNumber: rect.pageNumber,
    } as any);
  });

  if (selectTexts[0]) {
    selectTexts[0].text = rectStr.value;
  }

  pageSelectTexts.value = selectTexts;

  const div = overlay.find((item) => item.rectPageNumber === page)
    ?.div as HTMLDivElement;
  const bound = div.getBoundingClientRect();

  const { top, left, width, height } = bound;

  appendPopperElement(toolTip.value!, {
    x: left,
    y: top + window.scrollY,
    width,
    height,
  });

  handleOutClick();
}

const rectList = ref<RectOptions[]>([]);
const rectStr = ref('');
const pageNumber = ref(1);
const multiSegment = ref(false);
const newRect = ref(false);

const createHighlight = (id: string) => {
  annotateId.value = id;
};

const hideToolTip = () => {
  // 检查翻译弹窗是否被钉住，如果钉住则不关闭
  if (!createTranslateTippyVue.isPinned()) {
    createTranslateTippyVue.hide();
  }
  toolTip.value && (toolTip.value.style.display = 'none');
};

const clear = () => {
  document.body.removeEventListener('click', clear);

  hideToolTip();

  clearArrow();

  annotateId.value = '';
  toolTipType.value = ToolBarType.None;

  annotationStore.currentAnnotationId = '';

  destroyEditOverlay(props.pdfViewInstance);

  showTool.value.ai = true;
  showTool.value.comment = true;
  showTool.value.highlight = true;
};

const handleOutClick = async () => {
  await delay(100);
  document.body.addEventListener('click', clear);
};

const commonState = useCommentGlobalState();

const createRect = async (note?: string) => {
  const { newRectElement, newCanvasElement, newRectAnnotation } =
    props.clipAction.getClipPayload();
  const annotation = newRectAnnotation as AnnotationRect;
  const { color, fill } = rectStyleMap[commonState.value.styleId];
  setAttributes(newRectElement as SVGRectElement, {
    fill,
    color,
    stroke: color,
    lineColor: color,
  });

  Object.assign(annotation, {
    styleId: commonState.value.styleId,
    fill,
    color,
  });

  await saveRectAnnotation(
    props.pdfAnnotateInstance,

    store.state.user.userInfo,
    annotationStore,

    annotation,
    newRectElement as SVGRectElement,
    newCanvasElement as HTMLCanvasElement,

    note
  );

  if (checkNoteTabVisible()) {
    connectComment(annotation.uuid);
  }

  await uploadImageAndUpdateAnnotation(
    annotationStore,
    annotation,
    newCanvasElement as HTMLCanvasElement
  );

  props.clipAction.clearClip();

  return annotation.uuid;
};

const addOcrNote = async (translation: string) => {
  const uuid = await createRect(translation);
  annotationStore.currentAnnotationId = uuid;
};

const getOcrImage = async () => {
  let canvasElement: HTMLCanvasElement;

  if (!annotateId.value) {
    canvasElement = props.clipAction.getClipPayload()
      .newCanvasElement as HTMLCanvasElement;
  } else {
    canvasElement = await createCanvasElementFromRectAnnotation();
  }

  const picBase64 = compressImageBase64(canvasElement, 0.4, ImageMimeType.WEBP);

  if (picBase64.length > 500 * 1024) {
    throw new Error('图片过大，请重新截图');
  }

  return picBase64;
};

const createCanvasElementFromRectAnnotation = async () => {
  const canvasElement = document.createElement('canvas');
  const imgElement = document.createElement('img');
  imgElement.crossOrigin = 'anonymous';

  const rectAnnotation = annotationStore.controller.findAnnotation(
    annotateId.value as string
  )?.annotation as AnnotationRect;

  await new Promise((resolve, reject) => {
    imgElement.onload = resolve;
    imgElement.onerror = reject;
    imgElement.src = rectAnnotation.picUrl;
  });

  canvasElement.width = imgElement.width;
  canvasElement.height = imgElement.height;
  const context = canvasElement.getContext('2d') as CanvasRenderingContext2D;
  context.drawImage(imgElement, 0, 0);

  return canvasElement;
};

const annotationScreenShot = ({
  rect,
  pageNum,
  visibleBtns,
}: ScreenShotPayload) => {
  clear();

  multiSegment.value = false;

  rectList.value = [
    {
      pageNumber: pageNum,
      rotation: 0,
      ...rect,
    },
  ];

  rectStr.value = '';

  pageNumber.value = pageNum;

  toolTipType.value = ToolBarType.rect;

  newRect.value = true;

  annotateId.value = '';

  if (visibleBtns) {
    // @ts-ignore
    showTool.value = {
      ...visibleBtns,
    };
  }

  const { x, y, width, height } = rect;
  appendPopperElement(toolTip.value!, {
    y: y + window.scrollY,
    x,
    width,
    height,
  });

  handleOutClick();
};

const pageSelectTexts = ref<PageSelectText[]>([]);

const textSelect = (pageTexts: PageSelectText[]) => {
  console.log('pageTexts', pageTexts);

  if (activeTab.value === RightSideBarType.Copilot) {
    emitter.emit(COPILOT_ASK, pageTexts);
    // return
  }

  clear();

  const { pageNum } = pageTexts[pageTexts.length - 1];

  const container = props.pdfViewInstance?.getDocumentViewer().container!;

  multiSegment.value = checkMultiSegment(pageTexts);

  rectList.value = pageTexts.reduce((all: RectOptions[], current) => {
    const { rects, pageNum } = current;

    rects
      .map((item) => ({
        rotation: 0,
        ...item,
        pageNumber: pageNum,
      }))
      .forEach((item) => {
        all.push(item);
      });

    return all;
  }, []);

  rectStr.value = pageTexts.reduce((prev, current) => prev + current.text, '');

  pageNumber.value = pageNum;

  const bound = container
    .querySelector(`.page[data-page-number="${pageNum}"]`)!
    .getBoundingClientRect();

  const { left, top } = bound;

  const lastText = pageTexts.slice(-1).pop()!.rects.slice(-1).pop()!;

  const { x, y, width, height } = lastText;

  appendPopperElement(toolTip.value!, {
    y: y + window.scrollY + top,
    x: x + left,
    width,
    height,
  });

  toolTipType.value = ToolBarType.select;
  annotateId.value = '';

  handleOutClick();

  pageSelectTexts.value = pageTexts;
};

onMounted(() => {
  console.log('tooltip init events');
  UI.on(ANNOTATION_CLICK, annotationClick);
  UI.on(ANNOTATION_SCREENSHOT, annotationScreenShot);

  pdfViewer?.addEventListener(ViewerEvent.TEXT_SELECT, textSelect);
});

onUnmounted(() => {
  console.log('tooltip remove events');
  UI.off(ANNOTATION_CLICK, annotationClick);
  UI.off(ANNOTATION_SCREENSHOT, annotationScreenShot);

  pdfViewer.removeEventListener(ViewerEvent.TEXT_SELECT, textSelect);
});

store.watch(
  () => ({
    currentGroupId: currentGroupId.value,
    isOwner: isOwner.value,
  }),
  (newVal) => {
    if (newVal.isOwner && newVal.currentGroupId === SELF_NOTEINFO_GROUPID) {
      enableTranslateTippyVueToAddNote(true);
    } else {
      enableTranslateTippyVueToAddNote(false);
    }
    createTranslateTippyVue.hide();
  },
  {
    immediate: true,
  }
);

const getMouseEventPath = (event: MouseEvent) => {
  const path: Array<Element | Window | Document | null | undefined> = [];
  let currentElement: HTMLElement | null = event.target as HTMLElement | null;
  while (currentElement) {
    path.push(currentElement);
    currentElement = currentElement.parentElement;
  }
  if (!path.includes(window)) {
    if (!path.includes(document)) {
      path.push(document);
    }

    path.push(window);
  }

  return path;
};

const askCopilot = async () => {
  hideToolTip();

  let payload: ClipPayload | AskImagePayload;

  if (annotateId.value) {
    const result = annotationStore.controller.findAnnotation(annotateId.value);

    payload = {
      base64: (result.annotation as AnnotationRect).picUrl,
      pageNumber: result.pageNumber,
    };
  } else {
    payload = props.clipAction.getClipPayload();

    if (!payload) {
      return;
    }
  }

  switchAndShowTab(RightSideBarType.Copilot, {});

  await delay(100);
  await nextTick();

  emitter.emit(SUBMIT_CLIP_PAYLOAD_TO_COPILTO, payload);
};
</script>

<style lang="less" scoped>
.pdf-content-tooltip {
  background: #393c3e;
  color: #fff;
  font-weight: bold;
  font-size: 13px;
  position: absolute;
  top: -999px;
  display: none;
  outline: none;
  z-index: 999;

  border-radius: 4px;
}

.highlight-comment-options {
  display: flex;
  cursor: pointer;
  border-radius: 4px;
}

.color-options {
  padding: 6px 9px;
}
</style>
