<script setup lang="ts">
import {
  PolishRewriteRequest,
  ZhToEnRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { REQUEST_SERVICE_NAME_AI_POLISH } from '@common/api/const';
import { message } from 'ant-design-vue';
import copyTextToClipboard from 'copy-text-to-clipboard';
import { onUnmounted, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import DiffText from '@common/components/AITools/DiffText.vue';
import Typing from '@common/components/TypingTxt/index.vue';
import { ParagraphMode } from '@common/components/AIParagraph/type';
import { useTextSSE } from '@common/hooks/aitools/useClientSSE';
import { convertMode } from '@common/hooks/aitools/usePollingParagraphSSE';
import {
  reportEvent,
  EventCode,
  ModuleType,
  PageType,
} from '@common/utils/report';
import { TabKeyType } from '../../type';

const props = defineProps<{
  originalText: string;
  modifiedText: string;
  isDiff: boolean;
  modeValue?: ParagraphMode;
  requestId: string;
  uniqueId: string;
  type: TabKeyType.polish | TabKeyType.zhpolish | TabKeyType.translate;
}>();

const _modifiedText = ref(props.modifiedText);

const { t } = useI18n();

const handleCopy = () => {
  copyTextToClipboard(_modifiedText.value);
  message.success(t('common.tips.copied'));
  reportEvent(EventCode.readpaper_ai_polish_result_feedback_click, {
    element_name: 'copy',
    polish_type:
      props.type === TabKeyType.zhpolish
        ? ModuleType.POLISH_REWRITE_ZH
        : props.type === TabKeyType.polish
          ? ModuleType.POLISH_REWRITE
          : ModuleType.ZH_TO_EN,
    task_id: requestId.value || props.requestId,
    page_type: PageType.POLISH,
  });
};

const {
  loading,
  startPollingText,
  abortRequest,
  text: retryModifiedText,
  requestId,
  isResolvedError: error,
} = useTextSSE<PolishRewriteRequest | ZhToEnRequest>(
  props.modeValue
    ? `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/rewrite`
    : `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/zhToEn`,
  {
    polish_type: props.modeValue
      ? props.type === TabKeyType.zhpolish
        ? ModuleType.POLISH_REWRITE_ZH
        : ModuleType.POLISH_REWRITE
      : ModuleType.ZH_TO_EN,
    page_type: PageType.POLISH,
    event: EventCode.readpaper_ai_polish_response_time,
  }
);

const onTypingFinished = () => {
  if (!retryModifiedText.value) {
    return '';
  }
  _modifiedText.value = retryModifiedText.value;
};

const handleRetry = async () => {
  _modifiedText.value = '';
  const text = props.originalText;
  if (props.modeValue) {
    await startPollingText(
      {
        mode: convertMode(props.modeValue),
        text,
        uniqueId: props.uniqueId,
        isZH: props.type === TabKeyType.zhpolish,
      },
      true
    );
  } else {
    await startPollingText(
      {
        text,
        uniqueId: props.uniqueId,
      },
      true
    );
  }
};

onUnmounted(() => {
  abortRequest();
});
</script>
<template>
  <div class="bg-white rounded-2xl mb-4 p-4">
    <div v-if="isDiff">
      <DiffText
        v-if="_modifiedText"
        :originalText="originalText"
        :modifiedText="_modifiedText"
        :diff-chars="type === TabKeyType.zhpolish"
      />
      <Typing
        v-else
        :text="retryModifiedText"
        :is-pending="loading"
        @typing:finished="onTypingFinished"
      />
    </div>
    <div
      v-else
      class="leading-[22px]"
    >
      <div class="mb-4 text-rp-neutral-8">
        <span class="bg-rp-neutral-2">{{ originalText }}</span>
      </div>
      <div
        v-if="_modifiedText"
        class="font-bold"
      >
        {{ _modifiedText }}
      </div>
      <Typing
        v-else
        :text="retryModifiedText"
        :is-pending="loading"
        @typing:finished="onTypingFinished"
      />
      <div
        v-if="error"
        class="text-center mt-4 text-rp-neutral-8 cursor-pointer"
      >
        {{ error.message || 'unknown error' }}
      </div>
    </div>
    <div class="flex justify-end space-x-4 mt-3">
      <a-button
        shape="round"
        :loading="loading"
        :disabled="loading"
        @click="handleRetry"
      >
        重新生成
      </a-button>
      <a-button
        type="primary"
        shape="round"
        @click="handleCopy"
      >
        {{
          $t('common.text.copy')
        }}
      </a-button>
    </div>
  </div>
</template>
