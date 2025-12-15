<script setup lang="ts">
import {
  PolishRewriteRequest,
  ZhToEnRequest,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { inject, ref } from 'vue';
import copyTextToClipboard from 'copy-text-to-clipboard';
import { message } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';
import { REQUEST_SERVICE_NAME_AI_POLISH } from '@common/api/const';
import { RefundReasonScene } from '@common/api/aibeans';
import { useTextSSE } from '@common/hooks/aitools/useClientSSE';
import { convertMode } from '@common/hooks/aitools/usePollingParagraphSSE';
import { ParagraphMode } from '@common/components/AIParagraph/type';
import BeansRefund from '@common/components/AIBean/Refund.vue';
import Typing from '@common/components/TypingTxt/index.vue';
import DiffText from './DiffText.vue';
import wordsCount from 'words-count';
import {
  reportEvent,
  EventCode,
  ModuleType,
  PageType,
} from '@common/utils/report';
import { BizType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/PolishFeedbackInfo';
import { TabKeyType } from './type';

const { t } = useI18n();

const showDiffTxtHtml = ref(false);
const originText = ref('');

const textCount = ref(0);
const refundbale = inject('refundable') as boolean;

const props = defineProps<{
  type: TabKeyType.polish | TabKeyType.zhpolish | TabKeyType.translate;
  modeValue?: ParagraphMode;
  isDiff: boolean;
}>();

const {
  loading,
  isResolvedError: error,
  startPollingText,
  abortRequest,
  text: modifiedText,
  requestId,
  result,
} = props.type === TabKeyType.zhpolish || props.type === TabKeyType.polish
  ? useTextSSE<PolishRewriteRequest>(
      `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/rewrite`,
      {
        polish_type:
          props.type === TabKeyType.zhpolish
            ? ModuleType.POLISH_REWRITE_ZH
            : ModuleType.POLISH_REWRITE,
        page_type: PageType.POLISH,
        event: EventCode.readpaper_ai_polish_response_time,
      }
    )
  : useTextSSE<ZhToEnRequest>(
      `/api${REQUEST_SERVICE_NAME_AI_POLISH}/text/polish/zhToEn`,
      {
        polish_type: ModuleType.ZH_TO_EN,
        page_type: PageType.POLISH,
        event: EventCode.readpaper_ai_polish_response_time,
      }
    );

const startParagraph = async (text: string) => {
  showDiffTxtHtml.value = false;
  originText.value = text;
  try {
    if (props.modeValue) {
      await startPollingText({
        mode: convertMode(props.modeValue),
        text,
        isZH: props.type === TabKeyType.zhpolish,
      }).promise;
    } else {
      await (
        startPollingText as (chat: ZhToEnRequest) => {
          promise: Promise<void>;
          resolve: () => void;
          reject: (err: Error) => void;
        }
      )({
        text,
      }).promise;
    }
  } catch (error) {
    originText.value = '';
  }
};

const onTypingFinished = () => {
  if (!modifiedText.value) {
    return '';
  }
  textCount.value = wordsCount(modifiedText.value);
  showDiffTxtHtml.value = true;
};

const clear = () => {
  showDiffTxtHtml.value = false;
  originText.value = '';
  abortRequest();
};

const handleCopy = () => {
  copyTextToClipboard(modifiedText.value);
  message.success(t('common.tips.copied'));
  reportEvent(EventCode.readpaper_ai_polish_result_feedback_click, {
    element_name: 'copy',
    polish_type:
      props.type === TabKeyType.zhpolish
        ? ModuleType.POLISH_REWRITE_ZH
        : props.type === TabKeyType.polish
          ? ModuleType.POLISH_REWRITE
          : ModuleType.ZH_TO_EN,
    task_id: requestId.value,
    page_type: PageType.POLISH,
  });
};

const handleBeansRefunded = () => {
  if (result.value) {
    result.value.isRefundAiBean = true;
  }
};

defineExpose({
  startParagraph,
  clear,
});
</script>
<template>
  <div class="p-6 overflow-hidden">
    <div class="flex flex-col h-full">
      <div class="flex-1 overflow-auto">
        <div
          v-if="error"
          class="h-full flex items-center justify-center"
        >
          <a-result
            status="warning"
            title="请求失败"
            :sub-title="`Error: ${error.message || 'unknow error'}`"
          />
        </div>
        <Typing
          v-else-if="!showDiffTxtHtml && originText"
          :text="modifiedText"
          :is-pending="loading"
          @typing:finished="onTypingFinished"
        />
        <div v-else-if="showDiffTxtHtml">
          <DiffText
            v-if="isDiff"
            :original-text="originText"
            :modified-text="modifiedText"
            :diff-chars="type === TabKeyType.zhpolish"
          />
          <div v-else>
            {{ modifiedText }}
          </div>
        </div>
      </div>
      <div
        v-if="modifiedText && !loading"
        class="flex items-center mt-4"
      >
        <span class="flex-1 text-rp-neutral-6">{{ textCount }}
          {{ $t('common.text.words', textCount > 1 ? 2 : 1) }}</span>
        <span class="mr-6">
          <BeansRefund
            v-if="refundbale && result"
            :no-icon="true"
            :btn-props="{
              type: 'default',
              shape: 'round',
              size: 'middle',
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
        </span>
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
  </div>
</template>
