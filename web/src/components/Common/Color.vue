<template>
  <div
    ref="colorRef"
    class="color-container"
  >
    <div
      v-for="(option, key) in styleMap"
      :key="option.color"
      :class="{ 'option-container': true, active: option.color === color }"
      @click.stop.prevent="handleClick(key)"
      @mouseenter="handleMouseEnter(key)"
    >
      <div
        :style="{ background: option.color }"
        class="color"
      />
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { useMouseInElement } from '@vueuse/core';
import { styleMap } from '@/style/select';

const emit = defineEmits<{
  (e: 'change', option: any): void;
  (e: 'enterChange', option: any): void;
}>();
const props = defineProps<{ color: string; handleOutside: () => void }>();

const handleClick = (key: keyof typeof styleMap) => {
  emit('change', Number(key));
};

const handleMouseEnter = (key: keyof typeof styleMap) => {
  emit('enterChange', Number(key));
};

const colorRef = ref();

const { isOutside } = useMouseInElement(colorRef);

watch(isOutside, (val) => {
  if (val) {
    props.handleOutside();
  }
});
</script>

<style lang="less" scoped>
.color-container {
  display: flex;
  align-items: center;
  flex-direction: column;
  color: #6fc169;
}

.option-container {
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  width: 28px;
  height: 28px;

  &:hover {
    background: #52565a;
  }
}

.color {
  width: 10px;
  height: 10px;
  border-radius: 50%;
}

.active {
  background: #52565a;
}
</style>
