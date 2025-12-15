<script setup lang="ts">
import { watch, App, ref } from 'vue';
import {
  ViewerEvent,
  ScrollEvent,
  ViewerController,
} from '@idea/pdf-annotate-viewer';
import buildMarkers from '~/src/dom/markers';
import { createPDFSourceCopyrightVue } from '~/src/dom/copyright';

import { CopyrightProps } from '../Copyright/type';
import PDFWebviewVue from './PDFWebview.vue';
import BackToButtonVue from '../ToolBar/BackToButton.vue';
import ToolbarNote from '@/components/ToolBarNote/index.vue';
import ToolbarVue from '@/components/ToolBar/index.vue';
import { ToolbarLeftSideBtnEvent } from '../ToolBar/components/LeftSideBtn.vue';
import { ToolbarPageEvent } from '../ToolBar/components/Page.vue';
import { PDFWebviewSettings } from '~/src/hooks/UserSettings/usePDFWebviewSettings';
import {
  changePDFScale,
  changePDFScaleBySteps,
  goToPDFPage,
} from '~/src/dom/pdf';
import { SideTabCommonSettings } from '~/src/hooks/UserSettings/useSideTabSettings';
import { ToolbarScaleEvent } from '../ToolBar/components/Scale.vue';
import { ownNoteOrVisitSharedNote, pdfStatusInfo } from '~/src/store';
import { UserStatusEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/note/OpenNote';
import { footerNoteInfoHeight } from '~/src/hooks/UserSettings/useFooterNoteInfo';
import { useResizeObserver } from '@vueuse/core';
import { useClip } from '~/src/hooks/useHeaderScreenShot';

const props = defineProps<{
  pdfUrl: string;
  copyRightProps: CopyrightProps | null;
  pdfId: string;
  isGroupPdf: boolean;
  paperId: string;
  isPrivatePaper: boolean;
  pdfWebviewSettings: PDFWebviewSettings;
  leftSideTabSettings: SideTabCommonSettings;
  showAnnotations: boolean;
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
}>();

const emit = defineEmits<{
  (event: 'created', payload: { pdfWebview: ViewerController }): void;
  (event: 'finished', payload: { pdfWebview: ViewerController }): void;
  (event: 'toggleBackToButton', dir: 'up' | 'down' | false): void;
  (event: 'toolbar-event', payload: ToolbarLeftSideBtnEvent): void;
}>();

let pdfViewInstance: ViewerController;

let copyrightVueApp: App<Element>;

const toolBarNoteOffset = ref(0);

watch(
  () => props.isGroupPdf,
  (newVal) => {
    if (!copyrightVueApp) {
      return;
    }
    const iconSpan = copyrightVueApp._container?.querySelector(
      '.js-group-intro .js-pdf-group-info'
    ) as HTMLSpanElement;
    const pdfVersionsSpan = copyrightVueApp._container?.querySelector(
      '.js-pdf-versions .js-pdf-versions-switcher'
    ) as HTMLSpanElement;
    if (newVal) {
      if (iconSpan) {
        iconSpan.style.display = 'inline-block';
      }
      if (pdfVersionsSpan) {
        pdfVersionsSpan.style.display = 'none';
      }
    } else {
      if (iconSpan) {
        iconSpan.style.display = 'none';
      }
      if (pdfVersionsSpan) {
        pdfVersionsSpan.style.display = 'flex';
      }
    }
  }
);

const onToolBarNoteReLoc = () => {
  const docViewer = pdfViewInstance?.getDocumentViewer();
  const pdfViewer = docViewer?.getPdfViewer();
  const wrapperWidth = docViewer?.container?.offsetWidth;
  const curPage = pdfViewer.currentPageNumber - 1;
  const viewport = pdfViewer?.getPageView(curPage);
  const curPageWidth = viewport?.div?.offsetWidth;
  const offset = (wrapperWidth - curPageWidth) / 2 - 12;

  if (offset >= 44) {
    toolBarNoteOffset.value = offset;
  } else {
    toolBarNoteOffset.value = 0;
  }
};

const onPDFWebviewCreated = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  pdfWebview.addEventListener(ViewerEvent.PROGRESS_CHANGE, (percent) => {
    //console.log('loading percent', percent);
  });

  pdfWebview.addEventListener(ViewerEvent.PAGES_INIT, ({ source }) => {
    const div = document.createElement('div');
    source.container.firstElementChild!.appendChild(div);
    copyrightVueApp = createPDFSourceCopyrightVue({
      pageViewerDiv: div,
      pdfWebview,
      copyrightProps: props.copyRightProps,
      isGroupPdf: props.isGroupPdf,
      paperId: props.paperId,
      isPrivatePaper: props.isPrivatePaper,
    });
  });

  if (ownNoteOrVisitSharedNote.value) {
    // 要在下面的build之前调用buildMarkers，先绑定PAGE_RENDERED事件
    buildMarkers(pdfWebview, props.pdfId);
  }

  emit('created', { pdfWebview });
};

