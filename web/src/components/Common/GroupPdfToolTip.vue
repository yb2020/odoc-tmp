<template>
  <span class="js-pdf-group-info group-info-icon">
    <a-tooltip placement="bottomRight">
      <template #title>{{ $t('tasks.pdfTip2') }}</template>
      {{ $t('tasks.pdfTip1') }} <question-circle-outlined />
    </a-tooltip>
  </span>
</template>
<script setup lang="ts">
import { QuestionCircleOutlined } from '@ant-design/icons-vue';
</script>
<style lang="less" scoped>
.group-info-icon {
  color: #262625;
  font-size: 13px;
}
</style>
