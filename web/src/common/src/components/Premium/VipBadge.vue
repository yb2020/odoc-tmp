<template>
  <div
    :class="[
      'flex items-center px-3 py-1 rounded-3xl text-white cursor-pointer',
      `vip-badge-${PremiumVipPreferences[vipType]?.key || 'free'}`,
    ]"
  >
    <svg
      class="mr-1 mt-px"
      width="16"
      height="16"
      viewBox="0 0 16 16"
      fill="none"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M0.97902 2C0.662096 2 0.471044 2.351 0.643099 2.61716L7.66399 13.4779C7.82159 13.7217 8.17823 13.7217 8.33583 13.4779L15.3567 2.61715C15.5288 2.351 15.3377 2 15.0208 2H11.2104C11.0746 2 10.9482 2.06884 10.8745 2.18284L7.99991 6.62955L5.12535 2.18284C5.05166 2.06884 4.92518 2 4.78943 2H0.97902Z"
        :fill="colorR"
      />
      <path
        fill-rule="evenodd"
        clip-rule="evenodd"
        d="M11.2019 2C11.0669 2 10.941 2.06813 10.8671 2.18117L5.58594 10.2617L7.66393 13.4779C7.82153 13.7217 8.17817 13.7217 8.33577 13.4779L15.3567 2.61715C15.5287 2.351 15.3377 2 15.0207 2H11.2019Z"
        :fill="colorL"
      />
      <defs>
        <linearGradient
          id="vipfree"
          x1="7.99988"
          y1="1"
          x2="7.99988"
          y2="13.6607"
          gradientUnits="userSpaceOnUse"
        >
          <stop stop-color="#1F71E0" />
          <stop
            offset="1"
            stop-color="#1F71E0"
            stop-opacity="0"
          />
        </linearGradient>
        <linearGradient
          id="vipother"
          x1="7.99988"
          y1="1"
          x2="7.99988"
          y2="13.6607"
          gradientUnits="userSpaceOnUse"
        >
          <stop stop-color="white" />
          <stop
            offset="1"
            stop-color="white"
            stop-opacity="0"
          />
        </linearGradient>
      </defs>
    </svg>
    <span class="text-sm">
      {{
        vipType === VipType.FREE
          ? `${$t('common.premium.btns.get')}${
            isChineseLanguage ? '' : ' '
          }${
            isLogin ? $t('common.premium.senior') : $t(`common.premium.vip`)
          }`
          : $t(`common.premium.versions.${PremiumVipPreferences[vipType]?.key || 'free'}`)
      }}
    </span>
  </div>
</template>
<script setup lang="ts">
import { VipType } from '@common/stores/vip';
import { computed } from 'vue';
import { PremiumVipPreferences } from '@common/components/Premium/types';
import { Language } from 'go-sea-proto/gen/ts/lang/Language';
import { useI18n } from 'vue-i18n';

const props = defineProps<{
  vipType: VipType;
  isLogin: boolean;
}>();

const { locale } = useI18n();

// 判断当前是否为中文语言
const isChineseLanguage = computed(() => {
  return locale.value.startsWith('zh');
});

const colorL = computed(() => {
  return props.vipType === VipType.FREE ? '#1F71E0' : 'white';
});

const colorR = computed(() => {
  return props.vipType === VipType.FREE ? 'url(#vipfree)' : 'url(#vipother)';
});
</script>
<style lang="less">
.vip-badge {
  &-free {
    color: theme('colors.rp-blue-6') !important;
    background: linear-gradient(90deg, #e8f5ff 0%, #bcddf9 100%);
  }

  &-standard {
    background: linear-gradient(97.28deg, #6aa9ec 6.55%, #438de6 90.32%);
  }

  &-profession {
    background: linear-gradient(97.28deg, #5a77eb 6.55%, #3a57e6 90.32%);
  }

  &-enterprise,
  &-outstanding {
    background: linear-gradient(97.28deg, #4b5977 6.55%, #24304b 90.32%);
  }
}
</style>
