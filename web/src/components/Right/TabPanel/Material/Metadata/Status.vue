<template>
  <div class="metadata-field">
    <div class="metadata-field-name">
      {{ $t('viewer.parsedProgress') }}
    </div>
    <div class="metadata-field-value">
      <span>{{ statusText }}</span>
      <ReloadOutlined 
        v-if="isFailed || !isAllCompleted" 
        :class="['reparse-icon', { 'reparse-spinning': isReparsing }]"
        @click="reparse" 
      />
    </div>
  </div>
  <div class="metadata-field">
    <div class="metadata-field-name">
      {{ $t('viewer.fullTextTranslateStatus') }}
    </div>
    <div class="metadata-field-value">
      <CheckCircleFilled v-if="isContentParsed" class="status-icon status-success" />
      <CloseCircleFilled v-else class="status-icon status-error" />
    </div>
  </div>
  <div class="metadata-field">
    <div class="metadata-field-name">
      {{ $t('viewer.aiCopilotStatus') }}
    </div>
    <div class="metadata-field-value">
      <CheckCircleFilled v-if="isEmbeddingReady" class="status-icon status-success" />
      <CloseCircleFilled v-else class="status-icon status-error" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, watch, ref, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { message } from 'ant-design-vue';
import { CheckCircleFilled, CloseCircleFilled, ReloadOutlined } from '@ant-design/icons-vue';
import { isParseFailedStatus, calculateParsedProgress } from '~/src/utils/pdf-upload/statusMapper';
import { reParsePaper } from '~/src/api/parse';
import { getUserDocStatusByIds } from '~/src/api/material';
import { UserDocParsedStatusEnum } from 'go-sea-proto/gen/ts/doc/UserDocParsedStatus';
import { useDocStore } from '~/src/stores/docStore';

const props = defineProps<{
  docId: bigint | undefined;
  pdfId: bigint | undefined;
}>();

const { t } = useI18n();
const docStore = useDocStore();

// 从 API 获取的状态
const apiParsedStatus = ref<number | undefined>();
const apiEmbeddingStatus = ref<number | undefined>();

// 是否正在重新解析
const isReparsing = ref(false);

const isFailed = computed(() => {
  if (apiParsedStatus.value === undefined) {
    return false;
  }
  return isParseFailedStatus(apiParsedStatus.value);
});

// 判断两个状态是否都完成
const isAllCompleted = computed(() => {
  return apiParsedStatus.value === UserDocParsedStatusEnum.CONTENT_DATA_PARSED 
    && apiEmbeddingStatus.value === UserDocParsedStatusEnum.EMBEDDED;
});

const statusText = computed(() => {
  if (apiParsedStatus.value === undefined) {
    return '';
  }
  
  return calculateParsedProgress(apiParsedStatus.value, apiEmbeddingStatus.value || 0);
});

const isContentParsed = computed(() => {
  return apiParsedStatus.value === UserDocParsedStatusEnum.CONTENT_DATA_PARSED;
});

const isEmbeddingReady = computed(() => {
  return apiEmbeddingStatus.value === UserDocParsedStatusEnum.EMBEDDED;
});

// 判断是否需要继续轮询
const needsPolling = computed(() => {
  return !isContentParsed.value || !isEmbeddingReady.value;
});

// 轮询定时器
let pollingTimer: ReturnType<typeof setTimeout> | null = null;
const POLLING_INTERVAL = 3000; // 3秒轮询一次

// 获取文档状态
const fetchDocStatus = async () => {
  if (!props.docId) return;
  
  try {
    // @ts-ignore
    const response = await getUserDocStatusByIds({
      docIds: [props.docId]
    });
    
    // @ts-ignore - API 返回的字段是 items，不是 userDocStatusList
    if (response.items && response.items.length > 0) {
      // @ts-ignore
      const docStatus = response.items[0];
      // API 返回的字段是 status，不是 parsedStatus
      apiParsedStatus.value = docStatus.status;
      apiEmbeddingStatus.value = docStatus.embeddingStatus;
      
      // 同步状态到 docStore，保持数据一致性
      docStore.updateDocStatus(docStatus.status, docStatus.embeddingStatus);
      
      // 如果状态未完成，继续轮询
      if (needsPolling.value) {
        startPolling();
      } else {
        stopPolling();
        // 状态完成，停止旋转动画
        isReparsing.value = false;
      }
    }
  } catch (error) {
    console.error('Failed to fetch doc status:', error);
    // 出错时也继续轮询
    if (needsPolling.value) {
      startPolling();
    }
  }
};

// 开始轮询
const startPolling = () => {
  stopPolling(); // 先清除旧的定时器
  pollingTimer = setTimeout(() => {
    fetchDocStatus();
  }, POLLING_INTERVAL);
};

// 停止轮询
const stopPolling = () => {
  if (pollingTimer) {
    clearTimeout(pollingTimer);
    pollingTimer = null;
  }
};

// 监听 docId 变化，重新获取状态
watch(() => props.docId, (newDocId) => {
  stopPolling(); // 先停止旧的轮询
  if (newDocId) {
    fetchDocStatus();
  }
}, { immediate: true });

// 组件卸载时清理定时器
onUnmounted(() => {
  stopPolling();
});

const reparse = async () => {
  if (!props.pdfId) {
    message.error('文档ID无效');
    return;
  }

  if (isReparsing.value) {
    message.warning('该文档正在重新解析中，请稍候');
    return;
  }

  try {
    isReparsing.value = true;
    await reParsePaper({ pdfId: props.pdfId });
    message.success('重新解析请求已提交');
    // 重新获取状态
    await fetchDocStatus();
  } catch (error) {
    console.error('Reparse failed:', error);
    message.error('重新解析请求失败，请稍后重试');
    isReparsing.value = false;
  }
};
</script>

<style scoped lang="less">
.metadata-field-value {
  display: flex;
  align-items: center;
  gap: 8px; /* 给状态文本和链接之间增加一些间距 */
}

.reparse-icon {
  font-size: 16px;
  color: var(--site-theme-brand);
  cursor: pointer;
  margin-left: 8px;

  &:hover {
    opacity: 0.8;
  }

  &.reparse-spinning {
    animation: rotate 1s linear infinite;
  }
}

@keyframes rotate {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.status-icon {
  font-size: 16px;
}

.status-success {
  color: #52c41a;
}

.status-error {
  color: #ff4d4f;
}
</style>
