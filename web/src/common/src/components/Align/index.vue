<template>
  <div
    ref="targetEl"
    class="align-target"
  >
    <slot />
    <Teleport :to="to || 'body'">
      <div
        v-show="visible"
        ref="wrapperEl"
        class="align-wrapper w-fit"
        :class="alignClass"
      >
        <slot name="align" />
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import align from 'dom-align';
import { Align } from 'ant-design-vue/lib/dropdown/props';
import { useResizeObserver, useWindowSize } from '@vueuse/core';
import { computed, onMounted, onUpdated, ref, toRef } from 'vue';
import { debounce } from 'lodash-es';

const props = defineProps<{
  to?: HTMLElement;
  target?: HTMLElement | SVGElement;
  visible?: boolean;
  alignProps: Align;
  alignClass?: string;
}>();
const emit = defineEmits(['align']);

const visible = toRef(props, 'visible');
const alignProps = toRef(props, 'alignProps');
const targetEl = ref<undefined | HTMLElement>();
const wrapperEl = ref<undefined | HTMLElement>();

const realTarget = computed(() => props.target || targetEl.value);

const alignSelf = debounce(
  () => {
    if (visible.value && realTarget.value && wrapperEl.value) {
      const alignConfig: Align = {
        points: ['tl', 'tl'],
        offset: [0, 0],
        ...alignProps.value,
      };
      align(wrapperEl.value, realTarget.value, alignConfig);
      emit('align');
    }
  },
  100,
  {
    leading: false,
  }
);

onUpdated(alignSelf);
onMounted(() => {
  useResizeObserver(wrapperEl, alignSelf);
});
if (realTarget.value instanceof HTMLElement) {
  useResizeObserver(realTarget, alignSelf);
}
if (typeof document !== 'undefined') {
  useResizeObserver(document.body, alignSelf);
}

defineExpose({
  wrapperEl,
});
</script>
