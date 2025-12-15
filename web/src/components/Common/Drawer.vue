<template>
  <div
    ref="wrapperRef"
    :class="`drawer js-drawer-${dir}`"
    :style="{ width: `${width}px` }"
  >
    <slot @visibleChange="handleVisibleChange" />
    <span
      v-if="holder"
      class="drawer-btn js-interact-finder-allow absolute left-0 top-1/2"
    >
      <HolderOutlined />
    </span>
    <a-tooltip placement="left">
      <template #title>
        {{ $t('viewer.expand') }}<component :is="getShortcutTxt(shortcut)" />
      </template>
      <ButtonVue
        v-if="wrapperRef"
        :ele="wrapperRef"
        :dir="dir"
        :visible="visible"
        @visible-change="handleVisibleChange"
      />
    </a-tooltip>
  </div>
</template>
<script setup lang="ts">
import ButtonVue from '../Common/Button.vue';
import { onMounted, ref, watch } from 'vue';
import { HolderOutlined } from '@ant-design/icons-vue';
import isMobile from 'is-mobile';
import interact from 'interactjs';
import { getShortcutTxt } from '../../store/shortcuts';

const props = defineProps<{
  dir: 'left' | 'right';
  minWidth: number;
  initialWidth: number;
  visible: boolean;
  shortcut?: string;
}>();
const wrapperRef = ref();

let validOffsetWidth = Math.max(props.initialWidth, props.minWidth);

const emit = defineEmits<{
  (event: 'visibleChange', visible: boolean): void;
  (event: 'widthChange', width: number): void;
  (event: 'update:visible', visible: boolean): void;
}>();

const visible = ref(props.visible);
const holder = isMobile({ tablet: true });

const width = ref(
  (() => {
    if (!visible.value) {
      return 0;
    }
    return validOffsetWidth;
  })()
);

const handleVisibleChange = (val?: boolean) => {
  visible.value = val === undefined ? !visible.value : val;
  emit('update:visible', visible.value);
  emit('visibleChange', visible.value);
};

watch(
  () => props.visible,
  (newVal, oldVal) => {
    if (newVal) {
      width.value = Math.max(validOffsetWidth, props.minWidth);
    } else {
      width.value = 0;
    }
    visible.value = newVal;
    emit('visibleChange', newVal);
  }
);

onMounted(() => {
  if (wrapperRef.value) {
    let maxWidth = window.innerWidth;
    interact(wrapperRef.value).resizable({
      edges: props.dir === 'left' ? { right: true } : { left: true },
      ignoreFrom: '.js-interact-finder-ignore',
      allowFrom: holder ? '.js-interact-finder-allow' : undefined,
      listeners: {
        start: (event) => {
          (event.target as HTMLElement).style.transition = 'none';
          const edgeDrawerEle = document.body.querySelector(
            `.js-drawer-${props.dir === 'left' ? 'right' : 'left'}`
          ) as HTMLDivElement;
          if (edgeDrawerEle) {
            maxWidth = window.innerWidth - edgeDrawerEle.offsetWidth;
          }
        },
        move: function (event) {
          let { x, y } = event.target.dataset;

          x = (parseFloat(x) || 0) + event.deltaRect.left;
          y = (parseFloat(y) || 0) + event.deltaRect.top;

          if (event.rect.width < (props.minWidth || 200)) {
            return;
          }

          const drawerEle = (event.target as HTMLDivElement).querySelector(
            '.js-drawer'
          ) as HTMLDivElement;

          if (drawerEle) {
            drawerEle.style.width = `${event.rect.width}px`;
          }

          event.target.style.width = `${event.rect.width}px`;

          Object.assign(event.target.dataset, { x, y });
        },
        end: (event) => {
          const target = event.target as HTMLElement;
          target.style.transition = 'width 0.3s cubic-bezier(0.7, 0.3, 0.1, 1)';
          const curWidth = Math.min(
            Math.max(target.offsetWidth, props.minWidth || 200),
            maxWidth
          );
          const drawerEle = (event.target as HTMLDivElement).querySelector(
            '.js-drawer'
          ) as HTMLDivElement;
          if (drawerEle) {
            drawerEle.style.width = `${curWidth}px`;
          }
          target.style.width = `${curWidth}px`;
          width.value = curWidth;
          validOffsetWidth = curWidth;
          emit('widthChange', curWidth);
        },
      },
    });
  }
});

defineExpose({
  handleVisibleChange,
});
</script>
<style lang="less" scoped>
.drawer {
  height: 100%;
  position: relative;
  transition: width 0.3s cubic-bezier(0.7, 0.3, 0.1, 1);
}
.drawer-btn {
  transform: translate(0, -50%);
}
</style>
