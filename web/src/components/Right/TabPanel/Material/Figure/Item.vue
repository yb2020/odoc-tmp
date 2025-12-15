<template>
  <div class="figure-item">
    <div
      class="title"
      @click="goFigure(item)"
    >
      {{
        item.desc
          ? item.desc
          : item.refIdx
            ? item.refIdx
            : `Page ${item.pageNum}`
      }}
    </div>
    <div class="box">
      <img
        ref="triggerRef"
        class="img"
        :src="item.url"
        @click="addToView(item)"
      >
    </div>
  </div>
</template>
<script lang="ts" setup>
import { PdfFigureAndTableInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { ref, computed } from 'vue';
import { store, selfNoteInfo } from '~/src/store';
import { ParseMutationTypes } from '~/src/store/parse';
import { goBbox } from '~/src/util/bbox';
import { FigureTippyTriggerItem } from '~/src/store/parse/type';
import { usePdfStore } from '~/src/stores/pdfStore';

const props = defineProps<{
  item: PdfFigureAndTableInfo;
  pdfId: string;
}>();

const pdfStore = usePdfStore();

const pdfViewerRef = computed(() => {
  return pdfStore.getViewer(props.pdfId);
});

const triggerRef = ref();

const goFigure = (item: PdfFigureAndTableInfo) => {
  if (item.bbox && pdfViewerRef.value) {
    goBbox(item.pageNum, item.bbox, pdfViewerRef.value);
  }
};
const addToView = (item: PdfFigureAndTableInfo) => {
  const pdfId = selfNoteInfo.value?.pdfId || '';
  const payload: FigureTippyTriggerItem = {
    triggerEle: triggerRef.value,
    id: item.url,
    pdfId,
  };
  store.commit(`parse/${ParseMutationTypes.SET_FIGURE_VIEWER_ITEM}`, payload);
};
</script>
<style scoped lang="less">
.figure-item {
  margin: 0 -20px 12px;
  padding: 0 20px;

  &:hover {
    background-color: var(--site-theme-bg-hover);
  }

  .title {
    cursor: pointer;
    padding: 6px 0;
  }
}

.box {
  padding: 5px;
  background: var(--site-theme-pdf-panel-secondary);
  border-radius: 0px 0px 4px 4px;
}

.img {
  width: 100%;
  cursor: pointer;
}
</style>
