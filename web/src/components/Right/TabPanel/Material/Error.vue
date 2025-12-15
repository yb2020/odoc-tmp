<template>
  <div class="initial-container">
    <div
      class="detail"
      @click="handleRedo"
    >
      {{ error }}
      <RedoOutlined />
    </div>
  </div>
</template>
<script lang="ts" setup>
import { RedoOutlined } from '@ant-design/icons-vue';

const emit = defineEmits<{ (event: 'redo'): void }>()

const props = defineProps<{ error: string }>()

const handleRedo = () => {
  emit('redo');
};

</script>
<style lang="less" scoped>
.initial-container {
  text-align: center;
  margin: 20px auto;
  cursor: pointer;
  .detail {
    font-size: 13px;
    font-weight: 400;
    color: var(--site-theme-text-secondary);
  }
}
</style>
