<template>
  <Library />
</template>

<script setup lang="ts">
import { defineAsyncComponent } from 'vue';
// 使用defineAsyncComponent正确处理懒加载组件
const Library = defineAsyncComponent(() => import('@/components/Library/Library.vue'))
import { useReportVisitDuration } from '@common/utils/report'
import { PageType } from '@common/utils/report'
import { LIBRARY_CONTAINER_CLASSNAME } from '@/components/Library/helper'

// 使用简化版的报告函数
useReportVisitDuration(
  () => '',
  () => ({
    page_type: 'library', // 直接使用字符串而不是 PageType.library
    type_parameter: 'none',
  }),
  () => {
    const container = document.getElementsByClassName(LIBRARY_CONTAINER_CLASSNAME)[0]
    return !container || !container.parentElement
  }
)
</script>

<style lang="less">
@import '@/assets/css/antd.less';

// 避免重写全局antd样式
.library-wrap {
  .ant-btn.ant-btn-default {
    color: #262625;
    border-color: #eceef2;
    &:focus {
      color: #262625;
      border-color: #eceef2;
    }
  }
}

.recommend-popover {
  .ant-popover-content {
    margin-left: 2px;
    .ant-popover-inner {
      .ant-popover-inner-content {
        padding: 8px;
        font-size: 12px;
        font-weight: 400;
        color: #4e5969;
        line-height: 18px;
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        .recommend-title {
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
          flex-shrink: 0;
        }
      }
    }
  }
}
</style>
