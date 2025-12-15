<script setup lang="ts">
import { ref, watch } from 'vue';
import PDFWebview, {
  ViewerEvent,
  PDFViewerScale,
  ViewerController,
} from '@idea/pdf-annotate-viewer';
import ToolbarVue from '@/components/ToolBar/index.vue';

import PDFWebviewVue from './PDFWebview.vue';
import { ToolbarLeftSideBtnEvent } from '../ToolBar/components/LeftSideBtn.vue';
import { ToolbarPageEvent } from '../ToolBar/components/Page.vue';
import {
  changePDFScale,
  changePDFScaleBySteps,
  goToPDFPage,
} from '~/src/dom/pdf';
import usePDFWebviewSettings, {
  PDFWebviewSettings,
} from '~/src/hooks/UserSettings/usePDFWebviewSettings';
import { ToolbarScaleEvent } from '../ToolBar/components/Scale.vue';
import { SaveMode } from '~/src/hooks/UserSettings/useRemoteUserSettings';

const props = defineProps<{
  pdfUrl: string;
  scale?: PDFViewerScale;
  pdfId: string;
  showToolbar: boolean;
  mainPdfWebviewSettings: PDFWebviewSettings;
}>();

const emit = defineEmits<{
  (event: 'created', payload: { pdfWebview: ViewerController }): void;
  (event: 'finished', payload: { pdfWebview: ViewerController }): void;
}>();

const numPages = ref(0);

let pdfWebviewViewerController: ViewerController;

const isPDFWebviewFinished = ref(false);

const onPDFWebviewFinished = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  if (thumbnailWrapper.value) {
    pdfWebview.enableThumbnailViewer(thumbnailWrapper.value);
  }
  numPages.value =
    pdfWebview.getDocumentViewer().getPdfViewer().pdfDocument?.numPages || 0;
  pdfWebviewViewerController = pdfWebview;
  isPDFWebviewFinished.value = true;
  pdfWebview.getDocumentViewer().enableWheelToScale();
  emit('finished', { pdfWebview });
};

const thumbnailWrapper = ref<HTMLDivElement>();

const handleToolbarEvent = (
  payload: ToolbarPageEvent | ToolbarLeftSideBtnEvent | ToolbarScaleEvent
) => {
  if (!pdfWebviewViewerController) {
    return;
  }
  if (payload.type === 'toolbar:page') {
    goToPDFPage(payload.pageNumber, true, pdfWebviewViewerController);
  } else if (payload.type === 'toolbar:scale') {
    changePDFScale(payload.scaleValue, pdfWebviewViewerController);
  } else if (payload.type === 'toolbar:scale:increase') {
    changePDFScaleBySteps('increase', pdfWebviewViewerController);
  } else if (payload.type === 'toolbar:scale:decrease') {
    changePDFScaleBySteps('decrease', pdfWebviewViewerController);
  }
};

const { pdfWebviewSettings, setScaleSetting, setCurrentPageSetting } =
  usePDFWebviewSettings(false);

const onPDFWebviewCreated = ({
  pdfWebview,
}: {
  pdfWebview: ViewerController;
}) => {
  pdfWebview.addEventListener(ViewerEvent.SCALE_CHANGING, (scaleInfo) => {
    scaleInfo.presetValue = scaleInfo.presetValue || '';
    setScaleSetting(scaleInfo, SaveMode.local);
  });

  pdfWebview.addEventListener(ViewerEvent.PAGE_CHANGING, (pageInfo) => {
    setCurrentPageSetting(pageInfo.pageNumber, SaveMode.local);
  });
};

watch(
  () => props.showToolbar,
  (newVal) => {
    if (!isPDFWebviewFinished.value) {
      return;
    }
    if (newVal) {
      setTimeout(() => {
        pdfWebviewViewerController?.getThumbnailViewer()?.forceRendering();
      }, 500);
    } else {
      setTimeout(() => {
        const scaleValue =
          props.mainPdfWebviewSettings.scalePresetValue ||
          props.mainPdfWebviewSettings.scale;
        changePDFScale(
          scaleValue === 'page-width' ? 'page-width' : Number(scaleValue),
          pdfWebviewViewerController
        );
      }, 50);
    }
  }
);
</script>

<template>
  <div
    v-show="showToolbar"
    ref="thumbnailWrapper"
    class="pdf-reader-thumbnails"
  />
  <PDFWebviewVue
    :pdf-id="pdfId"
    :pdf-url="pdfUrl"
    :scale="scale"
    @finished="onPDFWebviewFinished"
    @created="onPDFWebviewCreated"
  />
  <ToolbarVue
    v-if="showToolbar && isPDFWebviewFinished"
    :page-props="{ numPages, currentPage: pdfWebviewSettings.currentPage }"
    :scale-props="{
      scale: pdfWebviewSettings.scale,
      scalePresetValue: pdfWebviewSettings.scalePresetValue,
    }"
    :pdf-id="pdfId"
    :is-login-user="true"
    :bottom-distance="0"
    @toolbar-event="handleToolbarEvent"
  />
</template>

<style scoped lang="less">
.pdf-reader-thumbnails {
  width: 160px;
}
</style>
