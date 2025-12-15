<template>
  <teleport to="body">
    <div
      v-if="vipStore.showPayDialog"
      class="vip-iframe-modal absolute z-[9999] w-full h-full top-0 left-0"
    >
      <a-spin
        v-if="loading"
        spinning
        class="!absolute w-full h-full top-0 left-0 !flex items-center justify-center"
      />
      <iframe
        width="100%"
        height="100%"
        :src="url"
        frameborder="0"
        @load="onLoad"
      />
    </div>
  </teleport>
</template>

<script setup lang="ts">
import { PayStatus } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { computed, onMounted, ref, watch } from 'vue';
import { useVipStore, VipType } from '@common/stores/vip';
import { getDomainOrigin } from '@common/utils/env';

const vipStore = useVipStore();
const props = computed(
  () =>
    vipStore.payDialogProps || {
      needVipType: VipType.STANDARD,
    }
);

const url = computed(
  () =>
    `${getDomainOrigin()}/dialog/vipPay#page=${
      props.value.reportParams?.page_type || 'note'
    }&type=${props.value.needVipType}`
);
const loading = ref(true);

enum Actions {
  PAY_SUCC = 'pay_succ',
  PAY_CANCEL = 'pay_cancel',
}

export interface Event<T = any> {
  action: Actions;
  payload: T;
}

onMounted(() => {
  window.addEventListener('message', (e: MessageEvent<Event>) => {
    switch (e.data?.action) {
      case Actions.PAY_SUCC:
        const { payStatus } = e.data.payload;
        vipStore.hideVipPayDialog(payStatus || PayStatus.PAY_SUCCESS);
        break;
      case Actions.PAY_CANCEL:
        vipStore.hideVipPayDialog(PayStatus.PAY_PRE);
        break;
      default:
    }
  });
});

watch(
  () => vipStore.showPayDialog,
  () => {
    if (!vipStore.showPayDialog) {
      onUnload();
    }
  }
);

const onLoad = () => (loading.value = false);
const onUnload = () => (loading.value = true);
</script>

<style lang="less">
.vip-iframe-modal {
  &.ant-modal {
    padding-bottom: 0;
  }
  .ant-modal-close {
    display: none;
  }
  .ant-modal-content,
  .ant-modal-body {
    padding: 0;
    height: 100%;
  }
}
</style>
