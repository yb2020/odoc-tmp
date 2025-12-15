<template>
  <p
    v-if="curAnswer.isOptimized"
    class="tag overflow-hidden absolute h-5 pl-2.5 pr-1.5 leading-5 top-0 right-0 text-xs rounded-bl-[20px] bg-rp-grass-8"
  >
    {{ $t('aiCopilot.optimized') }}
  </p>
  <div
    v-if="curAnswer.isOptimized"
    class="w-full text-rp-white-8 text-xs"
  >
    <p
      v-if="curAnswer.optimizedFeedbackResult"
      class="text-center"
    >
      {{ $t('aiCopilot.optimizedFeedback') }}
    </p>
    <div v-else>
      <div class="flex items-center mb-2">
        <span>{{ $t('aiCopilot.optimizedLabel') }}</span>
        <ul class="list-none p-0 mb-0 flex-1 flex flex-wrap gap-1.5">
          <li
            v-for="option in options"
            class="leading-1 flex items-center justify-center w-[52px] py-1.5 bg-rp-dark-8 hover:bg-rp-blue-6 text-rp-white-10 cursor-pointer"
            @click="onFeedback(option)"
          >
            {{ $t(`aiCopilot.${option}`) }}
          </li>
        </ul>
      </div>
      <p class="mb-0">
        {{ $t('aiCopilot.optimizedLabelTip') }}
      </p>
    </div>
  </div>
  <div
    v-else
    class="text-rp-white-10"
  >
    <BeansRefund
      v-if="curAnswer.needRefundAiBean !== false"
      :btn-props="{
        class: 'copilot-refund',
      }"
      :scene="RefundReasonScene.Copilot"
      :tid="curAnswer.id"
      :ctime="curAnswer.createDate"
      :withdrawn="!!curAnswer.isRefundAiBean"
      @ok="onOptimazeAnswer"
    />
  </div>
</template>

<script setup lang="ts">
import { AnswerInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiAssistReading';
import { Modal } from 'ant-design-vue';
import { useI18n } from 'vue-i18n';
import { RefundReasonScene } from '@common/api/aibeans';
import BeansRefund from '@common/components/AIBean/Refund.vue';
import { OptimizeFeedbackType, optimizeFeedback } from '~/src/api/copilot';
import { useRequest } from 'ahooks-vue';
import { ref } from 'vue';

const props = defineProps<{
  questionId: string;
  curAnswer: AnswerInfo;
}>();

const emit = defineEmits<{
  (e: 'refunded'): void;
  (e: 'feedbacked', x: string): void;
  (
    e: 'optimaze-answer',
    questionId: string,
    curAnswer: AnswerInfo,
    optimizeId: string
  ): void;
}>();

const i18n = useI18n();

const options = Object.values(OptimizeFeedbackType);

const onOptimazeAnswer = (optimizeId: string) => {
  emit('refunded');

  Modal.confirm({
    title: i18n.t('aiCopilot.confirmOptimize'),
    content: i18n.t('aiCopilot.confirmOptimizeTip'),
    okText: i18n.t('common.text.optimize'),
    onOk: async () => {
      emit('optimaze-answer', props.questionId, props.curAnswer, optimizeId);
    },
  });
};

const isFeedbacking = ref(false);
const { run: onFeedback } = useRequest(
  async (result: OptimizeFeedbackType) => {
    if (isFeedbacking.value) {
      return;
    }

    isFeedbacking.value = true;
    await optimizeFeedback({
      result,
      answerId: props.curAnswer.id,
    });

    emit('feedbacked', result);
    isFeedbacking.value = false;
  },
  {
    manual: true,
  }
);
</script>

<style lang="less">
.copilot-refund.btn-refund-txt.ant-btn {
  svg {
    fill: #ffffd8;
  }

  &:not([disabled]):hover {
    color: theme('colors.rp-darkblue-7');

    svg {
      fill: theme('colors.rp-darkblue-7');
    }
  }

  &[disabled] {
    color: theme('colors.rp-white-6');

    svg {
      fill: theme('colors.rp-white-6');
    }
  }
}
</style>