const onPDFWebviewFinished = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  pdfViewInstance = pdfWebview;
  numPages.value =
    pdfWebview.getDocumentViewer().getPdfViewer().pdfDocument?.numPages || 0;
  if (ownNoteOrVisitSharedNote.value) {
    pdfWebview
      .getScrollController()
      .addEventListener(
        ScrollEvent.LAST_POSITION_UPDATED,
        ({ lastPosition, source }) => {
          setTimeout(() => {
            const dir = lastPosition
              ? source.getPdfViewerContainer().scrollTop > lastPosition.top
                ? 'up'
                : 'down'
              : false;
            emit('toggleBackToButton', dir);
            toggleBackToButton(dir);
          }, 100);
        }
      );

    const docViewer = pdfWebview.getDocumentViewer();
    const container = docViewer.container;

    useResizeObserver(container, onToolBarNoteReLoc);
    pdfWebview.addEventListener(ViewerEvent.SCALE_CHANGING, onToolBarNoteReLoc);
    pdfWebview.addEventListener(ViewerEvent.PAGE_CHANGING, onToolBarNoteReLoc);

    docViewer.enableWheelToScale();
  }

  isPDFWebviewFinished.value = true;
  emit('finished', { pdfWebview });
};

// ---------START: 回到刚才位置---------
const backToButtonDir = ref<'up' | 'down' | false>(false);
const toggleBackToButton = (dir: 'up' | 'down' | false) => {
  backToButtonDir.value = dir;
};
// ---------END: 回到刚才位置---------

const isPDFWebviewFinished = ref(false);
const numPages = ref(0);

const handleToolbarEvent = (
  payload: ToolbarLeftSideBtnEvent | ToolbarPageEvent | ToolbarScaleEvent
) => {
  if (!pdfViewInstance) {
    return;
  }
  if (payload.type === 'toolbar:page') {
    goToPDFPage(payload.pageNumber, false, pdfViewInstance);
  } else if (payload.type === 'toolbar:leftside') {
    emit('toolbar-event', payload);
  } else if (payload.type === 'toolbar:scale') {
    changePDFScale(payload.scaleValue, pdfViewInstance);
  } else if (payload.type === 'toolbar:scale:increase') {
    changePDFScaleBySteps('increase', pdfViewInstance);
  } else if (payload.type === 'toolbar:scale:decrease') {
    changePDFScaleBySteps('decrease', pdfViewInstance);
  }
};
</script>

<template>
  <div class="pdf-note-webview">
    <PDFWebviewVue
      :pdf-id="pdfId"
      :pdf-url="pdfUrl"
      :scale="pdfWebviewSettings.scalePresetValue || pdfWebviewSettings.scale"
      :show-annotations="showAnnotations"
      @created="onPDFWebviewCreated"
      @finished="onPDFWebviewFinished"
    />
    <BackToButtonVue
      v-if="ownNoteOrVisitSharedNote"
      :dir="backToButtonDir"
      :pdfViewInstance="pdfViewInstance"
    />
    <ToolbarVue
      v-if="pdfStatusInfo.hasPdfAccessFlag && isPDFWebviewFinished"
      :page-props="{ numPages, currentPage: pdfWebviewSettings.currentPage }"
      :scale-props="{
        scale: pdfWebviewSettings.scale,
        scalePresetValue: pdfWebviewSettings.scalePresetValue,
      }"
      :pdf-id="pdfId"
      :left-btn-props="{
        sideTabSettings: leftSideTabSettings,
      }"
      :is-login-user="
        pdfStatusInfo.noteUserStatus === UserStatusEnum.OWNER ||
          pdfStatusInfo.noteUserStatus === UserStatusEnum.GUEST
      "
      :bottom-distance="footerNoteInfoHeight"
      @toolbar-event="handleToolbarEvent"
    />
    <ToolbarNote
      :style="
        toolBarNoteOffset > 0
          ? {
            transform: `translateX(100%)`,
            right: `${toolBarNoteOffset}px`,
          }
          : null
      "
      :clip-selecting="clipSelecting"
      :clip-action="clipAction"
    />
  </div>
</template>
<style lang="less">
.pdf-note-webview {
  flex: 1;
}
</style>
