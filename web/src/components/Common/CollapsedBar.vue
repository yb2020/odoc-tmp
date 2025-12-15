<template>
  <transition name="fade">
    <div
      v-show="isVisible && !isCollapsed"
      class="collapsed-bar absolute z-[100] py-1.5 px-3 left-1/2 -translate-x-1/2 rounded flex items-center"
      :class="class"
      :style="style"
    >
      <div
        class="btn-icon btn-collapse w-9 h-8 flex justify-center items-center cursor-pointer"
        @click="handleToggle(true)"
      >
        <slot name="icon-collapse">
          <double-left-outlined class="text-xs" />
        </slot>
      </div>
      <slot />
    </div>
  </transition>
  <div
    v-show="isCollapsed"
    class="collapsed-icon btn-icon btn-expand absolute z-[999] w-6 h-11 rounded-sm flex justify-center items-center cursor-pointer"
    :class="class"
    :style="iconStyle"
    @click="handleToggle(false)"
  >
    <slot name="icon-expand">
      <MenuUnfoldOutlined class="text-sm" />
    </slot>
  </div>
</template>

<script setup lang="ts">
import { useMouseInElement } from '@vueuse/core';
import { StyleValue, onMounted, ref, watch } from 'vue';
import { DoubleLeftOutlined, MenuUnfoldOutlined } from '@ant-design/icons-vue';
import { ElementClick, reportClick } from '~/src/api/report';

const props = defineProps<{
  collapsed?: boolean;
  class?: string;
  style?: StyleValue;
  iconStyle?: StyleValue;
  reportProps?: {
    element: string;
  };
}>();
const emit = defineEmits<{
  (e: 'toggled', x: boolean): void;
}>();

const toolbarRef = ref(null);
const isCollapsed = ref(props.collapsed ?? false);
const isVisible = ref(true);
const { elementX, elementY, isOutside } = useMouseInElement(toolbarRef);

watch(
  () => props.collapsed,
  () => {
    isCollapsed.value = props.collapsed ?? false;
  }
);

const handleToggle = (toggle: boolean) => {
  isCollapsed.value = toggle;

  if (props.reportProps?.element) {
    reportClick(
      props.reportProps.element as ElementClick,
      toggle ? 'on' : 'off'
    );
  }

  emit('toggled', toggle);
};

let timer: NodeJS.Timeout | null = null;
const handleMove = () => {
  isVisible.value = true;

  timer && clearTimeout(timer);
  if (isOutside.value) {
    timer = setTimeout(() => (isVisible.value = false), 3000);
  }
};

onMounted(() => {
  watch([elementX, elementY], handleMove);
});
</script>

<style scoped lang="less">
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.mobile-viewport {
  .collapsed-bar {
    display: none;
  }
}

.collapsed-bar {
  background-color: var(--site-theme-bg-primary);
  border: 1px solid var(--site-theme-divider);
  box-shadow: 3px 3px 12px 0px rgba(0, 0, 0, 0.3);
  
  .btn-collapse {
    &:hover {
      background-color: var(--site-theme-bg-hover);
    }
  }
}

.collapsed-icon {
  background-color: var(--site-theme-pdf-panel-collapsed);
  border: 1px solid var(--site-theme-divider);
  box-shadow: 3px 6px 12px rgba(0, 0, 0, 0.12), 1px 3px 6px rgba(0, 0, 0, 0.2);
  color: var(--site-theme-pdf-panel-text);
  
  &:hover {
    background-color: var(--site-theme-pdf-panel-collapsed-hover);
  }
}
</style>
