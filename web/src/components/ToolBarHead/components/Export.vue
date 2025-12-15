<template>
  <a-dropdown
    v-model:visible="visible"
    :trigger="['click']"
    destroyPopupOnHide
  >
    <div
      class="text-base cursor-pointer"
      @click.prevent
    >
      <a-tooltip :title="$t('toolbar.export')">
        <download-outlined />
      </a-tooltip>
      <!-- <span class="text">{{ $t("viewer.export") }}</span> -->
    </div>
    <template #overlay>
      <ExportPopover
        :pdf-id="pdfId"
        :note-id="noteId"
        :with-note="withNote"
        @hide-popover="hide"
        @open-file-choose-modal="openModal"
      />
    </template>
  </a-dropdown>
  <a-modal
    v-model:visible="modalVisible"
    :title="$t('viewer.exportNotesTip')"
    :footer="null"
    :width="240"
    wrapClassName="export-file-modal"
  >
    <div class="export-list -m-6">
      <div
        class="export-item !justify-center font-medium"
        @click="exportFile('md')"
      >
        Markdown
        <LoadingOutlined v-if="loading.md" />
      </div>
      <div
        class="export-item !justify-center font-medium"
        @click="exportFile('pdf')"
      >
        PDF
        <LoadingOutlined v-if="loading.pdf" />
      </div>
    </div>
  </a-modal>
</template>
<script lang="ts" setup>
import { DownloadOutlined, LoadingOutlined } from '@ant-design/icons-vue';
import { computed, ref } from 'vue';
import ExportPopover from './ExportPopover.vue';
import { useExportFile } from '~/src/hooks/useExport';
import { selfNoteInfo } from '~/src/store';

const props = defineProps<{
  pdfId: string;
  noteId: string;
  withNote: boolean;
}>();

const visible = ref(false);

const hide = () => {
  visible.value = false;
};

const modalVisible = ref(false);

const docNameDecoded = computed(() =>
  decodeURIComponent(
    selfNoteInfo.value?.docName || selfNoteInfo.value?.paperTitle || 'export'
  )
);

const { exportFile, loading } = useExportFile(props.noteId, docNameDecoded);

const openModal = () => {
  modalVisible.value = true;
  hide();
};
</script>
<style lang="less">
.export-file-modal {
  .ant-modal-content {
    background-color: #fff;
    .ant-modal-header {
      background-color: #fff;
      border: none;
      .ant-modal-title {
        color: theme('colors.rp-neutral-10');
      }
    }
  }
}
</style>
