<template>
  <div
    v-show="visible"
    ref="selectRef"
    class="select w-24"
  >
    <PerfectScrollbar
      ref="psRef"
      :style="{ maxHeight: '180px' }"
      :options="{
        suppressScrollX: true,
      }"
    >
      <div
        v-for="item in options"
        v-if="options.length > 0"
        :key="item.value"
        :class="['option', item.shortcut ? 'option__hot' : '']"
        @click="onChange(item.value)"
      >
        {{ item.title
        }}<!--<span v-for="btn in item.shortcut?.split('+') || []" :key="btn" class="option-key">{{
          btn in shortcutTxtMap ? shortcutTxtMap[btn] : btn
        }}</span>-->
      </div>
      <slot name="options" />
    </PerfectScrollbar>
    <slot />
  </div>
</template>
<script lang="ts" setup>
import { ref, watch, nextTick } from 'vue';
import { onClickOutside } from '@vueuse/core';
// import { shortcutTxtMap } from '../../../store/shortcuts';

const props = defineProps<{
  visible: boolean;
  options: { title: string; value: string; shortcut?: string }[];
}>();

const emit = defineEmits<{
  (event: 'update:visible', val: boolean): void;
  (event: 'selectChange', val: string): void;
}>();

const selectRef = ref();

onClickOutside(selectRef, (event) => {
  if ((event.target as HTMLElement)?.closest('.js_select')) {
    return;
  }
  emit('update:visible', false);
});

const onChange = (value: string) => {
  emit('selectChange', value);
};

const psRef = ref();

watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      nextTick(() => {
        psRef.value?.update();
      });
    }
  }
);
</script>
<style lang="less" scoped>
.select {
  background: #fff;
  position: absolute;
  bottom: 50px;
  font-size: 12px;
  line-height: 18px;
  box-shadow: 1px 1px 4px 0px rgba(0, 0, 0, 0.3);
  border-radius: 2px;
  border: 1px solid #d9d9d9;
  cursor: default;
  left: 0;
  // transform: translate(-50%, 0);

  .option {
    padding: 4px 12px;
    text-align: center;
    cursor: pointer;
    color: inherit;

    &:hover {
      background: #f5f5f5;
    }

    &:first-child {
      margin-top: 0px;
    }

    // &__hot {
    //   white-space: nowrap;
    //   border-top: 1px solid #e5e6eb;
    //   padding-top: 10px;
    // }

    &-key {
      padding: 0 4px;
      margin-left: 4px;
      line-height: 20px;
      background: rgba(0, 0, 0, 0.04);
      border: 1px solid rgba(0, 0, 0, 0.15);
      border-radius: 2px;

      &:first-child {
        margin-left: 8px;
      }
    }
  }
}
</style>
