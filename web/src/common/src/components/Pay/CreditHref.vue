<template>
  <a 
    href="javascript:void(0)" 
    :class="['credit-display', { 'loading': isLoading, 'low-credit': isCreditLow }]"
    @click.prevent="handleAddOnCreditPayment"
  >
    <span class="credit-icon"><IconBean /></span>
    <span class="credit-value">{{ isLoading ? loadingText : creditValue }}</span>
    <!-- <span v-if="errorMessage" class="error-tooltip">{{ errorMessage }}</span> -->
  </a>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useUserStore } from '@/common/src/stores/user';
import { handlePaymentClick, type PaymentOptions } from '@/api/stripe/payment';
import IconBean from '@/common/src/components/AIBean/Icon.vue';
import { useI18n } from 'vue-i18n';
import { message } from 'ant-design-vue';
import { OrderType } from 'go-sea-proto/gen/ts/order/OrderApi';


const { t } = useI18n();

const props = defineProps({
  loadingText: {
    type: String,
    default: '...'
  },
  lowCreditThreshold: {
    type: Number,
    default: 10
  },
  // highCreditThreshold: {
  //   type: Number,
  //   default: 100000
  // },
  numberCount: {
    type: Number,
    default: 1
  },
  openInNewTab: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['click']);

const isLoading = ref(false);
const errorMessage = ref('');
const userStore = useUserStore();

// 获取用户积分
const creditValue = computed(() => {
  const credits = userStore.getTotalCredits();
  return credits !== undefined ? credits : 0;
});



// 判断积分是否不足
const isCreditLow = computed(() => {
  return creditValue.value < props.lowCreditThreshold;
});

// 判断积分是否过多
const isCreditHigh = computed(() => {
  return creditValue.value > props.highCreditThreshold;
});

// 处理支付按钮点击
const handleAddOnCreditPayment = async () => {
  // 根据用户是否为Pro，动态设置订单类型
  // Pro用户是类型订阅加油包，免费用户是类型订阅会员
  const orderType = userStore.isPro() ? OrderType.SUB_PRO_ADD_ON_CREDIT : OrderType.SUB_PRO;

  const paymentOptions: PaymentOptions = {
    orderType: orderType,
    numberCount: props.numberCount,
    openInNewTab: props.openInNewTab
  };
  
  // if (isCreditHigh.value) {
  //   message.warning(t('aiCopilot.creditHigh'));
  //   return;
  // }
  await handlePaymentClick(paymentOptions, isLoading, errorMessage);
};
</script>

<style scoped lang="less">
.credit-display {
  display: inline-flex;
  align-items: center;
  border-radius: 16px;
  transition: all 0.3s;
  text-decoration: none;
  color: inherit;
  
  &:hover {
    color: var(--site-theme-brand);
  }
  
  .credit-icon {
    margin-right: 4px;
    display: flex;
    align-items: center;
  }
  
  .credit-value {
    font-weight: 500;
  }
  
  &.loading {
    opacity: 0.7;
  }
  
  &.low-credit {
    color: var(--site-theme-danger-color);
    animation: pulse 2s infinite;
  }
  
  .error-tooltip {
    position: absolute;
    bottom: -20px;
    left: 50%;
    transform: translateX(-50%);
    background: rgba(0, 0, 0, 0.8);
    color: white;
    padding: 4px 8px;
    border-radius: 4px;
    font-size: 12px;
    white-space: nowrap;
  }
}

@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(255, 77, 79, 0.4);
  }
  70% {
    box-shadow: 0 0 0 6px rgba(255, 77, 79, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(255, 77, 79, 0);
  }
}
</style>