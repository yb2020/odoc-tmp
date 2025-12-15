<template>
  <Tooltip
    :title="$t('message.markdownTip')"
    placement="bottomLeft"
    overlay-class-name="markdown-tooltip"
    :overlay-style="{
      whiteSpace: 'pre-line',
    }"
  >
    <div class="markdown-tip">
      <i class="aiknowledge-icon icon-markdown" />
    </div>
  </Tooltip>
</template>

<script setup lang="ts">
import { Tooltip } from 'ant-design-vue';
</script>

<style scoped lang="postcss">
.markdown-tip {
  height: 24px;
  width: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 2px;
  &:hover {
    background-color: var(--site-theme-bg-hover);
  }
  & > i {
    height: 18px;
    color: var(--site-theme-text-primary);
  }
}
</style>

<style lang="postcss">
.markdown-tooltip {
  .ant-tooltip-inner {
    background-color: #ffffff;
    font-size: 12px;
    color: #4e5969;
  }
  .ant-tooltip-arrow-content {
    background-color: #ffffff;
  }
}
</style>
