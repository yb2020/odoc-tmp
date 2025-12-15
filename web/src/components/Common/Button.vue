<template>
  <div
    class="button"
    :style="style"
    @click="handleClick"
  >
    <!-- <em :class="['iconfont', icon]"></em> -->
    <span
      v-if="dir === 'left'"
      class="text"
    >
      <!-- T<br />a<br />b<br />l<br />e -->
      <!-- 目<br />录 -->
      {{ $t('viewer.catalogue') }}
    </span>
    <template v-else>
      <DoubleLeftOutlined />
    </template>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { DoubleLeftOutlined } from '@ant-design/icons-vue';

const props = defineProps<{
  ele: HTMLDivElement;
  dir: 'right' | 'left';
  visible: boolean;
}>();

const icon = computed(() => {
  if (props.dir === 'right') {
    return props.visible
      ? 'iconxiangyousanjiaoxing'
      : 'iconxiangzuosanjiaoxing';
  }
  return props.visible ? 'iconxiangzuosanjiaoxing' : 'iconxiangyousanjiaoxing';
});

const style = ref(
  props.dir === 'right'
    ? {
        display: 'none',
        left: 0,
        transform: 'translate(-100%, -50%)',
      }
    : {
        display: 'none',
        right: 0,
        transform: 'translate(100%, -50%)',
      }
);

watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      style.value.display = 'none';
    } else {
      style.value.display = 'flex';
    }
  }
);

onMounted(() => {
  document.body.addEventListener('pointermove', (evt) => {
    if (props.visible) {
      style.value.display = 'none';
      return;
    }
    if (props.dir === 'right') {
      const offsetLeft = props.ele.offsetLeft;
      if (evt.pageX <= offsetLeft && evt.pageX >= offsetLeft - 60) {
        style.value.display = 'flex';
      } else {
        style.value.display = 'none';
      }
    } else {
      const offsetLeft = props.ele.offsetWidth;
      if (evt.pageX >= offsetLeft && evt.pageX <= offsetLeft + 60) {
        style.value.display = 'flex';
      } else {
        style.value.display = 'none';
      }
    }
  });
});

const emit = defineEmits<{ (event: 'visibleChange'): void }>();

const handleClick = () => {
  emit('visibleChange');
};
</script>
<style lang="less" scoped>
.button {
  top: 50%;
  position: absolute;
  width: 24px;
  // height: 48px;
  background: rgba(0, 0, 0, 0.35);
  // border: 1px solid #e4e7ed;
  text-align: center;

  position: absolute;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  display: flex;

  z-index: 1000;
  padding: 8px 2px;
  &:hover {
    background-color: rgba(0, 0, 0, 0.85);
  }
  .text {
    font-size: 12px;
    line-height: 16px;
    writing-mode: vertical-lr;
    letter-spacing: 2px;
  }
}
</style>
