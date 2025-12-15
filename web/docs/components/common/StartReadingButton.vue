<template>
  <a href="/workbench" target="_self" class="start-reading-link">
    <button class="download-btn">
      <span>{{ buttonText }}</span>
      <img :src="getImageUrl('/images/icons/arrow.svg')" class="arrow-icon" alt="arrow" />
    </button>
  </a>
</template>

<script setup lang="ts">
import { getImageUrl } from '../../utils/contextPath';
import { defineProps } from 'vue';

defineProps({
  buttonText: {
    type: String,
    required: true,
  }
});
</script>

<style scoped>
.start-reading-link {
  text-decoration: none;
}

.download-btn {
  min-width: 140px;
  height: 44px;
  font-size: 14px;
  background-color: #5B7CE2;
  border: 1px solid #5B7CE2;
  border-radius: 6px;
  color: white;
  font-weight: 500;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.3s ease;
  cursor: pointer;
}

.download-btn:hover {
  background-color: #4A6AD7;
  border-color: #4A6AD7;
  box-shadow: 0 4px 12px rgba(91, 124, 226, 0.3);
  transform: translateY(-1px);
}

.arrow-icon {
  width: 16px;
  height: 16px;
  margin-left: 8px;
  transition: transform 0.3s ease;
}

.download-btn:hover .arrow-icon {
  transform: translateX(2px);
}
</style>
