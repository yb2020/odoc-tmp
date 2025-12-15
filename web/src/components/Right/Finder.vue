<template>
  <div
    class="finder js-interact-finder-ignore js-pdf-finder"
    :style="{ top: (top || 0) + 'px' }"
  >
    <div ref="finderWrapper" />
  </div>
</template>
<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import interact from 'interactjs';

defineProps<{ top?: number }>()

const finderWrapper = ref<HTMLDivElement>()

onMounted(() => {
  if (finderWrapper.value) {
    interact(finderWrapper.value).resizable({
      edges: { left: true },
      listeners: {
        start: (event) => {
          (event.target as HTMLElement).style.transition = 'none'
        },
        move: function (event) {
          let { x, y } = event.target.dataset;

          x = (parseFloat(x) || 0) + event.deltaRect.left;
          y = (parseFloat(y) || 0) + event.deltaRect.top;

          if (event.rect.width < 200) {
            return;
          }

          Object.assign(event.target.style, {
            width: `${event.rect.width}px`,
          });

          Object.assign(event.target.dataset, { x, y });
        },
        end: (event) => {
          const target = event.target as HTMLElement;
          target.style.transition = 'width 0.3s cubic-bezier(0.7, 0.3, 0.1, 1)';
        },
      },
    });
  }
})

const getFinderWrapper = () => {
  return finderWrapper.value
}

defineExpose({
  getFinderWrapper,
})

</script>
<style lang="less" scoped>
.finder {
  position: relative;
  top: 0;
  right: 0;
  height: 100%;
  z-index: 10;
}
</style>
