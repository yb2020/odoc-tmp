<template>
  <div
    ref="thumbnailWrapper"
    class="thumbnails"
  />
</template>
<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { ViewerController } from '@idea/pdf-annotate-viewer';

const props = defineProps<{
  pdfId: string;
  pdfViewInstance: ViewerController;
}>();

const thumbnailWrapper = ref<HTMLDivElement>();

onMounted(() => {
  if (thumbnailWrapper.value) {
    props.pdfViewInstance.enableThumbnailViewer(thumbnailWrapper.value);
  }
});
</script>
<style lang="less" scoped>
.thumbnails {
  width: 100%;
  background-color: var(--site-theme-pdf-panel);
  
  :deep(.thumbnail) {
    border: 1px solid var(--site-theme-divider);
    
    &.selected {
      border: 2px solid var(--site-theme-brand);
    }
  }
  
  :deep(.thumbnail-label) {
    color: var(--site-theme-pdf-panel-text);
  }
}
</style>
