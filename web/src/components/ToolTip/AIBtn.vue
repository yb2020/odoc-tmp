<template>
  <a-tooltip
    :mouse-enter-delay="1"
    placement="bottomRight"
    :destroyTooltipOnHide="false"
  >
    <template
      v-if="vipStore.enabled"
      #title
    >
      {{ $t('aiCopilot.ocrTip')
      }}{{
        vipStore.role.vipPrivilege?.readingAskQuestion === 0
          ? $t('common.aibeans.freeuse')
          : $t('common.premium.wordings.beanspertime', [
            vipStore.role.vipPrivilege?.readingAskQuestion,
          ])
      }}<br>
      <span class="inline-flex items-center">
        {{ $t('common.aibeans.balance') }}{{ $t('common.symbol.colon')
        }}<AIBeans @click.stop />
      </span>
    </template>
    <span
      class="ai-button"
      @click="$emit('ask-copilot')"
    >
      <img
        src="@/assets/images/ai-tooltip-button.svg"
        style="height: 1rem; width: auto"
      >
    </span>
  </a-tooltip>
</template>

<script setup lang="ts">
import AIBeans from '@common/components/AIBean/index.vue';
// import { useAIBeans } from '@common/hooks/useAIBeans';
import { useVipStore } from '@common/stores/vip';

// const { beans } = useAIBeans();
const vipStore = useVipStore();

defineEmits(['ask-copilot']);
</script>

<style lang="less" scoped>
.ai-button {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  &:hover {
    background: #52565a;
  }
  .aiknowledge-icon {
    font-size: 16px;
  }
}
</style>
