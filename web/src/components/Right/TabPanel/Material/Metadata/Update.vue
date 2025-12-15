<script setup lang="ts">
import { onMounted, onUnmounted, ref } from 'vue';
import { LoadingOutlined } from '@ant-design/icons-vue';
import {
  PageType,
  reportReadpaperPopupPaperReferenceClick,
} from '~/src/api/report';
import { ElementClick } from '~/src/api/report';
import { reportElementClick } from '~/src/api/report';
import Tabs from '@/components/Meta/Tabs.vue';

const props = defineProps<{
  paperId: string;
  pdfId: string;
}>();

const loading = ref(false);
const showUpdateDialog = ref(false);

const handleOpenUpdateDialog = () => {
  reportElementClick({
    element_name: ElementClick.renew,
    page_type: PageType.note,
    type_parameter: 'none',
  });
  
  // 直接显示对话框，不再加载外部库
  showUpdateDialog.value = true;
  loading.value = true;
  emit('loading', loading.value);
  
  // 短暂延迟后关闭加载状态
  setTimeout(() => {
    loading.value = false;
    emit('loading', loading.value);
  }, 500);
};

const handleCloseDialog = () => {
  showUpdateDialog.value = false;
};

const handleUpdateSuccess = () => {
  showUpdateDialog.value = false;
  emit('update:success');
};

const emit = defineEmits<{
  (event: 'update:success'): void;
  (event: 'loading', loading: boolean): void;
}>();

onMounted(() => {
  window.addEventListener('message', onIframeMessage, false);
});

onUnmounted(() => {
  window.removeEventListener('message', onIframeMessage, false);
});

const onIframeMessage = (event: MessageEvent) => {
  const { data } = event;
  if (
    data.event === 'clickUpdateCitation' &&
    data.params.paperId === props.paperId
  ) {
    reportReadpaperPopupPaperReferenceClick({
      element_name: ElementClick.renew,
      page_type: PageType.note,
      type_parameter: 'none',
    });
    
    // 直接显示对话框
    showUpdateDialog.value = true;
    loading.value = true;
    emit('loading', loading.value);
    
    // 短暂延迟后关闭加载状态
    setTimeout(() => {
      loading.value = false;
      emit('loading', loading.value);
    }, 500);
  }
};
</script>
<template>
  <div>
    <div
      class="cite-wrap"
      style="margin-left: 4px"
      @click="handleOpenUpdateDialog()"
    >
      <loading-outlined v-if="loading" />
      <i
        v-else
        class="aiknowledge-icon icon-update-cite"
        aria-hidden="true"
      />
      {{ $t('info.update') }}
    </div>
    
    <a-modal
      v-model:visible="showUpdateDialog"
      :title="$t('meta.title')"
      :footer="null"
      :width="800"
      @cancel="handleCloseDialog"
    >
      <Tabs
        v-if="showUpdateDialog"
        :paper-id="paperId"
        :pdf-id="pdfId"
        :page-type="PageType.note"
        @update:success="handleUpdateSuccess"
        @cancel="handleCloseDialog"
      />
    </a-modal>
  </div>
</template>
<style scoped lang="less">
@import url('./Cite.less');
.cite-wrap {
  display: flex;
  justify-content: center;
  > * {
    margin-right: 10px;
  }
}
</style>
