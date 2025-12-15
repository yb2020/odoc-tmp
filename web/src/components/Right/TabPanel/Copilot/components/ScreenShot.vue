<template>
  <div class="py-1 mx-[10px] flex">
    <a-tooltip
      :title="$t('aiCopilot.screenshotTip')"
      placement="bottomLeft"
    >
      <div
        :class="[
          'screenshot-button',
          'cursor-pointer w-6 h-6 flex justify-center items-center',
          clipSelecting ? 'active' : '',
        ]"
        @click="handleScreenShot"
      >
        <i
          class="aiknowledge-icon icon-crop text-base"
          aria-hidden="true"
        />
      </div>
    </a-tooltip>
  </div>
</template>
<script lang="ts" setup>
import {
  useClip,
  ClipPayload,
  SUBMIT_CLIP_PAYLOAD_TO_COPILTO,
} from '@/hooks/useHeaderScreenShot';
import { onMounted, onUnmounted } from 'vue';
import { AskImagePayload } from '~/src/hooks/useCopilot';
import { AnnotationRect } from '~/src/stores/annotationStore/BaseAnnotationController';
import { emitter } from '~/src/util/eventbus';
import { ImageMimeType, compressImageBase64 } from '~/src/util/image';

const props = defineProps<{
  clipSelecting: boolean;
  clipAction: ReturnType<typeof useClip>['clipAction'];
  addImageToUpload?: (payload: { base64: string; pageNumber?: number }) => boolean;
}>();

const handleScreenShot = async () => {
  if (props.clipSelecting) {
    props.clipAction.cancelCut();
    return;
  }

  const payload = await props.clipAction.init();
  if (payload) {
    emitter.emit(SUBMIT_CLIP_PAYLOAD_TO_COPILTO, payload);
  }
};

const submit = async (payload: ClipPayload | AskImagePayload) => {
  // 如果是直接传入的 AskImagePayload
  if ((payload as AskImagePayload).base64) {
    // 如果有 addImageToUpload 方法，则添加到暂存区
    if (props.addImageToUpload) {
      props.addImageToUpload(payload as AskImagePayload);
    }
    return;
  }

  // 处理截图
  const { newRectElement, newCanvasElement, newRectAnnotation } =
    payload as ClipPayload;

  const base64 = compressImageBase64(
    newCanvasElement as HTMLCanvasElement,
    0.8,
    ImageMimeType.JPEG
  );

  // 如果有 addImageToUpload 方法，则添加到暂存区
  if (props.addImageToUpload) {
    props.addImageToUpload({
      base64,
      pageNumber: (newRectAnnotation as AnnotationRect).pageNumber,
    });
  }

  // 清理截图元素
  newRectElement?.remove();
  props.clipAction.clearClip();
};

onMounted(() => {
  emitter.on(SUBMIT_CLIP_PAYLOAD_TO_COPILTO, submit as () => Promise<void>);
});

onUnmounted(() => {
  emitter.off(SUBMIT_CLIP_PAYLOAD_TO_COPILTO, submit as () => Promise<void>);
});
</script>

<style lang="less" scoped>
.screenshot-button {
  border-radius: 4px;
  color: var(--site-theme-text-primary);

  &:hover,
  &.active {
    background-color: var(--site-theme-background-hover);
  }
}
</style>
