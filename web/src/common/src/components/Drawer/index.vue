<script setup lang="ts">
import { onMounted, ref, watch } from 'vue';
import {
  CloseOutlined,
  DoubleLeftOutlined,
  DoubleRightOutlined,
} from '@ant-design/icons-vue';
import interact from 'interactjs';

const visible = defineModel('visible', { default: true });

const props = defineProps<{
  initialWidth: number;
  maxWidth: number;
  minWidth: number;
  placement: 'left' | 'right';
  title?: string;
  closable?: boolean;
  resizable?: boolean;
}>();
const emit = defineEmits<{
  (event: 'resize', w: number): void;
}>();

const currentWidth = ref(visible.value ? props.initialWidth : 0);

watch(
  () => props.initialWidth,
  () => {
    currentWidth.value = visible.value ? props.initialWidth : 0;
  }
);

const toggleVisible = (v = !visible.value) => {
  visible.value = v;
};

const wrapperRef = ref<HTMLElement | null>(null);

onMounted(() => {
  if (wrapperRef.value && props.resizable) {
    let maxWidth = window.innerWidth;
    interact(wrapperRef.value as HTMLElement).resizable({
      edges: props.placement === 'left' ? { right: true } : { left: true },
      ignoreFrom: '.ps__rail-y,.ps__rail-x,.js-interact-drag-ignore',
      listeners: {
        start: (event) => {
          (event.target as HTMLElement).style.transition = 'none';
          const edgeDrawerEle = document.body.querySelector(
            `.js-drawer-${props.placement === 'left' ? 'right' : 'left'}`
          ) as HTMLDivElement;
          if (edgeDrawerEle) {
            maxWidth = window.innerWidth - edgeDrawerEle.offsetWidth;
          }
        },
        move: function (event) {
          let { x, y } = event.target.dataset;

          x = (parseFloat(x) || 0) + event.deltaRect.left;
          y = (parseFloat(y) || 0) + event.deltaRect.top;

          if (
            event.rect.width < (props.minWidth || 200) ||
            event.rect.width > props.maxWidth
          ) {
            return;
          }

          const target = event.target as HTMLDivElement;

          target.style.width = `${event.rect.width}px`;

          const inner = target.children[0] as HTMLDivElement;

          inner.style.width = `${event.rect.width}px`;

          Object.assign(event.target.dataset, { x, y });
        },
        end: (event) => {
          const target = event.target as HTMLElement;
          target.style.transition = 'width 0.3s cubic-bezier(0.7, 0.3, 0.1, 1)';
          const curWidth = Math.min(
            Math.max(target.offsetWidth, props.minWidth || 200),
            maxWidth
          );
          target.style.width = `${curWidth}px`;
          currentWidth.value = curWidth;
          emit('resize', curWidth);
        },
      },
    });
  }
});
</script>
<template>
  <section class="group h-full relative">
    <div
      ref="wrapperRef"
      class="drawer h-full"
      :style="{ width: (visible ? currentWidth : 0) + 'px' }"
    >
      <slot />
    </div>

    <slot
      v-if="closable"
      name="icon"
      class="z-10 opacity-0 transition-opacity duration-300 hover:opacity-100 group-hover:opacity-100"
      :visible="visible"
      :toggleVisible="toggleVisible"
    >
      <aside
        class="z-10 opacity-0 transition-opacity duration-300 hover:opacity-100 group-hover:opacity-100 absolute top-1/2 cursor-pointer"
        :class="placement === 'left' ? 'left-0' : 'right-0'"
        @click="toggleVisible()"
      >
        <template v-if="!visible">
          <DoubleRightOutlined v-if="placement === 'left'" />
          <DoubleLeftOutlined v-else />
        </template>
        <CloseOutlined v-else />
      </aside>
    </slot>
  </section>
</template>
<style lang="less" scoped>
.drawer {
  transition: width 0.3s cubic-bezier(0.7, 0.3, 0.1, 1);
}
</style>
