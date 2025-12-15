<template>
  <div>
    <ToolTip
      v-if="isPDFAnnotateActived"
      ref="toolTipRef"
      :pdf-id="pdfBaseInfo.pdfId"
      :note-id="noteInfo.noteId"
      :pdfViewInstance="pdfViewInstance!"
      :pdfAnnotateInstance="PDFAnnotateInstance!"
      :clip-action="clipAction"
    />
    <PDFNotesSvg
      v-if="isPDFAnnotateActived"
      :pdfViewInstance="pdfViewInstance!"
      :pdfAnnotateInstance="PDFAnnotateInstance!"
    />
    <PDFHotNotesSvg
      v-if="isPDFAnnotateActived && isSelfNote"
      :pdfViewInstance="pdfViewInstance!"
      :pdfAnnotateInstance="PDFAnnotateInstance!"
    />
    <PDFVocabularySvg
      v-if="isPDFAnnotateActived && isSelfNote"
      :pdfViewInstance="pdfViewInstance!"
      :pdfAnnotateInstance="PDFAnnotateInstance!"
    />
    <PDFAIHighlightsSvg
      v-if="isPDFAnnotateActived"
      :pdfViewInstance="pdfViewInstance!"
      :pdfAnnotateInstance="PDFAnnotateInstance!"
    />
  </div>
</template>

<script lang="ts" setup>
import { computed, shallowRef, watch } from 'vue';
import { delay } from 'lodash-es';
import isElVisible from 'element-visible';
import { currentGroupId, currentNoteInfo, isOwner, useStore } from '@/store';
import {
  PDFJSAnnotate,
  PDF_ANNOTATIONLAYER,
  renderAnnotate,
  renderAnnotationSvg,
  ANNOTATION_PAGE_RENDERED,
  ANNOTATION_PAGESVG_RENDERED,
  PDF_ANNOTATIONLAYER_GROUP_NOTES,
  PDF_ANNOTATE_SHAPE_DIV,
} from '@idea/pdf-annotate-core';
import type { AnnotationListGroupByVisible } from '@idea/pdf-annotate-core';
import { ViewerController, ViewerEvent } from '@idea/pdf-annotate-viewer';
import { SELF_NOTEINFO_GROUPID } from '~/src/store/base/type';
import ToolTip from '../ToolTip/index.vue';
import PDFNotesSvg from '~/src/components/PDFNotes/NotesSvg.vue';
import PDFHotNotesSvg from '~/src/components/PDFNotes/HotNotesSvg.vue';
import PDFVocabularySvg from '~/src/components/PDFNotes/VocabularySvg.vue';
import PDFAIHighlightsSvg from '~/src/components/PDFNotes/AIHighlightsSvg.vue';
import {
  createAvatar,
  clearAvatar,
  toggleAvatar,
  hideAvatar,
} from '~/src/dom/avatar';
import { PDFBaseInfo } from './PDFWebviewWithNote.vue';
import { PageType, reportModuleImpression } from '@/api/report';
import { ViewObserver } from '@/dom/isInViewport';

import { useAnnotationStore, noteBuffer } from '~/src/stores/annotationStore';
import { PDFPageView } from '@idea/pdfjs-dist/web/pdf_viewer';
import { AnnotationAll } from '~/src/stores/annotationStore/BaseAnnotationController';
import { useTextCallback, useShapeCallback } from './mouseCore';
import { convertTextAnnotation } from '~/src/api/annotations';
import { usePdfStore } from '~/src/stores/pdfStore';
import { useClip } from '~/src/hooks/useHeaderScreenShot';
import { emitter, ANNOTATION_CREATED } from '~/src/util/eventbus';

