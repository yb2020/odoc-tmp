<template>
  <div class="py-1 mx-[10px] flex">
    <a-tooltip
      :title="$t('translate.ocrTip')"
      placement="bottomLeft"
    >
      <div
        :class="[
          'cursor-pointer w-6 h-6 flex justify-center items-center hover:bg-gray-200',
          clipSelecting ? 'bg-gray-200' : '',
        ]"
        @click="handleScreenShot"
      >
        <i
          class="aiknowledge-icon icon-crop text-base leading-4 h-4"
          aria-hidden="true"
        />
      </div>
    </a-tooltip>
  </div>
</template>
<script lang="ts" setup>
import {
  ANNOTATION_SCREENSHOT,
  ScreenShotPayload,
  useClip,
} from '@/hooks/useHeaderScreenShot';
import { nextTick } from 'vue';
import { PDFJSAnnotate } from '@idea/pdf-annotate-core';
import {
  executeAndSetBeforeChangeTab,
  removeBeforeChangeTab,
} from '~/src/hooks/UserSettings/useSideTabSettings';
import { AnnotationRect } from '~/src/stores/annotationStore/BaseAnnotationController';

const props = defineProps<{
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
  pdfAnnotater?: PDFJSAnnotate;
}>();

const handleScreenShot = () => {
  if (props.clipSelecting) {
    cancelScreenShot();
  } else {
    startScreehShot();
  }
};

const startScreehShot = async () => {
  executeAndSetBeforeChangeTab(cancelScreenShot);

  await nextTick();

  const payload = await props.clipAction.init(true);
  if (!payload) {
    return;
  }

  const { newRectElement, newRectAnnotation } = payload;

  props.pdfAnnotater?.UI.emit(ANNOTATION_SCREENSHOT, {
    rect: newRectElement?.getBoundingClientRect(),
    pageNum: (newRectAnnotation as AnnotationRect).pageNumber,
    visibleBtns: {
      translate: true,
    },
  } as ScreenShotPayload);
};

const cancelScreenShot = () => {
  props.clipAction.cancelCut();
  removeBeforeChangeTab(cancelScreenShot);
};
</script>
