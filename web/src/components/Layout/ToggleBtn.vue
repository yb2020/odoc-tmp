<template>
  <button class="cursor-pointer">
    <span
      class="flex items-center justify-center w-3 h-10 rounded-[20px] text-xs text-[#73716f]"
      :class="['bg-[#ebeae8]', 'hover:bg-[#d8d8d8]']"
    >
      <CaretLeftOutlined v-if="visible" />
      <CaretRightOutlined v-else />
    </span>
  </button>
</template>

<script setup lang="ts">
import { CaretLeftOutlined, CaretRightOutlined } from '@ant-design/icons-vue'

defineProps<{
  visible: boolean
}>()
</script>