const props = defineProps<{
  noteInfo: { noteId: string; groupId: string };
  pdfBaseInfo: PDFBaseInfo;
  pdfViewInstance?: ViewerController;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const PDFAnnotateInstance = shallowRef<null | PDFJSAnnotate>(null);

const store = useStore();
const annotationStore = useAnnotationStore();
const userId = computed(() => store.state.user.userInfo?.id);
const isSelfNote = computed(
  () => props.noteInfo.groupId === SELF_NOTEINFO_GROUPID
);
const isPDFAnnotateActived = computed(
  () =>
    PDFAnnotateInstance.value &&
    props.noteInfo.noteId === currentNoteInfo.value?.noteId
);

const pageRender = async ({
  pageNumber,
  source,
}: {
  pageNumber: number;
  source: PDFPageView;
}) => {
  const svg = source.div.querySelector(
    `svg.${PDF_ANNOTATIONLAYER}`
  ) as SVGSVGElement;
  const g = svg?.querySelector<SVGGElement>(
    `.${PDF_ANNOTATIONLAYER_GROUP_NOTES}`
  );
  if (!g) {
    return;
  }

  for (let i = 0; i < g.children.length; i++) {
    const child = g.children[i] as SVGGElement;

    const { pdfAnnotateId } = child.dataset;
    const { top, left, width, height } = child.getBoundingClientRect();

    if (width !== 0 && height !== 0) {
      const { top: pageTop, left: pageLeft } =
        source.div.getBoundingClientRect();

      createAvatar(
        {
          top: top - pageTop + height / 2,
          left: left - pageLeft + width,
          pdfAnnotateId: pdfAnnotateId!,
          pageNumber: +pageNumber,
        },
        props.pdfViewInstance!
      );
    }
  }

  hideAvatar();
  toggleAvatar(
    userId.value as string,
    annotationStore.groupSelfVisible,
    annotationStore.groupOtherVisible,
    annotationStore.groupImageVisible
  );
};

let expoFlag = false;
/**
 * 会被多次调用
 */
const initExpoReport = () => {
  const pdfWebview = props.pdfViewInstance;
  if (!pdfWebview) {
    return;
  }
  // https://github.com/w3c/IntersectionObserver/issues/376
  // IntersectionObserver不适用svg元素
  const pageView = pdfWebview.getDocumentViewer().getPdfViewer();
  // 确保所有他人笔记都被监听
  const hotAnnotationLines = pageView.container.querySelectorAll('g[uuid]');
  ViewObserver.watchElements(
    [].slice.call(hotAnnotationLines),
    {
      listener: () => {
        if (!expoFlag) {
          expoFlag = true;
          reportModuleImpression({
            page_type: PageType.note,
            type_parameter: props.pdfBaseInfo.pdfId,
            module_type: 'others_note',
          });
        }
      },
    },
    pageView.container
  );
};

const initPDFAnnotate = () => {
  const { pdfViewInstance } = props;
  let instance: PDFJSAnnotate;

  if (!pdfViewInstance) {
    return;
  }

  const annotationStore = useAnnotationStore();
  PDFAnnotateInstance.value = instance = new PDFJSAnnotate(
    props.noteInfo.noteId,
    pdfViewInstance
  );
  annotationStore.instantiated = true;

  pdfViewInstance.addEventListener(
    ViewerEvent.TEXT_LAYER_RENDERED,
    ({ source }) => {
      renderPagesAnnotation(source);
    }
  );
  pdfViewInstance.addEventListener(
    ViewerEvent.PAGE_CHANGING,
    ({ source, pageNumber }) => {
      const pageView = source.getPageView(pageNumber - 1);
      renderPagesAnnotation(pageView, true);
    }
  );
  watch(
    () => [annotationStore.personVisible, annotationStore.hotVisible],
    renderVisiblePagesAnnotation
  );

  emitter.on(ANNOTATION_CREATED, renderVisiblePagesAnnotation);

  // watch(
  //   () => annotationStore.crossPageMap,
  //   () => {
  //     console.log('crossPageMap watcher triggered! Re-rendering annotations...');
  //     renderVisiblePagesAnnotation();
  //   },
  //   { deep: true }
  // );
  watch(currentNoteInfo, () => {
    if (props.pdfBaseInfo.pdfId === currentNoteInfo.value?.pdfId) {
      // 共用PDFViewer
      if (props.noteInfo.groupId !== currentGroupId.value) {
        // 代表当前笔记实例渲染到了，需要清除此实例笔记
        instance.clearPersonalNotes();
        instance.UI.disable();
      } else {
        // 恢复当前实例的笔记
        renderVisiblePagesAnnotation();
        instance.UI.enable();
      }
    }
  });

  clearAvatar();

  if (props.noteInfo.groupId !== SELF_NOTEINFO_GROUPID) {
    instance.UI.on(ANNOTATION_PAGESVG_RENDERED, pageRender);
  }
  instance.UI.on(ANNOTATION_PAGE_RENDERED, initExpoReport);

  renderVisiblePagesAnnotation();
};

const renderVisiblePagesAnnotation = () => {
  const { pdfViewInstance } = props;
  if (!pdfViewInstance) {
    return;
  }

  const pageView = pdfViewInstance.getDocumentViewer().getPdfViewer();
  const pages = (pageView._pages ?? []) as PDFPageView[];

  for (const page of pages) {
    renderPagesAnnotation(page);
  }
};

const renderPagesAnnotation = (page: PDFPageView, needCheck = false) => {
  delay(() => {
    if (!PDFAnnotateInstance.value) {
      return;
    }

    const isNotRendered =
      !needCheck || !page.div.querySelector(`.${PDF_ANNOTATE_SHAPE_DIV}`);
    const isLoaded = page.div.getAttribute('data-loaded') === 'true';
    const isVisible = isElVisible(page.div, 0.1);
    console.debug(
      `check rendering annotations for page: ${page.pdfPage.pageNumber}`,
      isNotRendered,
      isLoaded,
      isVisible
    );
    if (isNotRendered && isLoaded && isVisible) {
      loadAndRenderAnnotation(
        props.noteInfo.noteId,
        page.pdfPage.pageNumber,
        page,
        PDFAnnotateInstance.value
      );
      loadAndRenderHighlights(
        props.noteInfo.noteId,
        page.pdfPage.pageNumber,
        page,
        PDFAnnotateInstance.value
      );
    }
  }, 500);
};

const loadAndRenderAnnotation = async (
  documentId: string,
  pageNumber: number,
  source: PDFPageView,
  instance: PDFJSAnnotate
) => {
  // 小组笔记暂不支持非摘录类笔记
  if (
    !isSelfNote.value ||
    props.noteInfo.noteId !== currentNoteInfo.value?.noteId
  ) {
    return;
  }

  annotationStore.controller.loadAnnotationMap();

  renderAnnotate({
    documentId,
    pageNumber,
    source,
    instance,
    shapeTextBuffer: async () => {
      const [, textMap] = await noteBuffer.annotationBuffer!;
      convertTextAnnotation(textMap);
      return textMap;
    },
    handwriteBuffer: () => noteBuffer.handwriteBuffer,
    handwriteAndShapeVisible: annotationStore.personVisible,
    shapeEditable: isOwner.value,
    shapeBuffer: () => noteBuffer.shapeBuffer,
    textCallback: useTextCallback(annotationStore),
    shapeCallback: useShapeCallback(),
  });
};

const loadAndRenderHighlights = async (
  documentId: string,
  pageNumber: number,
  source: PDFPageView,
  instance: PDFJSAnnotate
) => {
  if (props.noteInfo.noteId !== currentNoteInfo.value?.noteId) {
    return;
  }

  const response1 = annotationStore.controller.loadAnnotationMap();
  const response2 = annotationStore.controller.loadHotAnnotationMap();
  await response1, await response2;
  const personAnnotations = annotationStore.crossPageMap[pageNumber] || [];
  const hotAnnotations = annotationStore.pageHotMap[pageNumber] || [];

  let annotationsExtractGrouped: AnnotationListGroupByVisible;
  if (isSelfNote.value) {
    annotationsExtractGrouped = [
      [personAnnotations as any[], annotationStore.personVisible],
      [hotAnnotations as any[], annotationStore.hotVisible],
    ];
  } else {
    const selfAnnotations: AnnotationAll[] = [];
    const otherAnnotations: AnnotationAll[] = [];
    personAnnotations.forEach((anno) => {
      if (anno.commentatorInfoView?.userId === userId.value) {
        selfAnnotations.push(anno);
      } else {
        otherAnnotations.push(anno);
      }
    });
    annotationsExtractGrouped = [
      [selfAnnotations as any[], annotationStore.groupSelfVisible],
      [otherAnnotations as any[], annotationStore.groupOtherVisible],
    ];
  }

  renderAnnotationSvg({
    documentId,
    pageNumber,
    source,
    viewport: source.viewport,
    instance,
    annotationsExtractGrouped,
  });
};

initPDFAnnotate();

const pdfStore = usePdfStore();
watch(
  PDFAnnotateInstance,
  () => {
    if (PDFAnnotateInstance.value) {
      pdfStore.setAnnotater(props.noteInfo.noteId, PDFAnnotateInstance.value);
    }
  },
  { immediate: true }
);
</script>

<style></style>
