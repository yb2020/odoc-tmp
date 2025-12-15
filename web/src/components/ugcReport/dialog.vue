<template>
  <a-modal
    :visible="ugcReportStore.showDialog"
    :title="$t('source.report.button')"
    :footer="null"
    centered
    destroy-on-close
    width="500px"
    wrap-class-name="ugc-report-modal-wrap"
    @cancel="handleCloseDialog"
  >
    <div class="ugc-report-modal-content">
      <iframe
        ref="ugcReportIframeRef"
        width="452"
        height="445"
        title="UgcReport"
        :src="ugcReportDialogUrl"
        class="ugc-report-iframe"
        :onload="handleLoad"
      />
    </div>
  </a-modal>
</template>
<script lang="ts" setup>
import { store } from '~/src/store';
import { useUgcReportStore } from '~/src/stores/ugcReport';
import { getDomainOrigin } from '~/src/util/env';
import { ref, computed, onMounted, onUnmounted } from 'vue';

const ugcReportStore = useUgcReportStore();

const ugcReportIframeRef = ref();

const ugcReportDialogUrl = computed(() => {
  return `${getDomainOrigin()}/dialog/ugcReport`;
});

const handleLoad = () => {
  ugcReportIframeRef.value.contentWindow.postMessage(
    {
      event: 'openUgcReportDialog',
      params: JSON.stringify(ugcReportStore.reportParams),
    },
    '*'
  );
};

const handleCloseDialog = () => {
  ugcReportStore.hideUgcReportDialog();
};

const callback = (event: MessageEvent) => {
  const { data } = event;

  if (data.event === 'closeDialog') {
    handleCloseDialog();
  }
};

onMounted(() => {
  window.addEventListener('message', callback, false);
});

onUnmounted(() => {
  window.removeEventListener('message', callback, false);
});
</script>
<style lang="less">
.ugc-report-modal-wrap {
  .ant-modal-content {
    .ant-modal-header {
      background: #fff;
      border: none;
      padding: 17px 24px;
      .ant-modal-title {
        font-weight: 600;
        font-size: 20px;
        line-height: 30px;
        color: #1d2229;
        text-align: center;
      }
    }
    .ant-modal-close {
      .ant-modal-close-x {
        height: 64px;
        line-height: 64px;
        font-size: 16px;
      }
    }
    .ant-modal-close-icon {
      color: rgba(0, 0, 0, 0.85);
    }
    .ant-modal-body {
      background: #fff;
    }
  }
}
.ugc-report-iframe {
  border: none;
}
</style>
