<template>
  <div class="err-section flex items-center justify-center overflow-auto">
    <a-result
      status="error"
      :title="title || 'Server Error'"
      :sub-title="subTitle"
    >
      <template #extra>
        <slot />
      </template>
    </a-result>
  </div>
</template>
<script setup lang="ts">
import { computed } from 'vue';

const props = defineProps<{ error: Error; scope?: string; title?: string }>();
const subTitle = computed(() => {
  if (props.scope) {
    return `[${props.scope}]: ${props.error.message}`;
  }
  return props.error.message;
});
</script>

<style lang="postcss" scoped>
.err-section {
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.65);
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(12px);
}
</style>
