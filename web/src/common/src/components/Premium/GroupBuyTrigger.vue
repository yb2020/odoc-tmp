<template>
  <template
    v-if="enabled || visible"
    class="vip-trigger flex justify-center"
  >
    <slot
      :enabled="enabled"
      :onBuyVip="onBuyVip"
    >
      <a-button
        v-bind="btnProps"
        class="vip-trigger-btn"
        shape="round"
        @click="onBuyVip"
      >
        {{ btnTxt || '立即开通' }}
      </a-button>
    </slot>
    <teleport
      v-if="isNavgating"
      to="body"
    >
      <a-spin />
    </teleport>
  </template>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue';
import { useVipStore, VipType, VipTypePayable } from '@common/stores/vip';
import {
  PayStatus,
  PaySourceType,
  VipPayType,
  GroupBuyItem,
} from '@idea/types-readpaper-proto/types/cn/edu/idea/cloud/aiKnowledge/vip/VipPayInfo';
import { message } from 'ant-design-vue';
import {
  getHostname,
  IS_MOBILE,
  IS_ELECTRON,
  isInOverseaseElectron,
  isDev,
} from '@common/utils/env';
import { getAuthorizationCode, getAuthorizationUrl } from '@common/api/auth';
import {
  VipType2ProductType,
  getPayOrigin,
  getPayURL,
  getQrCodeParams,
  GetGroupBuyScanQRCode,
} from '@common/api/vipPay';
import { identity, pickBy } from 'lodash';
import { useI18n } from 'vue-i18n';
import { ELECTRON_CHANNEL_EVENT_OPEN_URL, invoke } from '../../electron/bridge';
import { VipType2PayType } from './types';
import { useUrlSearchParams } from '@vueuse/core';
import { GroupProductItem } from '~common/src/components/Premium/gropbuy'; // 引入 gropbuy

const props = defineProps<{
  visible: boolean;
  src?: PaySourceType;
  btnTxt?: string;
  needVipType?: VipTypePayable;
  reportParams?: {
    pageType?: string;
    elementName: string;
  };
  btnProps?: object;
  disabled?: boolean;
  openSelf?: boolean;
  groupProductItem?: GroupProductItem[];
}>();
const emit = defineEmits(['paysucc']);

const vipStore = useVipStore();
const { channel } = useUrlSearchParams();
const { locale } = useI18n();
const enabled = !isInOverseaseElectron();
const isNavgating = ref(false);

watch(
  () => vipStore.payStatus,
  () => {
    if (vipStore.payStatus === PayStatus.PAY_SUCCESS) {
      // message.success('开通成功')
      emit('paysucc');
    }
  }
);

const onBuyVip = async (e: Event) => {
  // const target = props.openSelf ? '_self' : '_blank'
  // if (!vipStore.enabled) {
  //   window.open(`//${getHostname()}/vip`, target)
  // } else if (IS_MOBILE) {
  //   const type = props.needVipType || VipType.STANDARD

  //   const order = await getQrCodeParams({
  //     vipPayType: VipType2PayType[type as keyof typeof VipType2PayType],
  //   })

  //   const url = getPayURL(
  //     pickBy(
  //       {
  //         id: order.preOrderId,
  //         type: VipType2ProductType[type],
  //         src: `${props.src || ''}`,
  //         channel: typeof channel === 'string' ? channel : channel?.[0],
  //       },
  //       identity
  //     ),
  //     vipStore.env,
  //     vipStore.payOrigin || getPayOrigin()
  //   )

  //   isNavgating.value = true
  //   window.location.href = url
  // } else if (
  //   !vipStore.payByDialog &&
  //   window.location.origin !== (vipStore.payOrigin || getPayOrigin())
  // ) {
  //   const code = await getAuthorizationCode()
  //   const origin = vipStore.payOrigin || getPayOrigin()
  //   const params = new URLSearchParams(
  //     pickBy(
  //       {
  //         type: props.needVipType?.toString() || '',
  //         ...props.reportParams,
  //         env: vipStore.env,
  //         lang: locale.value,
  //       },
  //       identity
  //     )
  //   ).toString()
  //   const url = getAuthorizationUrl(
  //     {
  //       env: vipStore.env,
  //       authorizationCode: code,
  //       redirectUrl: `${origin}/home/vip.html${
  //         params.length ? '?' : ''
  //       }${params}`,
  //     },
  //     origin
  //   )

  //   e?.preventDefault()
  //   IS_ELECTRON
  //     ? invoke(ELECTRON_CHANNEL_EVENT_OPEN_URL, { url })
  //     : window.open(url, target)
  // } else if (enabled) {
  //   e.preventDefault()
  //   vipStore.showVipPayDialog({
  //     needVipType: props.needVipType,
  //     reportParams: {
  //       page_type: 'note',
  //       ...props.reportParams,
  //     },
  //   })
  // }

  // debugger
  // e.preventDefault()
  vipStore.showVipPayDialog({
    needVipType: props.needVipType,
    isGroupBuy: true,
    groupProductItem: props.groupProductItem,
    reportParams: {
      page_type: 'note',
      ...props.reportParams,
    },
  });
};
</script>
