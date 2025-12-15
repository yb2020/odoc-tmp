<script setup lang="ts">
defineProps<{
  text: string;
}>();
</script>
<template>
  <div class="tooltip flex justify-between items-center">
    <span class="text-rp-neutral-6 font-medium text mr-6">{{ text }}</span>
    <slot name="beans" />
  </div>
</template>
<style lang="less" scoped>
.text {
  font-size: 15px;
  line-height: 24px;
}
</style>
