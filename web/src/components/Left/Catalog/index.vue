<template>
  <div class="calalogs">
    <div
      v-if="loading"
      class="center"
    >
      <a-spin
        v-if="loading"
        :tip="$t('message.parsingCatalogue')"
      />
    </div>
    <Tree
      v-else-if="catalogs"
      :catalogs="catalogs"
      :pdfViewInstance="pdfViewInstance"
    />
    <div
      v-else
      class="center"
    >
      <div class="text">
        {{ $t('message.failToParseCatalogue') }}
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref } from 'vue';
import { ViewerController } from '@idea/pdf-annotate-viewer';
// import { PdfCatalogueInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/PdfParse';
import { PdfCatalogueInfo } from 'go-sea-proto/gen/ts/pdf/PdfParse';
import { getCatalog } from '~/src/api/category';
import { getPDFCatalog } from '~/src/util/pdf';
import Tree from './Tree.vue';

const loading = ref(true);
const catalogs = ref<PdfCatalogueInfo>();

const props = defineProps<{
  pdfId: string;
  pdfViewInstance: ViewerController;
}>();

const fetchData = async () => {
  const data = await getCatalog({
    pdfId: props.pdfId,
  }).catch(() => {});

  if (data && data.pdfCatalogue?.child?.length) {
    catalogs.value = data.pdfCatalogue;
  } else {
    const catalog = await getPDFCatalog(props.pdfViewInstance);

    if (catalog.child.length) {
      catalogs.value = catalog;
    }
  }

  loading.value = false;
};

fetchData();
</script>
<style lang="less" scoped>
.calalogs {
  height: calc(100% - 25px);
  width: 100%;
  background-color: var(--site-theme-pdf-panel);
  
  .center {
    height: 100%;
    text-align: center;
    display: flex;
    justify-content: center;
    align-items: center;
    flex-direction: column;
    overflow: hidden;
    
    .text {
      color: var(--site-theme-pdf-left-text);
      font-size: 13px;
    }
  }
}
</style>
