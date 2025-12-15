<template>
  <div
    class="pages"
    :style="{
      background: Object.keys(notes).length === 0 ? '' : '#2f3033',
    }"
  >
    <div
      v-for="(item, key) in notes"
      :key="key"
      :class="['page', `${key}` === `${activeKey}` ? 'active' : '']"
      @click="handleClick(`${key}`)"
    >
      {{ key }}
    </div>
  </div>
</template>

<script lang="ts" setup>
import scrollIntoView from 'scroll-into-view-if-needed';
import { onMounted } from 'vue';

const props = defineProps<{ notes: any; activeKey: number }>();
const emit = defineEmits(['setActiveKey']);

const handleClick = (pageNumber: string) => {
  const page = document.querySelector(
    `.notes-container .notes .page-${pageNumber}`
  )!;

  emit('setActiveKey', parseInt(pageNumber));

  scrollIntoView(page, {
    block: 'start',
    inline: 'start',
    behavior: 'smooth',
  });
};
</script>

<style lang="less" scoped>
.pages {
  width: 32px;
  height: 100%;

  border-radius: 0px 0px 8px 0px;

  flex: 0 0 auto;

  .page {
    width: 32px;
    height: 32px;

    font-family: 'Lato';
    font-style: normal;
    font-weight: 700;
    font-size: 12px;
    line-height: 18px;
    color: #a8adb3;
    opacity: 0.65;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;

    &.active {
      background: #4a4e52;

      color: #ffffff;
    }
  }
}
</style>
