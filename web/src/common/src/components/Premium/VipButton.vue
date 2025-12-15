<template>
  <Trigger
    v-if="userVipType !== null"
    visible
    class="flex items-center"
    :report-params="{
      elementName: ElementName.top_bar,
      ...reportParams,
    }"
  >
    <template #default="slotProps">
      <VipBadge
        :vip-type="userVipType"
        :is-login="isLogin"
        @click="openVip($event, slotProps.onBuyVip)"
      />
    </template>
  </Trigger>
</template>
<script setup lang="ts">
import Trigger from '@common/components/Premium/Trigger.vue';
import { useVipStore } from '@common/stores/vip';
import { computed } from 'vue';
import { ElementName } from '../../utils/report';
import VipBadge from './VipBadge.vue';

const props = defineProps<{
  isLogin: boolean;
  reportParams?: {
    pageType: string;
    elementName?: string;
  };
  showIfNotLogin?: boolean;
  openLogin?: () => void;
}>();

const vipStore = useVipStore();

vipStore.fetchVipProfile();

const userVipType = computed(() => vipStore.role.vipType);

const openVip = (evt: Event, onBuyVip: (e: Event) => void) => {
  if (!props.isLogin) {
    props.openLogin?.();
    return;
  }
  onBuyVip(evt);
};
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
