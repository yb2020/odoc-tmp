<template>
  <section class="px-6 flex gap-6 items-center">
    <PayQrcode
      class="mx-0"
      :params="{
        taskId: taskData?.id,
        reservedResourceId: taskData?.reservedResourceId || resId,
        vipPayType: VipPayType.PAPER_AI_REVIEW,
      }"
      @qrcode="onQrcode"
      @paysucc="onPaySucc"
    />

    <div class="flex-1 min-w-0 flex flex-col gap-2">
      <span
        class="w-fit bg-rp-red-6 text-xs tracking-tighter text-white py-1 px-2 rounded-tl-lg rounded-br-lg"
      >内测优惠</span>
      <span class="mb-1 text-xs text-rp-neutral-8">
        {{
          $t('common.aiReviewerQS.buyPatchTip', [
            config?.currentAiBeanPrice,
            taskData?.lockedBeanAmount,
          ])
        }}
      </span>
      <span class="text-rp-red-6 text-[40px] font-bold leading-[1]">￥<LoadingOutlined v-if="!price" /><span v-else>{{
        formatPrice(+price)
      }}</span></span>
      <PaySupport class="mt-4 pt-px !justify-start" />
      <Help>
        <p class="absolute right-0 bottom-0 m-0 text-rp-blue-6 cursor-help">
          {{ $t('common.pay.supportTip') }}
        </p>
      </Help>
    </div>
  </section>
</template>

<script setup lang="ts">
import { useRequest } from 'ahooks-vue';
import { getQSReviewConfig } from '@common/api/review';
import PayQrcode from '@common/components/Pay/Qrcode.vue';
import PaySupport from '@common/components/Pay/Support.vue';
import Help from '@common/components/Help.vue';
import {
  GetScanQRCodeResp,
  VipPayType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { ref } from 'vue';
import { formatPrice } from '../../utils/pay';
import { TaskInfo } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/review/AiReviewPaper';

defineProps<{
  resId?: string;
  taskData: TaskInfo;
}>();

const emit = defineEmits(['paysucc']);

const price = ref('');

const { data: config } = useRequest(async () => {
  const res = await getQSReviewConfig();

  return res;
}, {});

const onQrcode = (res: GetScanQRCodeResp) => {
  if (res.payTotalAmount) {
    price.value = res.payTotalAmount;
  }
};

const onPaySucc = () => emit('paysucc');
</script>

<style scoped></style>
