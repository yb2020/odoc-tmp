<script setup lang="ts">
import { computed } from 'vue';
import { isOwner, selfNoteInfo } from '~/src/store';
import PDFWebview from '../Main/PDFWebview.vue';
import NoteInfo from '../Private/NoteInfo.vue';
import Tip from './Tip.vue';

const pdfUrl = computed(() => selfNoteInfo.value?.pdfUrl || '');
const pdfId = computed(() => selfNoteInfo.value?.pdfId || '');

</script>

<template>
  <div class="wrap">
    <div class="main">
      <PDFWebview
        :pdf-url="pdfUrl"
        :is-allowed-to-take-note="false"
        scale="auto"
        :is-group-pdf="false"
        :pdf-id="pdfId"
        :paper-id="selfNoteInfo?.paperId || ''"
      />
    </div>
    <Tip v-if="isOwner" />
    <NoteInfo v-else />
  </div>
</template>

<style scoped lang="less">
.wrap {
  height: 100%;
  width: 100%;
  display: flex;
  .main {
    flex: 1;
    height: 100%;
    position: relative;
  }
}
</style>
