<template>
  <div class="dot-wrapper relative group">
    <span
      class="dot block w-2.5 h-2.5 rounded-full"
      :style="{ background: `${color}` }"
      @click.stop="() => {}"
    />
    <div
      v-if="editable"
      class="dot-option-list hidden group-hover:block absolute z-[3] w-full left-1/2 top-full py-2.5 rounded-[1px] border"
      :style="{
        backgroundColor: 'var(--site-theme-pdf-panel-secondary)',
        borderColor: 'var(--site-theme-divider)'
      }"
    >
      <div
        v-for="(option, key) in colors"
        :key="option.fill"
        class="dot-option h-[22px] flex items-center justify-center cursor-pointer dot-option-hover"
        @click.stop="handleClick(key)"
      >
        <span
          class="dot-option-dot w-2.5 h-2.5 rounded-full"
          :style="{ background: option.fill }"
        />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';
import { ColorStyle, styleMap } from '@/style/select';

const props = defineProps<{
  color: string;
  colorMap?: Record<number | string, Omit<ColorStyle, 'color'>>;
  editable?: boolean;
}>();

const emit = defineEmits<{
  (e: 'change', x: number | string): void;
}>();

const colors = computed(() => props.colorMap || styleMap);

const handleClick = async (key: number | string) => {
  emit('change', key);
};
</script>

<style lang="postcss" scoped>
.dot-option-list {
  transform: translate(-50%, 0);
}

.dot-option-hover:hover {
  background-color: var(--site-theme-background-hover);
}
</style>
