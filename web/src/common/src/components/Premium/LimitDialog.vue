<template>
  <a-modal
    :visible="visible"
    :title="null"
    :footer="null"
    @cancel="handleCancel"
  >
    <div>
      <p class="vip-limit-dialog__content">
        {{ msg }}
      </p>
      <div class="vip-limit-dialog__button">
        <span class="inline-block text-rp-blue-6">
          <a-button
            v-if="btnInfo"
            type="link"
          >
            <a
              target="_blank"
              :href="btnInfo.url"
            >{{ btnInfo.text }}</a>
          </a-button>
        </span>
        <Trigger
          v-if="needVipType && !exception?.notCurrentUser"
          visible
          class="flex items-center"
          :page="PageType.NOTE"
          :need-vip-type="needVipType"
          :report-params="reportParams"
        >
          <template #default="slotProps">
            <Help placement="bottom">
              <a
                v-if="clicked"
                class="text-rp-blue-6"
                @click.prevent="handleBuy(slotProps, $event)"
              >{{ $t('common.premium.btns.reget') }}</a>
            </Help>
            <a-button
              type="primary"
              @click="clicked ? handleCancel() : handleBuy(slotProps, $event)"
            >
              <template v-if="clicked">
                {{
                  $t('common.premium.btns.payed')
                }}
              </template>
              <template v-else>
                {{ $t('common.premium.btns.get') }}
                {{ vipStore.enabled ? $t('common.premium.senior') : 'VIP' }}
              </template>
            </a-button>
          </template>
        </Trigger>
        <a-button
          v-else
          type="primary"
          @click="handleCancel"
        >
          {{
            $t('common.premium.btns.known')
          }}
        </a-button>
      </div>
    </div>
  </a-modal>
</template>
<script lang="ts" setup>
import { mapKeys, camelCase } from 'lodash';
import { ref, computed, watch } from 'vue';
import { PageType, reportElementImpression } from '@common/utils/report';
import { useVipStore, VipTypePayable } from '@common/stores/vip';
import Trigger from '@common/components/Premium/Trigger.vue';
import Help from '@common/components/Help.vue';

const vipStore = useVipStore();
const visible = computed(() => vipStore.showLimitDialog === 'vip');
const msg = computed(() => vipStore.limitDialogMessage);
const exception = computed(() => vipStore.limitDialogProps.exception);
const needVipType = computed(
  () => exception.value?.needVipType as VipTypePayable
);
const btnInfo = computed(() => vipStore.limitDialogProps.leftBtn);
const reportParams = computed(
  () =>
    mapKeys(vipStore.limitDialogProps.reportParams, (_, k) => camelCase(k)) as {
      pageType?: PageType;
      elementName: string;
      typeParameter?: string;
      elementParameter?: string;
    }
);
const clicked = ref(false);

const handleBuy = ({ onBuyVip }: any, e: MouseEvent) => {
  if (vipStore.enabled && !vipStore.payByDialog) {
    clicked.value = true;
  }
  onBuyVip(e);
};

const handleCancel = () => {
  vipStore.hideVipLimitDialog();
};

watch(visible, () => {
  if (visible.value) {
    clicked.value = false;
    const { reportParams } = vipStore.limitDialogProps;
    reportElementImpression({
      page_type: PageType.NOTE,
      element_name: 'uppper_unknow_popup',
      type_parameter: 'none',
      element_parameter: 'none',
      module_type: 'none',
      ...reportParams,
    });
  }
});
</script>
<style lang="less" scoped>
.vip-limit-dialog {
  &__content {
    font-size: 16px;
    line-height: 26px;
    font-weight: 500;
    margin-top: 16px;
  }
  &__button {
    margin-top: 24px;
    text-align: right;
    margin-bottom: -8px;
    .ant-btn-primary {
      margin-left: 16px;
    }
  }
}
</style>
