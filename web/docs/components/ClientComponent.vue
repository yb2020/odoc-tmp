<script setup>
// 注意：这里不使用 TypeScript，避免可能的类型问题
import { ref, onMounted } from 'vue'

// 直接在 script setup 中定义变量和逻辑
const message = ref('这是一个纯客户端组件')
const count = ref(0)

// 所有逻辑都放在 onMounted 中，确保只在客户端执行
onMounted(() => {
  console.log('ClientComponent 已挂载')
  
  // 示例：每秒增加计数
  setInterval(() => {
    count.value++
  }, 1000)
})
</script>

<template>
  <div class="client-component">
    <h3>{{ message }}</h3>
    <p>计数器: {{ count }}</p>
  </div>
</template>

<style scoped>
.client-component {
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
  margin: 20px 0;
  background-color: #f9f9f9;
}
</style>
