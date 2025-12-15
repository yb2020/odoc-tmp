<template>
  <div class="export-list">
    <!-- 导出PDF -->
    <div
      class="export-item"
      @click="exportPDF"
    >
      {{ $t('viewer.exportPDF') }}
      <LoadingOutlined v-if="loading.pdf" />
    </div>
    <!-- 导出PDF带笔记 -->
    <!-- <Trigger
      visible
      :report-params="{
        pageType: PageType.note,
        elementName: ElementName.upperNotePDFDownloadPopup,
      }"
    >
      <template #default="slotProps">
        <div
          v-if="withNote"
          class="export-item"
          @click="exportPDFWithNote($event, slotProps)"
        >
          <span>{{ $t('viewer.exportPDF')
          }}<span class="export-tip">{{ $t('viewer.withNote') }}</span></span>
          <LoadingOutlined v-if="loading.note" />
          <UnlockFilled v-else-if="unlock" />
          <a-tooltip
            v-else
            :title="$t('viewer.exportTip', [vipStore.seniorTxt])"
          >
            <LockFilled />
          </a-tooltip>
        </div>
      </template>
    </Trigger> -->
    <!-- 导出笔记 -->
    <Trigger
      visible
      :report-params="{
        pageType: PageType.note,
        elementName: ElementName.upperNoteDownloadPopup,
      }"
    >
      <template #default="slotProps">
        <div
          v-if="withNote"
          class="export-item"
          @click="exportFilesWithNote($event, slotProps)"
        >
          <span>{{ $t('viewer.exportNotes') }}</span>
          <UnlockFilled v-if="unlock" />
          <a-tooltip
            v-else
            :title="$t('viewer.exportTip', [vipStore.seniorTxt])"
          >
            <LockFilled />
          </a-tooltip>
        </div>
      </template>
    </Trigger>
  </div>
</template>
<script lang="ts" setup>
import { ref } from 'vue';
import {
  LockFilled,
  UnlockFilled,
  LoadingOutlined,
} from '@ant-design/icons-vue';
import { message } from 'ant-design-vue';
import {
  getExportEnable,
  getExportPDFUrl,
  validExportCountIncr,
} from '@/api/export';
import { polling } from '~/src/util/polling';
import {
  ExportPdfRequest,
  ExportPdfResponse,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/pdf/Export';
import { downloadUrl } from '@common/utils/url';
import Trigger from '@common/components/Premium/Trigger.vue';
import { useVipStore } from '@common/stores/vip';
import { usePdfStore } from '@/stores/pdfStore';
import {
  ElementName,
  PageType,
  reportElementImpression,
} from '~/src/api/report';

const props = defineProps<{
  pdfId: string;
  noteId: string;
  withNote: boolean;
}>();

const emit = defineEmits<{
  (event: 'hidePopover'): void;
  (event: 'openFileChooseModal'): void;
}>();

const vipStore = useVipStore();
const pdfStore = usePdfStore();
const unlock = ref(false);

const init = async () => {
  // const enable = await getExportEnable(); // 老版本通过设置导出权限 request api /microservice-readpaper-pdf/note/isVip
  const enable = true;
  unlock.value = enable;

  if (!unlock.value) {
    reportElementImpression({
      page_type: PageType.note,
      element_name: ElementName.upperNoteDownloadPopup,
      type_parameter: 'none',
    });
  }
};

init();

const loading = ref({
  pdf: false,
  note: false,
});

const exportPDF = async () => {
  // if (loading.value.pdf) {
  //   return;
  // }
  // loading.value.pdf = true;
  // try {
  //   await validExportCountIncr({}); //老版本：这可能是用来验证或增加用户导出次数统计的接口
  //   loading.value.pdf = false;
  // } catch (error) {
  //   loading.value.pdf = false;
  //   return;
  // }
  const pdfViewer = pdfStore.getViewer(props.pdfId);
  if (pdfViewer) {
    pdfViewer.getDownloadController().download();
  } else {
    message.error('PDF加载解析中，请耐心等待');
  }
  emit('hidePopover');
};

const fetchExportPDFWithNoteStatus = (params: ExportPdfRequest) => {
  return polling<ExportPdfRequest, ExportPdfResponse>({
    fn: getExportPDFUrl,
    maxAttempts: 1,
    interval: 5000,
    params,
    validate: (res) => {
      if (res.needFetch) {
        return false;
      }
      return true;
    },
  });
};

const exportPDFWithNote = async (
  e: MouseEvent,
  slotProps?: {
    onBuyVip: (e: MouseEvent) => void;
  }
) => {
  try {
    if (unlock.value) {
      if (loading.value.note) {
        return;
      }
      loading.value.note = true;
      try {
        const res = fetchExportPDFWithNoteStatus({ noteId: props.noteId });
        const data = await res.request;
        downloadUrl(data.url!, 'export.pdf');
      } catch (error) {}
      loading.value.note = false;
    } else {
      slotProps?.onBuyVip(e);
    }
  } finally {
    emit('hidePopover');
  }
};

// 导出笔记
const exportFilesWithNote = (
  e: MouseEvent,
  slotProps?: {
    onBuyVip: (e: MouseEvent) => void;
  }
) => {
  if (unlock.value) {
    emit('openFileChooseModal');
  } else {
    slotProps?.onBuyVip(e);
    emit('hidePopover');
  }
};
</script>
<style lang="less">
.export-list {
  background-color: #fff;
  color: #1d2229;
  padding: 16px;
  // border: 1px solid #E4E7ED;
  box-shadow:
    0px 3px 6px rgba(0, 0, 0, 0.12),
    0px 6px 16px rgba(0, 0, 0, 0.08),
    0px 9px 28px rgba(0, 0, 0, 0.05);
  border-radius: 4px;
  & > .export-item {
    cursor: pointer;
    background-color: #f7f8fa;
    padding: 13px 16px;
    display: flex;
    justify-content: space-between;
    align-items: center;
    .export-tip {
      color: #86919c;
      font-size: 13px;
    }
    .anticon-lock:hover {
      color: #1f71e0;
    }
    .anticon-unlock {
      color: #52c41a;
    }
  }
  .export-item + .export-item {
    margin-top: 12px;
  }
}
</style>
