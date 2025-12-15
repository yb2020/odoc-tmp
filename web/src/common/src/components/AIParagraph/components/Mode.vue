<script setup lang="ts">
import { ParagraphMode } from '../type';
import { computed, watch } from 'vue';

const props = defineProps<{
  disabledMode?: ParagraphMode[];
}>();

const modeValue = defineModel('mode', {
  default: ParagraphMode.standard,
});

const modeList = computed(() => {
  return [
    ParagraphMode.improve,
    // ParagraphMode.standard,
    ParagraphMode.shorten,
    ParagraphMode.expand,
    ParagraphMode.simple,
    ParagraphMode.reduceSimilar,
  ].filter((mode) => {
    return !props.disabledMode?.includes(mode);
  });
});

watch(
  () => props.disabledMode,
  () => {
    if (props.disabledMode?.includes(modeValue.value)) {
      modeValue.value = modeList.value[0];
    }
  },
  {
    immediate: true,
  }
);
</script>
<template>
  <div>
    <span class="text-base font-medium pr-4 text-rp-neutral-10">{{
      $t('common.text.mode')
    }}</span>
    <a-radio-group
      v-model:value="modeValue"
      class="space-x-2"
    >
      <a-radio-button
        v-for="mode in modeList"
        :value="mode"
        class="!rounded-[20px]"
      >
        {{ $t(`common.aitools.polishModes.${mode}`) }}
      </a-radio-button>
    </a-radio-group>
  </div>
</template>
<style lang="postcss" scoped>
.ant-radio-button-wrapper {
  &::before {
    display: none !important;
  }
  border-left-width: 1px;
  border-radius: 20px !important;

  color: rgba(0, 0, 0, 0.65);
}
</style>
