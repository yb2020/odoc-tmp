<template>
  <section class="aibeans-buy pt-4 rounded-md bg-white text-rp-neutral-10">
    <p class="mb-0 text-sm flex gap-1.5 items-center justify-center">
      {{ $t('common.aibeans.balance') }}{{ $t('common.symbol.colon') }}
      <IconBean />
      <span class="text-2xl leading-9">{{ beans }}</span>
    </p>
    <header class="px-6 mb-4 flex flex-col items-center justify-center">
      <div
        v-if="!upgradable && upgradeInfo"
        class="mt-4"
      >
        <span class="mr-4 whitespace-nowrap">
          升级到{专业版/企业版}会员，获取更大的豆包容量
          <!-- {{
            $t('common.premium.upgradeBeanTip', [
              upgradeInfo.aiBeanCountEveryWeek,
            ])
          }} -->
        </span>
        <Trigger
          visible
          :btnProps="{ type: 'primary' }"
          :btnTxt="`${$t('common.premium.upgrade')}${$t(
            `common.premium.versions.version`
          )}`"
          :needVipType="
            PayType2VipType[upgradeInfo.vipPayType] as VipTypePayable
          "
          :report-params="{
            pageType: PageType.aibeans_pay,
            elementName: ElementName.update_version,
          }"
        />
      </div>
      <h1 class="mb-0 flex-1 text-inherit mt-4 text-lg font-medium text-center">
        {{ $t('common.aibeans.buyTt') }}
      </h1>
    </header>
    <div class="px-6 flex justify-center items-center pt-px pb-4">
      <template v-if="disabled">
        <div
          class="my-[130px] pt-px flex-1 text-rp-neutral-8 text-base text-center font-medium"
        >
          <img
            class="w-20 h-20 mb-2"
            src="@common/../assets/images/beans/pay-disabled.svg"
            alt="disabled"
          >
          <p class="m-0">
            {{ $t('common.aibeans.disabled') }}
          </p>
        </div>
      </template>
      <template v-else>
        <template v-if="upgradable && upgradeInfo">
          <PayWay :fill="preferences.color">
            <template #title>
              <IconPremium class="w-4 h-4 mr-2" />
              {{ $t('common.premium.upgradeTo') }}
              {{ $t(`common.premium.versions.${preferences.key}`) }}
            </template>
            <template #beans>
              {{ $t('common.premium.upgradePrice') }}
              <span class="text-[40px] leading-[1] text-rp-red-6">{{ formatPrice(upgradeInfo?.cost)
              }}{{ $t('common.units.yen') }}</span>
              <span class="text-sm leading-[1]">({{
                $t('common.premium.beansTip', [
                  upgradeInfo?.aiBeanCountEveryWeek,
                ])
              }})</span>
            </template>
            <template #desc>
              <span class="text-rp-neutral-6 text-xs">({{ $t('common.premium.beansTip2') }})</span>
            </template>
            <template #qrcode>
              <PayQrcode
                :disabled="disabled"
                :params="{
                  vipPayType: upgradeInfo.vipPayType,
                }"
                @paysucc="onPaySucc(upgradeInfo.vipPayType)"
              />
            </template>
          </PayWay>
          <span class="flex-grow text-center text-sm -mt-4">{{ $t('common.text.or') }}<br>/</span>
        </template>
        <PayWay
          :beans="+aiBeanPkg.aiBeanCount"
          :price="+aiBeanPkg.price"
          :expireDays="aiBeanPkg.expireDays"
        >
          <template
            v-if="disabledToBuy"
            #qrcode
          >
            <div
              class="relative mx-auto w-40 h-40 p-[9px] border border-solid border-rp-neutral-3"
            >
              <div
                class="h-full bg-rp-neutral-1 p-2 flex flex-col gap-2 justify-center items-center"
              >
                <span class="text-center">{{ buyConfirmInfo?.message }}</span>
                <div>
                  <a
                    href="/vip"
                    target="_blank"
                  ><a-button
                    v-if="buyConfirmInfo?.button === 'vip'"
                    type="primary"
                    size="small"
                  >
                    {{ $t('common.premium.btns.get')
                    }}{{ $t('common.premium.senior') }}
                  </a-button></a>
                  <a-button
                    v-if="buyConfirmInfo?.button === 'confirm'"
                    type="link"
                    @click="disabledToBuy = false"
                  >
                    {{ $t('common.aibeans.persistBuy') }}
                  </a-button>
                </div>
              </div>
            </div>
          </template>
          <template
            v-else
            #qrcode
          >
            <PayQrcode
              :disabled="disabled"
              :params="{
                vipPayType: VipPayType.AIBEAN,
              }"
              @paysucc="onPaySucc(VipPayType.AIBEAN)"
            />
          </template>
        </PayWay>
      </template>
    </div>
    <!-- <footer class="px-6 pt-3 pb-6 border-t border-rp-neutral-3 bg-rp-neutral-1">
      
      <ul class="list-none p-0 m-0 text-rp-neutral-8 text-xs leading-5">
        <li v-for="x in buyRules">{{ x }}</li>
        <li>
          {{ buyRules.length + 1 }}、{{ $t('common.aibeans.buyRules4LowVip') }}
        </li>
      </ul>
    </footer> -->
  </section>
</template>

<script setup lang="ts">
import { VipType } from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/user/VipUserInterface';
import {
  AiBeanPackage,
  UpgradeVipPro,
  VipPayType,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { PremiumVipPreferences } from '@common/stores/vip';
import { PayType2VipType, VipTypePayable } from '@common/api/vipPay';
import { computed, ref } from 'vue';
import { useI18n } from 'vue-i18n';
import Trigger from '@common/components/Premium/Trigger.vue';
import IconPremium from '~common/assets/images/premium/icon-premium.svg?component';
import IconBean from './Icon.vue';
import PayWay from './PayWay.vue';
import PayQrcode from '../Pay/Qrcode.vue';
import { formatPrice } from '../../utils/pay';
import { PageType, ElementName } from '../../utils/report';

const props = defineProps<{
  beans: number;
  disabled?: boolean;
  upgradable?: boolean;
  upgradeInfo?: UpgradeVipPro;
  aiBeanPkgInfo?: AiBeanPackage;
  reportParams?: any;
  buyConfirmInfo?: {
    disabled: boolean;
    message: string;
    button: 'none' | 'vip' | 'confirm';
  };
}>();
const emit = defineEmits<{
  (e: 'supplied', count: number): void;
  (e: 'upgraded', type: VipPayType): void;
}>();

const disabledToBuy = ref(!!props.buyConfirmInfo?.disabled);

const { t } = useI18n();

const preferences = computed(() => {
  const { vipPayType } = props.upgradeInfo || {};
  return vipPayType
    ? PremiumVipPreferences[PayType2VipType[vipPayType]]
    : VipType.UNRECOGNIZED;
});

const aiBeanPkg = computed(
  () =>
    props.aiBeanPkgInfo || {
      price: 3000,
      aiBeanCount: '200',
      expireDays: 30,
    }
);

const buyRules = computed(() => t('common.aibeans.buyRules').split('\n'));

const onPaySucc = (type: VipPayType) => {
  if (type === VipPayType.AIBEAN) {
    emit('supplied', +aiBeanPkg.value.aiBeanCount);
  } else {
    emit('upgraded', type);
  }
};
</script>
