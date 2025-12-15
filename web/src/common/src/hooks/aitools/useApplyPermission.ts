import {
  ApplyQuotaStatus,
  GetApplyQuotaResp,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/AiPolishText';
import { getApplyQuotaStatus, applyQuota } from '@common/api/revise';
import { useRequest } from 'ahooks-vue';
import { computed, ref } from 'vue';
import { createSharedComposable } from '@vueuse/core';
import { useVipStore } from '../../stores/vip';

const useApplyPermission = () => {
  const vipStore = useVipStore();
  const { data, run: getPermission } = useRequest(getApplyQuotaStatus, {
    manual: true,
  });
  const applyStatus = computed(
    () => data.value?.status ?? ApplyQuotaStatus.UNRECOGNIZED
  );

  const isValidApplyStatus = computed(() => {
    return applyStatus.value !== ApplyQuotaStatus.UNRECOGNIZED;
  });

  const applyingModalVisible = ref(false);

  const hasPermission = computed(() => {
    return applyStatus.value === ApplyQuotaStatus.ACCEPT || vipStore.enabled;
  });

  const initPermission = async (res?: GetApplyQuotaResp) => {
    if (res) {
      data.value = res;
      return res;
    }
    return getPermission();
  };

  const checkPermission = async () => {
    applyingModalVisible.value = !hasPermission.value;
    return !hasPermission.value;
  };

  const { loading: applyLoading, run: applyPermission } = useRequest(
    async () => {
      const res = await applyQuota();

      data.value = {
        ...data.value,
        status: res?.status,
      };
      applyingModalVisible.value = !(res?.status === ApplyQuotaStatus.ACCEPT);

      return res;
    },
    {
      manual: true,
    }
  );

  return {
    initPermission,
    checkPermission,
    applyPermission,
    hasPermission,
    applyStatus,
    applyLoading,
    applyingModalVisible,
    isValidApplyStatus,
  };
};

export default createSharedComposable(useApplyPermission);
