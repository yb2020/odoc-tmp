<template>
  <a-modal
    :visible="visible"
    :title="null"
    :footer="null"
    :closable="true"
    :maskClosable="false"
    width="424px"
    centered
    destroy-on-close
    class="ocr-modal"
    @cancel="handleCancel"
  >
    <div class="orc-modal-content">
      <div v-if="needVip">
        <img :src="PREMIUM_OCR">
        <div class="ocr-msg">
          {{ msg }}
        </div>
        <!-- <Trigger visible :page="PageType.note" :need-vip-type="vipNext[0]"> -->
        <Trigger
          visible
          :report-params="{
            pageType: PageType.note,
            elementName:
              vipStore.role.vipType > VipType.FREE
                ? ElementName.upperOcrTranslateLimitPopup
                : ElementName.upperOcrTranslatePopup,
          }"
          :need-vip-type="exception?.needVipType as VipTypePayable"
        >
          <!-- <template #default="slotProps">
            <a-button type="primary" @click="slotProps.onBuyVip">
              升级至{{ vipOrPremium }}
            </a-button>
          </template> -->
          <template #default="slotProps">
            <a-button
              type="primary"
              @click="slotProps.onBuyVip"
            >
              升级至{{
                vipStore.enabled ? $t('common.premium.senior') : 'VIP'
              }}
            </a-button>
          </template>
        </Trigger>
      </div>
      <div v-else>
        <img :src="OCR_LIMIT">
        <div>{{ msg }}</div>
      </div>
    </div>
  </a-modal>
</template>
<script setup lang="ts">
import PREMIUM_OCR from '@/assets/images/premium-ocr.svg';
import OCR_LIMIT from '@/assets/images/ocr-limit.svg';
import { computed, watch } from 'vue';
import {
  ElementName,
  PageType,
  reportElementImpression,
} from '~/src/api/report';
import { VipTypePayable, useVipStore, VipType } from '@common/stores/vip';
import Trigger from '@common/components/Premium/Trigger.vue';
// import useVipPayConfig from '@common/hooks/useVipPayConfig';
// const payStore = useVipPayConfig();

const vipStore = useVipStore();

const visible = computed(() => vipStore.showLimitDialog === 'ocr');

const msg = computed(() => vipStore.limitDialogMessage);
const exception = computed(() => vipStore.limitDialogProps.exception);
const needVip = computed(
  () => exception.value?.needVipType && !exception.value?.notCurrentUser
);

const handleCancel = () => {
  vipStore.hideVipLimitDialog();
};

watch(visible, (val) => {
  if (!val) {
    return;
  }

  reportElementImpression({
    page_type: PageType.note,
    type_parameter: 'none',
    element_name: needVip.value
      ? ElementName.upperOcrTranslatePopup
      : ElementName.upperOcrTranslateLimitPopup,
    element_parameter: 'none',
  });
});

// const vipNext = computed<
//   null | [VipType.STANDARD | VipType.PROFESSIONAL, string]
// >(() => {
//   const vipType = vipStore.roles[0]?.vipType ?? VipType.FREE;

//   switch (vipType) {
//     case VipType.FREE:
//     case VipType.UNRECOGNIZED:
//       return [VipType.STANDARD, '标准'];
//     case VipType.STANDARD:
//       return [VipType.PROFESSIONAL, '专业'];
//     case VipType.PROFESSIONAL:
//     case VipType.ENTERPRISE:
//       return null;
//   }
// });

// const vipOrPremium = computed(() => {
//   return vipStore.enabled ? `${vipNext.value?.[1]}版会员` : 'VIP';
// });

// const vipCount = computed(() => {
//   const vipType = vipNext.value?.[0] ?? VipType.STANDARD;

//   const count = payStore.data.value?.vipPayPrivilege.find(
//     (item) => item.vipType === vipType
//   )?.ocrTranslateCountLimit as number;

//   return count;
// });
</script>

<style lang="less">
.ocr-modal {
  border-radius: 2px;
  .ant-modal-body {
    padding: 0;
    background-color: var(--site-theme-bg-primary);
    // height: 262px;
    border-radius: 2px;
  }
  .ant-modal-close-x {
    height: 48px;
    width: 48px;
    margin-right: 8px;
    line-height: 48px;
    color: var(--site-theme-text-secondary);
  }
}
</style>
<style lang="less" scoped>
.orc-modal-content > div {
  padding-top: 44px;
  padding-bottom: 40px;
  display: flex;
  flex-direction: column;
  align-items: center;
  color: var(--site-theme-text-primary);
  img {
    display: block;
  }
  div {
    font-size: 18px;
    line-height: 28px;
  }
  span {
    font-size: 14px;
    line-height: 22px;
  }
  button {
    margin-top: 16px;
    height: 32px;
    width: 248px;
  }
}

.ocr-msg {
  word-wrap: break-word;
  word-break: break-all;
  padding: 0 24px;
}
</style>
