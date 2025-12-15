import {
  GetScanQRCodeReq,
  PayStatus,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { ComputedRef, MaybeRef, Ref, computed, ref, toRef, watch } from 'vue';
import { useRequest } from 'ahooks-vue';
import { getPayInfo, getQrCodeParams } from '@common/api/vipPay';
import { useI18n } from 'vue-i18n';

export function isWaiting(status?: PayStatus) {
  return [PayStatus.PAY_WAITING].includes(status!);
}

export function usePlaceOrder(
  params: MaybeRef<GetScanQRCodeReq>,
  disabled: Ref<boolean> | ComputedRef<boolean>
) {
  const { t } = useI18n();

  return useRequest(
    () => {
      if (disabled.value) {
        throw new Error(t('common.aibeans.disabled') as string);
      }

      return getQrCodeParams(toRef(params).value);
    },
    {
      manual: true,
    }
  );
}

export function useOrderInfo(
  orderId: Ref<string | undefined>,
  autoPolling = ref(true)
) {
  const isFirst = ref(true);
  const needPolling = ref(true);
  const ready = computed(
    () =>
      typeof document !== 'undefined' && needPolling.value && !!orderId.value
  );

  watch(orderId, () => (needPolling.value = true));

  return useRequest(
    async () => {
      try {
        const res = await getPayInfo({
          preOrderId: orderId.value!,
        });

        needPolling.value =
          [PayStatus.PAY_PRE, PayStatus.PAY_WAITING].includes(res.status) &&
          autoPolling.value;

        return res;
      } catch (e) {
        needPolling.value = !isFirst.value && autoPolling.value;

        throw e;
      } finally {
        isFirst.value = false;
      }
    },
    {
      ready,
      pollingInterval: 5000,
      pollingSinceLastFinished: true,
      pollingWhenHidden: false,
    }
  );
}
