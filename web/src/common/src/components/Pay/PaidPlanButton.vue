<template>
  <a 
    href="javascript:void(0)" 
    :class="['paid-plan-button', { 'loading': isLoading }]"
    @click.prevent="handlePayment"
  >
    {{ isLoading ? loadingText : buttonTextComputed }}
    <span v-if="errorMessage" class="error-tooltip">{{ errorMessage }}</span>
  </a>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { handlePaymentClick, type PaymentOptions } from '@/api/stripe/payment';
import { useI18n } from 'vue-i18n';

const { t } = useI18n();

const props = defineProps({
  buttonText: {
    type: String,
    default: ''
  },
  loadingText: {
    type: String,
    default: 'paying...'
  },
  orderType: {
    type: Number,
    default: 2
  },
  numberCount: {
    type: Number,
    default: 1
  },
  openInNewTab: {
    type: Boolean,
    default: true
  }
});

const isLoading = ref(false);
const errorMessage = ref<string | null>(null);

// 如果没有提供按钮文本，则使用国际化文本
const buttonTextComputed = computed(() => {
  return props.buttonText || t('common.pay.button');
});

// 处理支付按钮点击
const handlePayment = async () => {
  const paymentOptions: PaymentOptions = {
    orderType: props.orderType,
    numberCount: props.numberCount,
    openInNewTab: props.openInNewTab
  };
  
  await handlePaymentClick(paymentOptions, isLoading, errorMessage);
};
</script>

<style scoped>
.paid-plan-button:hover {
  background-color: var(--site-theme-brand-dark);
}

.paid-plan-button.loading {
  opacity: 0.7;
  cursor: wait;
}

.error-tooltip {
  position: absolute;
  bottom: -40px;
  left: 50%;
  transform: translateX(-50%);
  background-color: var(--site-theme-error);
  color: var(--site-theme-text-inverse);
  padding: 5px 10px;
  border-radius: 4px;
  font-size: 12px;
  white-space: nowrap;
  z-index: 10;
  animation: fadeIn 0.3s;
}

.error-tooltip::before {
  content: '';
  position: absolute;
  top: -5px;
  left: 50%;
  transform: translateX(-50%);
  border-left: 5px solid transparent;
  border-right: 5px solid transparent;
  border-bottom: 5px solid var(--site-theme-error);
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
</style>