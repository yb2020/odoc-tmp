import { BaseUseRequestOptions, useRequest } from 'ahooks-vue';
import { ref } from 'vue';
import { GetVipPrivilegeResp } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { getVipPayConfig } from '../api/vipPay';

type Config = GetVipPrivilegeResp & {
  payVipInAiBeanPage?: boolean;
};

const DEFAULT_CONFIG = {
  maxDays: 350,
  allowPlanLowVipDays: 90,
};

export default function useVipPayConfig(
  interval = 5 * 60 * 1000,
  opts?: Partial<BaseUseRequestOptions<GetVipPrivilegeResp>>
) {
  const needPolling = ref(true);
  const { data, ...rest } = useRequest(
    async () => {
      let res = {} as Config;
      try {
        // 已弃用：/pay/public/getVipPrivilege接口已不再使用
        res = await getVipPayConfig(); // 使用模拟实现

        // 未开启的时候不轮询
        // 开启时轮询，关闭后页面1分钟内响应渲染
        needPolling.value = !!res?.paySwitch && interval > 0;
      } catch (e) {
        needPolling.value = false;

        // throw e;
      }

      return {
        ...DEFAULT_CONFIG,
        ...res,
      };
    },
    {
      ...opts,
      pollingInterval: interval,
      pollingSinceLastFinished: true,
      pollingWhenHidden: false,
    }
  );

  return {
    data,
    ...rest,
  };
}
