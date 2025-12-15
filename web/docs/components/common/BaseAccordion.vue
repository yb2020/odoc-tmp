<template>
  <div class="base-accordion">
    <slot name="content" 
      :currentIndex="currentIndex" 
      :handleItemClick="handleItemClick"
      :pauseAutoRotation="pauseAutoRotation"
      :resumeAutoRotation="resumeAutoRotation">
    </slot>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount } from 'vue';

const props = defineProps({
  itemsCount: {
    type: Number,
    required: true
  },
  autoRotate: {
    type: Boolean,
    default: true
  },
  rotationInterval: {
    type: Number,
    default: 5000 // 默认5秒切换一次
  },
  initialIndex: {
    type: Number,
    default: 0
  }
});

const emit = defineEmits(['change']);

// 当前选中的索引
const currentIndex = ref(props.initialIndex);

// 自动轮播相关变量
let autoRotationTimer = null;
const isPaused = ref(false);

// 处理点击或悬停在项目上的事件
const handleItemClick = (idx) => {
  currentIndex.value = idx;
  emit('change', idx);
};

// 自动轮播到下一个项目
const rotateToNext = () => {
  if (isPaused.value) return;
  
  const nextIndex = (currentIndex.value + 1) % props.itemsCount;
  handleItemClick(nextIndex);
};

// 开始自动轮播
const startAutoRotation = () => {
  if (!props.autoRotate) return;
  
  // 清除可能存在的旧定时器
  if (autoRotationTimer) {
    clearInterval(autoRotationTimer);
  }
  
  // 设置新的定时器
  autoRotationTimer = setInterval(rotateToNext, props.rotationInterval);
};

// 暂停自动轮播
const pauseAutoRotation = () => {
  isPaused.value = true;
};

// 恢复自动轮播
const resumeAutoRotation = () => {
  isPaused.value = false;
};

// 组件挂载后初始化
onMounted(() => {
  // 开始自动轮播
  startAutoRotation();
});

// 组件卸载前清理定时器
onBeforeUnmount(() => {
  if (autoRotationTimer) {
    clearInterval(autoRotationTimer);
  }
});

// 暴露方法给父组件
defineExpose({
  currentIndex,
  handleItemClick,
  pauseAutoRotation,
  resumeAutoRotation
});
</script>

<style scoped>
.base-accordion {
  width: 100%;
}
</style>
