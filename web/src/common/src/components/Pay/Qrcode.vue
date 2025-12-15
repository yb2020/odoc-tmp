<template>
  <div
    class="relative mx-auto w-40 h-40 p-[9px] border border-solid border-rp-neutral-3"
  >
    <a-spin :spinning="loading" />
    <QrcodeVue
      v-if="orderId"
      :value="url"
      :size="140"
    />
    <div
      v-else
      class="h-full bg-rp-neutral-1 p-2 flex overflow-auto text-ellipsis break-all"
    >
      <span v-if="error">{{ error.message }}</span>
    </div>
    <p
      v-if="waiting"
      class="absolute w-full top-full left-0 m-0 text-xs text-center text-rp-red-6"
    >
      已扫描，等待付款中
    </p>
  </div>
</template>

<script setup lang="ts">
import {
  GetScanQRCodeReq,
  GetScanQRCodeResp,
  GetPayInfoResp,
  PayStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { computed, onMounted, toRef, watch } from 'vue';
import QrcodeVue from 'qrcode.vue';
import { getPayOrigin, VipPayType2ProductType } from '@common/api/vipPay';
import {
  isWaiting,
  usePlaceOrder,
  useOrderInfo,
} from '@common/hooks/useOrderInfo';
import { useVipStore } from '../../stores/vip';

const props = defineProps<{
  disabled?: boolean;
  params: GetScanQRCodeReq;
}>();
const emit = defineEmits<{
  (e: 'qrcode', o: GetScanQRCodeResp): void;
  (e: 'paysucc', s: PayStatus): void;
  (e: 'paytimeout', s: PayStatus): void;
}>();

const vipStore = useVipStore();
const params = toRef(props, 'params');
const disabled = toRef(props, 'disabled');

const {
  data: order,
  loading,
  error,
  run: onPlaceOrder,
} = usePlaceOrder(params, disabled);

watch(order, () => {
  if (order.value) {
    emit('qrcode', order.value);
  }
});

const orderId = computed(() => order.value?.preOrderId);
const { data: orderInfo } = useOrderInfo(orderId);

const url = computed(() => {
  const { vipPayType } = props.params;
  const params = new URLSearchParams({
    id: orderId.value || '',
    //@ts-ignore
    type: VipPayType2ProductType[vipPayType],
  });
  if (vipStore.env) {
    params.append('env', vipStore.env);
  }
  const origin = vipStore.payOrigin || getPayOrigin();
  return `${origin}/pay?${params.toString()}`;
});
const waiting = computed(() => isWaiting(orderInfo.value?.status));

onMounted(() => {
  onPlaceOrder();
});

watch(orderInfo, () => {
  const payStatus = orderInfo.value?.status;
  if (payStatus === PayStatus.PAY_SUCCESS) {
    emit('paysucc', payStatus);
  } else if (payStatus === PayStatus.PAY_TIMEOUT) {
    emit('paytimeout', payStatus);
    onPlaceOrder();
  }
});
</script>
