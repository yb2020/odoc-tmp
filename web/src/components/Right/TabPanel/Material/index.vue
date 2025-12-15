<template>
  <PerfectScrollbar class="material-wrap">
    <Metadata
      v-if="selfNoteInfo?.noteId"
      :with-partition="
        envStore.viewerConfig.PDFmetaReginData !== false && !isWebEN
      "
      :with-author-link="envStore.viewerConfig.authorLink !== false"
      :doc-info="docInfo"
      :doc-info-reload="() => docStore.fetchDocInfo(selfNoteInfo.noteId!)"
      :error="fetchError"
    />
    <Attachments
      v-if="docInfo?.docId"
      :doc-id="docInfo.docId"
      :calculateFileMD5="calculateFileMD5"
    />
    <FigureVue :pdf-id="pdfId" />
  </PerfectScrollbar>
</template>
<script setup lang="ts">
import { isOwner, selfNoteInfo } from '~/src/store';
import Metadata from './Metadata/index.vue';
import FigureVue from './Figure/index.vue';
import Attachments from './Attachments/index.vue';
import { useEnvStore } from '~/src/stores/envStore';
import { computed, watch } from 'vue';
import { useDocStore } from '~/src/stores/docStore';
import { storeToRefs } from 'pinia';
import { checkOpenPaper } from '~/src/api/helper';
import { calculateFileMD5 } from '~/src/util/md5';
import { useLanguage } from '@/hooks/useLanguage';

defineProps<{ paperId: string; pdfId: string }>();

const envStore = useEnvStore();

const { isEnUS } = useLanguage();
const isWebEN = isEnUS; // 保持向后兼容的命名

const docStore = useDocStore();
const { docInfo, error: fetchError } = storeToRefs(docStore);

</script>

<style scoped lang="less">
.material-wrap {
  height: 100%;
  padding: 20px;
}
</style>
