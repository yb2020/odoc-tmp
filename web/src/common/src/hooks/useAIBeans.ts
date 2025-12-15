import {
  BaseUseRequestOptions,
  useDocumentVisibility,
  useRequest,
} from 'ahooks-vue';
import { createSharedComposable } from '@vueuse/core';
import { h, ref, AppContext, onUnmounted, onMounted, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import { Modal, message } from 'ant-design-vue';
import type { ModalFunc } from 'ant-design-vue/lib/modal/Modal';
import {
  AiBeanCountResponse,
  BuyAiBeanResponse,
  NeedVipException,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import {
  RefundReasonScene,
  getBeansBuyWays,
  // getBeansInfo,
  getBeansRefundReasons,
} from '@common/api/aibeans';
import AIBeanPayIndex from '@common/components/AIBean/PayIndex.vue';
import { useVipStore } from '../stores/vip';
import { ResponseError } from '../api/type';
import { ERROR_CODE_BEANS_NOT_ENOUGH, ERROR_CODE_NEED_VIP } from '../api/const';
import { ReportParams, reportElementImpression } from '../utils/report';
import { getPayOrigin } from '../api/vipPay';
import { getAuthorizationCode, getAuthorizationUrl } from '../api/auth';
import { identity, pickBy } from 'lodash';
import { IS_ELECTRON } from '../utils/env';
import { ELECTRON_CHANNEL_EVENT_OPEN_URL, invoke } from '../electron/bridge';

export enum BeanScenes {
  COPILOT_ASK,
  POLISH,
  POLISH_ALL,
  REVIEW,
  TITLE_GEN,
  ABSTRACT_GEN,
  ABSTRACT_POLISH_1,
  ABSTRACT_POLISH_2,
  INTRO_POLISH_1,
  INTRO_POLISH_2,
  AI_TRANSLATE,
  REVIEWQS,
}

export enum BeanThresholds {
  BEANS_LEFT_LIMIT = 'beans_left_limit',
  VIP_DAYS_LEFT_LIMIT = 'vip_days_left_limit',
}

const DEFAULT_BEAN_UNIT: {
  [x in BeanScenes | BeanThresholds]: number;
} = {
  [BeanScenes.COPILOT_ASK]: 8,
  [BeanScenes.POLISH]: 2,
  [BeanScenes.POLISH_ALL]: 30,
  [BeanScenes.REVIEW]: 30,
  [BeanScenes.TITLE_GEN]: 6,
  [BeanScenes.ABSTRACT_GEN]: 6,
  [BeanScenes.ABSTRACT_POLISH_1]: 0,
  [BeanScenes.ABSTRACT_POLISH_2]: 8,
  [BeanScenes.INTRO_POLISH_1]: 3,
  [BeanScenes.INTRO_POLISH_2]: 5,
  [BeanScenes.REVIEWQS]: 2880,
  [BeanThresholds.BEANS_LEFT_LIMIT]: 50,
  [BeanThresholds.VIP_DAYS_LEFT_LIMIT]: 20,
  [BeanScenes.AI_TRANSLATE]: 1,
};

function useAIBeansRaw(
  appContext?: AppContext,
  opts?: Partial<BaseUseRequestOptions<AiBeanCountResponse>>
) {
  const beans = ref(-1);
  const threhold = ref(30);
  const confirmable = ref(true);
  const units = ref({
    ...DEFAULT_BEAN_UNIT,
  });
  const { t } = useI18n();
  const vipStore = useVipStore();
  const { data, run, refresh, ...rest } = useRequest(
    async () => {
      if (!vipStore.enabled) {
        return {} as AiBeanCountResponse;
      }

      // const res = await getBeansInfo();
      // beans.value = 0;
      // beans.value = +res.aiBeanCount;
      // units.value = {
      //   ...units.value,
      //   [BeanScenes.COPILOT_ASK]: +res.aiBeanCountCostEveryQuesion,
      //   [BeanScenes.POLISH]: +res.aiBeanCountCostEveryPolish,
      //   [BeanScenes.REVIEW]: +res.aiBeanCountCostEveryReview,
      //   [BeanScenes.TITLE_GEN]: +res.aiBeanCountCostEveryTitleGeneration,
      //   [BeanScenes.ABSTRACT_GEN]: +res.aiBeanCountCostEveryAbstractGeneration,
      //   [BeanScenes.ABSTRACT_POLISH_1]:
      //     +res.aiBeanCountCostEveryAbstractTruingFramework,
      //   [BeanScenes.ABSTRACT_POLISH_2]:
      //     +res.aiBeanCountCostEveryAbstractRewrite,
      //   [BeanThresholds.BEANS_LEFT_LIMIT]: +res.promptIfAiBeanLeft,
      //   [BeanThresholds.VIP_DAYS_LEFT_LIMIT]: +res.promptIfVipNotEnough,
      //   [BeanScenes.AI_TRANSLATE]: +res.aiBeanCountCostEveryPolish,
      // };
      // threhold.value = +res.promptIfAiBeanCostMoreThan;

      // return res;
      return {} as AiBeanCountResponse;
    },
    {
      manual: true,
      ...opts,
    }
  );
  watch(
    () => vipStore.enabled,
    (v) => v && run(),
    { immediate: true }
  );
  useDocumentVisibility({
    onVisible() {
      vipStore.enabled && refresh();
    },
  });

  const supplyBeans = (count = 0) => {
    if (count >= 0) {
      beans.value += count;
    }
    refresh();
  };

  const refundBeans = (scene: BeanScenes) => {
    if (typeof units.value[scene] === 'number') {
      beans.value += units.value[scene];
    }

    refresh();
  };

  const consumeBeans = async (scene: BeanScenes, price?: number) => {
    const cost = price || units.value[scene];
    if (vipStore.enabled && typeof cost === 'number') {
      const minus = () => {
        if (beans.value >= cost) {
          beans.value -= cost;
        }
      };
      if (cost >= threhold.value && confirmable.value) {
        return new Promise<boolean>((resolve) => {
          Modal.confirm({
            appContext,
            width: 480,
            icon: null,
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-ignore
            // footer: null,
            okText: t('common.text.ok'),
            closable: true,
            wrapClassName: 'aibeans-confirm',
            content: h(
              'p',
              {
                class: 'text-center text-rp-neutral-10 text-lg',
              },
              t('common.aibeans.confirmCost', [cost])
            ),
            onOk() {
              minus();
              resolve(true);
            },
            onCancel() {
              resolve(false);
            },
          });
        });
      }

      minus();
    }

    return Promise.resolve(true);
  };

  return {
    ...rest,
    data,
    units,
    beans,
    threhold,
    confirmable,
    refresh,
    supplyBeans,
    consumeBeans,
    refundBeans,
  };
}

export const useAIBeans = createSharedComposable(useAIBeansRaw);

function useAIBeansBuyRaw(
  appContext?: AppContext,
  props?: object,
  opts?: Partial<BaseUseRequestOptions<BuyAiBeanResponse>>
) {
  const instance = ref<ReturnType<ModalFunc>>();
  const vipStore = useVipStore();
  const { locale } = useI18n();
  const { beans, supplyBeans } = useAIBeans();
  const { data, run, ...rest } = useRequest(getBeansBuyWays, {
    ...opts,
    manual: true,
  });
  const onSucc = (beans?: number) => {
    supplyBeans(beans);
    hideBuyDialog();
  };

  if (!vipStore.payByDialog) {
    const handler = (
      e: MessageEvent<{
        type: string;
        beans: number;
      }>
    ) => {
      if (e.origin !== (vipStore.payOrigin || getPayOrigin())) {
        return;
      }
      switch (e.data.type) {
        case 'paySucc':
          onSucc(e.data.beans);
          break;
      }
    };
    onMounted(() => {
      window.addEventListener('message', handler);
    });
    onUnmounted(() => {
      window.removeEventListener('message', handler);
    });
  }

  const showBuyDialog = async () => {
    if (!vipStore.payByDialog) {
      const code = await getAuthorizationCode();
      const origin = vipStore.payOrigin || getPayOrigin();
      const s = new URLSearchParams(
        pickBy(
          {
            lang: locale.value,
            env: vipStore.env,
          },
          identity
        )
      ).toString();
      const url = getAuthorizationUrl(
        {
          env: vipStore.env,
          authorizationCode: code,
          redirectUrl: `${origin}/pdf-annotate/aibeans.html${s ? `?${s}` : ''}`,
        },
        origin
      );

      IS_ELECTRON
        ? invoke(ELECTRON_CHANNEL_EVENT_OPEN_URL, { url })
        : window.open(url, '_blank');
      return;
    }

    await run();
    instance.value = Modal.confirm({
      appContext,
      icon: null,
      // eslint-disable-next-line @typescript-eslint/ban-ts-comment
      // @ts-ignore
      footer: null,
      closable: true,
      width: 'fit-content',
      wrapClassName: 'aibeans-modal',
      content: h(AIBeanPayIndex, {
        beans: beans.value,
        disabled: !data.value?.paySwitch,
        upgradeInfo: data.value?.upgradeVipPro,
        aiBeanPkgInfo: data.value?.aiBeanPackage,
        onUpgraded: () => onSucc(),
        onSupplied: onSucc,
        ...props,
      }),
    });
  };
  const hideBuyDialog = () => {
    instance.value?.destroy();
  };
  const checkBeansErr = (err: unknown, reportParams: ReportParams) => {
    const e = err as ResponseError;
    if (e?.code === ERROR_CODE_NEED_VIP) {
      vipStore.showVipLimitDialog(e?.message, {
        exception: e?.extra as NeedVipException,
        reportParams,
      });
      return true;
    }
    if (e?.code === ERROR_CODE_BEANS_NOT_ENOUGH) {
      showBuyDialog();
      message.warn(e?.message);
      reportElementImpression(reportParams as any);
      return true;
    }

    return false;
  };

  return {
    ...rest,
    run,
    data,
    showBuyDialog,
    hideBuyDialog,
    checkBeansErr,
  };
}

export const useAIBeansBuy = createSharedComposable(useAIBeansBuyRaw);

export function useAIBeansRefundReasons(scene: RefundReasonScene) {
  return useRequest(
    async () => {
      const reasons = await getBeansRefundReasons(scene);

      return reasons.map(
        ({ name: text, internationalize: textEn, ...rest }) => {
          return {
            ...rest,
            text,
            textEn,
          };
        }
      );
    },
    {
      manual: true,
    }
  );
}

// export const useAIBeansRefundReasons = createSharedComposable(
//   useAIBeansRefundReasonsRaw
// )
