<template>
  <div 
    v-if="visible" 
    class="limit-dialog-overlay" 
    @click.self="close"
  >
    <div class="limit-dialog">
      <div class="limit-dialog-header">
        <button 
          class="close-button" 
          @click="close"
        >
          <CloseOutlined />
        </button>
      </div>
      <div class="limit-dialog-content">
        <div class="icon-wrapper">
          <ExclamationCircleFilled />
        </div>
        <p class="limit-message">{{ message }}</p>
      </div>
      <div class="limit-dialog-footer">
        <template v-if="isFree">
          <a-button 
            type="primary" 
            :loading="isLoading"
            @click="handleMembershipClick"
            class="upgrade-button"
          >
            {{ $t('common.dialogPay.dialog.upgradeMembership') }}
          </a-button>
        </template>
        <template v-else>
          <template v-if="buyCredit">
            <a-button 
              type="primary" 
              :loading="isLoading"
              @click="handleCreditClick"
              class="upgrade-button"
            >
              {{ $t('common.dialogPay.dialog.buyCredit') }}
            </a-button>
          </template>
          <template v-if="beyondVipLimit">
            <p class="membership-tip">
              {{ $t('common.dialogPay.dialog.membershipTip') }}
            </p>
          </template>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { Button as AButton } from 'ant-design-vue';
import { CloseOutlined, ExclamationCircleFilled } from '@ant-design/icons-vue';
import { handlePaymentClick, type PaymentOptions } from '@/api/stripe/payment';
import { useLimitDialogStore } from '@/common/src/stores/limitDialog';
import { useUserStore } from '@/common/src/stores/user';
const limitDialogStore = useLimitDialogStore();
const userStore = useUserStore();

// 如果是会员的情况下，状态码属于这个区间，则代表该会员已无法处理，需要联系客服
const beyondVipLimit = computed(() => limitDialogStore.code >= 4200 && limitDialogStore.code <= 4699);
// 是否会员
const isFree = computed(() => userStore.isFree());
// 如果是会员的情况下，状态码属于这个区间，则代表需要购买积分包
const buyCredit = computed(() => limitDialogStore.code >= 4007 && limitDialogStore.code <= 4008);

// 响应式状态
const isLoading = ref(false);
const errorMessage = ref<string | null>(null);

// 计算属性
const visible = computed(() => limitDialogStore.visible);
const message = computed(() => limitDialogStore.message);

// 方法
const close = () => {
  limitDialogStore.close();
};

const handleCreditClick = async () => {
  const paymentOptions: PaymentOptions = {
    orderType: 3, // 积分包订单类型
    numberCount: 1,
    openInNewTab: true,
  };
  
  await handlePaymentClick(paymentOptions, isLoading, errorMessage);
  // 支付完成后关闭对话框
  close();
};

const handleMembershipClick = async () => {
  const paymentOptions: PaymentOptions = {
    orderType: 2, // 会员订单类型
    numberCount: 1,
    openInNewTab: true,
  };
  
  await handlePaymentClick(paymentOptions, isLoading, errorMessage);
  // 支付完成后关闭对话框
  close();
};
</script>

<style scoped lang="less">
.limit-dialog-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.limit-dialog {
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
  display: flex;
  flex-direction: column;
  width: 480px;
  max-width: 90%;
  max-height: 80%;
  animation: dialogFadeIn 0.3s;
  position: relative;
  overflow: hidden;
}

.limit-dialog-header {
  padding: 16px 24px;
  display: flex;
  justify-content: flex-end;
  align-items: center;
  
  .close-button {
    background: none;
    border: none;
    font-size: 16px;
    cursor: pointer;
    color: rgba(0, 0, 0, 0.45);
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: 4px;
    
    &:hover {
      background-color: #f5f5f5;
      color: rgba(0, 0, 0, 0.75);
    }
  }
}

.limit-dialog-content {
  padding: 0 40px 24px;
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  
  .icon-wrapper {
    color: #faad14;
    font-size: 48px;
    margin-bottom: 20px;
  }
  
  .limit-message {
    margin: 0;
    line-height: 1.6;
    color: #333;
    font-size: 16px;
    font-weight: 500;
    text-align: center;
    margin-bottom: 24px;
  }
}

.membership-tip {
  margin-top: 16px;
  text-align: center;
  color: #666;
  font-size: 14px;
}

.limit-dialog-footer {
  padding: 0 40px 32px;
  display: flex;
  flex-direction: column;
  align-items: center;
  
  .upgrade-button {
    height: 44px;
    font-size: 16px;
    font-weight: 500;
    padding: 0 32px;
    border-radius: 6px;
    background: linear-gradient(90deg, #1890ff, #096dd9);
    border: none;
    box-shadow: 0 2px 8px rgba(24, 144, 255, 0.3);
    
    &:hover {
      background: linear-gradient(90deg, #40a9ff, #1890ff);
    }
  }
}

@keyframes dialogFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>