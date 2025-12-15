<script setup lang="ts">
import {
  PolishRewriteRequest,
  ZhToEnRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { BizType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/PolishFeedbackInfo';
import { inject, onUnmounted, ref } from 'vue';
import { ReloadOutlined } from '@ant-design/icons-vue';
import { REQUEST_SERVICE_NAME_AI_POLISH } from '@common/api/const';
import { RefundReasonScene } from '@common/api/aibeans';
import { useSentenceSSE } from '@common/hooks/aitools/useClientSSE';
import { convertMode } from '@common/hooks/aitools/usePollingParagraphSSE';
import { ParagraphMode } from '@common/components/AIParagraph/type';
import BeansRefund from '@common/components/AIBean/Refund.vue';
import Item from './components/Item.vue';
import { EventCode, ModuleType, PageType } from '@common/utils/report';
import { TabKeyType } from '../type';

const props = defineProps<{
  isDiff: boolean;
  modeValue?: ParagraphMode;
  type: TabKeyType.polish | TabKeyType.zhpolish | TabKeyType.translate;
}>();

const originalText = ref('');
const refundbale = inject('refundable') as boolean;

// const sentences = ref([
//   {
//     originalText: '原文1',
//     modifiedText: '修改1',
//   },
//   {
//     originalText: '原文1',
//     modifiedText: '修改1',
//   },
//   {
//     originalText: '原文1',
//     modifiedText: '修改1',
//   },
//   {
//     originalText: '原文1',
//     modifiedText: '修改1',
//   },
//   {
//     originalText: '原文1',
//     modifiedText: '修改1',
//   },
// ])

const isInProcess = defineModel('isInProcess', {
  type: Boolean,
  default: false,
});

const {
  sentences,
  startPollingSentences,
  loading,
  error,
  abortRequest,
  requestId,
  result,
} = props.modeValue
  ? useSentenceSSE<PolishRewriteRequest>(
      `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/rewriteSentence`,
      {
        polish_type: ModuleType.POLISH_REWRITE,
        page_type: PageType.POLISH,
        event: EventCode.readpaper_ai_polish_response_time,
      }
    )
  : useSentenceSSE<ZhToEnRequest>(
      `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/zhToEnSentence`,
      {
        polish_type: ModuleType.ZH_TO_EN,
        page_type: PageType.POLISH,
        event: EventCode.readpaper_ai_polish_response_time,
      }
    );

const startSentences = async (text: string) => {
  originalText.value = text;
  isInProcess.value = true;
  if (props.modeValue) {
    await startPollingSentences({
      text,
      mode: convertMode(props.modeValue),
      isZH: props.type === TabKeyType.zhpolish,
    }).promise;
  } else {
    await (
      startPollingSentences as (chat: ZhToEnRequest) => {
        promise: Promise<void>;
        resolve: () => void;
        reject: (err: Error) => void;
      }
    )({
      text,
    }).promise;
  }
};

const handleClear = () => {
  isInProcess.value = false;
  originalText.value = '';
  abortRequest();
};

const handleBack = () => {
  isInProcess.value = false;
};

const handleBeansRefunded = () => {
  if (result.value) {
    result.value.isRefundAiBean = true;
  }
};

onUnmounted(() => {
  abortRequest();
});

defineExpose({
  startSentences,
  clear() {
    abortRequest();
  },
});
</script>
<template>
  <div
    v-if="isInProcess"
    class="sentences flex flex-col"
  >
    <div
      class="border mx-6 flex border-rp-neutral-4 px-4 bg-rp-neutral-3 text-rp-neutral-6 rounded-2xl items-center mb-4"
    >
      <div class="flex-1 w-0 overflow-auto truncate py-4">
        {{ originalText }}
      </div>
      <div class="ml-6">
        <a-button
          v-if="error"
          shape="round"
          class="!bg-rp-neutral-3"
          @click="handleBack"
        >
          {{ $t('common.text.back') }}
        </a-button>
        <a-button
          v-else
          shape="round"
          class="!bg-rp-neutral-3"
          @click="handleClear"
        >
          {{ $t('common.aitools.clearTask') }}
        </a-button>
      </div>
    </div>
    <div class="flex-1 h-0 px-6 overflow-auto">
      <Item
        v-for="(s, idx) in sentences"
        :key="idx"
        :original-text="s.origin"
        :modified-text="s.target"
        :is-diff="isDiff"
        :mode-value="modeValue"
        :request-id="requestId"
        :unique-id="s.uniqueId"
        :type="type"
      />
      <div
        v-if="loading"
        class="text-center mt-4"
      >
        <a-spin />
      </div>
      <div
        v-if="error"
        class="text-center mt-4 text-rp-neutral-8 cursor-pointer"
      >
        <div @click="startSentences(originalText)">
          <ReloadOutlined class="mr-2" />{{ error.message || 'unknown error' }}
        </div>
      </div>
      <BeansRefund
        v-if="refundbale && result && !loading"
        :btn-props="{
          class: 'text-rp-neutral-8',
        }"
        :scene="
          modeValue
            ? RefundReasonScene.Paragraph
            : RefundReasonScene.Translation
        "
        :tid="requestId"
        :ttype="modeValue ? BizType.AI_POLISH : BizType.ZH_TRANSLATION"
        :ctime="result.createTime"
        :withdrawn="!!result.isRefundAiBean"
        @ok="handleBeansRefunded"
      />
    </div>
  </div>
</template>
<style lang="less" scoped>
.sentences {
  :deep(.btn-refund svg) {
    fill: theme('colors.rp-neutral-8');
  }
}
</style>
