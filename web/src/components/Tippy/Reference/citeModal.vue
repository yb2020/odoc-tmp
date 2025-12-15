<template>
  <div>
    <Update
      v-if="
        envStore.viewerConfig.updateMetaButton !== false &&
          store.getters['user/isLogin']
      "
      :pdf-id="pdfId || ''"
      :paper-id="paperData?.paperId || ''"
      @update:success="onUpdateSuccess"
      @loading="onUpdateLoading"
    />
    <Modal
      v-model:visible="citeDialogVisible"
      destroy-on-close
      :title="$t('viewer.citeModalTitle')"
      :footer="null"
      :width="560"
      wrap-class-name="cite-modal-wrap"
      :z-index="100"
      @cancel="handleCiteCancel"
    >
      <div class="cite-modal-content">
        <iframe
          ref="citeIframeRef"
          width="520"
          height="255"
          title="Citation"
          :src="citeDialogUrl"
          class="cite-iframe"
        />
      </div>
    </Modal>
    <Modal
      v-model:visible="manageDialogVisible"
      destroy-on-close
      :title="$t('info.manageCiteTitle')"
      :footer="null"
      :width="560"
      wrap-class-name="manage-citation-style-modal"
      @cancel="handleCancel"
    >
      <div class="cite-modal-content">
        <iframe
          width="520"
          height="314"
          title="manageCitationStyle"
          :src="`${getDomainOrigin()}/dialog/manage`"
          class="cite-iframe"
        />
      </div>
    </Modal>
  </div>
</template>
<script setup lang="ts">
import { Modal } from 'ant-design-vue';
import { PaperDetailInfo } from 'go-sea-proto/gen/ts/paper/Paper'
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { getDomainOrigin } from '~/src/util/env';
import { PageType, getPdfIdFromUrl } from '~/src/api/report';
import { selfNoteInfo, store } from '@/store';
import Update from '@/components/Right/TabPanel/Material/Metadata/Update.vue';
import { useEnvStore } from '~/src/stores/envStore';
const envStore = useEnvStore();

const props = defineProps<{
  paperData: PaperDetailInfo;
}>();

const emit = defineEmits<{
  (event: 'hide'): void;
  (event: 'update:success'): void;
}>();

const pdfId = selfNoteInfo.value?.pdfId ?? getPdfIdFromUrl();

const manageDialogVisible = ref(false);
const citeIframeRef = ref();

const citeDialogVisible = defineModel('citeDialogVisible', { default: false });
const citeDialogUrl = computed(() => {
  return `${getDomainOrigin()}/dialog/cite?pdfId=${pdfId}${
    props.paperData?.paperId ? `&paperId=${props.paperData?.paperId}` : ''
  }&ts=${Date.now()}&pageType=${PageType.note}&typeParameter=${pdfId}`;
});

const onUpdateSuccess = () => {
  emit('update:success');
  if (citeDialogVisible.value) {
    citeIframeRef.value?.contentWindow.postMessage({
      event: 'UpdateManageCitation',
      params: {},
    });
  }
};

const onUpdateLoading = (loading: boolean) => {
  if (!loading) {
    citeIframeRef.value?.contentWindow?.postMessage(
      {
        event: 'stopLoadingClickUpdateCitation',
        params: {},
      },
      '*'
    );
  }
};
const handleCancel = () => {
  manageDialogVisible.value = false;
  citeIframeRef.value?.contentWindow?.postMessage(
    {
      event: 'cancelManageCitationStyleModal',
      params: {},
    },
    '*'
  );
};
const callback = (event: MessageEvent) => {
  const { data } = event;
  if (data.event === 'clickManageCitationStyle') {
    manageDialogVisible.value = true;
  }
};
const handleCiteCancel = () => {
  citeDialogVisible.value = false;
};

onMounted(() => {
  window.addEventListener('message', callback, false);
});

onUnmounted(() => {
  window.removeEventListener('message', callback, false);
});
</script>

<style scoped lang="less">
.cite-wrap {
  height: 0;
}
</style>
