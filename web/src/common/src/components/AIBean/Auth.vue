<template>
  <Trigger
    visible
    class="flex items-center"
    :need-vip-type="lowestVipType"
    :report-params="{
      pageType: PageType.POLISH,
      elementName: ElementName.ai_bean,
    }"
  >
    <template #default="{ onBuyVip }">
      <AIBeans @click="onCheckVip($event, onBuyVip)" />
    </template>
  </Trigger>
</template>

<script setup lang="ts">
// import { computed } from 'vue'
import AIBeans from '@common/components/AIBean/index.vue';
import Trigger from '@common/components/Premium/Trigger.vue';
import { useVipStore, VipTypePayable } from '@common/stores/vip';
// import useVipPayConfig from '@common/hooks/useVipPayConfig'
import { ElementName, PageType } from '@common/utils/report';

const props = defineProps<{
  lowestVipType?: VipTypePayable;
}>();

const vipStore = useVipStore();
// const { data: vipConfig } = useVipPayConfig()

vipStore.fetchVipProfile();

// const lowestVipType = computed(() => {
//   const item = vipConfig.value?.vipPayPrivilege?.find((x) => x.polishEnable)

//   return (item?.vipType as VipTypePayable) ?? VipType.PROFESSIONAL
// })

const onCheckVip = (e: MouseEvent, onBuyVip: (e: MouseEvent) => any) => {
  if (props.lowestVipType && vipStore.role.vipType < props.lowestVipType) {
    e.preventDefault();
    onBuyVip(e);
  }
};
</script>
