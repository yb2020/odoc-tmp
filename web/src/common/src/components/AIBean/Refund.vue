<template>
  <component
    :is="btnProps?.type ? 'a-button' : 'span'"
    v-if="withdrawn"
    v-bind="
      btnProps?.type
        ? {
          size: 'small',
          ...btnProps,
          disabled: true,
        }
        : {}
    "
    :class="[
      'btn-refunded',
      btnProps?.type
        ? '!bg-transparent !text-rp-neutral-5 !border-rp-neutral-4'
        : '!text-rp-neutral-6',
      btnProps?.class,
    ]"
  >
    {{ $t('common.aibeans.refunded') }}
  </component>
  <a-tooltip
    v-else
    :trigger="[!isRefundable ? 'hover' : '']"
  >
    <template
      v-if="!isRefundable"
      #title
    >
      <div v-if="!data?.canCreditRefund">
        {{ $t(`common.aibeans.refundDisabled`) }}
      </div>
      <div v-else-if="isTimeout">
        {{ $t('common.aibeans.refundTimeout') }}
      </div>
      <div v-else-if="isUnlimitedCard">
        {{ $t(`common.aibeans.refundDisabledByCard`) }}
      </div>
    </template>
    <a-button
      v-bind="{
        type: 'text',
        size: 'small',
        ...btnProps,
      }"
      class="btn-refund"
      :class="[
        btnProps?.class,
        {
          'btn-refund-txt': !btnProps?.type,
        },
      ]"
      :disabled="!isRefundable"
      @click="onRefund"
    >
      <IconRefund
        v-if="!noIcon"
        class="inline align-middle -mt-0.5 mr-1"
        v-bind="iconProps"
      />{{ $t('common.aibeans.refund') }}
    </a-button>
  </a-tooltip>
  <RefundModal
    v-model:visible="isRefunding"
    :tid="tid"
    :ttype="ttype"
    :scene="scene"
    @ok="onRefunded"
  />
</template>

<script setup lang="ts">
import { BizType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/polish/PolishFeedbackInfo';
import { ActivityTypeEnum } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipActivitiesInfo';
import { SVGAttributes, computed, ref, watch, triggerRef } from 'vue';
// import { QuestionCircleOutlined } from '@ant-design/icons-vue'
import { ButtonProps } from 'ant-design-vue';
import IconRefund from '~common/assets/images/beans/refund.svg?component';
import { useAIBeans } from '../../hooks/useAIBeans';
import { useActivityCards } from '../../hooks/useActivityCards';
import { RefundReasonScene } from '../../api/aibeans';
import RefundModal from './RefundModal.vue';

const props = defineProps<{
  tid: string;
  ttype?: BizType;
  ctime: string;
  scene: RefundReasonScene;
  withdrawn: boolean;
  noIcon?: boolean;
  btnProps?: ButtonProps & { class?: string };
  iconProps?: SVGAttributes;
}>();

const emit = defineEmits<{
  (e: 'ok', x: string): void;
}>();

const { data, refresh } = useAIBeans();
const { cards } = useActivityCards();
const seconds = computed(
  () => parseInt(`${data.value?.refundTimeliness ?? 1800}`, 10) || 30 * 60
);
const tsvalid = computed(
  () => new Date(+props.ctime).getTime() + seconds.value * 1000
);

const isTimeout = computed(() => {
  return Date.now() > tsvalid.value;
});
const isUnlimitedCard = computed(
  () =>
    cards.value[ActivityTypeEnum.UN_LIMIT_POLISH_CARD] &&
    ![RefundReasonScene.Copilot].includes(props.scene)
);
const isRefundable = computed(
  () =>
    !props.withdrawn &&
    data.value?.canCreditRefund &&
    !isUnlimitedCard.value &&
    !isTimeout.value
);
const isRefunding = ref(false);

let timer: ReturnType<typeof setTimeout>;
watch(
  tsvalid,
  () => {
    if (!isTimeout.value) {
      timer = setTimeout(() => {
        triggerRef(tsvalid);
      }, tsvalid.value - Date.now());
    } else {
      clearInterval(timer);
    }
  },
  {
    immediate: true,
  }
);

const onRefund = () => {
  isRefunding.value = true;
};

const onRefunded = (optimizeId: string) => {
  isRefunding.value = false;

  refresh();
  emit('ok', optimizeId);
};
</script>

<style scoped lang="less">
.btn-refund-txt.ant-btn {
  @apply text-sm;
  @apply p-0;
  @apply transition-none;
  @apply bg-transparent;
  color: inherit;

  :deep(svg) {
    fill: theme('colors.rp-neutral-8');
  }

  &:not([disabled]):hover {
    color: theme('colors.rp-blue-6');

    :deep(svg) {
      fill: theme('colors.rp-blue-6');
    }
  }

  &[disabled] {
    color: theme('colors.rp-neutral-5');

    :deep(svg) {
      fill: theme('colors.rp-neutral-5');
    }

    :deep(span) {
      pointer-events: all;
    }
  }

  :deep(*) {
    span {
      @apply inline-flex;
      @apply items-center;
    }
  }
}

.btn-refund.ant-btn-default {
  color: theme('colors.rp-neutral-8');
  border-color: theme('colors.rp-neutral-4');
  background-color: transparent;

  &:not([disabled]):hover {
    color: theme('colors.rp-blue-6');
    border-color: theme('colors.rp-blue-6');
  }

  &[disabled] {
    cursor: not-allowed;
    color: theme('colors.rp-neutral-5');
    border-color: theme('colors.rp-neutral-4');
    background-color: transparent;
  }
}
</style>
