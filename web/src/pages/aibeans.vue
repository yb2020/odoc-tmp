<template>
  <main class="main py-4 mx-24 text-rp-neutral-10">
    <p
      v-if="licensee"
      class="relative mb-4 py-px text-center text-base"
    >
      <span class="relative z-10 px-3 bg-rp-neutral-1">ReadPaper {{ $t('common.pay.authorize', [licensee]) }}</span>
      <span
        class="absolute top-1/2 left-0 w-full h-px -translate-y-1/2 bg-rp-neutral-3"
      />
    </p>
    <div class="flex justify-center mt-16">
      <div class="">
        <h2 class="text-lg font-bold text-rp-neutral-10">
          关于AI豆
        </h2>
        <p class="text-sm text-rp-neutral-7">
          AI豆是AI能量豆的简称，在您使用AI相关的功能，例如：AI辅读、AI润色、AI翻译等功能的时候，都需要消耗一定数量的AI豆
        </p>
        <h2 class="text-lg font-bold text-rp-neutral-10">
          各会员版本AI豆包
        </h2>
        <p>
          我们给各版本会员提供了豆包，不同的会员版本，豆包容量不同。在每周的周一0点0分，豆包会自动补满
        </p>
        <ul>
          <li>普通会员：无</li>
          <li>
            标准版会员：提供标准版AI豆包，容量为200个AI豆，当周有效，每周一自动补满
          </li>
          <li>
            专业版会员：提供专业版AI豆包，容量为600个AI豆，当周有效，每周一自动补满
          </li>
          <li>
            企业版会员：提供定制化AI豆包，容量通常>1200个AI豆，当周有效，每周一自动补满
          </li>
        </ul>
        <h2 class="text-lg font-bold text-rp-neutral-10">
          每周豆包容量不够用，怎么办？
        </h2>
        <p class="text-sm text-rp-neutral-7">
          我们提供如下方式进行扩容：
        </p>
        <ul>
          <li>
            参与固定活动：<a
              href="https://readpaper.com/home/mine"
              target="_blank"
            >&lt;AI大咖活动&gt;</a>，每周用得越多，下周获赠容量越多
          </li>
          <li>
            参与限时活动：我们会有不定期得活动，参与后也能获得相应活动周期内的扩容
            <a
              href="https://readpaper.com/home/mine"
              target="_blank"
            ><a-button>去看看</a-button></a>
          </li>
          <li v-if="data?.upgradeVipPro">
            升级到{专业版/企业版}会员，获取更大的豆包容量
            <Trigger
              visible
              :btnProps="{ type: 'primary' }"
              :btnTxt="`${$t('common.premium.upgrade')}${$t(
                `common.premium.versions.version`
              )}`"
              :needVipType="
                PayType2VipType[
                  data?.upgradeVipPro.vipPayType
                ] as VipTypePayable
              "
              :reportParams="{
                pageType: PageType.aibeans_pay,
                elementName: ElementName.updateVersion,
              }"
            />
          </li>
          <li>购买“单次临时扩容豆包”</li>
        </ul>
        <h2 class="text-lg font-bold text-rp-neutral-10">
          AI豆包注意事项
        </h2>
        <ul>
          <li>非转让性：AI豆包绑定于用户账户，不能转让给其他用户</li>
          <li>
            使用限制：如果标准/专业/企业版会员过期，仅仅有AI豆无法让你使用对应的AI功能
          </li>
          <li>
            豆包容量决定上限：当期剩余的AI豆，并不会提升豆包容量，鼓励每期都尽量用完
          </li>
        </ul>
      </div>
      <div class="w-[1px] bg-rp-neutral-3 h-auto my-4 mx-8" />
      <AIBeanPayIndex
        v-if="buyEnable.init"
        class=""
        :beans="beans"
        :disabled="!config?.paySwitch"
        :upgradable="config?.payVipInAiBeanPage"
        :upgradeInfo="data?.upgradeVipPro"
        :aiBeanPkgInfo="data?.aiBeanPackage"
        :buyConfirmInfo="buyEnable"
        @upgraded="onSucc"
        @supplied="onSucc"
      />
    </div>
  </main>
</template>

<script setup lang="ts">
import { useUrlSearchParams } from '@vueuse/core';
import langs from '@/locals/lang.json';
import AIBeanPayIndex from '@common/components/AIBean/PayIndex.vue';
import { useVipStore, VipType } from '@common/stores/vip';
import { useEnvInject } from '@common/hooks/useEnvInject';
import useVipPayConfig from '@common/hooks/useVipPayConfig';
import {
  useAIBeans,
  useAIBeansBuy,
  BeanThresholds,
} from '@common/hooks/useAIBeans';
import { computed, onMounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { message } from 'ant-design-vue';
import { ElementName, PageType, useReportVisitDuration } from '../api/report';
import axios from '@/api/axios';
import Trigger from '@common/components/Premium/Trigger.vue';
import { PayType2VipType, VipTypePayable } from '@common/api/vipPay';

const params = useUrlSearchParams<{
  // type: string;
  // page: string;
  // param: string;
  lang: string;
}>();

// @ts-ignore
useEnvInject([axios]);

useReportVisitDuration(
  () => '',
  () => ({
    page_type: PageType.aibeans_pay,
    type_parameter: 'none',
  })
);

const { t, locale } = useI18n();
if (langs.lang.includes(params.lang)) {
  locale.value = params.lang;
}

const vipStore = useVipStore();
vipStore.fetchVipConfig();
vipStore.fetchVipProfile();

const { data: config } = useVipPayConfig();
const licensee = computed(() => config.value?.licensee);
const { beans, units } = useAIBeans();
const { data, run } = useAIBeansBuy();

run();

const buyEnable = computed<{
  init: boolean;
  disabled: boolean;
  message: string;
  button: 'none' | 'vip' | 'confirm';
}>(() => {
  if (
    !Object.keys(vipStore.privileges).length ||
    !vipStore.inited ||
    beans.value === -1
  ) {
    return {
      init: false,
      disabled: false,
      message: '',
      button: 'none',
    };
  }
  const maxDaysRole = vipStore.roles.sort(
    (a, b) => b.vipLeftDays - a.vipLeftDays
  )[0];
  console.log(maxDaysRole, vipStore.roles);
  const confs = {
    ['noVip']: {
      message: t('common.aibeans.noVipTip'),
      button: 'vip',
    },
    ['lessDays']: {
      message: t('common.aibeans.daysLessTip'),
      button: 'confirm',
    },
    ['moreBeans']: {
      message: t('common.aibeans.leftMoreTip'),
      button: 'confirm',
    },
  };
  let c = null;

  if (
    maxDaysRole.vipType === VipType.UNRECOGNIZED ||
    maxDaysRole.vipType === VipType.FREE ||
    !maxDaysRole.vipLeftDays
  ) {
    c = confs['noVip'];
  } else if (
    units.value[BeanThresholds.VIP_DAYS_LEFT_LIMIT] > 0 &&
    maxDaysRole.vipLeftDays <= units.value[BeanThresholds.VIP_DAYS_LEFT_LIMIT]
  ) {
    c = confs['lessDays'];
  } else if (
    units.value[BeanThresholds.BEANS_LEFT_LIMIT] > 0 &&
    beans.value >= units.value[BeanThresholds.BEANS_LEFT_LIMIT]
  ) {
    c = confs['moreBeans'];
  }
  return {
    init: true,
    disabled: !!c,
    message: c?.message || '',
    button: (c?.button || 'none') as 'none' | 'vip' | 'confirm',
  };
});

const onSucc = () => {
  message.success(
    `${t('common.aibeans.tipSucc')}${t('common.symbol.comma')}${t(
      'common.tips.closeSelf',
      [3]
    )}`
  );
  (window.opener as Window)?.postMessage(
    {
      type: 'paySucc',
      beans: beans.value,
    },
    '*'
  );

  setTimeout(() => {
    window.close();
  }, 3000);
};

// debug
// onMounted(() => {
//   setTimeout(onSucc, 3000);
// });
</script>

<style lang="postcss">
#app {
  background-color: theme('colors.rp-neutral-1');
  height: 100%;
  overflow-y: auto;
}
</style>

<style lang="less" scoped>
:deep(*) {
  .aibeans-buy {
    @apply pt-6;
    @apply pb-10;
    border: 1px solid #dee7f2;

    & > * {
      @apply mx-4;
    }

    .aibeans-way {
      @apply w-80;
    }

    footer {
      @apply mx-10;
      padding: theme('spacing.6') theme('spacing.9');
      border: 1px solid theme('colors.rp-neutral-3');
    }
  }
}
</style>
